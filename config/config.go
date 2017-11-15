package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kreuzwerker/awsu/log"
	"github.com/yawn/envmap"
	ini "gopkg.in/ini.v1"
)

type Config struct {
	Profiles map[string]string `ini:"profiles"`
}

// DetectWorkspace implements the workspace detection from Terraform
func (c *Config) DetectWorkspace() string {

	// this could be done with new(command.Meta).Workspace() as well but this
	// adds 90MB of Terraform to the build

	const (
		DefaultDataDir       = ".terraform"
		DefaultStateName     = "default"
		DefaultWorkspaceFile = "environment"
		WorkspaceNameEnvVar  = "TF_WORKSPACE"
	)

	if envVar := os.Getenv(WorkspaceNameEnvVar); envVar != "" {
		return envVar
	}

	if envVar := os.Getenv("AWSU_WORKSPACE"); envVar != "" {
		return envVar
	}

	envData, err := ioutil.ReadFile(filepath.Join(DefaultDataDir, DefaultWorkspaceFile))

	current := string(bytes.TrimSpace(envData))

	if err != nil || current == "" {
		return DefaultStateName
	}

	return current

}

func (c *Config) Get(workspace string) (string, error) {

	profile, ok := c.Profiles[workspace]

	if !ok {
		return "", fmt.Errorf("no profile configured for workspace %q", workspace)
	}

	return profile, nil

}

// keys returns a sorted list of workspace names
func (c *Config) workspaces() []string {

	var keys sort.StringSlice

	for k := range c.Profiles {
		keys = append(keys, k)
	}

	keys.Sort()

	return keys

}

// Detect will load configuration values from the environment
func Detect() *Config {

	const (
		prefix = "AWSU_WORKSPACE_"
	)

	config := &Config{
		Profiles: make(map[string]string),
	}

	for k, v := range envmap.Import() {

		if strings.HasPrefix(k, prefix) {
			k = strings.ToLower(strings.TrimPrefix(k, prefix))
			config.Profiles[k] = v
		}

	}

	if len(config.Profiles) == 0 {
		return nil
	}

	log.Log("loading config from environment (%s)", config.workspaces())

	return config

}

// Load will load configuration values from a file
func Load(dir string) (*Config, error) {

	path := filepath.Join(dir, ".awsu")

	cfg, err := ini.Load(path)

	if err != nil {
		return nil, err
	}

	sec, err := cfg.GetSection("profiles")

	if err != nil {
		return nil, err
	}

	config := &Config{
		Profiles: sec.KeysHash(),
	}

	log.Log("loading config from file %q (%s)", path, config.workspaces())

	return config, nil

}
