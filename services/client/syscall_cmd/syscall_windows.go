//go:build windows
// +build windows

package syscall_cmd

import (
	"fmt"
	"syscall"
)

func GetCmdSyscall(cmd string) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c "%s"`, cmd)}
}
