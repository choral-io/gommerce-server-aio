package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type RoleUser struct {
	bun.BaseModel `bun:"table:role_users,alias:role_user"`

	// Columns
	RoleId    string       `bun:"role_id,pk"`
	UserId    string       `bun:"user_id,pk"`
	Immutable bool         `bun:"immutable"`
	CreatedAt time.Time    `bun:"created_at"`
	UpdatedAt sql.NullTime `bun:"updated_at"`
	DeletedAt sql.NullTime `bun:"deleted_at,soft_delete,nullzero"`

	// Relations
	Role *Role `bun:"rel:belongs-to,join:role_id=id"`
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

func (m *RoleUser) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.Immutable = false
		m.CreatedAt = time.Now()
		m.UpdatedAt = sql.NullTime{Valid: false}
		m.DeletedAt = sql.NullTime{Valid: false}
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	}
	return nil
}
