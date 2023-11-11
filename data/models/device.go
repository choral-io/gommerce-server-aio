package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/choral-io/gommerce-server-core/data"
	"github.com/uptrace/bun"
)

type Device struct {
	bun.BaseModel `bun:"table:devices,alias:device"`

	// Columns
	Id        string            `json:"id" bun:"id,pk"`
	UserId    sql.NullString    `json:"user_id" bun:"user_id"`
	ClientId  sql.NullString    `json:"client_id" bun:"client_id"`
	CreatedAt time.Time         `json:"created_at" bun:"created_at"`
	UpdatedAt sql.NullTime      `json:"updated_at" bun:"updated_at"`
	TraceCode string            `json:"trace_code" bun:"trace_code"`
	PushToken sql.NullString    `json:"_" bun:"push_token"`
	Metadata  map[string]string `json:"metadata" bun:"metadata,json_use_number"`
}

func (m *Device) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if idw := data.IdWorkerFromContext(ctx); idw != nil {
			m.Id = idw.NextHex()
		}
		m.CreatedAt = time.Now()
		m.UpdatedAt = sql.NullTime{Valid: false}
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	}
	return nil
}
