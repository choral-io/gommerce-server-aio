package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/choral-io/gommerce-server-core/data"
	"github.com/uptrace/bun"
)

const (
	USER_PROFILE_DISPLAY_NAME_ATTRIBUTE = "profile.display_name"
	USER_PROFILE_AVATAR_URL_ATTRIBUTE   = "profile.avatar_url"
	USER_PROFILE_GENDER_ATTRIBUTE       = "profile.gender"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:user"`

	// Columns
	Id             string            `bun:"id,pk"`
	RealmId        string            `bun:"realm_id"`
	CreatorId      sql.NullString    `bun:"creator_id"`
	Disabled       bool              `bun:"disabled"`
	Approved       bool              `bun:"approved"`
	Verified       bool              `bun:"verified"`
	Immutable      bool              `bun:"immutable"`
	CreatedAt      time.Time         `bun:"created_at"`
	UpdatedAt      sql.NullTime      `bun:"updated_at"`
	DeletedAt      sql.NullTime      `bun:"deleted_at,soft_delete,nullzero"`
	ExpiresAt      sql.NullTime      `bun:"expires_at"`
	FirstLoginTime sql.NullTime      `bun:"first_login_time"`
	LastActiveTime sql.NullTime      `bun:"last_active_time"`
	Flags          int64             `bun:"flags"`
	Attributes     map[string]string `bun:"attributes,json_use_number"`
	PhoneNumber    sql.NullString    `bun:"phone_number"`
	EmailAddress   sql.NullString    `bun:"email_address"`
	Description    sql.NullString    `bun:"description"`

	// Relations
	Realm   *Realm   `bun:"rel:belongs-to,join:realm_id=id"`
	Profile *Profile `bun:"rel:has-one,join:id=id"`
	Creator *User    `bun:"rel:belongs-to,join:creator_id=id"`
}

func (m *User) BeforeAppendModel(_ context.Context, query bun.Query) error {
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
