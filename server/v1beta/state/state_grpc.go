package state_v1beta

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	state_pb "github.com/choral-io/gommerce-protobuf-go/state/v1beta"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/choral-io/gommerce-server-core/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/rueidis"
	"google.golang.org/grpc"
)

const (
	TtlInSecondsKey    = "ttlInSeconds"
	StorageKeyTemplate = "state:store:%s:%s"
)

type stateStoreServiceServer struct {
	state_pb.UnimplementedStateStoreServiceServer

	rdb rueidis.Client
}

func NewStateStoreServiceServer(rdb rueidis.Client) state_pb.StateStoreServiceServer {
	return &stateStoreServiceServer{
		rdb: rdb,
	}
}

func (s *stateStoreServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&state_pb.StateStoreService_ServiceDesc, s)
}

func (s *stateStoreServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return state_pb.RegisterStateStoreServiceHandler(ctx, mux, conn)
}

func (s *stateStoreServiceServer) Authorize(ctx context.Context, _ string) error {
	return secure.Authorize(ctx, secure.AuthFuncAuthenticated, secure.AuthFuncRequireSchema(secure.AuthSchemaBasic))
}

func (s *stateStoreServiceServer) GetState(ctx context.Context, req *state_pb.GetStateRequest) (*state_pb.GetStateResponse, error) {
	sub := secure.IdentityFromContext(ctx).Token().Subject()
	key := fmt.Sprintf(StorageKeyTemplate, sub, req.GetKey())
	cmd := s.rdb.B().Get().Key(key)
	data, err := s.rdb.Do(ctx, cmd.Build()).AsBytes()
	if err != nil && !errors.Is(err, rueidis.Nil) {
		return nil, err
	}
	return &state_pb.GetStateResponse{
		Data: data,
	}, nil
}

func (s *stateStoreServiceServer) SetState(ctx context.Context, req *state_pb.SetStateRequest) (*state_pb.SetStateResponse, error) {
	sub := secure.IdentityFromContext(ctx).Token().Subject()
	key := fmt.Sprintf(StorageKeyTemplate, sub, req.GetKey())
	cmd := s.rdb.B().Set().Key(key).Value(rueidis.BinaryString(req.GetData()))
	if val, ok := req.Metadata[TtlInSecondsKey]; ok {
		ttl, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return nil, validator.NewErrorWithCause("metadata."+TtlInSecondsKey, "metadata.ttlInSeconds must be an integer", err)
		}
		cmd.ExSeconds(ttl)
	}
	err := s.rdb.Do(ctx, cmd.Build()).Error()
	if err != nil {
		return nil, err
	}
	return &state_pb.SetStateResponse{}, nil
}

func (s *stateStoreServiceServer) DelState(ctx context.Context, req *state_pb.DelStateRequest) (*state_pb.DelStateResponse, error) {
	sub := secure.IdentityFromContext(ctx).Token().Subject()
	key := fmt.Sprintf(StorageKeyTemplate, sub, req.GetKey())
	cmd := s.rdb.B().Del().Key(key)
	err := s.rdb.Do(ctx, cmd.Build()).Error()
	if err != nil {
		return nil, err
	}
	return &state_pb.DelStateResponse{}, nil
}
