package repos

import (
	"context"
)

type RoleRepo interface {
	FindNamesByUser(ctx context.Context, userId string, sqts ...SelectQueryTransformer) ([]string, int64, error)
}
