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
	Id           string         `bun:"id,pk"`
	CreatedAt    time.Time      `bun:"created_at"`
	UpdatedAt    sql.NullTime   `bun:"updated_at"`
	DisplayName  string         `bun:"display_name"`
	AvatarUrl    sql.NullString `bun:"avatar_url"`
	Gender       sql.NullString `bun:"gender"`
	Birthdate    sql.NullTime   `bun:"birthdate"`
	Introduction sql.NullString `bun:"introduction"`

	// Relations
	User *User `bun:"rel:belongs-to,join:id=id"`
}

func (m *Profile) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = sql.NullTime{Valid: false}
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	}
	return nil
}
