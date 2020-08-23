// +build !windows,!linux

package config

// GetPlatformDefaultConfig gets the defaults for the platform
func GetPlatformDefaultConfig() []byte {
	return []byte(
		`os:
  openCommand: 'open {{filename}}'
  openLinkCommand: 'open {{link}}'`)
}
