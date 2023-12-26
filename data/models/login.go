package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/choral-io/gommerce-server-core/data"
	"github.com/uptrace/bun"
)

type Login struct {
	bun.BaseModel `bun:"table:logins,alias:login"`

	// Columns
	Id         string            `json:"id" bun:"id,pk"`
	UserId     string            `json:"user_id" bun:"user_id"`
	Disabled   bool              `json:"disabled" bun:"disabled"`
	Immutable  bool              `json:"immutable" bun:"immutable"`
	CreatedAt  time.Time         `json:"created_at" bun:"created_at"`
	UpdatedAt  sql.NullTime      `json:"updated_at" bun:"updated_at"`
	DeletedAt  sql.NullTime      `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`
	ExpiresAt  sql.NullTime      `json:"expires_at" bun:"expires_at"`
	Provider   string            `json:"provider" bun:"provider"`
	Identifier string            `json:"identifier" bun:"identifier"`
	Credential sql.NullString    `json:"-" bun:"credential"`
	Metadata   map[string]string `json:"metadata" bun:"metadata,json_use_number"`

	// Relations
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

func (m *Login) BeforeAppendModel(ctx context.Context, query bun.Query) error {
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
