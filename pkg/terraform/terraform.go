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

func GenerateTerraform(config *conf.Config) error {
	settings := map[string]struct {
		Main string
		Vars string
	}{
		"hetzner": {
			Main: hetznerMain,
			Vars: hetznerVars,
		},
	}

	tfSettings, ok := settings[config.CloudProviderConfig.Provider]
	if !ok {
		return fmt.Errorf("%s is not a supported cloud provider", config.CloudProviderConfig.Provider)
	}

	tmplMain, e := template.New("tf-main").Parse(tfSettings.Main)
	if e != nil {
		return e
	}
	var bufMain bytes.Buffer
	err := tmplMain.Execute(&bufMain, config)
	if err != nil {
		return err
	}

	tmplVars, e := template.New("tf-vars").Parse(tfSettings.Vars)
	if e != nil {
		return e
	}
	var bufVars bytes.Buffer
	err = tmplVars.Execute(&bufVars, config)
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
	err = os.WriteFile(filepath.Clean(filepath.Join(folder, "main.tf")), bufMain.Bytes(), 0600)
	if err != nil {
		return err
	}
	return nil

}
