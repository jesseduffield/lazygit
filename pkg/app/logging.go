// +build !windows

package app

import (
	"fmt"
	"github.com/aybabtme/humanlog"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"log"
	"os"
)

func TailLogs() {
	logFilePath, err := config.LogPath()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Tailing log file %s\n\n", logFilePath)

	opts := humanlog.DefaultOptions
	opts.Truncates = false

	_, err = os.Stat(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Log file does not exist. Run `lazygit --debug` first to create the log file")
		}
		log.Fatal(err)
	}

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
