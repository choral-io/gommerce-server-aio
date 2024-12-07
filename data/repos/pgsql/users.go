package repos_pgsql

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-aio/data/repos"
)

type userRepo struct {
	bdb bun.IDB
}

func (r *userRepo) FindByID(ctx context.Context, id string, sqts ...repos.SelectQueryTransformer) (*models.User, error) {
	user := new(models.User)
	if query, err := repos.TransformSelectQuery(ctx, r.bdb.NewSelect().Model(user).Where(`"user"."id" = ?`, id), sqts...); err != nil {
		return nil, err
	} else if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) Find(ctx context.Context, sqts ...repos.SelectQueryTransformer) ([]*models.User, int64, error) {
	var users []*models.User
	if query, err := repos.TransformSelectQuery(ctx, r.bdb.NewSelect().Model((*models.User)(nil)), sqts...); err != nil {
		return nil, 0, err
	} else if total, err := query.ScanAndCount(ctx, &users); err != nil {
		query.GetModel().Value()
		return nil, 0, err
	} else {
		query.GetModel().Value()
		return users, int64(total), nil
	}
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	return r.bdb.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(user).Exec(ctx); err != nil {
			return err
		}
		if user.Profile != nil {
			user.Profile.Id = user.Id
			if _, err := tx.NewInsert().Model(user.Profile).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *userRepo) UpdateLoginStatus(ctx context.Context, userId string, now time.Time) error {
	if _, err := r.bdb.NewUpdate().Model((*models.User)(nil)).
		Set(`"updated_at" = ?`, now).
		Set(`"first_login_time" = COALESCE("first_login_time", ?)`, now).
		Set(`"last_active_time" = ?`, now).
		Where(`"id" = ?`, userId).Exec(ctx); err != nil {
		return err
	}
	return nil
}
