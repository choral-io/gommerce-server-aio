package v1beta

import (
	"database/sql"

	iam "github.com/choral-io/gommerce-protobuf-go/iam/v1beta"
	gender "github.com/choral-io/gommerce-protobuf-go/types/v1/gender"
	sqlpb "github.com/choral-io/gommerce-protobuf-go/types/v1/sqlpb"
	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-core/secure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toUserPB(u models.User) *iam.User {
	// mask sensitive data
	// convert to proto
	r := &iam.User{
		Id:                 u.Id,
		Disabled:           u.Disabled,
		Approved:           u.Approved,
		Verified:           u.Verified,
		Immutable:          u.Immutable,
		CreatedAt:          timestamppb.New(u.CreatedAt),
		UpdatedAt:          sqlpb.FromNullTime(u.UpdatedAt),
		ExpiresAt:          sqlpb.FromNullTime(u.ExpiresAt),
		DeletedAt:          sqlpb.FromNullTime(u.DeletedAt),
		FirstLoginTime:     sqlpb.FromNullTime(u.FirstLoginTime),
		LastActiveTime:     sqlpb.FromNullTime(u.LastActiveTime),
		Flags:              u.Flags,
		Attributes:         u.Attributes,
		DisplayName:        sqlpb.FromNullString(u.DisplayName),
		Gender:             gender.FromSqlNullString(u.Gender),
		MaskedPhoneNumber:  sqlpb.FromNullString(sql.NullString{Valid: u.PhoneNumber.Valid, String: secure.MaskString(u.PhoneNumber.String)}),
		MaskedEmailAddress: sqlpb.FromNullString(sql.NullString{Valid: u.EmailAddress.Valid, String: secure.MaskString(u.EmailAddress.String)}),
		Description:        sqlpb.FromNullString(u.Description),
	}
	if u.Realm != nil {
		r.Realm = u.Realm.Name
	}
	if u.Creator != nil {
		r.Creator = toUserPB(*u.Creator)
	}
	return r
}
