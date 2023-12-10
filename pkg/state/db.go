package state

import (
	"database/sql"
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

func (d *Db) initDb() (*sql.DB, error) {
	return sql.Open("sqlite3", filepath.Join(d.folder, "state.db"))
}

func (d *Db) Sync(config *conf.Config, inventory *ansible.Inventory) error {
	db, err := d.initDb()
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
	for groupname, _ := range config.ServerGroups {
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

	for group, hostGroup := range inventory.All.Children {
		for _, host := range hostGroup.GetHosts() {
			if hostStruct, found := hostGroup.Hosts[host]; found {
				stmt, err := db.Prepare(`
				INSERT INTO servers(id, public_ip, private_ip, hostname, is_lb_target, instance_type, server_group_id)
				VALUES (?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT(id)
				DO UPDATE SET public_ip = excluded.public_ip, private_ip = excluded.private_ip, hostname = excluded.hostname, is_lb_target = excluded.is_lb_target, instance_type = excluded.instance_type, server_group_id = excluded.server_group_id;
				`)
				if err != nil {
					return err
				}
				defer stmt.Close() //nolint: all
				// instance type and lb target fields need to be calculted
				instanceType := config.ServerGroups[group].InstanceType
				isLbTarget := config.ServerGroups[group].LbTarget
				_, err = stmt.Exec(hostStruct.ID, hostStruct.PublicIP, hostStruct.PrivateIP, hostStruct.HostName, isLbTarget, instanceType, group)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
