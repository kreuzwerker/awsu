package main

import (
	"github.com/gesellix/awsu/command"
)

var (
	build   string
	time    string
	version string
)

func main() {
	command.Execute(version, build, time)
}
