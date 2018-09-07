package command

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const app = "awsu"

var this = Version{}

// Execute is the main entry point into the app
func Execute(version, build, time string) {

	this.Build = build
	this.Time = time
	this.Version = version

	if _, err := rootCmd.ExecuteC(); err != nil {

		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)

	}

}

// flag adds a flag with pointer value, default value, long and short flags,
// matching environment variable and matching description
func flag(fs *pflag.FlagSet, val, def interface{}, long, short, env, desc string) {

	switch t := val.(type) {

	case *bool:
		fs.BoolVarP(t, long, short, def.(bool), desc)
	case *time.Duration:
		fs.DurationVarP(t, long, short, def.(time.Duration), desc)
	case *string:
		fs.StringVarP(t, long, short, def.(string), desc)
	default:
		panic("unexpected value")
	}

	viper.BindPFlag(long, fs.Lookup(long))

	if env != "" {
		viper.BindEnv(long, env)
	}

}
