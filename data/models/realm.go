package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/choral-io/gommerce-server-core/data"
	"github.com/uptrace/bun"
)

type Realm struct {
	bun.BaseModel `bun:"table:realms,alias:realm"`

	// Columns
	Id          string         `json:"id" bun:"id,pk"`
	Disabled    bool           `json:"disabled" bun:"disabled"`
	Immutable   bool           `json:"immutable" bun:"immutable"`
	CreatedAt   time.Time      `json:"created_at" bun:"created_at"`
	UpdatedAt   sql.NullTime   `json:"updated_at" bun:"updated_at"`
	DeletedAt   sql.NullTime   `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`
	Name        string         `json:"name" bun:"name"`
	Title       string         `json:"title" bun:"title"`
	Description sql.NullString `json:"description" bun:"description"`
}

func (m *Realm) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if idw := data.IdWorkerFromContext(ctx); idw != nil {
			m.Id = idw.NextHex()
		}
		m.Immutable = false
		m.CreatedAt = time.Now()
		m.UpdatedAt = sql.NullTime{Valid: false}
		m.DeletedAt = sql.NullTime{Valid: false}
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	case *bun.DeleteQuery:
		m.UpdatedAt = m.DeletedAt
	}
	return nil
}
