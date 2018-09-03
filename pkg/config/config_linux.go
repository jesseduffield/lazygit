package config

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() []byte {
	return []byte(
		`os:
  openCommand: 'bash -c \"xdg-open {{filename}} &>/dev/null &\"'`)
}
