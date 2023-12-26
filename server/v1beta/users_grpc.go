package v1beta

import (
	"context"
	"database/sql"
	"errors"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-core/data"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	iam "github.com/choral-io/gommerce-protobuf-go/iam/v1beta"
	gender "github.com/choral-io/gommerce-protobuf-go/types/v1/gender"
	sqlpb "github.com/choral-io/gommerce-protobuf-go/types/v1/sqlpb"
)

type usersServiceServer struct {
	iam.UnimplementedUsersServiceServer

	bdb bun.IDB
}

func NewUsersServiceServer(bdb bun.IDB) iam.UsersServiceServer {
	return &usersServiceServer{bdb: bdb}
}

func (s *usersServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&iam.UsersService_ServiceDesc, s)
}

func (s *usersServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return iam.RegisterUsersServiceHandler(ctx, mux, conn)
}

func (s *usersServiceServer) Authorize(ctx context.Context, procedure string) error {
	if procedure == iam.UsersService_GetIdentity_FullMethodName {
		return secure.Authorize(ctx, secure.AuthFuncAuthenticated, secure.AuthFuncRequireSchema(secure.AUTH_SCHEMA_BEARER))
	}
	if procedure == iam.UsersService_ListUsers_FullMethodName {
		return secure.Authorize(ctx, secure.AuthFuncAuthenticated, secure.AuthFuncRequireRealm(REALM_ADMIN))
	}
	return nil
}

func (s *usersServiceServer) Register(ctx context.Context, req *iam.RegisterRequest) (*iam.RegisterResponse, error) {
	realm := &models.Realm{}
	if err := s.bdb.NewSelect().Model(realm).Where("", req.Realm).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.InvalidArgument, "realm %s not found", req.Realm)
		}
		return nil, status.Errorf(codes.Unknown, "error retrieving realm %s: %v", req.Realm, err)
	}
	if !realm.AllowRegistration() {
		return nil, status.Errorf(codes.PermissionDenied, "realm %s does not allow registration", req.Realm)
	}
	user := &models.User{
		RealmId:     realm.Id,
		Disabled:    false,
		Approved:    true,
		Verified:    true,
		Attributes:  map[string]string{},
		DisplayName: sqlpb.ToNullString(req.DisplayName),
		Gender:      gender.ToSqlNullString(req.Gender),
	}
	if user.DisplayName.Valid {
		user.Attributes["profile.display_name"] = user.DisplayName.String
	}
	if user.Gender.Valid {
		user.Attributes["profile.gender"] = user.Gender.String
	}
	login := &models.Login{
		Provider:   LOGIN_PROVIDER_FORM_PASSWORD,
		Identifier: req.Username,
		Metadata:   map[string]string{},
	}
	if hp, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost); err != nil {
		return nil, status.Errorf(codes.Unknown, "error hashing password: %v", err)
	} else {
		login.Credential = sql.NullString{Valid: true, String: string(hp)}
	}
	err := s.bdb.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(user).Exec(ctx); err != nil {
			return status.Errorf(codes.Unknown, "error creating user: %v", err)
		}
		login.UserId = user.Id
		if _, err := tx.NewInsert().Model(login).Exec(ctx); err != nil {
			return status.Errorf(codes.Unknown, "error creating login: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &iam.RegisterResponse{
		User: toUserPB(*user),
	}, nil
}

func (s *usersServiceServer) ListUsers(ctx context.Context, req *iam.ListUsersRequest) (*iam.ListUsersResponse, error) {
	var users []models.User
	query := s.bdb.NewSelect().Model(&users)
	total, err := query.Apply(data.WithPaging(req)).
		Relation("Realm", func(sq *bun.SelectQuery) *bun.SelectQuery { return sq.Column("name") }).
		Relation("Creator").
		Relation("Creator.Realm", func(sq *bun.SelectQuery) *bun.SelectQuery { return sq.Column("name") }).
		ScanAndCount(ctx)
	if err != nil {
		return nil, err
	}
	res := &iam.ListUsersResponse{
		Page:  req.Page,
		Size:  req.Size,
		Total: int64(total),
		Items: make([]*iam.User, len(users)),
	}
	for i, u := range users {
		res.Items[i] = toUserPB(u)
	}
	return res, nil
}

func (s *usersServiceServer) GetIdentity(ctx context.Context, req *iam.GetIdentityRequest) (*iam.GetIdentityResponse, error) {
	user := models.User{
		Id: secure.IdentityFromContext(ctx).Token().Subject(),
	}
	err := s.bdb.NewSelect().Model(&user).WherePK().
		Relation("Realm", func(sq *bun.SelectQuery) *bun.SelectQuery { return sq.Column("name") }).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &iam.GetIdentityResponse{
		User:  toUserPB(user),
		Scope: secure.IdentityFromContext(ctx).Token().Scope(),
	}, nil
}
