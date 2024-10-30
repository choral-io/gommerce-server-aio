package server_v1

import (
	"github.com/choral-io/gommerce-server-aio/server"
	utils_v1 "github.com/choral-io/gommerce-server-aio/server/v1/utils"
)

func ProviceRegistrations() []any {
	return server.AnnotateRegistrations(
		utils_v1.NewSequenceServiceServer,
		utils_v1.NewSnowflakeServiceServer,
		utils_v1.NewDateTimeServiceServer,
		utils_v1.NewPasswordServiceServer,
	)
}
