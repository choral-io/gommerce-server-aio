package repos_pgsql

import (
	"context"

	"github.com/redis/rueidis"
	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-aio/data/repos"
)

type realmsRepo struct {
	repos.BaseRepo

	bdb bun.IDB
	rdb rueidis.Client
}

func NewRealmsRepo(bdb bun.IDB, rdb rueidis.Client) repos.RealmsRepo {
	return &realmsRepo{bdb: bdb, rdb: rdb}
}

func (r *realmsRepo) WithDB(bdb bun.IDB) repos.RealmsRepo {
	if r.bdb == bdb {
		return r
	}
	n := *r
	n.bdb = bdb
	return &n
}

func (r *realmsRepo) FindOneByName(ctx context.Context, name string, sqts ...repos.SelectQueryTransformer) (*models.Realm, error) {
	realm := new(models.Realm)
	if query, err := repos.TransformSelectQuery(ctx, r.bdb.NewSelect().Model(realm).Where("name = ?", name), sqts...); err != nil {
		return nil, err
	} else if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return realm, nil
}
