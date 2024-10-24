package repos_pgsql

import (
	"context"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-aio/data/repos"
	"github.com/redis/rueidis"
	"github.com/uptrace/bun"
)

type clientsRepo struct {
	repos.BaseRepo

	bdb bun.IDB
	rdb rueidis.Client
}

func NewClientsRepo(bdb bun.IDB, rdb rueidis.Client) repos.ClientsRepo {
	return &clientsRepo{bdb: bdb, rdb: rdb}
}

func (r *clientsRepo) WithDB(bdb bun.IDB) repos.ClientsRepo {
	if r.bdb == bdb {
		return r
	}
	n := *r
	n.bdb = bdb
	return &n
}

func (r *clientsRepo) FindOneBySecretKey(ctx context.Context, secretKey string, sqts ...repos.SelectQueryTransformer) (*models.Client, error) {
	client := new(models.Client)
	if query, err := repos.TransformSelectQuery(
		ctx,
		r.bdb.NewSelect().Model(client).Where("secret_key = ?", secretKey),
		sqts...,
	); err != nil {
		return nil, err
	} else if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return client, nil
}
