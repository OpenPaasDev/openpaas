package state

import (
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var fs embed.FS

func Migrate(dbDir string) error {
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	db, err := sql.Open("sqlite3", filepath.Join(dbDir, "state.db"))
	if err != nil {
		return err
	}
	defer func() {
		e := db.Close()
		fmt.Printf("Error closing db %v", e)
	}()
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs",
		d,
		"sqlite3", driver)

	if err != nil {
		return err
	}
	return m.Up()

}
