package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/choral-io/gommerce-server-aio/data/drivers" // register db drivers
	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/schema"
	"golang.org/x/crypto/bcrypt"
)

const (
	ansi_reset     = "\033[0m"
	ansi_red       = "\033[31m"
	ansi_green     = "\033[32m"
	ansi_yellow    = "\033[33m"
	ansi_blue      = "\033[34m"
	base58_symbols = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

func init() {
	godotenv.Overload("./prisma/.env")
}

func main() {
	log.SetFlags(0)
	log.Printf("%sSeeding database...%s", ansi_blue, ansi_reset)
	if err := seed(context.Background()); err != nil {
		log.Printf("%sfailed to seed database: %s%v%s", ansi_yellow, ansi_red, err, ansi_reset)
		os.Exit(1)
	}
	log.Printf("%sDatabase seeded.%s", ansi_green, ansi_reset)
}

func seed(ctx context.Context) error {
	// load env vars
	driver := os.Getenv("GO_SQL_DATA_DRIVER")
	if driver == "" {
		driver = "pgsql"
	}
	source := os.Getenv("GO_SQL_DATA_SOURCE")
	if source == "" {
		return fmt.Errorf("GO_SQL_DATA_SOURCE is not set")
	}

	// create bun db
	var dialect schema.Dialect
	switch driver {
	case "pg", "pgsql":
		dialect = pgdialect.New()
	case "mysql":
		dialect = mysqldialect.New()
	case "mssql":
		dialect = mssqldialect.New()
	default:
		return fmt.Errorf("unsupported driver: %s", driver)
	}
	var bdb bun.IDB
	if sdb, err := sql.Open(driver, source); err != nil {
		return err
	} else {
		if err := sdb.Ping(); err != nil {
			return err
		}
		bdb = bun.NewDB(sdb, dialect, bun.WithDiscardUnknownColumns())
	}

	// insert data
	return bdb.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		adminRealm := models.Realm{
			Immutable: true,
			Flags:     0b0000,
			Name:      "admin",
			Title:     "Admin",
		}
		if _, err := tx.NewInsert().Model(&adminRealm).Exec(ctx); err != nil {
			return err
		}

		usersRealm := models.Realm{
			Immutable: true,
			Flags:     models.REALM_FLAGS_ALLOW_REGISTRATION,
			Name:      "users",
			Title:     "Users",
		}
		if _, err := tx.NewInsert().Model(&usersRealm).Exec(ctx); err != nil {
			return err
		}

		adminRole := models.Role{
			RealmId:     adminRealm.Id,
			Immutable:   true,
			Name:        "Admin",
			Description: sql.NullString{Valid: true, String: "Built-in admin role."},
		}
		if _, err := tx.NewInsert().Model(&adminRole).Exec(ctx); err != nil {
			return err
		}

		adminUser := models.User{
			RealmId:   adminRealm.Id,
			Approved:  true,
			Verified:  true,
			Immutable: true,
			Flags:     0b0000,
			Attributes: map[string]string{
				"profile.display_name": "Admin",
			},
			Description: sql.NullString{Valid: true, String: "Built-in admin user."},
		}
		if _, err := tx.NewInsert().Model(&adminUser).Exec(ctx); err != nil {
			return err
		}

		adminProfile := models.Profile{
			Id:          adminUser.Id,
			DisplayName: sql.NullString{Valid: true, String: "Admin"},
		}
		if _, err := tx.NewInsert().Model(&adminProfile).Exec(ctx); err != nil {
			return err
		}

		adminLogin := models.Login{
			UserId:     adminUser.Id,
			Immutable:  true,
			Provider:   models.LOGIN_PROVIDER_FORM_PASSWORD,
			Identifier: "admin",
			Metadata:   map[string]string{},
		}
		if pwd, err := secure.RandString(16, base58_symbols); err != nil {
			return err
		} else if hp, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost); err != nil {
			return err
		} else {
			adminLogin.Credential = sql.NullString{Valid: true, String: string(hp)}
			log.Printf("%susing randomly generated password for admin user:        %s%s%s", ansi_blue, ansi_yellow, pwd, ansi_reset)
		}
		if _, err := tx.NewInsert().Model(&adminLogin).Exec(ctx); err != nil {
			return err
		}

		roleUsers := []models.RoleUser{
			{RoleId: adminRole.Id, UserId: adminUser.Id},
		}
		if _, err := tx.NewInsert().Model(&roleUsers).Exec(ctx); err != nil {
			return err
		}

		consoleClient := models.Client{
			Immutable:   true,
			Description: sql.NullString{Valid: true, String: "Web-based console client."},
		}
		if pwd, err := secure.RandString(16, base58_symbols); err != nil {
			return err
		} else {
			consoleClient.SecretKey = pwd
			log.Printf("%susing sequence generated secret key for console client:  %s%s%s", ansi_blue, ansi_yellow, pwd, ansi_reset)
		}
		if pwd, err := secure.RandString(32, base58_symbols); err != nil {
			return err
		} else if hp, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost); err != nil {
			return err
		} else {
			consoleClient.SecretCode = sql.NullString{Valid: true, String: string(hp)}
			log.Printf("%susing randomly generated secret code for console client: %s%s%s", ansi_blue, ansi_yellow, pwd, ansi_reset)
		}
		if _, err := tx.NewInsert().Model(&consoleClient).Exec(ctx); err != nil {
			return err
		}

		clientUsers := []models.ClientUser{
			{ClientId: consoleClient.Id, UserId: adminUser.Id, Immutable: true},
		}
		if _, err := tx.NewInsert().Model(&clientUsers).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}
