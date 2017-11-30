package command

import (
	"fmt"
	"os"
)

const app = "awsu"

var (
	ctx  interface{}
	this = Version{}
)

func Execute(version, build, time string) {

	this.Build = build
	this.Time = time
	this.Version = version

	if _, err := rootCmd.ExecuteC(); err != nil {

		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)

	}

}
