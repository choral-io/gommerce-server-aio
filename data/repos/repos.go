package repos

import (
	"context"
	"database/sql"
)

type DataRepos interface {
	RunInTx(context.Context, *sql.TxOptions, func(context.Context, DataRepos) error) error
	Clients() ClientRepo
	Realms() RealmRepo
	Users() UserRepo
	Roles() RoleRepo
	Logins() LoginRepo
}
