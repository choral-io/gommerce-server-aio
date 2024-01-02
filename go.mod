module github.com/choral-io/gommerce-server-aio

go 1.21.0

replace github.com/choral-io/gommerce-protobuf-go v0.0.0 => ../gommerce-protobuf-go

replace github.com/choral-io/gommerce-server-core v0.0.0 => ../gommerce-server-core

require (
	github.com/choral-io/gommerce-protobuf-go v0.0.0
	github.com/choral-io/gommerce-server-core v0.0.0
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.1
	github.com/joho/godotenv v1.5.1
	github.com/redis/rueidis v1.0.26
	github.com/uptrace/bun v1.1.16
	github.com/uptrace/bun/dialect/mssqldialect v1.1.16
	github.com/uptrace/bun/dialect/mysqldialect v1.1.16
	github.com/uptrace/bun/dialect/pgdialect v1.1.16
	github.com/uptrace/bun/driver/pgdriver v1.1.16
	go.opentelemetry.io/otel/metric v1.21.0
	go.opentelemetry.io/otel/trace v1.21.0
	go.uber.org/fx v1.20.1
	golang.org/x/crypto v0.17.0
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.2 // indirect
	github.com/expr-lang/expr v1.15.7 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/rs/cors v1.10.1 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/uptrace/bun/extra/bunotel v1.1.16 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelsql v0.2.3 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.46.1 // indirect
	go.opentelemetry.io/otel v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v0.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.44.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.21.0 // indirect
	go.opentelemetry.io/otel/sdk v1.21.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.21.0 // indirect
	go.opentelemetry.io/proto/otlp v1.0.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240102182953-50ed04b92917 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240102182953-50ed04b92917 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mellium.im/sasl v0.3.1 // indirect
)
