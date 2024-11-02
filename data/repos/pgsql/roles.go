package repos_pgsql

import (
	"context"

	"github.com/redis/rueidis"
	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-aio/data/repos"
)

type rolesRepo struct {
	repos.BaseRepo

	bdb bun.IDB
	rdb rueidis.Client
}

func NewRolesRepo(bdb bun.IDB, rdb rueidis.Client) repos.RolesRepo {
	return &rolesRepo{bdb: bdb, rdb: rdb}
}

func (r *rolesRepo) WithDB(bdb bun.IDB) repos.RolesRepo {
	if r.bdb == bdb {
		return r
	}
	n := *r
	n.bdb = bdb
	return &n
}

func (r *rolesRepo) FindNamesForUser(ctx context.Context, userId string, sqts ...repos.SelectQueryTransformer) ([]string, int64, error) {
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
