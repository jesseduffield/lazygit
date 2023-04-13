package config

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() OSConfig {
	return OSConfig{
		Open:     `start "" {{filename}}`,
		OpenLink: `start "" {{link}}`,
	}
}
