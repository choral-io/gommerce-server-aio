package repos

import (
	"context"

	"github.com/choral-io/gommerce-server-aio/data/models"
)

type RealmRepo interface {
	FindByName(ctx context.Context, name string, sqts ...SelectQueryTransformer) (*models.Realm, error)
}
