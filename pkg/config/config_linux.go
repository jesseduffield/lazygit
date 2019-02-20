package config

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() []byte {
	return []byte(
		`os:
  openCommand: 'sh -c "xdg-open {{filename}} >/dev/null"'
  openLinkCommand: 'sh -c "xdg-open {{link}} >/dev/null"'`)
}
