package v1beta

import (
	"context"
	"fmt"
	"strconv"

	state "github.com/choral-io/gommerce-protobuf-go/state/v1beta"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/choral-io/gommerce-server-core/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/rueidis"
	"google.golang.org/grpc"
)

const (
	TTL_IN_SECONDS_KEY   = "ttlInSeconds"
	STORAGE_KEY_TEMPLATE = "state:store:%s:%s"
)

type stateStoreServiceServer struct {
	state.UnimplementedStateStoreServiceServer

	rdb rueidis.Client
}

func NewStateStoreServiceServer(rdb rueidis.Client) state.StateStoreServiceServer {
	return &stateStoreServiceServer{
		rdb: rdb,
	}
}

func (s *stateStoreServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&state.StateStoreService_ServiceDesc, s)
}

func (s *stateStoreServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return state.RegisterStateStoreServiceHandler(ctx, mux, conn)
}

func (s *stateStoreServiceServer) Authorize(ctx context.Context, procedure string) error {
	return secure.Authorize(ctx, secure.AuthFuncAuthenticated, secure.AuthFuncRequireSchema(secure.AUTH_SCHEMA_BASIC))
}

func (s *stateStoreServiceServer) GetState(ctx context.Context, req *state.GetStateRequest) (*state.GetStateResponse, error) {
	sub := secure.IdentityFromContext(ctx).Token().Subject()
	key := fmt.Sprintf(STORAGE_KEY_TEMPLATE, sub, req.GetKey())
	cmd := s.rdb.B().Get().Key(key)
	data, err := s.rdb.Do(ctx, cmd.Build()).ToString()
	if err != nil && err != rueidis.Nil {
		return nil, err
	}
	return &state.GetStateResponse{
		Data: []byte(data),
	}, nil
}

func (s *stateStoreServiceServer) SetState(ctx context.Context, req *state.SetStateRequest) (*state.SetStateResponse, error) {
	sub := secure.IdentityFromContext(ctx).Token().Subject()
	key := fmt.Sprintf(STORAGE_KEY_TEMPLATE, sub, req.GetKey())
	cmd := s.rdb.B().Set().Key(key).Value(string(req.GetData()))
	if val, ok := req.Metadata[TTL_IN_SECONDS_KEY]; ok {
		ttl, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return nil, validator.NewErrorWithCause("metadata.ttlInSeconds", "metadata.ttlInSeconds must be an integer", err)
		}
		cmd.ExSeconds(ttl)
	}
	err := s.rdb.Do(ctx, cmd.Build()).Error()
	if err != nil {
		return nil, err
	}
	return &state.SetStateResponse{}, nil
}

func (s *stateStoreServiceServer) DelState(ctx context.Context, req *state.DelStateRequest) (*state.DelStateResponse, error) {
	sub := secure.IdentityFromContext(ctx).Token().Subject()
	key := fmt.Sprintf(STORAGE_KEY_TEMPLATE, sub, req.GetKey())
	cmd := s.rdb.B().Del().Key(key)
	err := s.rdb.Do(ctx, cmd.Build()).Error()
	if err != nil {
		return nil, err
	}
	return &state.DelStateResponse{}, nil
}
