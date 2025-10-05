package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-core/data"
)

type ChatSession struct {
	bun.BaseModel `bun:"table:chat_sessions,alias:chat_session"`

	// Columns
	Id           string         `bun:"id,pk"`
	Readonly     bool           `bun:"readonly"`
	CreatedAt    time.Time      `bun:"created_at"`
	UpdatedAt    sql.NullTime   `bun:"updated_at"`
	DeletedAt    sql.NullTime   `bun:"deleted_at,nullzero,soft_delete"`
	IconUrl      sql.NullString `bun:"icon_url"`
	Title        string         `bun:"title"`
	Introduction sql.NullString `bun:"introduction"`

	// Relations
	Records []*ChatRecord `bun:"rel:has-many,join:id=session_id"`
	Members []*ChatMember `bun:"rel:has-many,join:id=session_id"`
}

func (m *ChatSession) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if m.Id == "" {
			m.Id = data.DefaultIDWorker().NextHex()
		}
		m.UpdatedAt = sql.NullTime{Valid: false}
		m.DeletedAt = sql.NullTime{Valid: false}
		m.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	}
	return nil
}
