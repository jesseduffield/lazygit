// configdir provides access to configuration folder in each platforms.
//
// System wide configuration folders:
//
//   - Windows: %PROGRAMDATA% (C:\ProgramData)
//   - Linux/BSDs: ${XDG_CONFIG_DIRS} (/etc/xdg)
//   - MacOSX: "/Library/Application Support"
//
// User wide configuration folders:
//
//   - Windows: %APPDATA% (C:\Users\<User>\AppData\Roaming)
//   - Linux/BSDs: ${XDG_CONFIG_HOME} (${HOME}/.config)
//   - MacOSX: "${HOME}/Library/Application Support"
//
// User wide cache folders:
//
//   - Windows: %LOCALAPPDATA% (C:\Users\<User>\AppData\Local)
//   - Linux/BSDs: ${XDG_CACHE_HOME} (${HOME}/.cache)
//   - MacOSX: "${HOME}/Library/Caches"
//
// configdir returns paths inside the above folders.

package configdir

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type ConfigType int

const (
	System ConfigType = iota
	Global
	All
	Existing
	Local
	Cache
)

// Config represents each folder
type Config struct {
	Path string
	Type ConfigType
}

func (c Config) Open(fileName string) (*os.File, error) {
	return os.Open(filepath.Join(c.Path, fileName))
}

func (c Config) Create(fileName string) (*os.File, error) {
	err := c.CreateParentDir(fileName)
	if err != nil {
		return nil, err
	}
	return os.Create(filepath.Join(c.Path, fileName))
}

func (c Config) ReadFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(c.Path, fileName))
}

// CreateParentDir creates the parent directory of fileName inside c. fileName
// is a relative path inside c, containing zero or more path separators.
func (c Config) CreateParentDir(fileName string) error {
	return os.MkdirAll(filepath.Dir(filepath.Join(c.Path, fileName)), 0755)
}

func (c Config) WriteFile(fileName string, data []byte) error {
	err := c.CreateParentDir(fileName)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(c.Path, fileName), data, 0644)
}

func (c Config) MkdirAll() error {
	return os.MkdirAll(c.Path, 0755)
}

func (c Config) Exists(fileName string) bool {
	_, err := os.Stat(filepath.Join(c.Path, fileName))
	return !os.IsNotExist(err)
}

// ConfigDir keeps setting for querying folders.
type ConfigDir struct {
	VendorName      string
	ApplicationName string
	LocalPath       string
}

func New(vendorName, applicationName string) ConfigDir {
	return ConfigDir{
		VendorName:      vendorName,
		ApplicationName: applicationName,
	}
}

func (c ConfigDir) joinPath(root string) string {
	if c.VendorName != "" && hasVendorName {
		return filepath.Join(root, c.VendorName, c.ApplicationName)
	}
	return filepath.Join(root, c.ApplicationName)
}

func (c ConfigDir) QueryFolders(configType ConfigType) []*Config {
	if configType == Cache {
		return []*Config{c.QueryCacheFolder()}
	}
	var result []*Config
	if c.LocalPath != "" && configType != System && configType != Global {
		result = append(result, &Config{
			Path: c.LocalPath,
			Type: Local,
		})
	}
	if configType != System && configType != Local {
		result = append(result, &Config{
			Path: c.joinPath(globalSettingFolder),
			Type: Global,
		})
	}
	if configType != Global && configType != Local {
		for _, root := range systemSettingFolders {
			result = append(result, &Config{
				Path: c.joinPath(root),
				Type: System,
			})
		}
	}
	if configType != Existing {
		return result
	}
	var existing []*Config
	for _, entry := range result {
		if _, err := os.Stat(entry.Path); !os.IsNotExist(err) {
			existing = append(existing, entry)
		}
	}
	return existing
}

func (c ConfigDir) QueryFolderContainsFile(fileName string) *Config {
	configs := c.QueryFolders(Existing)
	for _, config := range configs {
		if _, err := os.Stat(filepath.Join(config.Path, fileName)); !os.IsNotExist(err) {
			return config
		}
	}
	return nil
}

func (c ConfigDir) QueryCacheFolder() *Config {
	return &Config{
		Path: c.joinPath(cacheFolder),
		Type: Cache,
	}
}
