package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

const (
	ChatMemberPermissionOwner  int64 = 0b_0000_0001_0000_0011
	ChatMemberPermissionAdmin  int64 = 0b_0000_0001_0000_0010
	ChatMemberPermissionMember int64 = 0b_0000_0001_0000_0000
)

type ChatMember struct {
	bun.BaseModel `bun:"table:chat_members,alias:chat_member"`

	// Columns
	UserId      string         `bun:"user_id,pk"`
	SessionId   string         `bun:"session_id,pk"`
	CreatedAt   time.Time      `bun:"created_at"`
	UpdatedAt   sql.NullTime   `bun:"updated_at"`
	ReadCursor  sql.NullString `bun:"read_cursor"`
	Permission  int64          `bun:"permission"`
	DisplayName sql.NullString `bun:"display_name"`

	// Relations
	User    *User        `bun:"rel:belongs-to,join:user_id=id"`
	Profile *Profile     `bun:"rel:belongs-to,join:user_id=id"`
	Session *ChatSession `bun:"rel:belongs-to,join:session_id=id"`
}

func (m *ChatMember) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = sql.NullTime{Valid: false}
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	}
	return nil
}
