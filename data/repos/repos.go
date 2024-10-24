package repos

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type baseRepo interface {
	mustBeBaseRepo()
}

type BaseRepo struct{}

func (*BaseRepo) mustBeBaseRepo() {}

type DataRepos interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (bun.Tx, error)
	RunInTx(ctx context.Context, opts *sql.TxOptions, f func(ctx context.Context, drs DataRepos) error) error
	Clients() ClientsRepo
	Realms() RealmsRepo
	Users() UsersRepo
	Roles() RolesRepo
	Logins() LoginsRepo
}

type dataRepos struct {
	bdb     bun.IDB
	clients ClientsRepo
	realms  RealmsRepo
	users   UsersRepo
	roles   RolesRepo
	logins  LoginsRepo
}

func NewDataRepos(
	bdb bun.IDB,
	clients ClientsRepo,
	realms RealmsRepo,
	users UsersRepo,
	roles RolesRepo,
	logins LoginsRepo,
) DataRepos {
	return &dataRepos{
		bdb:     bdb,
		clients: clients,
		realms:  realms,
		users:   users,
		roles:   roles,
		logins:  logins,
	}
}

func (r *dataRepos) BeginTx(ctx context.Context, opts *sql.TxOptions) (bun.Tx, error) {
	return r.bdb.BeginTx(ctx, opts)
}

func (r *dataRepos) RunInTx(ctx context.Context, opts *sql.TxOptions, f func(ctx context.Context, drs DataRepos) error) error {
	return r.bdb.RunInTx(ctx, opts, func(ctx context.Context, tx bun.Tx) error {
		n := *r
		n.bdb = tx
		return f(ctx, &n)
	})
}

func (r *dataRepos) Clients() ClientsRepo {
	return r.clients.WithDB(r.bdb)
}

func (r *dataRepos) Realms() RealmsRepo {
	return r.realms.WithDB(r.bdb)
}

func (r *dataRepos) Users() UsersRepo {
	return r.users.WithDB(r.bdb)
}

func (r *dataRepos) Roles() RolesRepo {
	return r.roles.WithDB(r.bdb)
}

func (r *dataRepos) Logins() LoginsRepo {
	return r.logins.WithDB(r.bdb)
}
