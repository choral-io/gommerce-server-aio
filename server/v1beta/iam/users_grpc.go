package iam_v1beta

import (
	"context"
	"database/sql"
	"errors"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/choral-io/gommerce-server-core/secure"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-aio/data/repos"

	iam_pb "github.com/choral-io/gommerce-protobuf-go/iam/v1beta"
	sqlpb "github.com/choral-io/gommerce-protobuf-go/types/v1/sqlpb"
)

type usersServiceServer struct {
	iam_pb.UnimplementedUsersServiceServer

	drs repos.DataRepos
}

func NewUsersServiceServer(drs repos.DataRepos) iam_pb.UsersServiceServer {
	return &usersServiceServer{drs: drs}
}

func (s *usersServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&iam_pb.UsersService_ServiceDesc, s)
}

func (s *usersServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return iam_pb.RegisterUsersServiceHandler(ctx, mux, conn)
}

func (s *usersServiceServer) Authorize(ctx context.Context, procedure string) error {
	if procedure == iam_pb.UsersService_GetIdentity_FullMethodName {
		return secure.Authorize(ctx, secure.AuthFuncAuthenticated, secure.AuthFuncRequireSchema(secure.AuthSchemaBearer))
	}
	if procedure == iam_pb.UsersService_ListUsers_FullMethodName {
		return secure.Authorize(ctx, secure.AuthFuncAuthenticated, secure.AuthFuncRequireRealm(RealmAdmin))
	}
	return nil
}

func (s *usersServiceServer) Register(ctx context.Context, req *iam_pb.RegisterRequest) (*iam_pb.RegisterResponse, error) {
	realm, err := s.drs.Realms().FindByName(ctx, req.Realm)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.InvalidArgument, "realm %s not found", req.Realm)
		}
		return nil, status.Errorf(codes.Unknown, "error retrieving realm %s: %v", req.Realm, err)
	}
	if !realm.AllowRegistration() {
		return nil, status.Errorf(codes.PermissionDenied, "realm %s does not allow registration", req.Realm)
	}
	user := models.User{
		RealmId:    realm.Id,
		Disabled:   false,
		Approved:   true,
		Verified:   true,
		Attributes: map[string]string{},
		Profile: &models.Profile{
			DisplayName: req.DisplayName.GetValue(),
			AvatarUrl:   sqlpb.ToNullString(req.AvatarUrl),
			Gender:      sqlpb.EnumToNullName(req.Gender),
		},
	}
	if user.Profile.DisplayName == "" {
		user.Profile.DisplayName = req.Username
	}
	user.Attributes[models.USER_PROFILE_DISPLAY_NAME_ATTRIBUTE] = user.Profile.DisplayName
	if user.Profile.AvatarUrl.Valid {
		user.Attributes[models.USER_PROFILE_AVATAR_URL_ATTRIBUTE] = user.Profile.AvatarUrl.String
	}
	if user.Profile.Gender.Valid {
		user.Attributes[models.USER_PROFILE_GENDER_ATTRIBUTE] = user.Profile.Gender.String
	}
	login := models.Login{
		Provider:   LoginProviderFormPassword,
		Identifier: req.Username,
		Metadata:   map[string]string{},
	}
	if hp, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost); err != nil {
		return nil, status.Errorf(codes.Unknown, "error hashing password: %v", err)
	} else {
		login.Credential = sql.NullString{Valid: true, String: string(hp)}
	}
	err = s.drs.RunInTx(ctx, nil, func(ctx context.Context, drst repos.DataRepos) error {
		if err := drst.Users().Create(ctx, &user); err != nil {
			return status.Errorf(codes.Unknown, "error creating user: %v", err)
		}
		login.UserId = user.Id
		if err := drst.Logins().CreateLogin(ctx, &login); err != nil {
			return status.Errorf(codes.Unknown, "error creating login: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &iam_pb.RegisterResponse{
		User: toUserPB(&user),
	}, nil
}

func (s *usersServiceServer) ListUsers(ctx context.Context, req *iam_pb.ListUsersRequest) (*iam_pb.ListUsersResponse, error) {
	users, total, err := s.drs.Users().Find(
		ctx,
		repos.WithPagination(req),
		repos.WithRelation("Realm", "name"),
		repos.WithRelation("Creator"),
		repos.WithRelation("Creator.Realm", "name"),
	)
	if err != nil {
		return nil, err
	}
	res := &iam_pb.ListUsersResponse{
		Page:  req.Page,
		Size:  req.Size,
		Total: int64(total),
		Items: make([]*iam_pb.User, len(users)),
	}
	for i, u := range users {
		res.Items[i] = toUserPB(u)
	}
	return res, nil
}

func (s *usersServiceServer) GetIdentity(ctx context.Context, _ *iam_pb.GetIdentityRequest) (*iam_pb.GetIdentityResponse, error) {
	user, err := s.drs.Users().FindByID(
		ctx,
		secure.IdentityFromContext(ctx).Token().Subject(),
		repos.WithRelation("Realm", "name"),
	)
	if err != nil {
		return nil, err
	}
	return &iam_pb.GetIdentityResponse{
		User:  toUserPB(user),
		Scope: secure.IdentityFromContext(ctx).Token().Scope(),
	}, nil
}
