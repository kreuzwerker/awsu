package config

import (
	"fmt"
	"time"

	"github.com/kreuzwerker/awsu/log"
)

// Config is the central configuration struct of awsu
type Config struct {
	Duration              time.Duration
	Grace                 time.Duration
	ConfigFile            string
	Generator             string
	MFASerial             string
	NoCache               bool
	Profile               string
	Profiles              Profiles
	SharedCredentialsFile string
	Verbose               bool
}

// Init will perform post config initializations and validations
func (c *Config) Init() error {

	if c.Verbose {
		log.Debug = true
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
		return fmt.Errorf("invalid grace %q for duration %q", c.Grace.String(), c.Duration.String())
	}

	return nil

}
