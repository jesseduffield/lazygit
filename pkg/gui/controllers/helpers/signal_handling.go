//go:build !windows

package helpers

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func canSuspendApp() bool {
	return true
}

func sendStopSignal() error {
	return syscall.Kill(0, syscall.SIGSTOP)
}

// setForegroundPgrp sets the current process group as the foreground process group
// for the terminal, allowing the program to read input after resuming from suspension.
func setForegroundPgrp() error {
	fd, err := unix.Open("/dev/tty", unix.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)

	pgid := syscall.Getpgrp()

	return unix.IoctlSetPointerInt(fd, unix.TIOCSPGRP, pgid)
}

func handleResumeSignal(log *logrus.Entry, onResume func() error) {
	if err := setForegroundPgrp(); err != nil {
		log.Warning(err)
		return
	}

	if err := onResume(); err != nil {
		log.Warning(err)
	}
}

func installResumeSignalHandler(log *logrus.Entry, onResume func() error) {
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGCONT)

		for sig := range sigs {
			switch sig {
			case syscall.SIGCONT:
				handleResumeSignal(log, onResume)
			}
		}
	}()
}
