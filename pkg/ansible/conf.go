package ansible

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/OpenPaasDev/core/pkg/conf"
	"gopkg.in/yaml.v3"
)

type Inventory struct {
	All All `yaml:"all"`
}

type InventoryJson struct {
	Servers       map[string]HostValues `json:"servers"`
	ConsulVolumes Volumes               `json:"consul_volumes"`
	ClientVolumes Volumes               `json:"client_volumes"`
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
	Host      string `json:"host"`
	HostName  string `json:"host_name"`
	PrivateIP string `json:"private_ip"`
	ServerID  string `json:"server_id"`
}

type All struct {
	Children map[string]HostGroup `yaml:"children"`
}

type HostGroup struct {
	Hosts map[string]AnsibleHost `yaml:"hosts"`
}

type AnsibleHost struct {
	PrivateIP string            `yaml:"private_ip"`
	HostName  string            `yaml:"host_name"`
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
	return &config, nil
}

func GenerateInventory(config *conf.Config) (*Inventory, error) {
	jsonFile, err := os.Open(filepath.Clean(filepath.Join(config.BaseDir, "inventory-output.json")))
	if err != nil {
		return nil, err
	}
	defer func() {
		e := jsonFile.Close()
		fmt.Println(e)
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

	for k, v := range inventory.Servers {
		inv.All.Children[k] = HostGroup{Hosts: make(map[string]AnsibleHost)}
		for _, host := range v.Value {
			mounts := []Mount{}
			for _, vol := range inventory.ClientVolumes.Value {
				if fmt.Sprintf("%v", vol.ServerID) == host.ServerID {
					mounts = append(mounts, Mount{
						Name:      vol.Name,
						Path:      vol.Path,
						MountPath: vol.Mount,
						Owner:     config.CloudProviderConfig.User,
					})
				}
			}

			inv.All.Children[k].Hosts[host.Host] = AnsibleHost{
				PrivateIP: host.PrivateIP,
				HostName:  host.HostName,
				Mounts:    mounts,
			}
		}
	}

	bytes, err := yaml.Marshal(inv)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(filepath.Clean(filepath.Join(config.BaseDir, "inventory")), bytes, 0600)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}
