package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	KeyCacheTTL              = "cache-ttl"
	KeyConfigFile            = "config-file"
	KeyNoCache               = "no-cache"
	KeyProfile               = "profile"
	KeySessionTTL            = "session-ttl"
	KeySharedCredentialsFile = "shared-credentials-file"
	KeyVerbose               = "verbose"
)

type Config struct {
	CacheTTL              time.Duration
	ConfigFile            string
	NoCache               bool
	Profile               string
	Profiles              Profiles
	SessionTTL            time.Duration
	SharedCredentialsFile string
	Verbose               bool
}

func NewConfig() *Config {

	return &Config{
		CacheTTL:              viper.GetDuration(KeyCacheTTL),
		ConfigFile:            viper.GetString(KeyConfigFile),
		NoCache:               viper.GetBool(KeyNoCache),
		Profile:               viper.GetString(KeyProfile),
		SessionTTL:            viper.GetDuration(KeySessionTTL),
		SharedCredentialsFile: viper.GetString(KeySharedCredentialsFile),
		Verbose:               viper.GetBool(KeyVerbose),
	}

}
