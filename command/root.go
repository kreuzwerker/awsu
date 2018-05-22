package command

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/kreuzwerker/awsu/config"
	"github.com/kreuzwerker/awsu/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yawn/doubledash"
)

var rootConfig *config.Config

var rootFlags = struct {
	cacheTTL              time.Duration
	configFile            string
	noCache               bool
	profile               string
	sessionTTL            time.Duration
	sharedCredentialsFile string
	verbose               bool
}{}

var rootCmd = &cobra.Command{
	Use: app,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		rootConfig = config.NewConfig()

		if rootConfig.Verbose {
			log.Debug = true
		}

		var err error

		rootConfig.Profiles, err = config.Load(
			rootConfig.ConfigFile,
			rootConfig.SharedCredentialsFile,
		)

		return err

	},
	RunE: func(cmd *cobra.Command, args []string) error {

		creds, err := newSession(rootConfig)

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

	rootCmd.PersistentFlags().DurationP(config.KeyCacheTTL, "t", 8*time.Hour, "time to live for cached role credentials")
	viper.BindPFlag(config.KeyCacheTTL, rootCmd.PersistentFlags().Lookup(config.KeyCacheTTL))
	viper.BindEnv(config.KeyCacheTTL, "AWSU_CACHE_ROLE_TTL")

	viper.BindEnv(config.KeyConfigFile, "AWS_CONFIG_FILE")
	viper.SetDefault(config.KeyConfigFile, defaults.SharedCredentialsFilename())

	rootCmd.PersistentFlags().BoolVarP(&rootFlags.noCache, config.KeyNoCache, "n", false, "disable caching of short-term credentials")
	viper.BindPFlag(config.KeyNoCache, rootCmd.PersistentFlags().Lookup(config.KeyNoCache))
	viper.BindEnv(config.KeyNoCache, "AWSU_NO_CACHE")

	// TODO: enable config file workspace mapping here
	rootCmd.PersistentFlags().StringVarP(&rootFlags.profile, config.KeyProfile, "p", "", "shared config profile to use")
	viper.BindPFlag(config.KeyProfile, rootCmd.PersistentFlags().Lookup(config.KeyProfile))
	viper.BindEnv(config.KeyProfile, "AWS_PROFILE")
	viper.SetDefault(config.KeyProfile, "default")

	rootCmd.PersistentFlags().DurationP(config.KeySessionTTL, "s", 8*time.Hour, "time to live for cached session token credentials")
	viper.BindPFlag(config.KeySessionTTL, rootCmd.PersistentFlags().Lookup(config.KeySessionTTL))
	viper.BindEnv(config.KeySessionTTL, "AWSU_CACHE_SESSION_TOKEN_TTL")

	viper.BindEnv(config.KeySharedCredentialsFile, "AWS_SHARED_CREDENTIALS_FILE")
	viper.SetDefault(config.KeySharedCredentialsFile, defaults.SharedConfigFilename())

	rootCmd.PersistentFlags().BoolP(config.KeyVerbose, "v", false, "enable verbose operations")
	viper.BindPFlag(config.KeyVerbose, rootCmd.PersistentFlags().Lookup(config.KeyVerbose))
	viper.BindEnv(config.KeyVerbose, "AWSU_VERBOSE")

}
