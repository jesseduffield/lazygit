package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/utils/yaml_utils"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

const (
	themesDirName         = "themes"
	themeFileExtension    = ".yml"
	selectedThemeFileName = "selected_theme.yml"
)

// ThemeFileError reports a theme file that couldn't be loaded.
type ThemeFileError struct {
	Path string
	Err  error
}

func (e *ThemeFileError) Error() string {
	return fmt.Sprintf("The theme file `%s` couldn't be loaded.\n%s", e.Path, e.Err)
}

func (e *ThemeFileError) Unwrap() error {
	return e.Err
}

// themeGuiConfig lists the only user config keys that a theme file may set.
type themeGuiConfig struct {
	Theme      ThemeConfig `yaml:"theme"`
	DarkTheme  ThemeConfig `yaml:"darkTheme"`
	LightTheme ThemeConfig `yaml:"lightTheme"`
}

type themeFile struct {
	Gui themeGuiConfig `yaml:"gui"`
}

// validateThemeFileContent checks that a theme file sets no other keys than
// the ones in themeGuiConfig, and that none of its values is empty (see
// checkThemeValues). Decoding strictly into these named types reports any
// other key together with its line, while the keys of the maps (author names
// and branch name patterns) stay free. An empty file is a valid theme that
// changes nothing.
func validateThemeFileContent(content []byte) error {
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)
	var decoded themeFile
	if err := decoder.Decode(&decoded); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return fmt.Errorf("A theme file may only set gui.theme, gui.darkTheme and gui.lightTheme.\n%w", err)
	}

	if err := checkThemeValues(content); err != nil {
		return fmt.Errorf("A theme file can only set colors, not remove them, so none of its values may be empty.\n%w", err)
	}

	return nil
}

// themeFileValues holds the values of a theme file's settings as YAML nodes,
// which tell a missing or empty value apart from one that sets a color. They
// are grouped by the key of gui that they are in: theme, darkTheme or
// lightTheme.
type themeFileValues struct {
	Gui map[string]map[string]yaml.Node `yaml:"gui"`
}

// checkThemeValues returns an error for the first setting of a theme file, in
// the order of their names, whose value is null, an empty list or an empty
// map, or contains one of these. None of them sets a color, which is all that
// a theme is for, and a null or an empty list even clears what the config
// files loaded before the theme have set. The content must already have
// passed the strict decoding in validateThemeFileContent, so every key is
// known.
func checkThemeValues(content []byte) error {
	var values themeFileValues
	if err := yaml.Unmarshal(content, &values); err != nil {
		return err
	}

	// Each key of gui.theme, gui.darkTheme and gui.lightTheme is a setting of
	// its own
	settings := map[string]yaml.Node{}
	for guiKey, themeSettings := range values.Gui {
		for key, node := range themeSettings {
			settings["gui."+guiKey+"."+key] = node
		}
	}

	checked := map[*yaml.Node]bool{}
	for _, path := range slices.Sorted(maps.Keys(settings)) {
		node := settings[path]
		if err := checkThemeValue(&node, path, checked); err != nil {
			return err
		}
	}

	return nil
}

// checkThemeValue returns an error if node, the value at path, is null or an
// empty list or map, or contains such a value. It follows aliases and checks
// every node only once, so that an alias inside its own anchor can't make it
// loop.
func checkThemeValue(node *yaml.Node, path string, checked map[*yaml.Node]bool) error {
	if checked[node] {
		return nil
	}
	checked[node] = true

	switch node.Kind {
	case yaml.AliasNode:
		return checkThemeValue(node.Alias, path, checked)
	case yaml.SequenceNode:
		if len(node.Content) == 0 {
			return fmt.Errorf("%s is an empty list", path)
		}
		for i, item := range node.Content {
			if err := checkThemeValue(item, fmt.Sprintf("%s[%d]", path, i), checked); err != nil {
				return err
			}
		}
	case yaml.MappingNode:
		if len(node.Content) == 0 {
			return fmt.Errorf("%s is an empty map", path)
		}
		for i := 0; i < len(node.Content); i += 2 {
			key, value := node.Content[i], node.Content[i+1]
			if err := checkThemeValue(value, fmt.Sprintf("%s[%q]", path, key.Value), checked); err != nil {
				return err
			}
		}
	default:
		if node.ShortTag() == "!!null" {
			return fmt.Errorf("%s has no value", path)
		}
	}

	return nil
}

