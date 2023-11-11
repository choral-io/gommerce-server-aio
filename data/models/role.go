package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/choral-io/gommerce-server-core/data"
	"github.com/uptrace/bun"
)

type Role struct {
	bun.BaseModel `bun:"table:roles,alias:role"`

	// Columns
	Id          string         `json:"id" bun:"id,pk"`
	RealmId     string         `json:"realm_id" bun:"realm_id"`
	Disabled    bool           `json:"disabled" bun:"disabled"`
	Immutable   bool           `json:"immutable" bun:"immutable"`
	CreatedAt   time.Time      `json:"created_at" bun:"created_at"`
	UpdatedAt   sql.NullTime   `json:"updated_at" bun:"updated_at"`
	DeletedAt   sql.NullTime   `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`
	Name        string         `json:"name" bun:"name"`
	Description sql.NullString `json:"description" bun:"description"`
}

func (m *Role) BeforeAppendModel(ctx context.Context, query bun.Query) error {
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
	}
	return nil
}
