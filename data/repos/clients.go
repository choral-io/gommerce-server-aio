package repos

import (
	"context"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/uptrace/bun"
)

type ClientsRepo interface {
	baseRepo

	WithDB(bdb bun.IDB) ClientsRepo
	FindOneBySecretKey(ctx context.Context, secretKey string, sqts ...SelectQueryTransformer) (*models.Client, error)
}
