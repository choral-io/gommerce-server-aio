package repos_pgsql

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-aio/data/repos"
)

type clientRepo struct {
	bdb bun.IDB
}

func (r *clientRepo) FindBySecretKey(ctx context.Context, secretKey string, sqts ...repos.SelectQueryTransformer) (*models.Client, error) {
	client := new(models.Client)
	if query, err := repos.TransformSelectQuery(
		ctx,
		r.bdb.NewSelect().Model(client).Where(`"secret_key" = ?`, secretKey),
		sqts...,
	); err != nil {
		return nil, err
	} else if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return client, nil
}
