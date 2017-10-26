package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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

	return &Config{
		Profiles: sec.KeysHash(),
	}, nil

}
