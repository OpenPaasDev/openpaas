package terraform

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenPaasDev/openpaas/pkg/conf"
	"github.com/OpenPaasDev/openpaas/pkg/util"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

type tfConfig struct {
	Variable []Variable `hcl:"variable,block"`
}

type Variable struct {
	Type      *hcl.Attribute `hcl:"type"`
	Name      string         `hcl:"name,label"`
	Default   *cty.Value     `hcl:"default,optional"`
	Sensitive bool           `hcl:"sensitive,optional"`
}

// parse the main.tf file to extract data for testing
type TfMain struct {
	Terraform []Terraform `hcl:"terraform,block"`
	Raw       hcl.Body    `hcl:",remain"`
}

type Terraform struct {
	Backend []Backend `hcl:"backend,block"`
	Raw     hcl.Body  `hcl:",remain"`
}

type Backend struct {
	Name   string         `hcl:"name,label"`
	Config hcl.Attributes `hcl:",remain"` // Dynamic configuration attributes
}

func TestGenerateTerraform(t *testing.T) {
	config, err := conf.Load("../testdata/config.yaml")
	require.NoError(t, err)

	folder := util.RandString(8)
	config.BaseDir = folder
	defer func() {
		e := os.RemoveAll(filepath.Clean(folder))
		require.NoError(t, e)
	}()

	err = GenerateTerraform(config)
	require.NoError(t, err)

	parser := hclparse.NewParser()
	f, parseDiags := parser.ParseHCLFile(filepath.Clean(filepath.Join(folder, "terraform", "vars.tf")))
	assert.False(t, parseDiags.HasErrors())

	_, parseDiags = parser.ParseHCLFile(filepath.Clean(filepath.Join(folder, "terraform", "main.tf")))
	assert.False(t, parseDiags.HasErrors())

	var conf tfConfig
	decodeDiags := gohcl.DecodeBody(f.Body, nil, &conf)
	assert.False(t, decodeDiags.HasErrors())

	vars := []struct {
		name       string
		tpe        string
		defaultVal cty.Value
	}{
		{name: "hcloud_token", tpe: "string", defaultVal: cty.NullVal(cty.String)},
		//		{name: "ssh_keys", tpe: "list", defaultVal: cty.TupleVal([]cty.Value{cty.StringVal("123456")})},
		{name: "base_server_name", tpe: "string", defaultVal: cty.StringVal("nomad-srv")},
		{name: "firewall_name", tpe: "string", defaultVal: cty.StringVal("dev_firewall")},
		{name: "network_name", tpe: "string", defaultVal: cty.StringVal("dev_network")},
		{name: "allow_ips", tpe: "list", defaultVal: cty.TupleVal([]cty.Value{cty.StringVal("85.4.84.201/32")})},
		{name: "location", tpe: "string", defaultVal: cty.StringVal("nbg1")},
	}

	expectedMap := make(map[string]string)
	for _, v := range conf.Variable {
		for _, expected := range vars {
			if expected.name == v.Name {
				expectedMap[expected.name] = expected.name
				assert.Equal(t, expected.tpe, v.Type.Expr.Variables()[0].RootName())
				if expected.defaultVal != cty.NullVal(cty.String) && !strings.Contains(expected.name, "_count") {
					assert.Equal(t, expected.defaultVal, *v.Default)
				}
				if strings.Contains(expected.name, "_count") {
					assert.Equal(t, expected.defaultVal.AsBigFloat().String(), v.Default.AsBigFloat().String())
				}
			}
		}
	}

	assert.Equal(t, len(expectedMap), len(vars))
}

func TestGenerateTerraformWithLocal(t *testing.T) {
	config, err := conf.Load("../testdata/config.yaml")
	require.NoError(t, err)

	folder := util.RandString(8)
	config.BaseDir = folder
	defer func() {
		e := os.RemoveAll(filepath.Clean(folder))
		require.NoError(t, e)
	}()

	// set the backend to local before generating the terraform files
	config.TfState.Backend = "local"
	err = GenerateTerraform(config)
	require.NoError(t, err)

	parser := hclparse.NewParser()

	f, parseDiags := parser.ParseHCLFile(filepath.Clean(filepath.Join(folder, "terraform", "main.tf")))
	assert.False(t, parseDiags.HasErrors())

	var tfConf TfMain
	decodeDiags := gohcl.DecodeBody(f.Body, nil, &tfConf)
	assert.False(t, decodeDiags.HasErrors())
	assert.Empty(t, tfConf.Terraform[0].Backend)
}

func TestGenerateTerraformWithS3(t *testing.T) {
	config, err := conf.Load("../testdata/config.yaml")
	require.NoError(t, err)

	folder := util.RandString(8)
	config.BaseDir = folder
	defer func() {
		e := os.RemoveAll(filepath.Clean(folder))
		require.NoError(t, e)
	}()

	// load a remote s3 config from the test data
	err = GenerateTerraform(config)
	require.NoError(t, err)

	parser := hclparse.NewParser()

	tfVars, parseDiags := parser.ParseHCLFile(filepath.Clean(filepath.Join(folder, "terraform", "vars.tf")))
	assert.False(t, parseDiags.HasErrors())

	f, parseDiags := parser.ParseHCLFile(filepath.Clean(filepath.Join(folder, "terraform", "main.tf")))
	assert.False(t, parseDiags.HasErrors())

	// generate a context based on the vars
	var varsConf tfConfig
	decodeDiags := gohcl.DecodeBody(tfVars.Body, nil, &varsConf)
	assert.False(t, decodeDiags.HasErrors())

	ctx := &hcl.EvalContext{
		Variables: make(map[string]cty.Value),
	}
	for _, value := range varsConf.Variable {
		var defaultValue cty.Value
		if value.Default != nil {
			defaultValue = *value.Default
		} else {
			defaultValue = cty.StringVal("null pointer")
		}
		ctx.Variables[value.Name] = defaultValue
	}

	// parse main.tf into the struct
	var tfConf TfMain
	decodeDiags = gohcl.DecodeBody(f.Body, nil, &tfConf)
	assert.False(t, decodeDiags.HasErrors())

	vars := map[string]cty.Value{
		"endpoint":   cty.StringVal("endpoint_to_s3_compatible_storage"),
		"bucket":     cty.StringVal("bucket_name"),
		"region":     cty.StringVal("auto"),
		"access_key": cty.StringVal("env_var_access_key"),
		"secret_key": cty.StringVal("env_var_secret_key"),
		"key":        cty.StringVal("openpaas/terraform.tfstate"),
	}

	for key, attr := range tfConf.Terraform[0].Backend[0].Config {
		var value cty.Value
		value, _ = attr.Expr.Value(ctx)

		expected, ok := vars[key]
		// only check keys we are testing for
		if ok {
			assert.Equal(t, expected.Type(), value.Type())
			assert.Equal(t, expected, value)
		}
	}
}
