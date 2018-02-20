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

// TODO: perform the workspace to profile mapping here
var rootFlags = struct {
	cacheTTL              time.Duration
	configFile            string
	noCache               bool
	profile               string
	profiles              config.Profiles
	sharedCredentialsFile string
	verbose               bool
}{}

var rootCmd = &cobra.Command{
	Use: app,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		// TODO: apply this pattern to other commands as well

		rootFlags.cacheTTL = viper.GetDuration("cache-ttl")
		rootFlags.configFile = viper.GetString("config-file")
		rootFlags.noCache = viper.GetBool("no-cache")
		rootFlags.profile = viper.GetString("profile")
		rootFlags.sharedCredentialsFile = viper.GetString("shared-credentials-file")
		rootFlags.verbose = viper.GetBool("verbose")

		if rootFlags.verbose {
			log.Debug = true
		}

		var err error

		rootFlags.profiles, err = config.Load(
			rootFlags.configFile,
			rootFlags.sharedCredentialsFile,
		)

		return err

	},
	RunE: func(cmd *cobra.Command, args []string) error {
		sess, err := newSession(rootFlags.noCache,
			rootFlags.profile,
			rootFlags.profiles)
		if err != nil {
			return err
		}

		if len(doubledash.Xtra) > 0 {
			return sess.Exec(doubledash.Xtra[0], doubledash.Xtra)
		}

		fmt.Println(sess.String())

		return nil

	},
}

func init() {

	os.Args = doubledash.Args

	rootCmd.PersistentFlags().DurationP("cache-ttl", "t", 15*time.Minute, "time to live for cached short-term credentials")
	viper.BindPFlag("cache-ttl", rootCmd.PersistentFlags().Lookup("cache-ttl"))
	viper.BindEnv("cache-ttl", "AWSU_CACHE_TTL")

	viper.BindEnv("config-file", "AWS_CONFIG_FILE")
	viper.SetDefault("config-file", defaults.SharedCredentialsFilename())

	// TODO: enable config file workspace mapping here
	rootCmd.PersistentFlags().StringVarP(&rootFlags.profile, "profile", "p", "", "shared config profile to use")
	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))
	viper.BindEnv("profile", "AWS_PROFILE")
	viper.SetDefault("profile", "default")

	rootCmd.PersistentFlags().BoolVarP(&rootFlags.noCache, "no-cache", "n", false, "disable caching of short-term credentials")
	viper.BindPFlag("no-cache", rootCmd.PersistentFlags().Lookup("no-cache"))
	viper.BindEnv("no-cache", "AWSU_NO_CACHE")

	viper.BindEnv("shared-credentials-file", "AWS_SHARED_CREDENTIALS_FILE")
	viper.SetDefault("shared-credentials-file", defaults.SharedConfigFilename())


	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose operations")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindEnv("verbose", "AWSU_VERBOSE")

}
