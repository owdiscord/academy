// Package migrations contains our embedded dbmate backed migration files
package migrations

import (
	"embed"
	"log/slog"
	"net/url"
	"strings"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
)

//go:embed *.sql
var fs embed.FS

func Migrate(databaseURI string) error {
	if !strings.HasPrefix(databaseURI, "mysql://") {
		databaseURI = "mysql://" + databaseURI
	}

	url, err := url.Parse(databaseURI)
	if err != nil {
		return err
	}

	db := dbmate.New(url)
	db.FS = fs
	db.MigrationsDir = []string{"."}

	slog.Default().Info("running database migrations", "user", url.User.Username(), "database", url.Path)
	migrations, err := db.FindMigrations()
	if err != nil {
		return err
	}

	for _, m := range migrations {
		slog.Default().Info("migrating", "version", m.Version, "path", m.FilePath, "applied", m.Applied)
	}

	err = db.CreateAndMigrate()
	if err != nil {
		return err
	}

	slog.Default().Info("migrations complete")

	return nil
}
