package repos

import (
	"context"
)

type RoleRepo interface {
	FindNamesForUser(ctx context.Context, userId string, sqts ...SelectQueryTransformer) ([]string, int64, error)
}
