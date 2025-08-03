package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/schema"
	"golang.org/x/crypto/bcrypt"

	"github.com/choral-io/gommerce-server-aio/data/models"
	"github.com/choral-io/gommerce-server-core/secure"
)

func init() {
	// alias pg to pgsql
	sql.Register("pgsql", pgdriver.NewDriver())
}

const (
	ansiReset  = "\033[0m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiBlue   = "\033[34m"
	b58Chars   = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

func main() {
	_ = godotenv.Load("prisma/.env")
	os.Setenv("DATA_SEEDING_MODE", "true")
	log.SetFlags(0)
	log.Printf("%sSeeding database...%s", ansiBlue, ansiReset)
	if err := seed(context.Background()); err != nil {
		log.Printf("%sfailed to seed database: %s%v%s", ansiYellow, ansiRed, err, ansiReset)
		os.Exit(1)
	}
	log.Printf("%sDatabase seeded.%s", ansiGreen, ansiReset)
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
		if err := models.RegisterModels(bdb); err != nil {
			return err
		}
	}

	// insert data
	return bdb.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		systemRealm := models.Realm{
			Immutable: true,
			Flags:     0b0000,
			Name:      "system",
			Title:     "System",
		}
		if _, err := tx.NewInsert().Model(&systemRealm).Exec(ctx); err != nil {
			return err
		}

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
			Flags:     models.RealmFlagsAllowRegistration,
			Name:      "users",
			Title:     "Users",
		}
		if _, err := tx.NewInsert().Model(&usersRealm).Exec(ctx); err != nil {
			return err
		}

		systemUser := models.User{
			RealmId:   systemRealm.Id,
			Approved:  true,
			Verified:  true,
			Immutable: true,
			Flags:     0b0000,
			Attributes: map[string]string{
				models.USER_PROFILE_DISPLAY_NAME_ATTRIBUTE: "$SYSTEM",
			},
			Description: sql.NullString{Valid: true, String: "Built-in system user."},
		}
		if _, err := tx.NewInsert().Model(&systemUser).Exec(ctx); err != nil {
			return err
		}

		systemProfile := models.Profile{
			Id:          systemUser.Id,
			DisplayName: "$SYSTEM",
		}
		if _, err := tx.NewInsert().Model(&systemProfile).Exec(ctx); err != nil {
			return err
		}

		adminUser := models.User{
			RealmId:   adminRealm.Id,
			Approved:  true,
			Verified:  true,
			Immutable: true,
			Flags:     0b0000,
			Attributes: map[string]string{
				models.USER_PROFILE_DISPLAY_NAME_ATTRIBUTE: "Admin",
			},
			Description: sql.NullString{Valid: true, String: "Built-in admin user."},
		}
		if _, err := tx.NewInsert().Model(&adminUser).Exec(ctx); err != nil {
			return err
		}

		adminProfile := models.Profile{
			Id:          adminUser.Id,
			DisplayName: "Admin",
		}
		if _, err := tx.NewInsert().Model(&adminProfile).Exec(ctx); err != nil {
			return err
		}

		adminLogin := models.Login{
			RealmId:    adminUser.RealmId,
			UserId:     adminUser.Id,
			Immutable:  true,
			Provider:   models.LoginProviderFormPassword,
			Identifier: "admin",
			Metadata:   map[string]string{},
		}
		if pwd, err := secure.RandString(16, b58Chars); err != nil {
			return err
		} else if hp, err := bcrypt.GenerateFromPassword([]byte(pwd), 12); err != nil {
			return err
		} else {
			adminLogin.Credential = sql.NullString{Valid: true, String: string(hp)}
			log.Printf("%susing randomly generated password for admin user:        %s%s%s", ansiBlue, ansiYellow, pwd, ansiReset)
		}
		if _, err := tx.NewInsert().Model(&adminLogin).Exec(ctx); err != nil {
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

		roleUsers := []models.RoleUser{
			{
				RoleId:    adminRole.Id,
				UserId:    adminUser.Id,
				Immutable: true,
			},
		}
		if _, err := tx.NewInsert().Model(&roleUsers).Exec(ctx); err != nil {
			return err
		}

		consoleClient := models.Client{
			Immutable:   true,
			Description: sql.NullString{Valid: true, String: "Management console client."},
		}
		if pwd, err := secure.RandString(16, b58Chars); err != nil {
			return err
		} else {
			consoleClient.SecretKey = pwd
			log.Printf("%susing randomly generated secret key for console client:  %s%s%s", ansiBlue, ansiYellow, pwd, ansiReset)
		}
		if pwd, err := secure.RandString(32, b58Chars); err != nil {
			return err
		} else if hp, err := bcrypt.GenerateFromPassword([]byte(pwd), 12); err != nil {
			return err
		} else {
			consoleClient.SecretCode = sql.NullString{Valid: true, String: string(hp)}
			log.Printf("%susing randomly generated secret code for console client: %s%s%s", ansiBlue, ansiYellow, pwd, ansiReset)
		}
		if _, err := tx.NewInsert().Model(&consoleClient).Exec(ctx); err != nil {
			return err
		}

		portalClient := models.Client{
			Immutable:   true,
			Description: sql.NullString{Valid: true, String: "User portal client."},
		}
		if pwd, err := secure.RandString(16, b58Chars); err != nil {
			return err
		} else {
			portalClient.SecretKey = pwd
			log.Printf("%susing randomly generated secret key for portal client:   %s%s%s", ansiBlue, ansiYellow, pwd, ansiReset)
		}
		if pwd, err := secure.RandString(32, b58Chars); err != nil {
			return err
		} else if hp, err := bcrypt.GenerateFromPassword([]byte(pwd), 12); err != nil {
			return err
		} else {
			portalClient.SecretCode = sql.NullString{Valid: true, String: string(hp)}
			log.Printf("%susing randomly generated secret code for portal client:  %s%s%s", ansiBlue, ansiYellow, pwd, ansiReset)
		}
		if _, err := tx.NewInsert().Model(&portalClient).Exec(ctx); err != nil {
			return err
		}

		clientUsers := []models.ClientUser{
			{
				ClientId:  consoleClient.Id,
				UserId:    adminUser.Id,
				Immutable: true,
			},
		}
		if _, err := tx.NewInsert().Model(&clientUsers).Exec(ctx); err != nil {
			return err
		}

		chatSession := models.ChatSession{
			Readonly: true,
			Title:    "$SYSTEM",
		}
		if _, err := tx.NewInsert().Model(&chatSession).Exec(ctx); err != nil {
			return err
		}

		chatMembers := []models.ChatMember{
			{
				UserId:     systemUser.Id,
				SessionId:  chatSession.Id,
				Permission: models.ChatMemberPermissionOwner,
			},
			{
				UserId:     adminUser.Id,
				SessionId:  chatSession.Id,
				Permission: models.ChatMemberPermissionMember,
			},
		}
		if _, err := tx.NewInsert().Model(&chatMembers).Exec(ctx); err != nil {
			return err
		}

		chatRecord := models.ChatRecord{
			SessionId: chatSession.Id,
			CreatorId: systemUser.Id,
			Version:   "0.0.1",
			Headers: map[string]string{
				"type": "text",
			},
			Content: []byte(`"Welcome to Gommerce."`),
		}
		if _, err := tx.NewInsert().Model(&chatRecord).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}
