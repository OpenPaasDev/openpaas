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
	TfState             TerraformState         `yaml:"terraform_state"`
	Providers           map[string]interface{} `yaml:"providers"`
	CloudProviderConfig CloudProvider          `yaml:"cloud_provider_config"`
	ServerGroups        map[string]ServerGroup `yaml:"server_groups"`
	Services            map[string]interface{} `yaml:"services"`
}

type TerraformState struct {
	Backend string            `yaml:"backend"`
	Config  map[string]string `yaml:"config"`
}

type ServerGroup struct {
	Num          int      `yaml:"num"`
	InstanceType string   `yaml:"instance_type"`
	Volumes      []Volume `yaml:"volumes"`
	LbTarget     bool     `yaml:"lb_target"`
	Aliases      []string `yaml:"aliases"`
	Image        string   `yaml:"os_image"`
	SubnetID     int      `yaml:"subnet_id"`
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
	GithubIds        []string               `yaml:"github_ids"`
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
	ResourceNames    HetznerResourceNames `yaml:"resource_names"`
	LoadBalancerType string               `yaml:"load_balancer_type"`
}

type TFVarsConfig struct {
	ServerGroups   map[string]ServerGroup
	ProviderConfig interface{}
}

func Load(file string) (*Config, error) {
	fmt.Println("Loading config from file", file)
	bytes, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	// replace any env vars in the Tf State config with their values
	config.TfState = updateTerraformStateVarsConfig(config.TfState)

	return &config, nil
}

func LoadTFExecVars() *tfexec.VarOption {
	token := os.Getenv("HETZNER_TOKEN")
	return tfexec.Var(fmt.Sprintf("hcloud_token=%s", token))
}

func updateTerraformStateVarsConfig(tfState TerraformState) TerraformState {
	for key, value := range tfState.Config {
		envValue, exists := os.LookupEnv(value)
		// If the environment variable exists and is not an empty string, replace the value for that key in the map
		if exists && envValue != "" {
			tfState.Config[key] = envValue
		}
	}
	return tfState
}
