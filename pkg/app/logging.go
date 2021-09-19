//go:build !windows
// +build !windows

package app

import (
	"github.com/aybabtme/humanlog"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"log"
	"os"
)

func TailLogsForPlatform(logFilePath string, opts *humanlog.HandlerOptions) {
	cmd := secureexec.Command("tail", "-f", logFilePath)

	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if err := humanlog.Scanner(stdout, os.Stdout, opts); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
