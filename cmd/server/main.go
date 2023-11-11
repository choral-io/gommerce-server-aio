package main

import (
	"database/sql"
	"net/http"

	"github.com/choral-io/gommerce-server-core/config"
	"github.com/choral-io/gommerce-server-core/data"
	"github.com/choral-io/gommerce-server-core/logging"
	"github.com/choral-io/gommerce-server-core/otel"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/choral-io/gommerce-server-core/server"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/uptrace/bun"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	_ "github.com/choral-io/gommerce-server-aio/data/drivers" // register db drivers
	srv "github.com/choral-io/gommerce-server-aio/server"
	srv_v1 "github.com/choral-io/gommerce-server-aio/server/v1"
	srv_v1b "github.com/choral-io/gommerce-server-aio/server/v1beta"
)

var (
	grpc_server_fx_tag = `group:"grpc_servers"`
	grpc_server_anns   = []fx.Annotation{fx.As(new(any)), fx.ResultTags(grpc_server_fx_tag)}
)

func main() {
	fx.New(
		fx.Provide(config.LoadRootConfig, config.ExtractSections), // load and extract config sections
		fx.Provide(logging.NewLogger),                             // create logger
		fx.Provide(otel.NewServerResource),                        // create server resource for opentelemetry
		fx.Provide(otel.NewTracerProvider),                        // create tracer provider for opentelemetry
		fx.Provide(otel.NewMeterProvider),                         // create meter provider for opentelemetry
		fx.Provide(data.NewRedisClient, data.NewRedisSeq),         // create redis client and redis seq
		fx.Provide(data.NewIdWorker),                              // create id worker
		fx.Provide(data.NewBunDB, func(bdb *bun.DB) (*sql.DB, bun.IDB) { // create bun.DB and expose sql.DB and bun.IDB from it
			return bdb.DB, bdb
		}),
		fx.Provide(secure.NewTokenStore, srv_v1b.NewBasicTokenStore), // create token stores
		fx.Provide(func(uts secure.TokenStore, cts *srv_v1b.BasicTokenStore) *secure.ServerAuthorizer { // create server authorizer
			return secure.NewServerAuthorizer(map[string]secure.TokenStore{
				secure.AUTH_SCHEMA_BEARER: uts,
				secure.AUTH_SCHEMA_BASIC:  cts,
			})
		}),
		fx.Provide(srv.NewSelectorMatcher),
		fx.Provide(fx.Annotate(server.NewHealthServiceServer, grpc_server_anns...)),
		fx.Provide(fx.Annotate(srv_v1.NewSequenceServiceServer, grpc_server_anns...)),
		fx.Provide(fx.Annotate(srv_v1.NewSnowflakeServiceServer, grpc_server_anns...)),
		fx.Provide(fx.Annotate(srv_v1.NewPasswordServiceServer, grpc_server_anns...)),
		fx.Provide(fx.Annotate(srv_v1.NewDateTimeServiceServer, grpc_server_anns...)),
		fx.Provide(fx.Annotate(srv_v1b.NewTokensServiceServer, grpc_server_anns...)),
		fx.Provide(fx.Annotate(srv_v1b.NewUsersServiceServer, grpc_server_anns...)),
		fx.Provide(fx.Annotate(func(
			cfg config.ServerHTTPConfig, l logging.Logger, tp trace.TracerProvider, mp metric.MeterProvider,
			auth *secure.ServerAuthorizer, matcher selector.Matcher, regs []any,
		) (http.Handler, error) {
			return server.NewGRPCHandler(cfg,
				server.WithOTELStatsHandler(tp, mp),
				server.WithLoggingInterceptor(l),
				server.WithRecoveryInterceptor(nil),
				server.WithSecureInterceptor(auth, matcher),
				server.WithRegistrations(regs...),
			)
		}, fx.ParamTags(``, ``, ``, ``, ``, ``, grpc_server_fx_tag))),
		fx.Provide(server.NewHTTPServer), // create http server
		fx.Invoke(func(srv *server.HTTPServer, lc fx.Lifecycle) { // register http server to lifecycle
			lc.Append(fx.Hook{OnStart: srv.Start, OnStop: srv.Stop})
		}),
		fx.WithLogger(func(l logging.Logger) fxevent.Logger { // create logger for fx
			return logging.NewFxeventLogger(l).UseEventLevel(logging.LevelDebug)
		}),
	).Run()
}
