package command

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const app = "awsu"

var this = Version{}

// Execute is the main entry point into the app
func Execute(version, build, time string) {

	whitespace := regexp.MustCompile(`(\s{2,})`)

	this.Build = build
	this.Time = time
	this.Version = version

	if _, err := rootCmd.ExecuteC(); err != nil {

		fmt.Fprintf(os.Stderr, "error: %s\n",
			whitespace.ReplaceAllString(err.Error(), ""))

		os.Exit(-1)

	}

}

// flag adds a with default value, long and short flags a matching environment
// (if not empty) and matching description
func flag(fs *pflag.FlagSet, def interface{}, long, short, env, desc string) {

	switch t := def.(type) {

	case bool:
		fs.BoolP(long, short, t, desc)
	case time.Duration:
		fs.DurationP(long, short, t, desc)
	case string:
		fs.StringP(long, short, t, desc)
	default:
		panic("unexpected default value")
	}

	viper.BindPFlag(long, fs.Lookup(long))

	if env != "" {
		viper.BindEnv(long, env)
	}

	viper.SetDefault(long, fs.Lookup(long).DefValue)

}
