//go:build !windows
// +build !windows

package kill

import (
	"os/exec"
	"syscall"
)

// Kill kills a process. If the process has Setpgid == true, then we have anticipated that it might spawn its own child processes, so we've given it a process group ID (PGID) equal to its process id (PID) and given its child processes will inherit the PGID, we can kill that group, rather than killing the process itself.
func Kill(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		// You can't kill a person with no body
		return nil
	}

	if cmd.SysProcAttr != nil && cmd.SysProcAttr.Setpgid {
		// minus sign means we're talking about a PGID as opposed to a PID
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	return cmd.Process.Kill()
}

// PrepareForChildren ensures that child processes of this parent process will share the same group id
// as the parent, meaning when the call Kill on the parent process, we'll kill
// the whole group, parent and children both. Gruesome when you think about it.
func PrepareForChildren(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}
