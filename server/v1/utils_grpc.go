package v1

import (
	"context"
	"errors"
	"strconv"
	"time"

	tsoffset "github.com/choral-io/gommerce-protobuf-go/types/v1/tsoffset"
	utils "github.com/choral-io/gommerce-protobuf-go/utils/v1"
	"github.com/choral-io/gommerce-server-core/data"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/rueidis"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"

	"golang.org/x/crypto/bcrypt"
)

type sequenceServiceServer struct {
	utils.UnimplementedSequenceServiceServer

	seq data.Seq
}

func NewSequenceServiceServer(seq data.Seq) utils.SequenceServiceServer {
	return &sequenceServiceServer{seq: seq}
}

func (s *sequenceServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&utils.SequenceService_ServiceDesc, s)
}

func (s *sequenceServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return utils.RegisterSequenceServiceHandler(ctx, mux, conn)
}

func (s *sequenceServiceServer) NextValue(_ context.Context, req *utils.NextValueRequest) (*utils.NextValueResponse, error) {
	value, err := s.seq.Next(req.Key, req.MinValue, req.MaxValue)
	if err != nil {
		return nil, err
	}
	return &utils.NextValueResponse{
		Key:   req.Key,
		Value: value,
	}, nil
}

type snowflakeServiceServer struct {
	utils.UnimplementedSnowflakeServiceServer

	idw data.IdWorker
}

func NewSnowflakeServiceServer(idw data.IdWorker) utils.SnowflakeServiceServer {
	return &snowflakeServiceServer{idw: idw}
}

func (s *snowflakeServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&utils.SnowflakeService_ServiceDesc, s)
}

func (s *snowflakeServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return utils.RegisterSequenceServiceHandler(ctx, mux, conn)
}

func (s *snowflakeServiceServer) NextHex(ctx context.Context, _ *utils.NextHexRequest) (*utils.NextHexResponse, error) {
	return &utils.NextHexResponse{
		Value: s.idw.NextHex(),
	}, nil
}

func (s *snowflakeServiceServer) NextInt64(ctx context.Context, _ *utils.NextInt64Request) (*utils.NextInt64Response, error) {
	return &utils.NextInt64Response{
		Value: s.idw.NextInt64(),
	}, nil
}

type passwordServiceServer struct {
	utils.UnimplementedPasswordServiceServer
}

func NewPasswordServiceServer() utils.PasswordServiceServer {
	return &passwordServiceServer{}
}

func (s *passwordServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&utils.PasswordService_ServiceDesc, s)
}

func (s *passwordServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return utils.RegisterPasswordServiceHandler(ctx, mux, conn)
}

func (p *passwordServiceServer) GeneratePassword(_ context.Context, req *utils.GeneratePasswordRequest) (*utils.GeneratePasswordResponse, error) {
	if req.Symbols == "" {
		req.Symbols = secure.DEFAULT_PASSWORD_SYMBOLS
	}
	if req.Length <= 0 {
		req.Length = 16
	}
	pwd, err := secure.RandString(int(req.Length), req.Symbols)
	if err != nil {
		return nil, err
	}
	return &utils.GeneratePasswordResponse{
		Value: pwd,
	}, nil
}

func (p *passwordServiceServer) HashPassword(_ context.Context, req *utils.HashPasswordRequest) (*utils.HashPasswordResponse, error) {
	if len(req.Value) == 0 {
		err := errors.New("provided password msut not be empty")
		return nil, err
	}
	if value, err := bcrypt.GenerateFromPassword([]byte(req.Value), bcrypt.DefaultCost); err == nil {
		return &utils.HashPasswordResponse{
			Value: string(value),
		}, nil
	} else {
		return nil, err
	}
}

func (p *passwordServiceServer) ValidatePassword(_ context.Context, req *utils.ValidatePasswordRequest) (*utils.ValidatePasswordResponse, error) {
	if len(req.HashedPassword) == 0 {
		return nil, errors.New("provided password msut not be empty")
	}
	if len(req.ProvidedPassword) == 0 {
		return nil, errors.New("hashed password msut not be empty")
	}
	err := bcrypt.CompareHashAndPassword([]byte(req.HashedPassword), []byte(req.ProvidedPassword))
	return &utils.ValidatePasswordResponse{
		Valid: err == nil,
	}, nil
}

type dateTimeServiceServer struct {
	utils.UnimplementedDateTimeServiceServer

	bdb bun.IDB
	rdb rueidis.Client
}

func NewDateTimeServiceServer(bdb bun.IDB, rdb rueidis.Client) utils.DateTimeServiceServer {
	return &dateTimeServiceServer{
		bdb: bdb,
		rdb: rdb,
	}
}

func (s *dateTimeServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&utils.DateTimeService_ServiceDesc, s)
}

func (s *dateTimeServiceServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return utils.RegisterDateTimeServiceHandler(ctx, mux, conn)
}

func (d *dateTimeServiceServer) GetDBNow(ctx context.Context, _ *utils.GetDBNowRequest) (*utils.GetDBNowResponse, error) {
	var now time.Time
	if err := d.bdb.QueryRowContext(ctx, "SELECT NOW()").Scan(&now); err != nil {
		return nil, err
	}
	return &utils.GetDBNowResponse{
		Value: tsoffset.New(now),
	}, nil
}

func (d *dateTimeServiceServer) GetRedisNow(ctx context.Context, _ *utils.GetRedisNowRequest) (*utils.GetRedisNowResponse, error) {
	strs, err := d.rdb.Do(ctx, d.rdb.B().Time().Build()).AsStrSlice()
	if err != nil {
		return nil, err
	}
	sec, _ := strconv.ParseInt(strs[0], 10, 64)
	msec, _ := strconv.ParseInt(strs[1], 10, 64)
	return &utils.GetRedisNowResponse{
		Value: tsoffset.New(time.Unix(sec, msec*1000)),
	}, nil
}

func (d *dateTimeServiceServer) GetUTCNow(context.Context, *utils.GetUTCNowRequest) (*utils.GetUTCNowResponse, error) {
	return &utils.GetUTCNowResponse{
		Value: tsoffset.Now().UTC(),
	}, nil
}

func (d *dateTimeServiceServer) GetLocalNow(context.Context, *utils.GetLocalNowRequest) (*utils.GetLocalNowResponse, error) {
	return &utils.GetLocalNowResponse{
		Value: tsoffset.Now().Local(),
	}, nil
}

func (d *dateTimeServiceServer) WatchLocalNow(_ *utils.WatchLocalNowRequest, srv utils.DateTimeService_WatchLocalNowServer) error {
	if err := srv.Send(&utils.WatchLocalNowResponse{Value: tsoffset.Now()}); err != nil {
		return err
	}
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := srv.Send(&utils.WatchLocalNowResponse{Value: tsoffset.Now()}); err != nil {
				return err
			}
		case <-srv.Context().Done():
			if err := srv.Context().Err(); err == context.Canceled {
				return nil
			} else {
				return err
			}
		}
	}
}
