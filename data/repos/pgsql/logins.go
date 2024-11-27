package repos_pgsql

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-aio/data/repos"
)

type loginRepo struct {
	bdb bun.IDB
}

func (r *loginRepo) CreateLogin(ctx context.Context, login *models.Login) error {
	if _, err := r.bdb.NewInsert().Model(login).Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *loginRepo) FindOneByIdentifier(ctx context.Context, realmId, provider, identifier string, sqts ...repos.SelectQueryTransformer) (*models.Login, error) {
	login := new(models.Login)
	if query, err := repos.TransformSelectQuery(
		ctx,
		r.bdb.NewSelect().Model(login).
			Where(`"login"."provider" = ?`, provider).
			Where(`"login"."identifier" = ?`, identifier).
			Where(`"user"."realm_id" = ?`, realmId),
		sqts...,
	); err != nil {
		return nil, err
	} else if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return login, nil
}
