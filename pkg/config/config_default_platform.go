//go:build !windows && !linux
// +build !windows,!linux

package config

const DefaultEditCommandTemplate = `{{editor}} {{filename}}`

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() OSConfig {
	return OSConfig{
		EditCommand:         ``,
		EditCommandTemplate: DefaultEditCommandTemplate,
		OpenCommand:         "open {{filename}}",
		OpenLinkCommand:     "open {{link}}",
	}
}
