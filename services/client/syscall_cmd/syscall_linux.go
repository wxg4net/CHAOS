package syscall_cmd

import "syscall"

func GetCmdLinuxSyscall(cmd string) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
