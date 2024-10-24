package repos

import (
	"context"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/uptrace/bun"
)

type LoginsRepo interface {
	baseRepo

	WithDB(bdb bun.IDB) LoginsRepo
	CreateLogin(ctx context.Context, login *models.Login) error
	FindOneByIdentifier(ctx context.Context, realmId, provider, identifier string, sqts ...SelectQueryTransformer) (*models.Login, error)
}
