package repos_pgsql

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-aio/data/repos"
)

type roleRepo struct {
	bdb bun.IDB
}

func (r *roleRepo) FindNamesForUser(ctx context.Context, userId string, sqts ...repos.SelectQueryTransformer) ([]string, int64, error) {
	var roles []string
	if total, err := r.bdb.NewSelect().Model((*models.RoleUser)(nil)).
		Relation("Role", func(sq *bun.SelectQuery) *bun.SelectQuery { return sq.ExcludeColumn("*") }).
		Column("role.name").
		Where(`"role_user"."user_id" = ?`, userId).ScanAndCount(ctx, &roles); err != nil {
		return nil, 0, err
	} else {
		return roles, int64(total), nil
	}
}
