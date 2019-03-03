// +build !windows

package commands

import (
	"runtime"
)

func getPlatform() *Platform {
	return &Platform{
		os:                   runtime.GOOS,
		shell:                "bash",
		shellArg:             "-c",
		escapedQuote:         "'",
		openCommand:          "open {{filename}}",
		openLinkCommand:      "open {{link}}",
		fallbackEscapedQuote: "\"",
	}
}
