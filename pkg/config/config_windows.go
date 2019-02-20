package config

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() []byte {
	return []byte(
		`os:
  openCommand: 'cmd /c "start "" {{filename}}"'
  openLinkCommand: 'cmd /c "start "" {{link}}"'`)
}
