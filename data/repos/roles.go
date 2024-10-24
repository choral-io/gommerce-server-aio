package repos

import (
	"context"

	"github.com/uptrace/bun"
)

type RolesRepo interface {
	baseRepo

	WithDB(bdb bun.IDB) RolesRepo
	FindNamesForUser(ctx context.Context, userId string, sqts ...SelectQueryTransformer) ([]string, int64, error)
}
