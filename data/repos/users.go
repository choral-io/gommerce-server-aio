package repos

import (
	"context"
	"time"

	"github.com/choral-io/gommerce-server-aio/data/models"
)

type UserRepo interface {
	FindAll(ctx context.Context, sqts ...SelectQueryTransformer) ([]*models.User, int64, error)
	FindOneByID(ctx context.Context, id string, sqts ...SelectQueryTransformer) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateLoginStatus(ctx context.Context, userId string, now time.Time) error
}
