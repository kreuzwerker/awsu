package command

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	aki = "AWS_ACCESS_KEY_ID"
	ask = "AWS_SECRET_ACCESS_KEY"
	ast = "AWS_SESSION_TOKEN"
)

var rootFlags = struct {
	configPath string
	prefix     string
	profile    string
}{}

var rootCmd = &cobra.Command{
	Use: app,
	PersistentPreRunE: config.PreRun,
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
