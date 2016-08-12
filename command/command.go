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

func Execute(version, build string) {

	this.Build = build
	this.Version = version

	if _, err := rootCmd.ExecuteC(); err != nil {

		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)

	}

}
