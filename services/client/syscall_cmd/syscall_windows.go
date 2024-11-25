package syscall_cmd

import (
	"fmt"
	"syscall"
)

func GetCmdWindowsSyscall(cmd string) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c "%s"`, cmd)}
}
