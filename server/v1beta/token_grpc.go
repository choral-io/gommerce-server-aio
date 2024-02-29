package v1beta

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	iam "github.com/choral-io/gommerce-protobuf-go/iam/v1beta"
	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-core/config"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/choral-io/gommerce-server-core/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/uptrace/bun"
)

func (p *FormPasswordLoginProvider) Validate(req *iam.CreateTokenRequest) error {
	if req.GetUsername().GetValue() == "" {
		return validator.NewError("username", "username is required when using form password login provider")
	}
	if req.GetPassword().GetValue() == "" {
		return validator.NewError("password", "password is required when using form password login provider")
	}
	return nil
}

type tokensServiceServer struct {
	iam.UnimplementedTokensServiceServer

	cfg config.TokenConfig
	bdb bun.IDB
	ts  secure.TokenStore
	lps map[string]LoginProvider
}

func NewTokensServiceServer(cfg config.TokenConfig, bdb bun.IDB, ts secure.TokenStore) iam.TokensServiceServer {
	s := &tokensServiceServer{
		cfg: cfg,
		bdb: bdb,
		ts:  ts,
		lps: make(map[string]LoginProvider, 2),
	}

	s.lps[LOGIN_PROVIDER_FORM_PASSWORD] = NewFormPasswordLoginProvider(bdb)
	s.lps[LOGIN_PROVIDER_SMS_OTP_CODE] = NewSMSOTPCodeLoginProvider()

	return s
}

func (s *tokensServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&iam.TokensService_ServiceDesc, s)
}

func (s *tokensServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return iam.RegisterTokensServiceHandler(ctx, mux, conn)
}

func (s *tokensServiceServer) Authorize(ctx context.Context, procedure string) error {
	if procedure == iam.TokensService_CreateToken_FullMethodName || procedure == iam.TokensService_RefreshToken_FullMethodName {
		return secure.Authorize(ctx, secure.AuthFuncAuthenticated, secure.AuthFuncRequireSchema(secure.AUTH_SCHEMA_BASIC))
	}
	return nil
}

func (s *tokensServiceServer) CreateToken(ctx context.Context, req *iam.CreateTokenRequest) (*iam.CreateTokenResponse, error) {
	now := time.Now()
	provider, ok := s.lps[strings.ToUpper(req.Provider)]
	if !ok {
		return nil, errors.New("login provider not found")
	}
	if provider == nil {
		return nil, errors.New("login provider not implemented")
	}
	if v, ok := provider.(interface {
		Validate(*iam.CreateTokenRequest) error
	}); ok {
		if err := v.Validate(req); err != nil {
			return nil, err
		}
	}
	var realm models.Realm
	if err := s.bdb.NewSelect().Model(&realm).Where(`"realm"."name" = ?`, req.Realm).Scan(ctx); err != nil {
		return nil, fmt.Errorf("realm with name %s not found", req.Realm)
	}
	login, err := provider.Login(ctx, realm.Id, req.Username.GetValue(), req.Password.GetValue(), req.IdToken.GetValue(), nil)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, validator.NewError("username", "username not found")
	} else if err != nil {
		return nil, validator.NewError("", err.Error())
	}
	if login.User == nil {
		return nil, validator.NewError("username", "username not found")
	}
	if login.User.ExpiresAt.Valid && !login.User.ExpiresAt.Time.After(time.Now()) {
		return nil, errors.New("user expired")
	}
	if login.User.Disabled {
		return nil, errors.New("user disabled")
	}
	if !login.User.Approved {
		return nil, errors.New("user not approved")
	}
	if !login.User.Verified {
		return nil, errors.New("user not verified")
	}
	if login.Disabled {
		return nil, errors.New("login disabled")
	}
	if login.ExpiresAt.Valid && !login.ExpiresAt.Time.After(time.Now()) {
		return nil, errors.New("login expired")
	}
	var roles []string
	if err := s.bdb.NewSelect().Model((*models.RoleUser)(nil)).Relation("Role", func(sq *bun.SelectQuery) *bun.SelectQuery {
		return sq.ExcludeColumn("*")
	}).Column("role.name").Where(`"role_user"."user_id" = ?`, login.User.Id).Scan(ctx, &roles); err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	scope := make([]string, len(roles))
	for i, r := range roles {
		scope[i] = "ROLE_" + strings.ToUpper(r)
	}
	uat, err := s.ts.Issue(secure.NewToken(secure.TOKEN_TYPE_BEARER, realm.Name, secure.IdentityFromContext(ctx).Token().Subject(), login.User.Id, scope), s.cfg.GetAccessTokenTTL())
	if err != nil {
		return nil, err
	}
	urt, err := s.ts.Issue(secure.NewToken(secure.TOKEN_TYPE_REFRESH, realm.Name, secure.IdentityFromContext(ctx).Token().Subject(), login.User.Id, scope), s.cfg.GetRefreshTokenTTL())
	if err != nil {
		return nil, err
	}
	if err := s.bdb.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().Model((*models.User)(nil)).
			Set(`"updated_at" = ?`, now).
			Set(`"first_login_time" = COALESCE("user"."first_login_time", ?)`, now).
			Set(`"last_active_time" = ?`, now).
			Where(`"id" = ?`, login.UserId).Exec(ctx)
		if err != nil {
			return errors.New("failed to update user")
		}
		if erc, err := res.RowsAffected(); err != nil {
			return errors.New("failed to update user")
		} else if erc == 0 {
			return errors.New("failed to update user")
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &iam.CreateTokenResponse{
		TokenType:    secure.TOKEN_TYPE_BEARER,
		ExpiresIn:    int32(time.Until(now.Add(s.cfg.GetAccessTokenTTL())).Seconds()),
		AccessToken:  uat,
		RefreshToken: urt,
	}, nil
}

func (s *tokensServiceServer) RefreshToken(ctx context.Context, req *iam.RefreshTokenRequest) (*iam.RefreshTokenResponse, error) {
	now := time.Now()
	uat, err := s.ts.Renew(req.GetRefreshToken(), s.cfg.GetAccessTokenTTL())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token: %s", err)
	}
	token, _ := s.ts.Verify(uat)
	urt, err := s.ts.Issue(secure.NewToken(secure.TOKEN_TYPE_REFRESH, token.Realm(), token.Client(), token.Subject(), token.Scope()), s.cfg.GetRefreshTokenTTL())
	if err != nil {
		return nil, err
	}
	return &iam.RefreshTokenResponse{
		TokenType:    secure.TOKEN_TYPE_BEARER,
		ExpiresIn:    int32(time.Until(now.Add(s.cfg.GetAccessTokenTTL())).Seconds()),
		AccessToken:  uat,
		RefreshToken: urt,
	}, nil
}
