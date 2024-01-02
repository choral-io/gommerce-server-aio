package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type Profile struct {
	bun.BaseModel `bun:"table:profiles,alias:profile"`

	// Columns
	Id           string         `json:"id" bun:"id,pk"`
	CreatedAt    time.Time      `json:"created_at" bun:"created_at"`
	UpdatedAt    sql.NullTime   `json:"updated_at" bun:"updated_at"`
	DisplayName  sql.NullString `json:"display_name" bun:"display_name"`
	AvatarUrl    sql.NullString `json:"avatar_url" bun:"avatar_url"`
	Gender       sql.NullString `json:"gender" bun:"gender"`
	Birthdate    sql.NullTime   `json:"birthdate" bun:"birthdate"`
	Introduction sql.NullString `json:"introduction" bun:"introduction"`

	// Relations
	User *User `bun:"rel:belongs-to,join:id=id"`
}

func (m *Profile) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = sql.NullTime{Valid: false}
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	}
	return nil
}
