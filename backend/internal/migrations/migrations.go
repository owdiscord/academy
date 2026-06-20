// Package migrations contains our embedded dbmate backed migration files
package migrations

import (
	"embed"
	"fmt"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
)

//go:embed *.sql
var fs embed.FS

func Migrate(databaseURI string) error {
	url, err := url.Parse(databaseURI)
	if err != nil {
		return err
	}

	db := dbmate.New(url)
	db.FS = fs
	db.MigrationsDir = []string{"."}

	fmt.Println("Migrations:")
	migrations, err := db.FindMigrations()
	if err != nil {
		return err
	}
	for _, m := range migrations {
		fmt.Print(m.Version, m.FilePath)
	}

	fmt.Println("\nApplying...")
	err = db.CreateAndMigrate()
	if err != nil {
		return err
	}

	return nil
}
