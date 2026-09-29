package theme

import (
	"os"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
)

// writeThemeFile creates a theme file in the themes folder of the config dir.
// Call it from SetupConfig: the config dir is recreated after SetupRepo runs.
func writeThemeFile(cfg *config.AppConfig, name string, content string) {
	if err := writeThemeFileTo(cfg.GetUserConfigDir(), name, content); err != nil {
		panic(err)
	}
}

// writeThemeFileTo creates or replaces the theme file with the given name in
// the themes folder of configDir.
func writeThemeFileTo(configDir string, name string, content string) error {
	themesDir := filepath.Join(configDir, "themes")
	if err := os.MkdirAll(themesDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(themesDir, name+".yml"), []byte(content), 0o644)
}

// selectThemeAtStartup makes lazygit start with the given theme selected, as
// if it had been chosen in an earlier session.
func selectThemeAtStartup(cfg *config.AppConfig, name string) {
	writeSelectedThemeFile(cfg, "name: "+name+"\n")
}

// writeSelectedThemeFile writes the file in which lazygit remembers the
// selected theme. It lives next to state.yml, which is in the config dir
// because integration tests run lazygit with its own config dir (CONFIG_DIR).
// Call it from SetupConfig, like writeThemeFile.
func writeSelectedThemeFile(cfg *config.AppConfig, content string) {
	path := filepath.Join(cfg.GetUserConfigDir(), "selected_theme.yml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		panic(err)
	}
}