// convertOldThemeLayout returns the content of a theme file with the settings
// that themes for older versions of lazygit have directly in gui moved into
// gui.theme, by the same steps that migrate them in config files (see
// migrateThemeKeys), and whether it had to move any. It applies none of the
// other migrations, since those reject the YAML aliases that a theme may use
// to share colors. Content that needs no conversion is returned as it is.
func convertOldThemeLayout(content []byte) ([]byte, bool, error) {
	var rootNode yaml.Node
	if err := yaml.Unmarshal(content, &rootNode); err != nil || !hasGuiMap(&rootNode) {
		// Only a gui map can have the old layout. validateThemeFileContent
		// reports what is wrong with anything else, and accepts a gui without
		// a value, which migrateThemeKeys would reject.
		return content, false, nil
	}

	changes := NewChangesSet()
	if err := migrateThemeKeys(&rootNode, "Couldn't convert the theme to the current layout", changes); err != nil {
		return nil, false, err
	}
	if changes.Len() == 0 {
		return content, false, nil
	}

	convertedContent, err := yaml_utils.YamlMarshal(&rootNode)
	if err != nil {
		return nil, false, err
	}
	return convertedContent, true, nil
}

// convertedThemeHint follows an error in the content of a theme that
// convertOldThemeLayout has converted. The error's lines and keys are those of
// the converted content, and moving settings to the end of gui.theme can put
// an alias ahead of its anchor.
const convertedThemeHint = "This theme has authorColors, branchColorPatterns or branchColors " +
	"directly in gui, as themes for older versions of lazygit do, so it was converted in memory, " +
	"with these settings moved to the end of gui.theme; the lines and keys above refer to the " +
	"converted layout. Moving authorColors and branchColorPatterns into gui.theme yourself, and " +
	"replacing branchColors with gui.theme.branchColorPatterns, avoids the conversion."

// hasGuiMap reports whether the parsed YAML document is a map whose gui key is
// a map too.
func hasGuiMap(rootNode *yaml.Node) bool {
	if len(rootNode.Content) == 0 || rootNode.Content[0].Kind != yaml.MappingNode {
		return false
	}

	_, guiNode := yaml_utils.LookupKey(rootNode.Content[0], "gui")
	return guiNode != nil && guiNode.Kind == yaml.MappingNode
}

// readThemeFile reads a theme file, converts it from the old layout (see
// convertOldThemeLayout), and checks that it only sets theme keys. Theme files
// are often downloaded, or symlinked from a dotfiles repo, so unlike the
// user's own config files they are converted only in memory, and never
// rewritten.
func readThemeFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content, converted, err := convertOldThemeLayout(content)
	if err != nil {
		return nil, err
	}

	if err := validateThemeFileContent(content); err != nil {
		if converted {
			return nil, fmt.Errorf("%w\n%s", err, convertedThemeHint)
		}
		return nil, err
	}

	return content, nil
}

// isValidThemeName reports whether name can be the name of a file in the
// themes folder. The name is read from a file that users may edit by hand, so
// it must not be able to reach outside that folder; and names starting with a
// dot are editor lock files and the like, not themes.
func isValidThemeName(name string) bool {
	return name != "" && filepath.Base(name) == name && !strings.HasPrefix(name, ".")
}

// selectedThemeState is the content of the file that remembers the selected
// theme. It is kept out of AppState because every running lazygit rewrites
// state.yml from memory, which would undo a theme selected in another one.
type selectedThemeState struct {
	Name string `yaml:"name"`
}

func selectedThemeFilePath() (string, error) {
	return stateSiblingFilePath(selectedThemeFileName)
}

