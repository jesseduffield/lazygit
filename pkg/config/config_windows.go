package config

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() OSConfig {
	return OSConfig{
		OpenCommand:     `start "" {{filename}}`,
		OpenLinkCommand: `start "" {{link}}`,
	}
}
