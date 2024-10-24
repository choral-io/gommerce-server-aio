package repos_pgsql

import (
	"context"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-aio/data/repos"
	"github.com/redis/rueidis"
	"github.com/uptrace/bun"
)

type loginsRepo struct {
	repos.BaseRepo

	bdb bun.IDB
	rdb rueidis.Client
}

func NewLoginsRepo(bdb bun.IDB, rdb rueidis.Client) repos.LoginsRepo {
	return &loginsRepo{bdb: bdb, rdb: rdb}
}

func (r *loginsRepo) WithDB(bdb bun.IDB) repos.LoginsRepo {
	if r.bdb == bdb {
		return r
	}
	n := *r
	n.bdb = bdb
	return &n
}

func (r *loginsRepo) CreateLogin(ctx context.Context, login *models.Login) error {
	if _, err := r.bdb.NewInsert().Model(login).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *loginsRepo) FindOneByIdentifier(ctx context.Context, realmId, provider, identifier string, sqts ...repos.SelectQueryTransformer) (*models.Login, error) {
	login := new(models.Login)
	if query, err := repos.TransformSelectQuery(
		ctx,
		r.bdb.NewSelect().Model(login).
			Where(`"login"."provider" = ?`, provider).
			Where(`"login"."identifier" = ?`, identifier).
			Where(`"user"."realm_id" = ?`, realmId),
		sqts...,
	); err != nil {
		return nil, err
	} else if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return login, nil
}
