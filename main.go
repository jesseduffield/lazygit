package main

import (
	"github.com/jesseduffield/lazygit/pkg/app"
)

// These values may be set by the build script via the LDFLAGS argument
var (
	commit      string
	date        string
	version     string
	buildSource = "unknown"
)

func main() {
	ldFlagsBuildInfo := &app.BuildInfo{
		Commit:      commit,
		Date:        date,
		Version:     version,
		BuildSource: buildSource,
	}

	app.Start(ldFlagsBuildInfo, nil)
}
