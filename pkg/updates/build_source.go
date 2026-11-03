package updates

// The build sources that lazygit knows about. The build source is set with
// `-X main.buildSource=...` when building lazygit.
const (
	// Set by mergeBuildInfo for builds made with
	// `go install github.com/jesseduffield/lazygit@<version>`
	BuildSourceGoInstall = "goInstall"
)
