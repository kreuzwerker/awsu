package command

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/kreuzwerker/awsu/config"
	"github.com/spf13/cobra"
	"github.com/yawn/doubledash"
)

var rootFlags = new(config.Config)

var rootCmd = &cobra.Command{
	Use: app,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return rootFlags.Init()
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		creds, err := strategy.Apply(config)

		if err != nil {
			return err
		}

		if len(doubledash.Xtra) > 0 {
			return creds.Exec(doubledash.Xtra[0], doubledash.Xtra)
		}

		fmt.Println(creds.String())

		return nil

	},
}

func init() {

	os.Args = doubledash.Args

	flag(rootCmd.PersistentFlags(),
		&rootFlags.ConfigFile,
		defaults.SharedCredentialsFilename(),
		"config-config",
		"c",
		"AWS_CONFIG_FILE",
		"sets the config file",
	)

	flag(rootCmd.PersistentFlags(),
		&rootFlags.Duration,
		1*time.Hour,
		"duration",
		"d",
		"AWSU_DURATION",
		"duration to use for session tokens and roles",
	)

	flag(rootCmd.PersistentFlags(),
		&rootFlags.ConfigFile,
		"",
		"mfa-serial",
		"m",
		"AWSU_MFA_SERIAL",
		"set or override MFA serial",
	)

	flag(rootCmd.PersistentFlags(),
		&rootFlags.Generator,
		"yubikey",
		"generator",
		"g",
		"AWSU_TOKEN_GENERATOR",
		"configure the token generator to 'yubikey' or 'manual'",
	)

	flag(rootCmd.PersistentFlags(),
		&rootFlags.Grace,
		45*time.Minute,
		"grace",
		"r",
		"AWSU_GRACE",
		"distance to the duration before a cache credential is considered expired",
	)

	flag(rootCmd.PersistentFlags(),
		&rootFlags.NoCache,
		false,
		"no-cache",
		"n",
		"AWSU_NO_CACHE",
		"disable caching of short-term credentials",
	)

	flag(rootCmd.PersistentFlags(),
		&rootFlags.Profile,
		"default",
		"profile",
		"p",
		"AWS_PROFILE",
		"shared config profile to use",
	)

	flag(rootCmd.PersistentFlags(),
		&rootFlags.SharedCredentialsFile,
		defaults.SharedCredentialsFilename(),
		"shared-credentials-file",
		"s",
		"AWS_SHARED_CREDENTIALS_FILE",
		"shared credentials file to use",
	)

	flag(rootCmd.PersistentFlags(),
		&rootFlags.Verbose,
		false,
		"verbose",
		"v",
		"AWSU_VERBOSE",
		"enable verbose logging",
	)

}
