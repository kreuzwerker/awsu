package command

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

var config = new(prerunConfig)

type prerunConfig struct {
	config  *ini.File
	section *ini.Section
}

func (p *prerunConfig) PreRun(cmd *cobra.Command, args []string) error {

	var err error

	if p.config, err = ini.Load(rootFlags.configPath); err != nil {
		return fmt.Errorf("error while opening config at %q: %s",
			rootFlags.configPath,
			err)
	}

	if p.section, err = p.config.GetSection(rootFlags.profile); err != nil {
		return fmt.Errorf("error while fetching profile %q from config: %s",
			rootFlags.profile,
			err)
	}

	return nil

}
