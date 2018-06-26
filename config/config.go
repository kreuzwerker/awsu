package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	KeyCacheTTL              = "cache-ttl"
	KeyConfigFile            = "config-file"
	KeyGenerator             = "generator"
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
	Generator             string
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
		ConfigFile:            viper.GetString(KeyConfigFile),
		Generator:             viper.GetString(KeyGenerator),
		MFASerial:             viper.GetString(KeyMFASerial),
		NoCache:               viper.GetBool(KeyNoCache),
		Profile:               viper.GetString(KeyProfile),
		SessionTTL:            viper.GetDuration(KeySessionTTL),
		SharedCredentialsFile: viper.GetString(KeySharedCredentialsFile),
		Verbose:               viper.GetBool(KeyVerbose),
	}

}
