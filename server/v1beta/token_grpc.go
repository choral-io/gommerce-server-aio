package v1beta

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	iam "github.com/choral-io/gommerce-protobuf-go/iam/v1beta"
	"github.com/choral-io/gommerce-server-aio/data/repos"
	"github.com/choral-io/gommerce-server-core/config"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/choral-io/gommerce-server-core/validator"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p *formPasswordLoginProvider) Validate(req *iam.CreateTokenRequest) error {
	if req.GetUsername().GetValue() == "" {
		return validator.NewError("username", "username is required when using FORM_PASSWORD login provider")
	}
	if req.GetPassword().GetValue() == "" {
		return validator.NewError("password", "password is required when using FORM_PASSWORD login provider")
	}
	return nil
}

func (p *smsOTPCodeLoginProvider) Validate(req *iam.CreateTokenRequest) error {
	if req.GetUsername().GetValue() == "" {
		return validator.NewError("username", "username is required when using SMS_OTP_CODE login provider")
	}
	if req.GetPassword().GetValue() == "" {
		return validator.NewError("password", "password is required when using SMS_OTP_CODE login provider")
	}
	return nil
}

type tokensServiceServer struct {
	iam.UnimplementedTokensServiceServer

	cfg config.SecureTokenConfig
	drs repos.DataRepos
	ts  secure.TokenStore
	lps map[string]LoginProvider
}

func NewTokensServiceServer(cfg config.SecureTokenConfig, drs repos.DataRepos, ts secure.TokenStore) iam.TokensServiceServer {
	s := &tokensServiceServer{
		cfg: cfg,
		drs: drs,
		ts:  ts,
		lps: make(map[string]LoginProvider, 2),
	}

	s.lps[LoginProviderFormPassword] = NewFormPasswordLoginProvider(drs)
	s.lps[LoginProviderSmsOtpCode] = NewSMSOTPCodeLoginProvider()

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
		return secure.Authorize(ctx, secure.AuthFuncAuthenticated, secure.AuthFuncRequireSchema(secure.AuthSchemaBasic))
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
	realm, err := s.drs.Realms().FindOneByName(ctx, req.Realm)
	if err != nil {
		return nil, fmt.Errorf("realm with name %s not found: %w", req.Realm, err)
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
	roles, _, err := s.drs.Roles().FindNamesForUser(ctx, login.User.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	scope := make([]string, len(roles))
	for i, r := range roles {
		scope[i] = "ROLE_" + strings.ToUpper(r)
	}
	uat, err := s.ts.Issue(secure.NewToken(secure.TokenTypeBearer, realm.Name, secure.IdentityFromContext(ctx).Token().Subject(), login.User.Id, scope), s.cfg.GetAccessTokenTTL())
	if err != nil {
		return nil, err
	}
	urt, err := s.ts.Issue(secure.NewToken(secure.TokenTypeRefresh, realm.Name, secure.IdentityFromContext(ctx).Token().Subject(), login.User.Id, scope), s.cfg.GetRefreshTokenTTL())
	if err != nil {
		return nil, err
	}
	if err := s.drs.Users().UpdateLoginStatus(ctx, login.User.Id, now); err != nil {
		return nil, fmt.Errorf("failed to update login status: %w", err)
	}
	return &iam.CreateTokenResponse{
		TokenType:    secure.TokenTypeBearer,
		ExpiresIn:    int32(time.Until(now.Add(s.cfg.GetAccessTokenTTL())).Seconds()),
		AccessToken:  uat,
		RefreshToken: urt,
	}, nil
}

func (s *tokensServiceServer) RevokeToken(ctx context.Context, req *iam.RevokeTokenRequest) (*iam.RevokeTokenResponse, error) {
	splits := strings.SplitN(metadata.ExtractIncoming(ctx).Get(secure.AuthHeaderKey), " ", 2)
	if len(splits) == 2 && strings.EqualFold(splits[0], secure.AuthSchemaBearer) {
		if _, err := s.ts.Revoke(splits[1]); err != nil {
			return nil, err
		}
	}
	return &iam.RevokeTokenResponse{}, nil
}

func (s *tokensServiceServer) RefreshToken(_ context.Context, req *iam.RefreshTokenRequest) (*iam.RefreshTokenResponse, error) {
	now := time.Now()
	uat, err := s.ts.Renew(req.GetRefreshToken(), s.cfg.GetAccessTokenTTL())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token: %s", err)
	}
	token, _ := s.ts.Verify(uat)
	urt, err := s.ts.Issue(secure.NewToken(secure.TokenTypeRefresh, token.Realm(), token.Client(), token.Subject(), token.Scope()), s.cfg.GetRefreshTokenTTL())
	if err != nil {
		return nil, err
	}
	return &iam.RefreshTokenResponse{
		TokenType:    secure.TokenTypeBearer,
		ExpiresIn:    int32(time.Until(now.Add(s.cfg.GetAccessTokenTTL())).Seconds()),
		AccessToken:  uat,
		RefreshToken: urt,
	}, nil
}
