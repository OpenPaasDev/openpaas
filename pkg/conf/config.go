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
	CloudProviderConfig CloudProvider          `yaml:"cloud_provider_config"`
	IngressConfig       IngressConfig          `yaml:"ingress"`
	ServerGroups        map[string]ServerGroup `yaml:"server_groups"`
}

type ServerGroup struct {
	Num          int      `yaml:"num"`
	InstanceType string   `yaml:"instance_type"`
	Volumes      []Volume `yaml:"volumes"`
	IsPublic     bool     `yaml:"is_public"`
	Aliases      []string `yaml:"aliases"`
}

type Volume struct {
	Name  string `yaml:"name"`
	Path  string `yaml:"path"`
	Size  int    `yaml:"size"`
	Owner string `yaml:"owner"`
}

type IngressConfig struct {
	ManagementDomain string `yaml:"management_domain"`
}

type ClientVolume struct {
	Name   string `yaml:"name"`
	Client string `yaml:"client"`
	Path   string `yaml:"path"`
	Size   int    `yaml:"size"`
}

type CloudProvider struct {
	User             string                 `yaml:"sudo_user"`
	Dir              string                 `yaml:"sudo_dir"`
	NetworkInterface string                 `yaml:"internal_network_interface_name"`
	Provider         string                 `yaml:"provider"`
	ProviderSettings map[string]interface{} `yaml:"provider_settings"`
	AllowedIPs       []string               `yaml:"allowed_ips"`
}

type HetznerResourceNames struct {
	BaseServerName string `yaml:"base_server_name"`
	FirewallName   string `yaml:"firewall_name"`
	NetworkName    string `yaml:"network_name"`
}

type HetznerSettings struct {
	Location         string               `yaml:"location"`
	SSHKeys          []string             `yaml:"ssh_keys"`
	ResourceNames    HetznerResourceNames `yaml:"resource_names"`
	LoadBalancerType string               `yaml:"load_balancer_type"`
}

type TFVarsConfig struct {
	ServerGroups   map[string]ServerGroup
	ProviderConfig interface{}
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

func LoadTFVarsConfig(config Config) (*TFVarsConfig, error) {
	var providerConfig interface{}
	if config.CloudProviderConfig.Provider == "hetzner" {
		var hetznerConfig HetznerSettings
		bytes, err := yaml.Marshal(config.CloudProviderConfig.ProviderSettings)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(bytes, &hetznerConfig)
		if err != nil {
			return nil, err
		}
		providerConfig = hetznerConfig
	}

	for _, group := range config.ServerGroups {
		if group.Volumes == nil {
			group.Volumes = []Volume{}
		}
	}

	return &TFVarsConfig{
		ServerGroups:   config.ServerGroups,
		ProviderConfig: providerConfig,
	}, nil
}

func LoadTFExecVars() *tfexec.VarOption {
	token := os.Getenv("HETZNER_TOKEN")
	return tfexec.Var(fmt.Sprintf("hcloud_token=%s", token))
}
