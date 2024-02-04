package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DC                  string                 `yaml:"dc_name"`
	BaseDir             string                 `yaml:"base_dir"`
	OrgName             string                 `yaml:"org_name"`
	Providers           []ProviderConfig       `yaml:"providers"`
	CloudProviderConfig CloudProvider          `yaml:"cloud_provider_config"`
	ServerGroups        map[string]ServerGroup `yaml:"server_groups"`
	Services            map[string]interface{} `yaml:"services"`
}

type ProviderConfig struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}

type ServerGroup struct {
	Num          int      `yaml:"num"`
	InstanceType string   `yaml:"instance_type"`
	Volumes      []Volume `yaml:"volumes"`
	LbTarget     bool     `yaml:"lb_target"`
	Aliases      []string `yaml:"aliases"`
}

type Volume struct {
	Name  string `yaml:"name"`
	Path  string `yaml:"path"`
	Size  int    `yaml:"size"`
	Owner string `yaml:"owner"`
}

type ClientVolume struct {
	Name   string `yaml:"name"`
	Client string `yaml:"client"`
	Path   string `yaml:"path"`
	Size   int    `yaml:"size"`
}

type CloudProvider struct {
	User             string                 `yaml:"sudo_user"`
	GithubIds        []string               `yaml:"ssh_key_github_ids"`
	Provider         string                 `yaml:"provider"`
	ProviderSettings map[string]interface{} `yaml:"provider_settings"`
	AllowedIPs       []string               `yaml:"allowed_ips"`
	SSHKey           string                 `yaml:"private_ssh_key"`
	GithubKeys       []GithubKey            `yaml:"-"` // Ignore in yaml, added to propagate information during config changes
}

func Load(file string) (*Config, error) {
	bytes, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadTFExecVars() *tfexec.VarOption {
	token := os.Getenv("HETZNER_TOKEN")
	return tfexec.Var(fmt.Sprintf("hcloud_token=%s", token))
}
