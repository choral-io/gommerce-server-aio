package v1beta

import (
	"context"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-core/data"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"

	iam "github.com/choral-io/gommerce-protobuf-go/iam/v1beta"
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

func (s *usersServiceServer) Authorize(ctx context.Context, method string) error {
	if method == iam.UsersService_ListUsers_FullMethodName {
		return secure.Authorize(ctx, secure.AuthFuncAuthenticated)
	}
	if method == iam.UsersService_GetIdentity_FullMethodName {
		return secure.Authorize(ctx, secure.AuthFuncAuthenticated, secure.AuthFuncRequireSchema(secure.AUTH_SCHEMA_BEARER))
	}
	return secure.Authorize(ctx, secure.AuthFuncAuthenticated)
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
	if err != nil {
		return nil, err
	}
	res := &iam.ListUsersResponse{
		Items: make([]*iam.User, len(users)),
		Page:  req.Page,
		Size:  req.Size,
		Total: int64(total),
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
