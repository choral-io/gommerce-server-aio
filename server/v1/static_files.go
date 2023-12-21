package v1

import (
	"context"
	"net/http"

	"github.com/choral-io/gommerce-server-aio/static"
	"github.com/choral-io/gommerce-server-core/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type staticFilesServer struct {
}

func NewStaticFilesServer() server.GatewayClientRegister {
	return &staticFilesServer{}
}

func (s *staticFilesServer) RegisterGatewayClient(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	fileServer := http.FileServer(http.FS(static.GetStaticFS()))
	return mux.HandlePath("GET", "/**", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		fileServer.ServeHTTP(w, r)
	})
}
