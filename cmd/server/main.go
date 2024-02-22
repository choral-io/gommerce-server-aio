package main

import (
	"context"
	"net/http"

	"github.com/choral-io/gommerce-server-core/config"
	"github.com/choral-io/gommerce-server-core/data"
	"github.com/choral-io/gommerce-server-core/events"
	"github.com/choral-io/gommerce-server-core/logging"
	"github.com/choral-io/gommerce-server-core/otel"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/choral-io/gommerce-server-core/server"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/nats-io/nats.go"
	"github.com/uptrace/bun"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	_ "github.com/choral-io/gommerce-server-aio/data/drivers" // register db drivers
	srv "github.com/choral-io/gommerce-server-aio/server"
	srv_v1 "github.com/choral-io/gommerce-server-aio/server/v1"
	srv_v1b "github.com/choral-io/gommerce-server-aio/server/v1beta"
	"github.com/choral-io/gommerce-server-aio/static"
)

var (
	grpc_server_fx_tag = `group:"grpc_servers"`
	grpc_servers_anns  = []fx.Annotation{fx.As(new(any)), fx.ResultTags(grpc_server_fx_tag)}
	grpc_handler_anns  = fx.ParamTags(``, grpc_server_fx_tag, ``, ``, ``, ``, ``)
)

func main() {
	fx.New(
		fx.Provide(config.LoadYamlConfig, config.ExtractSections), // load and extract config sections
		fx.Provide(logging.NewLogger),                             // create logger
		fx.Provide(otel.NewServerResource),                        // create server resource for opentelemetry
		fx.Provide(otel.NewTracerProvider),                        // create tracer provider for opentelemetry
		fx.Provide(otel.NewMeterProvider),                         // create meter provider for opentelemetry
		fx.Provide(data.NewRedisClient, data.NewRedisSeq),         // create redis client and redis seq
		fx.Provide(data.NewIdWorker),                              // create id worker
		fx.Provide(data.NewBunDB),                                 // create bun db
		fx.Provide(secure.NewTokenStore, srv.NewBasicTokenStore),  // create token stores
		fx.Provide(srv.NewServerAuthorizer),                       // create server authorizer
		fx.Provide(srv.NewSelectorMatcher),                        // create selector matcher
		fx.Provide(server.NewHTTPServer),                          // create http server
		fx.Provide(events.NewNATSConn),                            // create nats connection
		fx.Provide( // register grpc servers
			fx.Annotate(server.NewHealthServiceServer, grpc_servers_anns...),
			fx.Annotate(srv_v1.NewSequenceServiceServer, grpc_servers_anns...),
			fx.Annotate(srv_v1.NewSnowflakeServiceServer, grpc_servers_anns...),
			fx.Annotate(srv_v1.NewPasswordServiceServer, grpc_servers_anns...),
			fx.Annotate(srv_v1.NewDateTimeServiceServer, grpc_servers_anns...),
			fx.Annotate(srv_v1b.NewTokensServiceServer, grpc_servers_anns...),
			fx.Annotate(srv_v1b.NewUsersServiceServer, grpc_servers_anns...),
			fx.Annotate(srv_v1b.NewStateStoreServiceServer, grpc_servers_anns...),
		),
		fx.Provide( // create grpc handler
			fx.Annotate(func(cfg config.ServerHTTPConfig, regs []any,
				logger logging.Logger, tp trace.TracerProvider, mp metric.MeterProvider,
				auth *secure.ServerAuthorizer, matcher selector.Matcher,
			) (http.Handler, error) {
				fs := http.FS(static.GetStaticFS())
				return server.NewGRPCHandler(cfg,
					server.WithOTELStatsHandler(tp, mp),         // add opentelemetry stats handler
					server.WithLoggingInterceptor(logger),       // add logging interceptor
					server.WithRecoveryInterceptor(nil),         // add recovery interceptor
					server.WithSecureInterceptor(auth, matcher), // add secure interceptor
					server.WithValidatorInterceptor(),           // add validator interceptor
					server.WithRegistrations(regs...),           // add registrations
					server.WithStaticFileHandler("/**", fs),     // add static file handler
				)
			}, grpc_handler_anns)),
		fx.Invoke(data.SetDefaultIdWorker), // set default id worker
		fx.Invoke( // register db connection to lifecycle
			func(bdb bun.IDB, lc fx.Lifecycle) {
				lc.Append(fx.Hook{OnStop: func(ctx context.Context) error {
					return bdb.NewSelect().DB().Close()
				}})
			}),
		fx.Invoke( // register nats connection to lifecycle
			func(nc *nats.Conn, lc fx.Lifecycle) {
				lc.Append(fx.Hook{OnStop: func(ctx context.Context) error {
					return nc.Drain()
				}})
			}),
		fx.Invoke( // register http server to lifecycle
			func(srv *server.HTTPServer, lc fx.Lifecycle) {
				lc.Append(fx.Hook{OnStart: srv.Start, OnStop: srv.Stop})
			}),
		fx.WithLogger( // create logger for fx
			func(l logging.Logger) fxevent.Logger {
				return logging.NewFxeventLogger(l).UseEventLevel(logging.LevelDebug)
			}),
	).Run()
}
