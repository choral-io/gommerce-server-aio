package server_v1beta

import (
	"github.com/choral-io/gommerce-server-aio/server"
	chats_v1b "github.com/choral-io/gommerce-server-aio/server/v1beta/chats"
	iam_v1b "github.com/choral-io/gommerce-server-aio/server/v1beta/iam"
	oss_v1b "github.com/choral-io/gommerce-server-aio/server/v1beta/oss"
	state_v1b "github.com/choral-io/gommerce-server-aio/server/v1beta/state"
)

type (
	ObjectStoreService = oss_v1b.ObjectStoreService
)

var (
	NewObjectStoreService = oss_v1b.NewObjectStoreService
)

func ProviceRegistrations() []any {
	return server.AnnotateRegistrations(
		iam_v1b.NewTokensServiceServer,
		iam_v1b.NewUsersServiceServer,
		chats_v1b.NewChatsServiceServer,
		state_v1b.NewStateStoreServiceServer,
	)
}
