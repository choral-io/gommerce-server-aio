package repos_pgsql

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-aio/data/repos"
)

type dataRepos struct {
	bdb bun.IDB
}

func (r *dataRepos) RunInTx(ctx context.Context, opts *sql.TxOptions, f func(ctx context.Context, drs repos.DataRepos) error) error {
	return r.bdb.RunInTx(ctx, opts, func(ctx context.Context, tx bun.Tx) error {
		n := *r
		n.bdb = tx
		return f(ctx, &n)
	})
}

func (r *dataRepos) Clients() repos.ClientRepo {
	return &clientRepo{bdb: r.bdb}
}

func (r *dataRepos) Realms() repos.RealmRepo {
	return &realmRepo{bdb: r.bdb}
}

func (r *dataRepos) Users() repos.UserRepo {
	return &userRepo{bdb: r.bdb}
}

func (r *dataRepos) Roles() repos.RoleRepo {
	return &roleRepo{bdb: r.bdb}
}

func (r *dataRepos) Logins() repos.LoginRepo {
	return &loginRepo{bdb: r.bdb}
}

func NewDataRepos(bdb bun.IDB) repos.DataRepos {
	return &dataRepos{bdb: bdb}
}
