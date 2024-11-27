package repos

import (
	"context"

	"github.com/choral-io/gommerce-server-aio/data/models"
)

type LoginRepo interface {
	CreateLogin(ctx context.Context, login *models.Login) error
	FindOneByIdentifier(ctx context.Context, realmId, provider, identifier string, sqts ...SelectQueryTransformer) (*models.Login, error)
}
