package state_v1beta

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/rueidis"
	"google.golang.org/grpc"

	"github.com/choral-io/gommerce-server-core/config"
	"github.com/choral-io/gommerce-server-core/logging"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/choral-io/gommerce-server-core/validator"

	state_pb "github.com/choral-io/gommerce-protobuf-go/state/v1beta"
)

const (
	TtlInSecondsKey    = "ttlInSeconds"
	StorageKeyTemplate = "gommerce-server-aio:state-store:%s:%s"
)

type stateStoreServiceServer struct {
	state_pb.UnimplementedStateStoreServiceServer

	skt    string // storage key template
	rdb    rueidis.Client
	logger logging.Logger
}

func NewStateStoreServiceServer(rdb rueidis.Client, cfg config.RootConfig, logger logging.Logger) state_pb.StateStoreServiceServer {
	svc := &stateStoreServiceServer{
		skt:    StorageKeyTemplate, // default storage key template
		rdb:    rdb,
		logger: logger,
	}
	if err := cfg.GetValue("$.services.state.storage.key-template", &svc.skt); err != nil {
		logger.Error(context.Background(), "failed to read storage key template from config", "error", err)
		svc.skt = StorageKeyTemplate
		logger.Warn(context.Background(), "using default storage key template", "template", svc.skt)
	} else if svc.skt == "" {
		logger.Warn(context.Background(), "storage key template is empty, using default", "default_template", StorageKeyTemplate)
		svc.skt = StorageKeyTemplate // fallback to default if empty
	} else {
		logger.Info(context.Background(), "using custom storage key template", "template", svc.skt)
	}
	return svc
}

func (s *stateStoreServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&state_pb.StateStoreService_ServiceDesc, s)
}

func (s *stateStoreServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return state_pb.RegisterStateStoreServiceHandler(ctx, mux, conn)
}

func (s *stateStoreServiceServer) Authorize(ctx context.Context, _ string) error {
	return secure.Authorize(ctx, secure.AuthFuncAuthenticated)
}

func (s *stateStoreServiceServer) GetState(ctx context.Context, req *state_pb.GetStateRequest) (*state_pb.GetStateResponse, error) {
	sub := secure.IdentityFromContext(ctx).Token().Subject()
	key := fmt.Sprintf(s.skt, sub, req.GetKey())
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
	key := fmt.Sprintf(s.skt, sub, req.GetKey())
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
	key := fmt.Sprintf(s.skt, sub, req.GetKey())
	cmd := s.rdb.B().Del().Key(key)
	err := s.rdb.Do(ctx, cmd.Build()).Error()
	if err != nil {
		return nil, err
	}
	return &state_pb.DelStateResponse{}, nil
}