// loadSelectedThemeName returns the name of the selected theme, or "" if no
// theme is selected, which is also the case when the file is missing or empty.
func loadSelectedThemeName() (string, error) {
	path, err := selectedThemeFilePath()
	if err != nil {
		if os.IsPermission(err) {
			// apparently when people have read-only permissions they prefer us to fail silently
			return "", nil
		}
		return "", err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	var state selectedThemeState
	if err := yaml.Unmarshal(content, &state); err != nil {
		return "", fmt.Errorf("The file `%s` couldn't be parsed.\n%w", path, err)
	}

	return state.Name, nil
}

// GetThemesDir returns the folder that holds the theme files, or "" if there
// is no config dir (as with NewDummyAppConfig), in which case the folder would
// resolve relative to the working directory.
func (c *AppConfig) GetThemesDir() string {
	if c.userConfigDir == "" {
		return ""
	}

	return filepath.Join(c.userConfigDir, themesDirName)
}

// ListThemes returns the names of the themes in the themes folder, sorted
// case-insensitively. A theme is a regular file, or a link to one, whose name
// ends in .yml and doesn't start with a dot; its name is the file name without
// that extension. A missing themes folder has no themes.
func (c *AppConfig) ListThemes() ([]string, error) {
	themesDir := c.GetThemesDir()
	if themesDir == "" {
		return nil, nil
	}

	entries, err := os.ReadDir(themesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		name, ok := strings.CutSuffix(entry.Name(), themeFileExtension)
		if !ok || !isValidThemeName(name) {
			continue
		}

		// os.Stat follows links, so that a link to a theme file (e.g. into a
		// dotfiles repo) counts as a theme, while a broken link or a folder
		// doesn't
		info, err := os.Stat(filepath.Join(themesDir, entry.Name()))
		if err != nil || !info.Mode().IsRegular() {
			continue
		}

		names = append(names, name)
	}

	// os.ReadDir sorts by bytes, which puts "Zed" before "alpha"
	slices.SortStableFunc(names, func(a, b string) int {
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	})
	return names, nil
}

// resolveThemeName returns the name of the listed theme that the given name,
// as read from selected_theme.yml, refers to. A name edited by hand may differ
// in case from the file name: a file system that ignores case would still find
// the file, but the name wouldn't match the listed one, and other file systems
// wouldn't find the file at all. A name that no listed theme matches, which
// includes every name when the themes folder can't be listed, is returned as
// it is.
func (c *AppConfig) resolveThemeName(name string) string {
	if name == "" {
		return ""
	}

	themes, _ := c.ListThemes()
	if slices.Contains(themes, name) {
		return name
	}
	if listedName, ok := lo.Find(themes, func(theme string) bool {
		return strings.EqualFold(theme, name)
	}); ok {
		return listedName
	}

	return name
}

// GetSelectedTheme returns the name of the selected theme, or "" if none is
// selected. It returns the name even when the theme's file doesn't exist or
// couldn't be loaded; a name that matches a listed theme is spelled like it
// (see resolveThemeName).
func (c *AppConfig) GetSelectedTheme() string {
	return c.selectedTheme
}

// GetAppliedTheme returns the name of the theme that the user config in use
// was loaded with, or "": when no theme is selected, when the selected theme's
// file was missing at the last successful load, or when it couldn't be loaded
// (see GetThemeLoadError).
func (c *AppConfig) GetAppliedTheme() string {
	return c.appliedTheme
}

// loadedThemeName returns the name of the selected theme if its file is among
// the user config files and existed when they were last loaded, or "". Call it
// only after a load that succeeded: a load that fails has already recorded
// which files exist, but the config from before it stays in use.
func (c *AppConfig) loadedThemeName() string {
	// The only theme file among the user config files is the selected theme's
	isLoadedThemeFile := func(f *ConfigFile) bool { return f.isTheme && f.exists }
	if !lo.SomeBy(c.userConfigFiles, isLoadedThemeFile) {
		return ""
	}

	return c.selectedTheme
}

// GetThemeLoadError returns the error that kept the last
// ReloadUserConfigForRepo from applying the selected theme, or nil. It is
// either the theme file's error, or the error reading which theme is
// selected, in which case GetSelectedTheme returns "".
func (c *AppConfig) GetThemeLoadError() error {
	return c.themeLoadError
}

// themeConfigFile returns the config file of the theme with the given name, or
// nil if there can't be one. A missing theme file is skipped, not reported:
// the file may come back (e.g. with the next dotfiles sync), and as part of the
// user config files it is then picked up by the reload on focus.
func (c *AppConfig) themeConfigFile(name string) *ConfigFile {
	themesDir := c.GetThemesDir()
	if themesDir == "" || !isValidThemeName(name) {
		return nil
	}

	return &ConfigFile{
		Path:    filepath.Join(themesDir, name+themeFileExtension),
		Policy:  ConfigFilePolicySkipIfMissing,
		isTheme: true,
	}
}

// composeUserConfigFiles returns the files that make up the user config, in
// order of increasing precedence: the global files, the theme (if any), and
// the repo files. The theme comes after the global files so that selecting it
// overrides the colors in config.yml, and before the repo files so that a repo
// can still have colors of its own.
func (c *AppConfig) composeUserConfigFiles(theme *ConfigFile, repoConfigFiles []*ConfigFile) []*ConfigFile {
	if theme == nil {
		return slices.Concat(c.globalUserConfigFiles, repoConfigFiles)
	}

	return slices.Concat(c.globalUserConfigFiles, []*ConfigFile{theme}, repoConfigFiles)
}
