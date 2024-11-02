package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"

	"github.com/choral-io/gommerce-server-core/data"
)

type Device struct {
	bun.BaseModel `bun:"table:devices,alias:device"`

	// Columns
	Id        string            `bun:"id,pk"`
	UserId    sql.NullString    `bun:"user_id"`
	ClientId  sql.NullString    `bun:"client_id"`
	CreatedAt time.Time         `bun:"created_at"`
	UpdatedAt sql.NullTime      `bun:"updated_at"`
	TraceCode string            `bun:"trace_code"`
	PushToken sql.NullString    `bun:"push_token"`
	Metadata  map[string]string `bun:"metadata,json_use_number"`
}

func (m *Device) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if m.Id == "" {
			m.Id = data.DefaultIdWorker().NextHex()
		}
		m.CreatedAt = time.Now()
		m.UpdatedAt = sql.NullTime{Valid: false}
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	}
	return nil
}
