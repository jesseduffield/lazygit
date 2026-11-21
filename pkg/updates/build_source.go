package updates

// The build sources that lazygit knows about. The build source is set with
// `-X main.buildSource=...` when building lazygit.
const (
	// Set by .goreleaser.yml for our release archives
	BuildSourceRelease = "binaryRelease"
)
