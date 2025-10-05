package utils_v1

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/rueidis"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"

	"github.com/choral-io/gommerce-server-core/data"
	"github.com/choral-io/gommerce-server-core/secure"

	tsoffset "github.com/choral-io/gommerce-protobuf-go/types/v1/tsoffset"
	utils_pb "github.com/choral-io/gommerce-protobuf-go/utils/v1"
)

type SequenceServiceServer struct {
	utils_pb.UnimplementedSequenceServiceServer

	seq data.Seq
}

func NewSequenceServiceServer(seq data.Seq) utils_pb.SequenceServiceServer {
	return &SequenceServiceServer{seq: seq}
}

func (s *SequenceServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&utils_pb.SequenceService_ServiceDesc, s)
}

func (s *SequenceServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return utils_pb.RegisterSequenceServiceHandler(ctx, mux, conn)
}

func (s *SequenceServiceServer) NextValue(_ context.Context, req *utils_pb.NextValueRequest) (*utils_pb.NextValueResponse, error) {
	value, err := s.seq.Next(req.Key, req.MinValue, req.MaxValue)
	if err != nil {
		return nil, err
	}
	return &utils_pb.NextValueResponse{
		Key:   req.Key,
		Value: value,
	}, nil
}

type SnowflakeServiceServer struct {
	utils_pb.UnimplementedSnowflakeServiceServer

	idw data.IDWorker
}

func NewSnowflakeServiceServer(idw data.IDWorker) utils_pb.SnowflakeServiceServer {
	return &SnowflakeServiceServer{idw: idw}
}

func (s *SnowflakeServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&utils_pb.SnowflakeService_ServiceDesc, s)
}

func (s *SnowflakeServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return utils_pb.RegisterSequenceServiceHandler(ctx, mux, conn)
}

func (s *SnowflakeServiceServer) NextHex(_ context.Context, _ *utils_pb.NextHexRequest) (*utils_pb.NextHexResponse, error) {
	return &utils_pb.NextHexResponse{
		Value: s.idw.NextHex(),
	}, nil
}

func (s *SnowflakeServiceServer) NextInt64(_ context.Context, _ *utils_pb.NextInt64Request) (*utils_pb.NextInt64Response, error) {
	return &utils_pb.NextInt64Response{
		Value: s.idw.NextInt64(),
	}, nil
}

type PasswordServiceServer struct {
	utils_pb.UnimplementedPasswordServiceServer
}

func NewPasswordServiceServer() utils_pb.PasswordServiceServer {
	return &PasswordServiceServer{}
}

func (s *PasswordServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&utils_pb.PasswordService_ServiceDesc, s)
}

func (s *PasswordServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return utils_pb.RegisterPasswordServiceHandler(ctx, mux, conn)
}

func (s *PasswordServiceServer) GeneratePassword(_ context.Context, req *utils_pb.GeneratePasswordRequest) (*utils_pb.GeneratePasswordResponse, error) {
	if req.Symbols == "" {
		req.Symbols = secure.DefaultPasswordSymbols
	}
	if req.Length <= 0 {
		req.Length = 16
	}
	pwd, err := secure.RandString(int(req.Length), req.Symbols)
	if err != nil {
		return nil, err
	}
	return &utils_pb.GeneratePasswordResponse{
		Value: pwd,
	}, nil
}

func (s *PasswordServiceServer) HashPassword(_ context.Context, req *utils_pb.HashPasswordRequest) (*utils_pb.HashPasswordResponse, error) {
	if len(req.Value) == 0 {
		err := errors.New("provided password must not be empty")
		return nil, err
	}
	if value, err := bcrypt.GenerateFromPassword([]byte(req.Value), 12); err == nil {
		return &utils_pb.HashPasswordResponse{
			Value: string(value),
		}, nil
	} else {
		return nil, err
	}
}

func (s *PasswordServiceServer) ValidatePassword(_ context.Context, req *utils_pb.ValidatePasswordRequest) (*utils_pb.ValidatePasswordResponse, error) {
	if len(req.HashedPassword) == 0 {
		return nil, errors.New("provided password must not be empty")
	}
	if len(req.ProvidedPassword) == 0 {
		return nil, errors.New("hashed password must not be empty")
	}
	err := bcrypt.CompareHashAndPassword([]byte(req.HashedPassword), []byte(req.ProvidedPassword))
	return &utils_pb.ValidatePasswordResponse{
		Valid: err == nil,
	}, nil
}

type DateTimeServiceServer struct {
	utils_pb.UnimplementedDateTimeServiceServer

	bdb bun.IDB
	rdb rueidis.Client
}

func NewDateTimeServiceServer(bdb bun.IDB, rdb rueidis.Client) utils_pb.DateTimeServiceServer {
	return &DateTimeServiceServer{
		bdb: bdb,
		rdb: rdb,
	}
}

func (s *DateTimeServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&utils_pb.DateTimeService_ServiceDesc, s)
}

func (s *DateTimeServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return utils_pb.RegisterDateTimeServiceHandler(ctx, mux, conn)
}

func (s *DateTimeServiceServer) GetDBNow(ctx context.Context, _ *utils_pb.GetDBNowRequest) (*utils_pb.GetDBNowResponse, error) {
	var now time.Time
	if err := s.bdb.QueryRowContext(ctx, "SELECT NOW()").Scan(&now); err != nil {
		return nil, err
	}
	return &utils_pb.GetDBNowResponse{
		Value: tsoffset.New(now),
	}, nil
}

func (s *DateTimeServiceServer) GetRedisNow(ctx context.Context, _ *utils_pb.GetRedisNowRequest) (*utils_pb.GetRedisNowResponse, error) {
	strs, err := s.rdb.Do(ctx, s.rdb.B().Time().Build()).AsStrSlice()
	if err != nil {
		return nil, err
	}
	sec, _ := strconv.ParseInt(strs[0], 10, 64)
	msec, _ := strconv.ParseInt(strs[1], 10, 64)
	return &utils_pb.GetRedisNowResponse{
		Value: tsoffset.New(time.Unix(sec, msec*1000)),
	}, nil
}

func (s *DateTimeServiceServer) GetUTCNow(context.Context, *utils_pb.GetUTCNowRequest) (*utils_pb.GetUTCNowResponse, error) {
	return &utils_pb.GetUTCNowResponse{
		Value: tsoffset.Now().UTC(),
	}, nil
}

func (s *DateTimeServiceServer) GetLocalNow(context.Context, *utils_pb.GetLocalNowRequest) (*utils_pb.GetLocalNowResponse, error) {
	return &utils_pb.GetLocalNowResponse{
		Value: tsoffset.Now().Local(),
	}, nil
}

func (s *DateTimeServiceServer) WatchLocalNow(_ *utils_pb.WatchLocalNowRequest, srv utils_pb.DateTimeService_WatchLocalNowServer) error {
	if err := srv.Send(&utils_pb.WatchLocalNowResponse{Value: tsoffset.Now()}); err != nil {
		return err
	}
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := srv.Send(&utils_pb.WatchLocalNowResponse{Value: tsoffset.Now()}); err != nil {
				return err
			}
		case <-srv.Context().Done():
			return srv.Context().Err()
		}
	}
}
