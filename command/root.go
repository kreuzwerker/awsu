package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"

	"github.com/spf13/cobra"
)

var (
	config  *ini.File
	section *ini.Section
)

var rootFlags = struct {
	configPath string
	prefix     string
	profile    string
}{}

var rootCmd = &cobra.Command{
	Use: app,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		cfg, err := ini.Load(rootFlags.configPath)

		if err != nil {
			return fmt.Errorf("error while opening config at %q: %s",
				rootFlags.configPath,
				err)
		}

		config = cfg

		sec, err := config.GetSection(rootFlags.profile)

		if err != nil {
			return fmt.Errorf("error while fetching profile %q from config: %s",
				rootFlags.profile,
				err)
		}

		section = sec

		return nil

	},
}

func init() {

	config := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")

	if config == "" {

		for _, home := range []string{os.Getenv("HOME"), os.Getenv("USERPROFILE")} {

			if home != "" {
				config = filepath.Join(home, ".aws", "credentials")
				break
			}

		}

	}

	log.SetOutput(os.Stderr)

	rootCmd.PersistentFlags().StringVarP(&rootFlags.configPath,
		"config",
		"c",
		config,
		"AWS credentials file location")

	rootCmd.PersistentFlags().StringVarP(&rootFlags.prefix,
		"prefix",
		"x",
		"_",
		"Prefix token for shadowing environment variables")

	rootCmd.PersistentFlags().StringVarP(&rootFlags.profile,
		"profile",
		"p",
		os.Getenv("AWS_PROFILE"),
		"Profile to pick")

}
