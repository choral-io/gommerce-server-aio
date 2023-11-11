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
	RoleId    string       `json:"role_id" bun:"role_id,pk"`
	UserId    string       `json:"user_id" bun:"user_id,pk"`
	Immutable bool         `json:"immutable" bun:"immutable"`
	CreatedAt time.Time    `json:"created_at" bun:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at" bun:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`

	// Relations
	Role *Role `bun:"rel:belongs-to,join:role_id=id"`
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

func (m *RoleUser) BeforeAppendModel(ctx context.Context, query bun.Query) error {
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
