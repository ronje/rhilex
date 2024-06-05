//go:build windows

package crontask

import (
	"os/exec"
	"syscall"
)

func GetSysProcAttr() *syscall.SysProcAttr {
	return nil
}

func KillProcess(proc *exec.Cmd) error {
	return proc.Process.Kill()
}
