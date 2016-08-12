package main

import "github.com/kreuzwerker/awsu/command"

var (
	build   string
	version string
)

func main() {
	command.Execute(version, build)
}
