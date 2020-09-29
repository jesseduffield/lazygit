// +build !windows

package oscommands

import (
	"runtime"
)

func getPlatform() *Platform {
	return &Platform{
		OS:                   runtime.GOOS,
		CatCmd:               "cat",
		Shell:                "bash",
		ShellArg:             "-c",
		EscapedQuote:         "'",
		OpenCommand:          "open {{filename}}",
		OpenLinkCommand:      "open {{link}}",
		FallbackEscapedQuote: "\"",
	}
}
