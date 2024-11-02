package repos

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-aio/data/models"
)

type ClientsRepo interface {
	baseRepo

	WithDB(bdb bun.IDB) ClientsRepo
	FindOneBySecretKey(ctx context.Context, secretKey string, sqts ...SelectQueryTransformer) (*models.Client, error)
}
