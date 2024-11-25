package syscall_cmd

import "syscall"

func GetCmdDarwinSyscall(cmd string) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
