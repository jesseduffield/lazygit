package config

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() OSConfig {
	return OSConfig{
		OpenCommand:     `cmd /c "start "" {{filename}}"`,
		OpenLinkCommand: `cmd /c "start "" {{link}}"`,
	}
}
