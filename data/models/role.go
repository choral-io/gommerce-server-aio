package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-core/data"
)

type Role struct {
	bun.BaseModel `bun:"table:roles,alias:role"`

	// Columns
	Id          string         `bun:"id,pk"`
	RealmId     string         `bun:"realm_id"`
	Disabled    bool           `bun:"disabled"`
	Immutable   bool           `bun:"immutable"`
	CreatedAt   time.Time      `bun:"created_at"`
	UpdatedAt   sql.NullTime   `bun:"updated_at"`
	DeletedAt   sql.NullTime   `bun:"deleted_at,soft_delete,nullzero"`
	Name        string         `bun:"name"`
	Description sql.NullString `bun:"description"`

	// Relations
	Realm *Realm `bun:"rel:belongs-to,join:realm_id=id"`
}

func (m *Role) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if m.Id == "" {
			m.Id = data.DefaultIdWorker().NextHex()
		}
		m.CreatedAt = time.Now()
		m.UpdatedAt = sql.NullTime{Valid: false}
		m.DeletedAt = sql.NullTime{Valid: false}
		if !SeedingMode(ctx) {
			m.Immutable = false
		}
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
		if m.Immutable {
			return ErrImmutableModel
		}
	case *bun.DeleteQuery:
		if m.Immutable {
			return ErrImmutableModel
		}
	}
	return nil
}
