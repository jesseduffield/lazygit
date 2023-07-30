//go:build !windows
// +build !windows

package tail

import (
	"log"
	"os"
	"os/exec"

	"github.com/aybabtme/humanlog"
)

func tailLogsForPlatform(logFilePath string, opts *humanlog.HandlerOptions) {
	cmd := exec.Command("tail", "-f", logFilePath)

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
