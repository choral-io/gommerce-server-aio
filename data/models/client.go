package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/choral-io/gommerce-server-core/data"
	"github.com/uptrace/bun"
)

type Client struct {
	bun.BaseModel `bun:"table:clients,alias:client"`

	// Columns
	Id          string         `bun:"id,pk"`
	Disabled    bool           `bun:"disabled"`
	Immutable   bool           `bun:"immutable"`
	CreatedAt   time.Time      `bun:"created_at"`
	UpdatedAt   sql.NullTime   `bun:"updated_at"`
	DeletedAt   sql.NullTime   `bun:"deleted_at,soft_delete,nullzero"`
	ExpiresAt   sql.NullTime   `bun:"expires_at"`
	SecretKey   string         `bun:"secret_key"`
	SecretCode  sql.NullString `bun:"secret_code"`
	Description sql.NullString `bun:"description"`
}

func (m *Client) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if m.Id == "" {
			m.Id = data.DefaultIdWorker().NextHex()
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
