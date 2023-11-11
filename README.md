# Gommerce Server AIO

All in one server of project `gommerce`.

```sh
# config
cp ../gommerce-server-core/config/example.yaml ./config/app-local.yaml

# jwt
openssl genrsa 2048 | tee >(openssl rsa -pubout 2>/dev/null)

# build
CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/server ./cmd/server

# run
GOMMERCE_CONFIG_PATH=./config/app-local.yaml go run ./cmd/server

# test
go test ./... -v -cover
```
