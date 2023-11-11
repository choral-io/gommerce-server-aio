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
	ClientId  string       `json:"client_id" bun:"client_id,pk"`
	UserId    string       `json:"user_id" bun:"user_id,pk"`
	Immutable bool         `json:"immutable" bun:"immutable"`
	CreatedAt time.Time    `json:"created_at" bun:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at" bun:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`
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
