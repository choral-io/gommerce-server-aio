package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/choral-io/gommerce-server-core/config"
	"github.com/choral-io/gommerce-server-core/data"
	"github.com/choral-io/gommerce-server-core/dlock"
	"github.com/choral-io/gommerce-server-core/events"
	"github.com/choral-io/gommerce-server-core/logging"
	"github.com/choral-io/gommerce-server-core/otel"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/choral-io/gommerce-server-core/server"

	models "github.com/choral-io/gommerce-server-aio/data/models"
	repos "github.com/choral-io/gommerce-server-aio/data/repos/pgsql"
	srv "github.com/choral-io/gommerce-server-aio/server"
	srv1 "github.com/choral-io/gommerce-server-aio/server/v1"
	srv1b "github.com/choral-io/gommerce-server-aio/server/v1beta"
	static "github.com/choral-io/gommerce-server-aio/static"
)

func init() {
	// alias pg to pgsql
	sql.Register("pgsql", pgdriver.NewDriver())
}

func main() {
	env, ok := os.LookupEnv("GOMMERCE_ENVIRONMENT")
	if !ok {
		env = "development"
		_ = os.Setenv("GOMMERCE_ENVIRONMENT", env)
	}
	_ = godotenv.Load(fmt.Sprintf(".env.%s.local", env))
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load(fmt.Sprintf(".env.%s", env))
	_ = godotenv.Load(".env")
	fx.New(
		fx.Provide(config.LoadYamlConfig, config.ExtractSections), // load and extract config sections
		fx.Provide(logging.NewLogger),                             // provide logger
		fx.Provide(otel.NewServerResource),                        // provide server resource for opentelemetry
		fx.Provide(otel.NewTracerProvider),                        // provide tracer provider for opentelemetry
		fx.Provide(otel.NewMeterProvider),                         // provide meter provider for opentelemetry
		fx.Provide(data.NewRedisClient),                           // provide redis client
		fx.Provide(data.NewRedisSeq),                              // provide redis seq
		fx.Provide(dlock.NewRedisLocker),                          // provide redis locker
		fx.Provide(data.NewIdWorker),                              // provide id worker
		fx.Provide(data.NewBunDB),                                 // provide bun db
		fx.Provide(repos.NewDataRepos),                            // provide data repos
		fx.Provide(secure.NewTokenStore, srv.NewBasicTokenStore),  // provide token stores
		fx.Provide(srv.NewServerAuthorizer),                       // provide server authorizer
		fx.Provide(srv.NewSelectorMatcher),                        // provide selector matcher
		fx.Provide(server.NewHTTPServer),                          // provide http server
		fx.Provide(events.NewNATSConn),                            // provide nats connection
		fx.Provide(srv1b.NewObjectStoreService),                   // provide object store service
		fx.Provide(srv.ProviceRegistrations()...),                 // provide grpc servers
		fx.Provide(srv1.ProviceRegistrations()...),                // provide grpc servers
		fx.Provide(srv1b.ProviceRegistrations()...),               // provide grpc servers
		fx.Provide( // create grpc handler
			fx.Annotate(func(regs []any, cfg config.ServerHTTPConfig,
				logger logging.Logger, tp trace.TracerProvider, mp metric.MeterProvider,
				auth *secure.ServerAuthorizer, matcher selector.Matcher, oss *srv1b.ObjectStoreService,
			) (http.Handler, error) {
				return server.NewGRPCHandler(cfg,
					server.WithCorsOptions(cfg.GetCors()),               // add cors options
					server.WithOTELStatsHandler(tp, mp),                 // add opentelemetry stats handler
					server.WithLoggingInterceptor(logger),               // add logging interceptor
					server.WithRecoveryInterceptor(nil),                 // add recovery interceptor
					server.WithSecureInterceptor(auth, matcher),         // add secure interceptor
					server.WithValidatorInterceptor(),                   // add validator interceptor
					server.WithRegistrations(regs...),                   // add registrations
					server.WithServeMuxRoutes(oss.ServerMuxRoutes()...), // add oss mux routes
					server.WithStaticFileHandler("/**", static.FS()),    // add static file handler
				)
			}, fx.ParamTags(srv.ServerRegistrationsTag)),
		),
		fx.Invoke(logging.SetDefaultLogger), // set default logger
		fx.Invoke(data.SetDefaultIdWorker),  // set default id worker
		fx.Invoke(models.RegisterModels),    // register models
		fx.Invoke( // register db connection to lifecycle
			func(bdb bun.IDB, lc fx.Lifecycle) {
				lc.Append(fx.Hook{OnStop: func(ctx context.Context) error {
					return bdb.NewSelect().DB().Close()
				}})
			}),
		fx.Invoke( // register distributed locker to lifecycle
			func(locker dlock.Locker, lc fx.Lifecycle) {
				lc.Append(fx.Hook{OnStop: func(ctx context.Context) error {
					locker.Close()
					return nil
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
				return logging.NewFxeventLogger(l, logging.LevelInfo, logging.LevelError)
			}),
	).Run()
}
