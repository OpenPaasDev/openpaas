package ansible

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/OpenPaasDev/openpaas/pkg/conf"
	"gopkg.in/yaml.v3"
)

type Inventory struct {
	Path string `yaml:"-"`
	All  All    `yaml:"all"`
}

type InventoryJson struct {
	Servers HostValues `json:"servers"`
	Volumes Volumes    `json:"volumes"`
}
type Volumes struct {
	Value []Volume `json:"value"`
}

type Volume struct {
	Mount    string `json:"mount"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	ServerID int    `json:"server_id"`
}

type HostValues struct {
	Value []Host `json:"value"`
}

type Host struct {
	Group     string `json:"group"`
	Host      string `json:"host"`
	HostName  string `json:"host_name"`
	PrivateIP string `json:"private_ip"`
	ServerID  string `json:"server_id"`
	Image     string `json:"image"`
}

type All struct {
	Children map[string]HostGroup `yaml:"children"`
}

type HostGroup struct {
	Hosts map[string]AnsibleHost `yaml:"hosts"`
}

type AnsibleHost struct {
	PrivateIP string            `yaml:"private_ip"`
	PublicIP  string            `yaml:"public_ip"`
	HostName  string            `yaml:"host_name"`
	ID        string            `yaml:"id"`
	Mounts    []Mount           `yaml:"mounts"`
	ExtraVars map[string]string `yaml:"extra_vars"`
}

type Mount struct {
	Name      string `yaml:"name"`
	Path      string `yaml:"path"`
	MountPath string `yaml:"mount_path"`
	Owner     string `yaml:"owner"`
}

func (group *HostGroup) GetHosts() []string {
	res := []string{}
	for k := range group.Hosts {
		res = append(res, k)
	}
	return res
}

func (group *HostGroup) GetPrivateHosts() []string {
	res := []string{}
	for _, v := range group.Hosts {
		res = append(res, v.PrivateIP)
	}
	return res
}

func (group *HostGroup) GetPrivateHostNames() []string {
	res := []string{}
	for _, v := range group.Hosts {
		res = append(res, v.HostName)
	}
	return res
}

func (inv *Inventory) GetAllPrivateHosts() []string {
	hosts := []string{}
	rawHosts := []HostGroup{}
	seenHosts := make(map[string]string)

	for _, hostGroup := range inv.All.Children {
		rawHosts = append(rawHosts, hostGroup)
	}

	for _, hostGroup := range rawHosts {
		for _, host := range hostGroup.GetPrivateHosts() {
			if _, ok := seenHosts[host]; !ok {
				hosts = append(hosts, host)
				seenHosts[host] = host
			}
		}
		for _, host := range hostGroup.GetPrivateHostNames() {
			if _, ok := seenHosts[host]; !ok {
				hosts = append(hosts, host)
				seenHosts[host] = host
			}
		}
	}

	return hosts
}

func LoadInventory(file string) (*Inventory, error) {
	bytes, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return nil, err
	}
	var config Inventory
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	config.Path = file
	return &config, nil
}

func GenerateInventory(config *conf.Config) (*Inventory, error) {
	jsonFile, err := os.Open(filepath.Clean(filepath.Join(config.BaseDir, "inventory-output.json")))
	if err != nil {
		return nil, err
	}
	defer func() {
		e := jsonFile.Close()
		if e != nil {
			panic(e)
		}
	}()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var inventory InventoryJson

	err = json.Unmarshal(byteValue, &inventory)
	if err != nil {
		return nil, err
	}

	inv := Inventory{
		All: All{
			Children: make(map[string]HostGroup),
		},
	}

	for _, v := range inventory.Servers.Value {
		if _, ok := inv.All.Children[v.Group]; !ok {
			inv.All.Children[v.Group] = HostGroup{Hosts: make(map[string]AnsibleHost)}
		}

		volumes := []Mount{}
		for _, volume := range inventory.Volumes.Value {
			if fmt.Sprintf("%v", volume.ServerID) == v.ServerID {
				owner := "root"
				for _, vol := range config.ServerGroups[v.Group].Volumes {
					if strings.HasSuffix(volume.Name, vol.Name) {
						owner = vol.Owner
					}
				}

				volumes = append(volumes, Mount{
					Name:      volume.Name,
					Path:      volume.Path,
					MountPath: volume.Mount,
					Owner:     owner,
				})
			}

		}
		inv.All.Children[v.Group].Hosts[v.Host] = AnsibleHost{
			HostName:  v.HostName,
			PrivateIP: v.PrivateIP,
			PublicIP:  v.Host,
			Mounts:    volumes,
			ID:        v.ServerID,
			ExtraVars: map[string]string{
				"datacenter": config.DC,
				"os":         v.Image,
			},
		}

		for _, alias := range config.ServerGroups[v.Group].Aliases {
			inv.All.Children[alias] = HostGroup{Hosts: make(map[string]AnsibleHost)}
			inv.All.Children[alias] = inv.All.Children[v.Group]
		}
		// also map mounts
	}

	bytes, err := yaml.Marshal(inv)
	if err != nil {
		return nil, err
	}

	invPath := filepath.Clean(filepath.Join(config.BaseDir, "inventory"))
	err = os.WriteFile(invPath, bytes, 0600)
	if err != nil {
		return nil, err
	}

	inv.Path = invPath
	return &inv, nil
}
