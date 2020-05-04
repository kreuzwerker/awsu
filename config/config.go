package config

import (
	"fmt"
	"time"

	"github.com/gesellix/awsu/log"
)

const (
	errInvalidDuration = "invalid grace %q for duration %q"
)

// Config is the central configuration struct of awsu
type Config struct {
	ConfigFile            string        `mapstructure:"config-file"`
	Console               *Console      `mapstructure:"-"`
	Duration              time.Duration `mapstructure:"duration"`
	Generator             string        `mapstructure:"generator"`
	Grace                 time.Duration `mapstructure:"grace"`
	MFASerial             string        `mapstructure:"mfa-serial"`
	NoCache               bool          `mapstructure:"no-cache"`
	Profile               string        `mapstructure:"profile"`
	Profiles              Profiles      `mapstructure:"-"`
	Register              *Register     `mapstructure:"-"`
	SharedCredentialsFile string        `mapstructure:"shared-credentials-file"`
	Verbose               bool          `mapstructure:"verbose"`
}

// Init will perform post config initializations and validations
func (c *Config) Init() error {

	if c.Verbose {
		log.Verbose = true
	}

	profiles, err := Load(
		c.ConfigFile,
		c.SharedCredentialsFile,
	)

	if err != nil {
		return err
	}

	c.Profiles = profiles

	if c.Duration.Seconds() <= c.Grace.Seconds() {
		return fmt.Errorf(errInvalidDuration, c.Grace.String(), c.Duration.String())
	}

	return nil

}
