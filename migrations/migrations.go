package migrations

import (
	"embed"
	"fmt"
	"regexp"

	"github.com/TakeAway-Inc/platform/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func Migrate(log *logger.Logger, fs *embed.FS, dbUrl string) {
	source, err := iofs.New(fs, "migrations")
	if err != nil {
		log.Error("failed to read migrations source", err)

		return
	}

	instance, err := migrate.NewWithSourceInstance("iofs", source, makeMigrateUrl(dbUrl))
	if err != nil {
		log.Error("failed to initialization the migrate instance", err)

		return
	}

	err = instance.Up()

	switch err {
	case nil:
		log.Info("the migration schema successfully upgraded!")
	case migrate.ErrNoChange:
		log.Info("the migration schema not changed")
	default:
		log.Error("could not apply the migration schema", err)
	}
}

func makeMigrateUrl(dbUrl string) string {
	urlRe := regexp.MustCompile("^[^\\?]+")
	url := urlRe.FindString(dbUrl)

	sslModeRe := regexp.MustCompile("(sslmode=)[a-zA-Z0-9]+")
	sslMode := sslModeRe.FindString(dbUrl)

	return fmt.Sprintf("%s?%s", url, sslMode)
}
