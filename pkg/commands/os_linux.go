package commands

import (
	"runtime"
)

func getPlatform() *Platform {
	return &Platform{
		os:           runtime.GOOS,
		shell:        "bash",
		shellArg:     "-c",
		escapedQuote: "\"",
		openCommand:  "bash -c \"xdg-open {{filename}} &>/dev/null &\"",
	}
}
