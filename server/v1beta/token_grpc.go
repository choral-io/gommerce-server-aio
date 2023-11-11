package v1beta

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	iam "github.com/choral-io/gommerce-protobuf-go/iam/v1beta"
	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/uptrace/bun"
)

type tokensServiceServer struct {
	iam.UnimplementedTokensServiceServer

	ts  secure.TokenStore
	bdb bun.IDB
	lps map[string]LoginProvider
}

func NewTokensServiceServer(bdb bun.IDB, ts secure.TokenStore) iam.TokensServiceServer {
	s := &tokensServiceServer{
		bdb: bdb,
		ts:  ts,
		lps: make(map[string]LoginProvider, 2),
	}

	s.lps["FORM_PASSWORD"] = NewFormPasswordLoginProvider(bdb)
	s.lps["SMS_OTP_CODE"] = NewSMSOTPCodeLoginProvider()

	return s
}

func (s *tokensServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&iam.TokensService_ServiceDesc, s)
}

func (s *tokensServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return iam.RegisterTokensServiceHandler(ctx, mux, conn)
}

func (s *tokensServiceServer) CreateToken(ctx context.Context, req *iam.CreateTokenRequest) (*iam.CreateTokenResponse, error) {
	now := time.Now()
	ttl := 7 * 24 * 60 * 60 * time.Second
	provider, ok := s.lps[strings.ToUpper(req.Provider)]
	if !ok {
		return nil, errors.New("login provider not found")
	}
	if provider == nil {
		return nil, errors.New("login provider not implemented")
	}
	var realm models.Realm
	if err := s.bdb.NewSelect().Model(&realm).Where(`"realm"."name" = ?`, req.Realm).Scan(ctx); err != nil {
		return nil, fmt.Errorf("realm with name %s not found", req.Realm)
	}
	login, err := provider.Login(ctx, realm.Id, req.Username.GetValue(), req.Password.GetValue(), req.IdToken.GetValue(), nil)
	if err != nil {
		return nil, err
	}
	if login.User == nil {
		return nil, errors.New("user not found")
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
		scope[i] = "ROLE_" + r
	}
	accessToken, _ := s.ts.Issue(secure.NewToken(secure.TOKEN_TYPE_BEARER, realm.Name, secure.IdentityFromContext(ctx).Token().Subject(), login.User.Id, scope), ttl)
	refreshToken, _ := s.ts.Issue(secure.NewToken(secure.TOKEN_TYPE_REFRESH, realm.Name, secure.IdentityFromContext(ctx).Token().Subject(), login.User.Id, scope), ttl*10)
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
		TokenType:    "Bearer",
		ExpiresIn:    int32(time.Until(now.Add(ttl)).Seconds()),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *tokensServiceServer) RefreshToken(context.Context, *iam.RefreshTokenRequest) (*iam.RefreshTokenResponse, error) {
	return nil, errors.New("refresh token not implemented")
}