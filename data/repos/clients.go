package repos

import (
	"context"

	"github.com/choral-io/gommerce-server-aio/data/models"
)

type ClientRepo interface {
	FindBySecretKey(ctx context.Context, secretKey string, sqts ...SelectQueryTransformer) (*models.Client, error)
}
