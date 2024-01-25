package terraform

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/OpenPaasDev/openpaas/pkg/conf"
)

//go:embed templates/hetzner/main.tf
var hetznerMain string

//go:embed templates/hetzner/vars.tf
var hetznerVars string

//go:embed templates/hetzner/cloud-init.yml
var cloudInit string

func GenerateTerraform(config *conf.Config) error {
	settings := map[string]struct {
		Main      string
		Vars      string
		CloudInit string
	}{
		"hetzner": {
			Main:      hetznerMain,
			Vars:      hetznerVars,
			CloudInit: cloudInit,
		},
	}

	tfSettings, ok := settings[config.CloudProviderConfig.Provider]
	if !ok {
		return fmt.Errorf("%s is not a supported cloud provider", config.CloudProviderConfig.Provider)
	}

	tmplVars, e := template.New("tf-vars").Parse(tfSettings.Vars)
	if e != nil {
		return e
	}
	var bufVars bytes.Buffer

	//TODO do we use this allowedIps at all??
	allowedIps := []string{}
	config.CloudProviderConfig.ProviderSettings["https_allowed_ips"] = allowedIps

	err := tmplVars.Execute(&bufVars, config)
	if err != nil {
		return err
	}

	tmplCloudInit, e := template.New("cloud-init").Parse(tfSettings.CloudInit)
	if e != nil {
		return e
	}
	var bufCloudInit bytes.Buffer
	err = tmplCloudInit.Execute(&bufCloudInit, config)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Clean(filepath.Join(config.BaseDir, "terraform")), 0750)
	if err != nil {
		return err
	}
	folder := filepath.Join(config.BaseDir, "terraform")

	err = os.WriteFile(filepath.Clean(filepath.Join(folder, "vars.tf")), bufVars.Bytes(), 0600)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Clean(filepath.Join(folder, "main.tf")), []byte(hetznerMain), 0600)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Clean(filepath.Join(folder, "cloud-init.yml")), bufCloudInit.Bytes(), 0600)
	if err != nil {
		return err
	}
	return nil

}
