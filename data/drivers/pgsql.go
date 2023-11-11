package data

import (
	"database/sql"

	"github.com/uptrace/bun/driver/pgdriver"
)

func init() {
	// alias pg to pgsql
	sql.Register("pgsql", pgdriver.NewDriver())
}
