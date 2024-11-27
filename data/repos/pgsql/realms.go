package repos_pgsql

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-aio/data/repos"
)

type realmRepo struct {
	bdb bun.IDB
}

func (r *realmRepo) FindOneByName(ctx context.Context, name string, sqts ...repos.SelectQueryTransformer) (*models.Realm, error) {
	realm := new(models.Realm)
	if query, err := repos.TransformSelectQuery(ctx, r.bdb.NewSelect().Model(realm).Where(`"name" = ?`, name), sqts...); err != nil {
		return nil, err
	} else if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return realm, nil
}
