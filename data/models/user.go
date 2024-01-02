package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/choral-io/gommerce-server-core/data"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:user"`

	// Columns
	Id             string            `json:"id" bun:"id,pk"`
	RealmId        string            `json:"realm_id" bun:"realm_id"`
	CreatorId      sql.NullString    `json:"creator_id" bun:"creator_id"`
	Disabled       bool              `json:"disabled" bun:"disabled"`
	Approved       bool              `json:"approved" bun:"approved"`
	Verified       bool              `json:"verified" bun:"verified"`
	Immutable      bool              `json:"immutable" bun:"immutable"`
	CreatedAt      time.Time         `json:"created_at" bun:"created_at"`
	UpdatedAt      sql.NullTime      `json:"updated_at" bun:"updated_at"`
	DeletedAt      sql.NullTime      `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`
	ExpiresAt      sql.NullTime      `json:"expires_at" bun:"expires_at"`
	FirstLoginTime sql.NullTime      `json:"first_login_time" bun:"first_login_time"`
	LastActiveTime sql.NullTime      `json:"last_active_time" bun:"last_active_time"`
	Flags          int64             `json:"flags" bun:"flags"`
	Attributes     map[string]string `json:"attributes" bun:"attributes,json_use_number"`
	PhoneNumber    sql.NullString    `json:"-" bun:"phone_number"`
	EmailAddress   sql.NullString    `json:"-" bun:"email_address"`
	Description    sql.NullString    `json:"description" bun:"description"`

	// Relations
	Realm   *Realm   `bun:"rel:belongs-to,join:realm_id=id"`
	Profile *Profile `bun:"rel:has-one,join:id=id"`
	Creator *User    `bun:"rel:belongs-to,join:creator_id=id"`
}

func (m *User) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.Id = data.DefaultIdWorker().NextHex()
		m.Immutable = false
		m.CreatedAt = time.Now()
		m.UpdatedAt = sql.NullTime{Valid: false}
		m.DeletedAt = sql.NullTime{Valid: false}
	case *bun.UpdateQuery:
		m.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	}
	return nil
}
