//go:build linux
// +build linux

package syscall_cmd

import "syscall"

func GetCmdSyscall(cmd string) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
