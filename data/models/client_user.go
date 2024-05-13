package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type ClientUser struct {
	bun.BaseModel `bun:"table:client_users,alias:client_user"`

	// Columns
	ClientId  string       `bun:"client_id,pk"`
	UserId    string       `bun:"user_id,pk"`
	Immutable bool         `bun:"immutable"`
	CreatedAt time.Time    `bun:"created_at"`
	UpdatedAt sql.NullTime `bun:"updated_at"`
	DeletedAt sql.NullTime `bun:"deleted_at,soft_delete,nullzero"`
}

func (m *ClientUser) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.Immutable = false
		m.UpdatedAt = sql.NullTime{Valid: false}
		m.DeletedAt = sql.NullTime{Valid: false}
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	}
	return nil
}
