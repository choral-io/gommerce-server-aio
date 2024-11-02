package repos

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-aio/data/models"
)

type RealmsRepo interface {
	baseRepo

	WithDB(bdb bun.IDB) RealmsRepo
	FindOneByName(ctx context.Context, name string, sqts ...SelectQueryTransformer) (*models.Realm, error)
}
