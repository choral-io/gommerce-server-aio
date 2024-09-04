module github.com/choral-io/gommerce-server-aio

go 1.22.0

replace github.com/choral-io/gommerce-protobuf-go v0.0.0 => ../gommerce-protobuf-go

replace github.com/choral-io/gommerce-server-core v0.0.0 => ../gommerce-server-core

require (
	github.com/choral-io/gommerce-protobuf-go v0.0.0
	github.com/choral-io/gommerce-server-core v0.0.0
	github.com/coder/websocket v1.8.12
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.1.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0
	github.com/joho/godotenv v1.5.1
	github.com/minio/minio-go/v7 v7.0.76
	github.com/nats-io/nats.go v1.37.0
	github.com/redis/rueidis v1.0.45
	github.com/uptrace/bun v1.2.3
	github.com/uptrace/bun/dialect/mssqldialect v1.2.3
	github.com/uptrace/bun/dialect/mysqldialect v1.2.3
	github.com/uptrace/bun/dialect/pgdialect v1.2.3
	github.com/uptrace/bun/driver/pgdriver v1.2.3
	go.opentelemetry.io/otel/metric v1.29.0
	go.opentelemetry.io/otel/trace v1.29.0
	go.uber.org/fx v1.22.2
	golang.org/x/crypto v0.26.0
	google.golang.org/grpc v1.66.0
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.1.0 // indirect
	github.com/expr-lang/expr v1.16.9 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/puzpuzpuz/xsync/v3 v3.4.0 // indirect
	github.com/redis/rueidis/rueidisotel v1.0.45 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/uptrace/bun/extra/bunotel v1.2.3 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelsql v0.3.1 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.54.0 // indirect
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.29.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.29.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.29.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.29.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.29.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.29.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.29.0 // indirect
	go.opentelemetry.io/otel/sdk v1.29.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.29.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	go.uber.org/dig v1.18.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/mod v0.20.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mellium.im/sasl v0.3.1 // indirect
)
