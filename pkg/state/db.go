package state

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/OpenPaasDev/core/pkg/ansible"
	"github.com/OpenPaasDev/core/pkg/conf"
	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	folder string
}

func Init(folder string) *Db {
	return &Db{folder: folder}
}

func (d *Db) Sync(config *conf.Config, inventory *ansible.Inventory) error {
	db, err := sql.Open("sqlite3", filepath.Join(d.folder, "state.db"))
	if err != nil {
		return err
	}
	defer db.Close() //nolint: all

	statement, err := db.Prepare(`
        INSERT INTO datacenters(id, region) 
        VALUES (?, ?)
        ON CONFLICT(id) 
        DO UPDATE SET region = excluded.region;
    `)
	if err != nil {
		return err
	}
	defer statement.Close() //nolint: all
	_, err = statement.Exec(config.DC, config.CloudProviderConfig.ProviderSettings["location"])
	if err != nil {
		return err
	}
	for groupname, group := range config.ServerGroups {
		fmt.Println(groupname)
		fmt.Println(group)
		stmt, err := db.Prepare(`
			INSERT INTO server_groups(id, dc_id)
			VALUES (?, ?)
			ON CONFLICT(id)
			DO UPDATE SET dc_id = excluded.dc_id;
		`)
		if err != nil {
			return err
		}
		defer stmt.Close() //nolint: all
		_, err = stmt.Exec(groupname, config.DC)
		if err != nil {
			return err
		}

	}

	return nil
}
