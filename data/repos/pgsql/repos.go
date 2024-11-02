package repos_pgsql

import (
	"github.com/redis/rueidis"
	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-aio/data/repos"
)

func NewDataRepos(bdb bun.IDB, rdb rueidis.Client) repos.DataRepos {
	return repos.NewDataRepos(
		bdb,
		NewClientsRepo(bdb, rdb),
		NewRealmsRepo(bdb, rdb),
		NewUsersRepo(bdb, rdb),
		NewRolesRepo(bdb, rdb),
		NewLoginsRepo(bdb, rdb),
	)
}
