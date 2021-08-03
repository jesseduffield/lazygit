package config

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() OSConfig {
	return OSConfig{
		Editor:          ``,
		EditCommand:     `{{editor}} {{filename}}`,
		OpenCommand:     `cmd /c "start "" {{filename}}"`,
		OpenLinkCommand: `cmd /c "start "" {{link}}"`,
	}
}
