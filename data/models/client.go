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
	Id          string         `json:"id" bun:"id,pk"`
	Disabled    bool           `json:"disabled" bun:"disabled"`
	Immutable   bool           `json:"immutable" bun:"immutable"`
	CreatedAt   time.Time      `json:"created_at" bun:"created_at"`
	UpdatedAt   sql.NullTime   `json:"updated_at" bun:"updated_at"`
	DeletedAt   sql.NullTime   `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`
	ExpiresAt   sql.NullTime   `json:"expires_at" bun:"expires_at"`
	SecretKey   string         `json:"secret_key" bun:"secret_key"`
	SecretCode  sql.NullString `json:"_" bun:"secret_code"`
	Description sql.NullString `json:"description" bun:"description"`
}

func (m *Client) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.Id = data.DefaultIdWorker().NextHex()
		m.Immutable = false
		m.CreatedAt = time.Now()
		m.UpdatedAt = sql.NullTime{Valid: false}
		m.DeletedAt = sql.NullTime{Valid: false}
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	}
	return nil
}
