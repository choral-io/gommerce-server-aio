package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type UserDevice struct {
	bun.BaseModel `bun:"table:user_devices,alias:user_device"`

	// Columns
	UserId    string       `bun:"user_id,pk"`
	DeviceId  string       `bun:"device_id,pk"`
	CreatedAt time.Time    `bun:"created_at"`
	UpdatedAt sql.NullTime `bun:"updated_at"`
	DeletedAt sql.NullTime `bun:"deleted_at,soft_delete,nullzero"`
}

func (m *UserDevice) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.CreatedAt = time.Now()
		m.UpdatedAt = sql.NullTime{Valid: false}
		m.DeletedAt = sql.NullTime{Valid: false}
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	}
	return nil
}
