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
	UserId    string       `json:"user_id" bun:"user_id,pk"`
	DeviceId  string       `json:"device_id" bun:"device_id,pk"`
	CreatedAt time.Time    `json:"created_at" bun:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at" bun:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`
}

func (m *UserDevice) BeforeAppendModel(ctx context.Context, query bun.Query) error {
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
