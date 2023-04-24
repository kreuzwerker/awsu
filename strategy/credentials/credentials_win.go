// +build windows

package credentials

import (
	"os"
	"os/exec"
)

func exec_(cmd string, args []string, env []string) error {
	c := exec.Command(cmd, args[1:]...)
	c.Env = env
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
