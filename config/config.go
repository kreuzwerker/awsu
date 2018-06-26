package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	KeyCacheTTL              = "cache-ttl"
	KeyConfigFile            = "config-file"
	KeyMFASerial             = "mfa-serial"
	KeyNoCache               = "no-cache"
	KeyProfile               = "profile"
	KeySessionTTL            = "session-ttl"
	KeySharedCredentialsFile = "shared-credentials-file"
	KeyVerbose               = "verbose"
)

type Config struct {
	CacheTTL              time.Duration
	ConfigFile            string
	MFASerial             string
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
		MFASerial:             viper.GetString(KeyMFASerial),
		ConfigFile:            viper.GetString(KeyConfigFile),
		NoCache:               viper.GetBool(KeyNoCache),
		Profile:               viper.GetString(KeyProfile),
		SessionTTL:            viper.GetDuration(KeySessionTTL),
		SharedCredentialsFile: viper.GetString(KeySharedCredentialsFile),
		Verbose:               viper.GetBool(KeyVerbose),
	}

}
