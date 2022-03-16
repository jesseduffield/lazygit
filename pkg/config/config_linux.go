package config

const DefaultEditCommandTemplate = `{{editor}} {{filename}}`

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() OSConfig {
	return OSConfig{
		EditCommand:         ``,
		EditCommandTemplate: DefaultEditCommandTemplate,
		OpenCommand:         `xdg-open {{filename}} >/dev/null`,
		OpenLinkCommand:     `xdg-open {{link}} >/dev/null`,
	}
}
