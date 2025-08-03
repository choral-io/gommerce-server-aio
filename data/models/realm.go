package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-core/data"
)

const (
	RealmFlagsAllowRegistration int64 = 0b_0000_0001
)

type Realm struct {
	bun.BaseModel `bun:"table:realms,alias:realm"`

	// Columns
	Id          string         `bun:"id,pk"`
	Disabled    bool           `bun:"disabled"`
	Immutable   bool           `bun:"immutable"`
	CreatedAt   time.Time      `bun:"created_at"`
	UpdatedAt   sql.NullTime   `bun:"updated_at"`
	DeletedAt   sql.NullTime   `bun:"deleted_at,soft_delete,nullzero"`
	Flags       int64          `bun:"flags"`
	Name        string         `bun:"name"`
	Title       string         `bun:"title"`
	Description sql.NullString `bun:"description"`
}

func (m *Realm) BeforeAppendModel(ctx context.Context, query bun.Query) error {
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

func (m *Realm) AllowRegistration() bool {
	return m.Flags&RealmFlagsAllowRegistration != 0
}
