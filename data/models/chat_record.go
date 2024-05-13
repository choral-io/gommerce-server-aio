package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/choral-io/gommerce-server-core/data"
	"github.com/uptrace/bun"
)

type ChatRecord struct {
	bun.BaseModel `bun:"table:chat_records,alias:chat_record"`

	// Columns
	Id        string            `bun:"id,pk"`
	SessionId string            `bun:"session_id"`
	CreatorId string            `bun:"creator_id"`
	CreatedAt time.Time         `bun:"created_at"`
	UpdatedAt sql.NullTime      `bun:"updated_at"`
	DeletedAt sql.NullTime      `bun:"deleted_at,nullzero,soft_delete"`
	Version   string            `bun:"version"`
	Headers   map[string]string `bun:"headers,json_use_number"`
	Content   json.RawMessage   `bun:"content,json_use_number"`

	// Relations
	Session *ChatSession `bun:"rel:belongs-to,join:session_id=id"`
	Creator *User        `bun:"rel:belongs-to,join:creator_id=id"`
}

func (m *ChatRecord) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if m.Id == "" {
			m.Id = data.DefaultIdWorker().NextHex()
		}
		m.UpdatedAt = sql.NullTime{Valid: false}
		m.DeletedAt = sql.NullTime{Valid: false}
		m.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	}
	if m.Headers == nil {
		m.Headers = map[string]string{}
	}
	return nil
}
