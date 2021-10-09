package config

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() OSConfig {
	return OSConfig{
		EditCommand:         ``,
		EditCommandTemplate: `{{editor}} {{filename}}`,
		OpenCommand:         `xdg-open {{filename}} >/dev/null`,
		OpenLinkCommand:     `xdg-open {{link}} >/dev/null`,
	}
}
