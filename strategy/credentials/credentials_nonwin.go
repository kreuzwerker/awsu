// +build !windows

package credentials

import "syscall"

func exec_(cmd string, args []string, env []string) error {
	return syscall.Exec(cmd, args, env)
}
