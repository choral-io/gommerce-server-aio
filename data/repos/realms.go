package repos

import (
	"context"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/uptrace/bun"
)

type RealmsRepo interface {
	baseRepo

	WithDB(bdb bun.IDB) RealmsRepo
	FindOneByName(ctx context.Context, name string, sqts ...SelectQueryTransformer) (*models.Realm, error)
}
