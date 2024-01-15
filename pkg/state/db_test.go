package state

import (
	"os"
	"strings"
	"testing"

	"github.com/OpenPaasDev/core/pkg/ansible"
	"github.com/OpenPaasDev/core/pkg/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSync(t *testing.T) {
	inv, err := ansible.LoadInventory("testdata/inventory")
	require.NoError(t, err)
	assert.NotNil(t, inv)

	cnf, err := conf.Load("testdata/config.yaml")
	require.NoError(t, err)
	assert.NotNil(t, cnf)

	defer os.Remove("testdata/foo.db") //nolint

	db := InitWithName("testdata", "foo.db")
	err = db.Sync(cnf, inv)
	require.NoError(t, err)
	dbx, err := db.initDb()
	require.NoError(t, err)

	rows, err := dbx.Query("SELECT * FROM datacenters")
	require.NoError(t, rows.Err())
	require.NoError(t, err)

	i := 0
	for rows.Next() {
		i++
		var id, region string // Adjust the types according to your table columns
		// Scan the row's columns
		err = rows.Scan(&id, &region) // Modify the number of arguments & types according to your table
		require.NoError(t, err)
		assert.Equal(t, "hetzner", id)
		assert.Equal(t, "fsn1", region)
	}
	assert.Equal(t, 1, i)
	rows.Close() //nolint

	rows, err = dbx.Query("SELECT id FROM server_groups")
	require.NoError(t, rows.Err())
	require.NoError(t, err)
	serverGroups := []string{}
	for rows.Next() {
		i++
		var id string // Adjust the types according to your table columns
		// Scan the row's columns
		err = rows.Scan(&id) // Modify the number of arguments & types according to your table
		require.NoError(t, err)
		serverGroups = append(serverGroups, id)
	}

	assert.Contains(t, serverGroups, "server")
	assert.Contains(t, serverGroups, "agent")
	rows.Close() //nolint

	rows, err = dbx.Query("SELECT id, public_ip, private_ip, hostname, is_lb_target, instance_type, server_group_id FROM servers") //nolint
	require.NoError(t, err)
	i = 0
	for rows.Next() {
		i++
		var id, publicIP, privateIP, hostname, instanceType, serverGroup string // Adjust the types according to your table columns
		var isLBTarget bool
		// Scan the row's columns
		err = rows.Scan(&id, &publicIP, &privateIP, &hostname, &isLBTarget, &instanceType, &serverGroup) // Modify the number of arguments & types according to your table
		require.NoError(t, err)
		assert.True(t, strings.HasPrefix(privateIP, "10.0"))
		assert.Equal(t, "cx21", instanceType)
		serverGroups = append(serverGroups, id)
		if serverGroup == "server" {
			assert.True(t, strings.HasPrefix(hostname, "prod-server"))
			assert.False(t, isLBTarget)
		} else {
			assert.True(t, strings.HasPrefix(hostname, "prod-agent"))
			assert.True(t, isLBTarget)
		}
	}

	assert.Equal(t, 5, i)
	rows.Close() //nolint
}
