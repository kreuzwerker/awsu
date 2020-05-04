package command

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yawn/doubledash"

	"github.com/gesellix/awsu/config"
	"github.com/gesellix/awsu/strategy"
)

var conf config.Config

var rootCmd = &cobra.Command{
	Use:           app,
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		if err := viper.Unmarshal(&conf); err != nil {
			return err
		}

		return conf.Init()

	},
	RunE: func(cmd *cobra.Command, args []string) error {

		creds, err := strategy.Apply(&conf)

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
		defaults.SharedConfigFilename(),
		"config-file",
		"c",
		"AWS_CONFIG_FILE",
		"sets the config file",
	)

	flag(rootCmd.PersistentFlags(),
		1*time.Hour,
		"duration",
		"d",
		"AWSU_DURATION",
		"duration to use for session tokens and roles",
	)

	flag(rootCmd.PersistentFlags(),
		"",
		"mfa-serial",
		"m",
		"AWSU_MFA_SERIAL",
		"set or override MFA serial",
	)

	flag(rootCmd.PersistentFlags(),
		"yubikey",
		"generator",
		"g",
		"AWSU_TOKEN_GENERATOR",
		"configure the token generator to 'yubikey' or 'manual'",
	)

	flag(rootCmd.PersistentFlags(),
		45*time.Minute,
		"grace",
		"r",
		"AWSU_GRACE",
		"distance to the duration before a cache credential is considered expired",
	)

	flag(rootCmd.PersistentFlags(),
		false,
		"no-cache",
		"n",
		"AWSU_NO_CACHE",
		"disable caching of short-term credentials",
	)

	flag(rootCmd.PersistentFlags(),
		"default",
		"profile",
		"p",
		"AWS_PROFILE",
		"shared config profile to use",
	)

	flag(rootCmd.PersistentFlags(),
		defaults.SharedCredentialsFilename(),
		"shared-credentials-file",
		"s",
		"AWS_SHARED_CREDENTIALS_FILE",
		"shared credentials file to use",
	)

	flag(rootCmd.PersistentFlags(),
		false,
		"verbose",
		"v",
		"AWSU_VERBOSE",
		"enable verbose logging",
	)

}
