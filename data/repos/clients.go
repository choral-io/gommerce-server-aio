package repos

import (
	"context"

	"github.com/choral-io/gommerce-server-aio/data/models"
)

type ClientRepo interface {
	FindOneBySecretKey(ctx context.Context, secretKey string, sqts ...SelectQueryTransformer) (*models.Client, error)
}
