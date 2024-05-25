package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/choral-io/gommerce-server-core/data"
	"github.com/uptrace/bun"
)

const (
	LoginProviderFormPassword = "FORM_PASSWORD"
	LoginProviderSmsOtpCode   = "SMS_OTP_CODE"
)

type Login struct {
	bun.BaseModel `bun:"table:logins,alias:login"`

	// Columns
	Id         string            `bun:"id,pk"`
	UserId     string            `bun:"user_id"`
	Disabled   bool              `bun:"disabled"`
	Immutable  bool              `bun:"immutable"`
	CreatedAt  time.Time         `bun:"created_at"`
	UpdatedAt  sql.NullTime      `bun:"updated_at"`
	DeletedAt  sql.NullTime      `bun:"deleted_at,soft_delete,nullzero"`
	ExpiresAt  sql.NullTime      `bun:"expires_at"`
	Provider   string            `bun:"provider"`
	Identifier string            `bun:"identifier"`
	Credential sql.NullString    `bun:"credential"`
	Metadata   map[string]string `bun:"metadata,json_use_number"`

	// Relations
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

func (m *Login) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if m.Id == "" {
			m.Id = data.DefaultIdWorker().NextHex()
		}
		m.Immutable = false
		m.CreatedAt = time.Now()
		m.UpdatedAt = sql.NullTime{Valid: false}
		m.DeletedAt = sql.NullTime{Valid: false}
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	}
	return nil
}
