//go:build !windows && !linux
// +build !windows,!linux

package config

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() OSConfig {
	return OSConfig{
		Open:     "open -- {{filename}}",
		OpenLink: "open {{link}}",
	}
}
