package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

const (
	allowedThemeKeysMessage = "may only set gui.theme, gui.darkTheme and gui.lightTheme"
	emptyThemeValueMessage  = "none of its values may be empty"
	convertedThemeMessage   = "so it was converted in memory"
)

func TestValidateThemeFileContent(t *testing.T) {
	scenarios := []struct {
		name           string
		content        string
		expectedErrors []string
	}{
		{
			name:    "All allowed keys",
			content: "gui: {theme: {activeBorderColor: ['#ff00ff', bold], defaultFgColor: [white], authorColors: {'*': '#b4befe'}, branchColorPatterns: {'^feature/': green}}, darkTheme: {selectedLineBgColor: [blue]}, lightTheme: {selectedLineBgColor: ['#ccd0da'], authorColors: {'*': '#7287fd'}}}",
		},
		{
			name:    "Empty file",
			content: "",
		},
		{
			name:    "Only a comment",
			content: "# no colors yet\n",
		},
		{
			name:    "Gui without a value",
			content: "gui:\n",
		},
		{
			name:    "Themes without a value",
			content: "gui:\n  theme:\n  darkTheme:\n  lightTheme:\n",
		},
		{
			name:    "Anchors and aliases within the allowed keys",
			content: "gui: {theme: {activeBorderColor: &accent ['#ff00ff', bold]}, darkTheme: {searchingActiveBorderColor: *accent}}",
		},
		{
			name:    "Merge keys within the allowed keys",
			content: "gui: {theme: {authorColors: &palette {'*': '#b4befe'}, branchColorPatterns: {<<: *palette, '^feature/': green}}}",
		},
		{
			name:    "Merge keys across themes",
			content: "gui: {theme: {authorColors: &palette {'*': '#b4befe'}}, lightTheme: {authorColors: {<<: *palette, John: green}}}",
		},
		{
			// yaml skips the entry whose key is null, so this only sets John's
			// color; checking the values must not follow the alias forever
			name:    "Alias inside its own anchor",
			content: "gui: {theme: {authorColors: &authors {John: red, ~: *authors}}}",
		},
		{
			name:           "Unknown top-level key",
			content:        "git: {autoFetch: false}",
			expectedErrors: []string{allowedThemeKeysMessage, "field git not found"},
		},
		{
			name:           "Gui key that isn't about colors",
			content:        "gui: {nerdFontsVersion: '3'}",
			expectedErrors: []string{allowedThemeKeysMessage, "field nerdFontsVersion not found"},
		},
		{
			name:           "Color scheme",
			content:        "gui: {colorScheme: light, lightTheme: {activeBorderColor: [blue]}}",
			expectedErrors: []string{allowedThemeKeysMessage, "field colorScheme not found"},
		},
		{
			// convertOldThemeLayout moves these into gui.theme before the
			// content is validated
			name:           "Theme setting directly in gui",
			content:        "gui: {authorColors: {'*': '#b4befe'}}",
			expectedErrors: []string{allowedThemeKeysMessage, "field authorColors not found"},
		},
		{
			name:           "Unknown theme key",
			content:        "gui: {theme: {selectedRangeBgColor: [blue]}}",
			expectedErrors: []string{allowedThemeKeysMessage, "field selectedRangeBgColor not found"},
		},
		{
			name:           "Branch color patterns that aren't a map",
			content:        "gui: {theme: {branchColorPatterns: ['^feature/']}}",
			expectedErrors: []string{allowedThemeKeysMessage, "cannot unmarshal !!seq into map[string]string"},
		},
		{
			name:           "Unknown dark theme key",
			content:        "gui: {darkTheme: {colorScheme: dark}}",
			expectedErrors: []string{allowedThemeKeysMessage, "field colorScheme not found"},
		},
		{
			name:           "Theme keys without gui",
			content:        "theme: {activeBorderColor: [red]}",
			expectedErrors: []string{allowedThemeKeysMessage, "field theme not found"},
		},
		{
			name:           "Invalid yaml",
			content:        "gui: {theme: [",
			expectedErrors: []string{allowedThemeKeysMessage, "yaml:"},
		},
		{
			name:           "Author colors without a value",
			content:        "gui:\n  theme:\n    authorColors:\n",
			expectedErrors: []string{emptyThemeValueMessage, "gui.theme.authorColors has no value"},
		},
		{
			name:           "Border color set to null",
			content:        "gui:\n  theme:\n    activeBorderColor: null\n",
			expectedErrors: []string{emptyThemeValueMessage, "gui.theme.activeBorderColor has no value"},
		},
		{
			name:           "Empty list of colors",
			content:        "gui: {theme: {activeBorderColor: []}}",
			expectedErrors: []string{emptyThemeValueMessage, "gui.theme.activeBorderColor is an empty list"},
		},
		{
			name:           "Empty map of branch colors",
			content:        "gui: {theme: {branchColorPatterns: {}}}",
			expectedErrors: []string{emptyThemeValueMessage, "gui.theme.branchColorPatterns is an empty map"},
		},
		{
			name:           "Null in a list of colors",
			content:        "gui: {theme: {activeBorderColor: ['#ff00ff', ~]}}",
			expectedErrors: []string{emptyThemeValueMessage, "gui.theme.activeBorderColor[1] has no value"},
		},
		{
			name:           "Author without a color",
			content:        "gui: {theme: {authorColors: {John: ~}}}",
			expectedErrors: []string{emptyThemeValueMessage, `gui.theme.authorColors["John"] has no value`},
		},
		{
			name:           "Branch pattern without a color in the dark theme",
			content:        "gui: {darkTheme: {branchColorPatterns: {'^feature/': ~}}}",
			expectedErrors: []string{emptyThemeValueMessage, `gui.darkTheme.branchColorPatterns["^feature/"] has no value`},
		},
		{
			name:           "Empty list of colors in the light theme",
			content:        "gui: {lightTheme: {selectedLineBgColor: []}}",
			expectedErrors: []string{emptyThemeValueMessage, "gui.lightTheme.selectedLineBgColor is an empty list"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			err := validateThemeFileContent([]byte(s.content))
			if len(s.expectedErrors) == 0 {
				assert.NoError(t, err)
				return
			}
			for _, expectedError := range s.expectedErrors {
				assert.ErrorContains(t, err, expectedError)
			}
		})
	}
}

func TestConvertOldThemeLayout(t *testing.T) {
	scenarios := []struct {
		name              string
		content           string
		expected          string
		expectedConverted bool
		expectedError     string
	}{
		{
			// Marshaling would indent this by two spaces
			name:     "New layout",
			content:  "gui:\n    theme:\n        authorColors:\n            '*': '#b4befe'\n",
			expected: "gui:\n    theme:\n        authorColors:\n            '*': '#b4befe'\n",
		},
		{
			name:     "Empty file",
			content:  "",
			expected: "",
		},
		{
			// MoveYamlKey rejects a gui that isn't a map
			name:     "Gui without a value",
			content:  "gui:\n",
			expected: "gui:\n",
		},
		{
			name:     "Gui that isn't a map is left for validation",
			content:  "gui: [authorColors]\n",
			expected: "gui: [authorColors]\n",
		},
		{
			name:     "Invalid yaml is left for validation",
			content:  "gui: {authorColors: [",
			expected: "gui: {authorColors: [",
		},
		{
			name: "Author colors and branch color patterns move into the theme",
			content: "gui:\n" +
				"  authorColors:\n" +
				"    '*': '#b4befe'\n" +
				"  theme:\n" +
				"    activeBorderColor:\n" +
				"      - '#89b4fa'\n" +
				"      - bold\n" +
				"  branchColorPatterns:\n" +
				"    '^feature/': green\n",
			expected: "gui:\n" +
				"  theme:\n" +
				"    activeBorderColor:\n" +
				"      - '#89b4fa'\n" +
				"      - bold\n" +
				"    authorColors:\n" +
				"      '*': '#b4befe'\n" +
				"    branchColorPatterns:\n" +
				"      '^feature/': green\n",
			expectedConverted: true,
		},
		{
			name: "Branch colors become branch color patterns in the theme",
			content: "gui:\n" +
				"  branchColors:\n" +
				"    feature: green\n",
			expected: "gui:\n" +
				"  theme:\n" +
				"    branchColorPatterns:\n" +
				"      ^feature(/|$): green\n",
			expectedConverted: true,
		},
		{
			name: "Anchors and aliases",
			content: "gui:\n" +
				"  theme:\n" +
				"    activeBorderColor:\n" +
				"      - &accent '#ff00ff'\n" +
				"  authorColors: &authors\n" +
				"    '*': *accent\n" +
				"  darkTheme:\n" +
				"    authorColors:\n" +
				"      <<: *authors\n" +
				"      John: green\n",
			expected: "gui:\n" +
				"  theme:\n" +
				"    activeBorderColor:\n" +
				"      - &accent '#ff00ff'\n" +
				"    authorColors: &authors\n" +
				"      '*': *accent\n" +
				"  darkTheme:\n" +
				"    authorColors:\n" +
				"      !!merge <<: *authors\n" +
				"      John: green\n",
			expectedConverted: true,
		},
		{
			// The moved setting goes to the end of gui.theme, so the result
			// doesn't parse; readThemeFile explains that
			name: "Alias ahead of the anchor that moves",
			content: "gui:\n" +
				"  authorColors:\n" +
				"    '*': &lavender '#b4befe'\n" +
				"  theme:\n" +
				"    activeBorderColor:\n" +
				"      - *lavender\n",
			expected: "gui:\n" +
				"  theme:\n" +
				"    activeBorderColor:\n" +
				"      - *lavender\n" +
				"    authorColors:\n" +
				"      '*': &lavender '#b4befe'\n",
			expectedConverted: true,
		},
		{
			name: "Author colors both in gui and in the theme",
			content: "gui:\n" +
				"  authorColors:\n" +
				"    '*': '#b4befe'\n" +
				"  theme:\n" +
				"    authorColors:\n" +
				"      John: green\n",
			expectedError: "Couldn't convert the theme to the current layout for key gui.authorColors: new key `authorColors' already exists",
		},
		{
			name: "Branch colors in gui and branch color patterns in the theme",
			content: "gui:\n" +
				"  branchColors:\n" +
				"    feature: green\n" +
				"  theme:\n" +
				"    branchColorPatterns:\n" +
				"      '^docs/': blue\n",
			expectedError: "for key gui.branchColorPatterns: new key `branchColorPatterns' already exists",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			actual, converted, err := convertOldThemeLayout([]byte(s.content))
			if s.expectedError != "" {
				assert.ErrorContains(t, err, s.expectedError)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, s.expected, string(actual))
			assert.Equal(t, s.expectedConverted, converted)
		})
	}
}

func TestIsValidThemeName(t *testing.T) {
	scenarios := []struct {
		name     string
		expected bool
	}{
		{name: "dark", expected: true},
		{name: "Catppuccin Mocha", expected: true},
		{name: "", expected: false},
		{name: ".", expected: false},
		{name: "..", expected: false},
		{name: ".hidden", expected: false},
		{name: "../dark", expected: false},
		{name: "sub/dark", expected: false},
	}

	for _, s := range scenarios {
		assert.Equal(t, s.expected, isValidThemeName(s.name), "theme name %q", s.name)
	}
}

func TestSelectedThemeFilePath(t *testing.T) {
	stateDir := t.TempDir()
	t.Setenv("CONFIG_DIR", stateDir)

	path, err := selectedThemeFilePath()

	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(stateDir, selectedThemeFileName), path)
}

func TestLoadSelectedThemeName(t *testing.T) {
	stateDir := t.TempDir()
	t.Setenv("CONFIG_DIR", stateDir)
	path := filepath.Join(stateDir, selectedThemeFileName)

	name, err := loadSelectedThemeName()
	assert.NoError(t, err)
	assert.Equal(t, "", name, "missing file")

	writeThemeTestFile(t, path, "")
	name, err = loadSelectedThemeName()
	assert.NoError(t, err)
	assert.Equal(t, "", name, "empty file")

	writeThemeTestFile(t, path, "name: dark\n")
	name, err = loadSelectedThemeName()
	assert.NoError(t, err)
	assert.Equal(t, "dark", name)

	writeThemeTestFile(t, path, "name: [dark\n")
	_, err = loadSelectedThemeName()
	assert.ErrorContains(t, err, path)
}

func TestGetThemesDir(t *testing.T) {
	assert.Equal(t, "", NewDummyAppConfig().GetThemesDir())
	assert.Nil(t, NewDummyAppConfig().themeConfigFile("dark"))

	appConfig, configDir := newThemeTestAppConfig(t, "")
	assert.Equal(t, filepath.Join(configDir, "themes"), appConfig.GetThemesDir())
}

// Config files passed in LG_CONFIG_FILE replace config.yml, but not the config
// dir, so the themes are still found in it. A theme overrides all of those
// files, and the repo config still overrides the theme.
func TestThemesWithConfigFilesFromEnvironment(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("CONFIG_DIR", configDir)
	customConfigDir := t.TempDir()
	firstConfigPath := filepath.Join(customConfigDir, "first.yml")
	writeThemeTestFile(t, firstConfigPath, `gui:
  theme:
    inactiveBorderColor:
      - blue
`)
	secondConfigPath := filepath.Join(customConfigDir, "second.yml")
	writeThemeTestFile(t, secondConfigPath, `gui:
  theme:
    activeBorderColor:
      - red
`)
	t.Setenv("LG_CONFIG_FILE", firstConfigPath+","+secondConfigPath)
	writeThemeTestFile(t, filepath.Join(configDir, "themes", "pink.yml"), `gui:
  theme:
    activeBorderColor:
      - '#ff00ff'
    optionsTextColor:
      - '#00ff00'
`)
	selectThemeForTest(t, configDir, "pink")
	repoConfigPath := filepath.Join(t.TempDir(), "lazygit.yml")
	writeThemeTestFile(t, repoConfigPath, `gui:
  theme:
    optionsTextColor:
      - yellow
`)
	repoConfigFiles := []*ConfigFile{{Path: repoConfigPath, Policy: ConfigFilePolicySkipIfMissing}}
	appConfig, err := NewAppConfig("lazygit", "unversioned", "", "", "", false, t.TempDir())
	assert.NoError(t, err)

	assert.Equal(t, filepath.Join(configDir, "themes"), appConfig.GetThemesDir())
	themes, err := appConfig.ListThemes()
	assert.NoError(t, err)
	assert.Equal(t, []string{"pink"}, themes)

	err = appConfig.ReloadUserConfigForRepo(repoConfigFiles)

	assert.NoError(t, err)
	assert.NoError(t, appConfig.GetThemeLoadError())
	assert.Equal(t, "pink", appConfig.GetAppliedTheme())
	themeConfig := appConfig.GetUserConfig().Gui.Theme
	assert.Equal(t, []string{"blue"}, themeConfig.InactiveBorderColor, "the first file still applies where the theme is silent")
	assert.Equal(t, []string{"#ff00ff"}, themeConfig.ActiveBorderColor, "the theme overrides the last file")
	assert.Equal(t, []string{"yellow"}, themeConfig.OptionsTextColor, "the repo config overrides the theme")
}

func TestListThemes(t *testing.T) {
	themes, err := NewDummyAppConfig().ListThemes()
	assert.NoError(t, err)
	assert.Nil(t, themes, "no config dir")

	appConfig, configDir := newThemeTestAppConfig(t, "")

	themes, err = appConfig.ListThemes()
	assert.NoError(t, err)
	assert.Nil(t, themes, "missing themes folder")

	themesDir := filepath.Join(configDir, "themes")
	for _, fileName := range []string{
		"zed.yml",
		"Dark.yml",
		"alpha beta.yml",
		"light.yaml",
		"shout.YML",
		".hidden.yml",
		".yml",
		"notes.txt",
	} {
		writeThemeTestFile(t, filepath.Join(themesDir, fileName), "")
	}
	assert.NoError(t, os.Mkdir(filepath.Join(themesDir, "folder.yml"), 0o755))

	themes, err = appConfig.ListThemes()
	assert.NoError(t, err)
	assert.Equal(t, []string{"alpha beta", "Dark", "zed"}, themes)
}

func TestListThemesFollowsLinks(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	themesDir := filepath.Join(configDir, "themes")
	assert.NoError(t, os.MkdirAll(themesDir, 0o755))
	target := filepath.Join(t.TempDir(), "mocha.yml")
	writeThemeTestFile(t, target, "")
	if err := os.Symlink(target, filepath.Join(themesDir, "mocha.yml")); err != nil {
		t.Skipf("can't create symlinks here: %v", err)
	}
	assert.NoError(t, os.Symlink(filepath.Join(themesDir, "missing.yml"), filepath.Join(themesDir, "broken.yml")))

	themes, err := appConfig.ListThemes()

	assert.NoError(t, err)
	assert.Equal(t, []string{"mocha"}, themes)
}

func TestReloadUserConfigForRepoAppliesThemeBetweenGlobalAndRepoConfig(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, `gui:
  theme:
    activeBorderColor:
      - red
    inactiveBorderColor:
      - blue
    authorColors:
      John: red
      '*': white
`)
	writeThemeTestFile(t, filepath.Join(configDir, "themes", "pink.yml"), `gui:
  theme:
    activeBorderColor:
      - '#ff00ff'
    optionsTextColor:
      - '#00ff00'
    authorColors:
      '*': '#ff00ff'
`)
	selectThemeForTest(t, configDir, "pink")
	repoConfigPath := filepath.Join(t.TempDir(), "lazygit.yml")
	writeThemeTestFile(t, repoConfigPath, `gui:
  theme:
    optionsTextColor:
      - yellow
`)
	repoConfigFiles := []*ConfigFile{{Path: repoConfigPath, Policy: ConfigFilePolicySkipIfMissing}}

	err := appConfig.ReloadUserConfigForRepo(repoConfigFiles)

	assert.NoError(t, err)
	assert.NoError(t, appConfig.GetThemeLoadError())
	assert.Equal(t, "pink", appConfig.GetSelectedTheme())
	assert.Equal(t, "pink", appConfig.GetAppliedTheme())
	themeConfig := appConfig.GetUserConfig().Gui.Theme
	assert.Equal(t, []string{"#ff00ff"}, themeConfig.ActiveBorderColor, "the theme overrides the global config")
	assert.Equal(t, []string{"blue"}, themeConfig.InactiveBorderColor, "the global config stays where the theme is silent")
	assert.Equal(t, []string{"yellow"}, themeConfig.OptionsTextColor, "the repo config overrides the theme")
	assert.Equal(t,
		map[string]string{"John": "red", "*": "#ff00ff"},
		themeConfig.AuthorColors,
		"maps are merged key by key",
	)
	assert.Equal(t, repoConfigFiles, appConfig.repoUserConfigFiles)
	assert.Len(t, appConfig.globalUserConfigFiles, 1)
}

// A later file's branch color patterns come before those of the files loaded
// before it, so that they are tried first.
func TestBranchColorPatternsOfThemeComeBetweenRepoAndGlobalOnes(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, `gui:
  theme:
    branchColorPatterns:
      '^main$': red
      '^docs/': blue
`)
	writeThemeTestFile(t, filepath.Join(configDir, "themes", "pink.yml"), `gui:
  theme:
    branchColorPatterns:
      '^feature/': green
      '^docs/': '#ff00ff'
`)
	selectThemeForTest(t, configDir, "pink")
	repoConfigPath := filepath.Join(t.TempDir(), "lazygit.yml")
	writeThemeTestFile(t, repoConfigPath, `gui:
  theme:
    branchColorPatterns:
      '^fix/': yellow
`)
	repoConfigFiles := []*ConfigFile{{Path: repoConfigPath, Policy: ConfigFilePolicySkipIfMissing}}

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(repoConfigFiles))

	assert.Equal(t,
		ColorPatterns{
			{Pattern: "^fix/", Color: "yellow"},
			{Pattern: "^feature/", Color: "green"},
			{Pattern: "^docs/", Color: "#ff00ff"},
			{Pattern: "^main$", Color: "red"},
		},
		appConfig.GetUserConfig().Gui.Theme.BranchColorPatterns,
	)
}

// gui.darkTheme and gui.lightTheme apply after all config files are merged, so
// the ones in config.yml override the selected theme's gui.theme on their
// background.
func TestConfigDarkThemeOverridesSelectedThemeOnDarkBackground(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, `gui:
  darkTheme:
    activeBorderColor:
      - red
`)
	writeThemeTestFile(t, filepath.Join(configDir, "themes", "pink.yml"), `gui:
  theme:
    activeBorderColor:
      - '#ff00ff'
`)
	selectThemeForTest(t, configDir, "pink")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	guiConfig := appConfig.GetUserConfig().Gui
	assert.Equal(t, []string{"red"}, guiConfig.ThemeForBackground(false, "").ActiveBorderColor)
	assert.Equal(t, []string{"#ff00ff"}, guiConfig.ThemeForBackground(true, "").ActiveBorderColor)
}

func TestThemeFileDarkAndLightThemesOverrideItsTheme(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, `gui:
  darkTheme:
    activeBorderColor:
      - red
`)
	writeThemeTestFile(t, filepath.Join(configDir, "themes", "pink.yml"), `gui:
  theme:
    activeBorderColor:
      - '#ff00ff'
    inactiveBorderColor:
      - '#00ff00'
  darkTheme:
    activeBorderColor:
      - '#00ffff'
  lightTheme:
    activeBorderColor:
      - '#0000ff'
`)
	selectThemeForTest(t, configDir, "pink")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	assert.NoError(t, appConfig.GetThemeLoadError())
	guiConfig := appConfig.GetUserConfig().Gui
	darkTheme := guiConfig.ThemeForBackground(false, "")
	assert.Equal(t, []string{"#00ffff"}, darkTheme.ActiveBorderColor, "the theme's darkTheme also overrides the one in config.yml")
	assert.Equal(t, []string{"#00ff00"}, darkTheme.InactiveBorderColor)
	lightTheme := guiConfig.ThemeForBackground(true, "")
	assert.Equal(t, []string{"#0000ff"}, lightTheme.ActiveBorderColor)
	assert.Equal(t, []string{"#00ff00"}, lightTheme.InactiveBorderColor)
}

func TestReloadUserConfigForRepoWithoutSelectedTheme(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	writeThemeTestFile(t, filepath.Join(configDir, "themes", "pink.yml"), `gui:
  theme:
    branchColorPatterns:
      master: '#ff00ff'
`)

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	assert.Equal(t, "", appConfig.GetSelectedTheme())
	assert.Equal(t, "", appConfig.GetAppliedTheme())
	assert.NoError(t, appConfig.GetThemeLoadError())
	assert.Empty(t, appConfig.GetUserConfig().Gui.Theme.BranchColorPatterns)
	assert.False(t, lo.SomeBy(appConfig.userConfigFiles, func(f *ConfigFile) bool { return f.isTheme }))
}

func TestReloadUserConfigForRepoFallsBackToNoThemeWhenThemeFileIsInvalid(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, `gui:
  theme:
    branchColorPatterns:
      main: red
`)
	themePath := filepath.Join(configDir, "themes", "broken.yml")
	writeThemeTestFile(t, themePath, `gui:
  nerdFontsVersion: "3"
  theme:
    branchColorPatterns:
      master: '#ff00ff'
`)
	selectThemeForTest(t, configDir, "broken")

	err := appConfig.ReloadUserConfigForRepo(nil)

	assert.NoError(t, err)
	assert.Equal(t, ColorPatterns{{Pattern: "main", Color: "red"}}, appConfig.GetUserConfig().Gui.Theme.BranchColorPatterns)
	assert.Equal(t, "broken", appConfig.GetSelectedTheme())
	assert.Equal(t, "", appConfig.GetAppliedTheme())
	assert.False(t, lo.SomeBy(appConfig.userConfigFiles, func(f *ConfigFile) bool { return f.isTheme }))

	themeLoadError := appConfig.GetThemeLoadError()
	var themeFileError *ThemeFileError
	if assert.ErrorAs(t, themeLoadError, &themeFileError) {
		assert.Equal(t, themePath, themeFileError.Path)
	}
	assert.ErrorContains(t, themeLoadError, "The theme file `"+themePath+"` couldn't be loaded.")
	assert.ErrorContains(t, themeLoadError, "field nerdFontsVersion not found")
	assert.ErrorContains(t, themeLoadError, allowedThemeKeysMessage)
	assert.NotContains(t, themeLoadError.Error(), convertedThemeMessage, "the theme needed no conversion")
	assert.Equal(t, themeLoadError, appConfig.GetThemeLoadError(), "reading the error doesn't clear it")
}

func TestThemeLoadErrorIsClearedOnceTheThemeLoads(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	themePath := filepath.Join(configDir, "themes", "pink.yml")
	writeThemeTestFile(t, themePath, `gui:
  nerdFontsVersion: "3"
`)
	selectThemeForTest(t, configDir, "pink")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))
	assert.Error(t, appConfig.GetThemeLoadError())

	// The next reload (e.g. a repo switch) reads the still broken file again
	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))
	assert.Error(t, appConfig.GetThemeLoadError())

	writeThemeTestFile(t, themePath, `gui:
  theme:
    branchColorPatterns:
      master: '#ff00ff'
`)
	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	assert.NoError(t, appConfig.GetThemeLoadError())
	assert.Equal(t, "pink", appConfig.GetAppliedTheme())
	assert.Equal(t, ColorPatterns{{Pattern: "master", Color: "#ff00ff"}}, appConfig.GetUserConfig().Gui.Theme.BranchColorPatterns)
}

func TestReloadUserConfigForRepoFallsBackWhenThemeFileIsUnreadable(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	// A folder can be stat'ed but not read as a file
	themePath := filepath.Join(configDir, "themes", "pink.yml")
	assert.NoError(t, os.MkdirAll(themePath, 0o755))
	selectThemeForTest(t, configDir, "pink")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	assert.Equal(t, "", appConfig.GetAppliedTheme())
	var themeFileError *ThemeFileError
	if assert.ErrorAs(t, appConfig.GetThemeLoadError(), &themeFileError) {
		assert.Equal(t, themePath, themeFileError.Path)
	}
}

func TestReloadUserConfigForRepoSurvivesUnreachableThemeFile(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, `gui:
  theme:
    activeBorderColor:
      - red
`)
	// With a file where the themes folder should be, stat'ing the theme file
	// fails with an error other than "doesn't exist" on Linux and macOS, while
	// Windows reports the theme file as missing. Either way lazygit starts
	// without the theme.
	writeThemeTestFile(t, filepath.Join(configDir, "themes"), "")
	selectThemeForTest(t, configDir, "pink")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))
	assert.Equal(t, []string{"red"}, appConfig.GetUserConfig().Gui.Theme.ActiveBorderColor)
	assert.Equal(t, "pink", appConfig.GetSelectedTheme())
	assert.Equal(t, "", appConfig.GetAppliedTheme())
}

func TestMissingThemeFileIsAppliedOnceItAppears(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, `gui:
  theme:
    activeBorderColor:
      - red
`)
	selectThemeForTest(t, configDir, "later")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))
	assert.NoError(t, appConfig.GetThemeLoadError())
	assert.Equal(t, "later", appConfig.GetSelectedTheme())
	assert.Equal(t, "", appConfig.GetAppliedTheme(), "the theme file is missing")
	assert.Equal(t, []string{"red"}, appConfig.GetUserConfig().Gui.Theme.ActiveBorderColor)

	writeThemeTestFile(t, filepath.Join(configDir, "themes", "later.yml"), `gui:
  theme:
    activeBorderColor:
      - '#ff00ff'
`)
	err, didChange := appConfig.ReloadChangedUserConfigFiles()

	assert.NoError(t, err)
	assert.True(t, didChange)
	assert.Equal(t, "later", appConfig.GetAppliedTheme())
	assert.Equal(t, []string{"#ff00ff"}, appConfig.GetUserConfig().Gui.Theme.ActiveBorderColor)
}

func TestAppliedThemeIsKeptWhenReloadingTheThemeFileFails(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	selectThemeForTest(t, configDir, "later")
	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))
	assert.Equal(t, "", appConfig.GetAppliedTheme(), "the theme file is missing")

	themePath := filepath.Join(configDir, "themes", "later.yml")
	rewriteThemeTestFile(t, themePath, `gui:
  nerdFontsVersion: "3"
  theme:
    branchColorPatterns:
      master: '#ff00ff'
`)
	err, didChange := appConfig.ReloadChangedUserConfigFiles()

	var themeFileError *ThemeFileError
	if assert.ErrorAs(t, err, &themeFileError) {
		assert.Equal(t, themePath, themeFileError.Path)
	}
	assert.False(t, didChange)
	// The config from before the reload stays in use, and it has no theme
	assert.Empty(t, appConfig.GetUserConfig().Gui.Theme.BranchColorPatterns)
	assert.Equal(t, "", appConfig.GetAppliedTheme())
	assert.NoError(t, appConfig.GetThemeLoadError(), "the reload on focus reports its error itself")
}

func TestAppliedThemeIsClearedWhenTheThemeFileIsDeleted(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	themePath := filepath.Join(configDir, "themes", "pink.yml")
	writeThemeTestFile(t, themePath, `gui:
  theme:
    branchColorPatterns:
      master: '#ff00ff'
`)
	selectThemeForTest(t, configDir, "pink")
	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))
	assert.Equal(t, "pink", appConfig.GetAppliedTheme())

	assert.NoError(t, os.Remove(themePath))
	err, didChange := appConfig.ReloadChangedUserConfigFiles()

	assert.NoError(t, err)
	assert.True(t, didChange)
	assert.Equal(t, "pink", appConfig.GetSelectedTheme(), "the selection is kept")
	assert.Equal(t, "", appConfig.GetAppliedTheme())
	assert.Empty(t, appConfig.GetUserConfig().Gui.Theme.BranchColorPatterns)
}

func TestReloadUserConfigForRepoIgnoresInvalidThemeName(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	// This is the file that "../x" would point to if it were used as a path
	writeThemeTestFile(t, filepath.Join(configDir, "x.yml"), `gui:
  theme:
    branchColorPatterns:
      master: '#ff00ff'
`)
	selectThemeForTest(t, configDir, "../x")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))
	assert.NoError(t, appConfig.GetThemeLoadError())
	assert.Equal(t, "../x", appConfig.GetSelectedTheme())
	assert.Equal(t, "", appConfig.GetAppliedTheme())
	assert.Empty(t, appConfig.GetUserConfig().Gui.Theme.BranchColorPatterns)
}

func TestReloadUserConfigForRepoMatchesThemeNameIgnoringCase(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	themePath := filepath.Join(configDir, "themes", "pink.yml")
	writeThemeTestFile(t, themePath, `gui:
  theme:
    branchColorPatterns:
      master: '#ff00ff'
`)
	// As if edited by hand
	selectThemeForTest(t, configDir, "Pink")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	assert.NoError(t, appConfig.GetThemeLoadError())
	assert.Equal(t, "pink", appConfig.GetSelectedTheme(), "the name is spelled as listed")
	assert.Equal(t, "pink", appConfig.GetAppliedTheme())
	assert.Equal(t, ColorPatterns{{Pattern: "master", Color: "#ff00ff"}}, appConfig.GetUserConfig().Gui.Theme.BranchColorPatterns)
	loadedThemeFile, found := lo.Find(appConfig.userConfigFiles, func(f *ConfigFile) bool { return f.isTheme })
	if assert.True(t, found) {
		assert.Equal(t, themePath, loadedThemeFile.Path)
	}
}

func TestReloadUserConfigForRepoPrefersThemeNameWithSameCase(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	themesDir := filepath.Join(configDir, "themes")
	writeThemeTestFile(t, filepath.Join(themesDir, "Pink.yml"), `gui:
  theme:
    branchColorPatterns:
      master: '#00ff00'
`)
	writeThemeTestFile(t, filepath.Join(themesDir, "pink.yml"), `gui:
  theme:
    branchColorPatterns:
      master: '#ff00ff'
`)
	if themes, _ := appConfig.ListThemes(); len(themes) != 2 {
		t.Skip("the file system ignores case, so both names are the same file")
	}
	selectThemeForTest(t, configDir, "pink")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	assert.Equal(t, "pink", appConfig.GetSelectedTheme())
	assert.Equal(t, ColorPatterns{{Pattern: "master", Color: "#ff00ff"}}, appConfig.GetUserConfig().Gui.Theme.BranchColorPatterns)
}

func TestReloadUserConfigForRepoReportsUnparsableThemeSelection(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	selectionPath := filepath.Join(configDir, selectedThemeFileName)
	writeThemeTestFile(t, selectionPath, "name: [unclosed\n")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))
	assert.Equal(t, "", appConfig.GetSelectedTheme())
	assert.Equal(t, "", appConfig.GetAppliedTheme())
	assert.ErrorContains(t, appConfig.GetThemeLoadError(), selectionPath)

	writeThemeTestFile(t, selectionPath, "")
	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))
	assert.NoError(t, appConfig.GetThemeLoadError(), "the next reload clears the error")
}

func TestThemeFileIsWatchedButNotOfferedForEditing(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	themePath := filepath.Join(configDir, "themes", "pink.yml")
	writeThemeTestFile(t, themePath, `gui:
  theme:
    activeBorderColor:
      - '#ff00ff'
`)
	selectThemeForTest(t, configDir, "pink")
	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	assert.Equal(t, []string{filepath.Join(configDir, ConfigFilename)}, appConfig.GetUserConfigPaths())

	rewriteThemeTestFile(t, themePath, `gui:
  theme:
    activeBorderColor:
      - '#00ff00'
`)
	err, didChange := appConfig.ReloadChangedUserConfigFiles()

	assert.NoError(t, err)
	assert.True(t, didChange)
	assert.Equal(t, []string{"#00ff00"}, appConfig.GetUserConfig().Gui.Theme.ActiveBorderColor)
}

func TestLoadingThemeFileDoesNotRewriteIt(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	themePath := filepath.Join(configDir, "themes", "palette.yml")
	// Migrations that walk the whole file can't handle aliases, so this also
	// fails if the theme goes through them
	themeContent := `# One accent color for two keys
gui:
  theme:
    activeBorderColor: &accent
      - '#ff00ff'
      - bold
    searchingActiveBorderColor: *accent
`
	writeThemeTestFile(t, themePath, themeContent)
	selectThemeForTest(t, configDir, "palette")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	assert.NoError(t, appConfig.GetThemeLoadError())
	assert.Equal(t, []string{"#ff00ff", "bold"}, appConfig.GetUserConfig().Gui.Theme.SearchingActiveBorderColor)
	actualContent, err := os.ReadFile(themePath)
	assert.NoError(t, err)
	assert.Equal(t, themeContent, string(actualContent))
}

// Themes written for older versions of lazygit have the author colors and the
// branch color patterns directly in gui.
func TestOldLayoutThemeIsConvertedWithoutRewritingIt(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, `gui:
  theme:
    authorColors:
      John: red
`)
	themePath := filepath.Join(configDir, "themes", "mocha.yml")
	themeContent := `gui:
  theme:
    activeBorderColor:
      - '#89b4fa'
      - bold
  authorColors:
    '*': '#b4befe'
  branchColorPatterns:
    '^feature/': green
`
	writeThemeTestFile(t, themePath, themeContent)
	selectThemeForTest(t, configDir, "mocha")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	assert.NoError(t, appConfig.GetThemeLoadError())
	assert.Equal(t, "mocha", appConfig.GetAppliedTheme())
	themeConfig := appConfig.GetUserConfig().Gui.Theme
	assert.Equal(t, []string{"#89b4fa", "bold"}, themeConfig.ActiveBorderColor)
	assert.Equal(t, map[string]string{"John": "red", "*": "#b4befe"}, themeConfig.AuthorColors)
	assert.Equal(t, ColorPatterns{{Pattern: "^feature/", Color: "green"}}, themeConfig.BranchColorPatterns)
	actualContent, err := os.ReadFile(themePath)
	assert.NoError(t, err)
	assert.Equal(t, themeContent, string(actualContent))
}

func TestOldLayoutThemeWithBranchColorsIsConverted(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	themePath := filepath.Join(configDir, "themes", "mocha.yml")
	themeContent := `gui:
  branchColors:
    feature: green
`
	writeThemeTestFile(t, themePath, themeContent)
	selectThemeForTest(t, configDir, "mocha")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	assert.NoError(t, appConfig.GetThemeLoadError())
	assert.Equal(t, "mocha", appConfig.GetAppliedTheme())
	assert.Equal(t,
		ColorPatterns{{Pattern: "^feature(/|$)", Color: "green"}},
		appConfig.GetUserConfig().Gui.Theme.BranchColorPatterns,
	)
	actualContent, err := os.ReadFile(themePath)
	assert.NoError(t, err)
	assert.Equal(t, themeContent, string(actualContent))
}

func TestLoadingOldLayoutThemeFileWithAliasesDoesNotRewriteIt(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	themePath := filepath.Join(configDir, "themes", "palette.yml")
	// Migrations that walk the whole file can't handle aliases, so this also
	// fails if the theme goes through them
	themeContent := `# One accent color for the border and the authors
gui:
  theme:
    activeBorderColor:
      - &accent '#ff00ff'
  authorColors:
    '*': *accent
`
	writeThemeTestFile(t, themePath, themeContent)
	selectThemeForTest(t, configDir, "palette")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	assert.NoError(t, appConfig.GetThemeLoadError())
	assert.Equal(t, map[string]string{"*": "#ff00ff"}, appConfig.GetUserConfig().Gui.Theme.AuthorColors)
	actualContent, err := os.ReadFile(themePath)
	assert.NoError(t, err)
	assert.Equal(t, themeContent, string(actualContent))
}

func TestOldLayoutThemeWithMergeKeysInBranchColorPatternsIsConverted(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	writeThemeTestFile(t, filepath.Join(configDir, "themes", "palette.yml"), `gui:
  authorColors: &palette
    '*': '#b4befe'
  branchColorPatterns:
    <<: *palette
    '^feature/': green
`)
	selectThemeForTest(t, configDir, "palette")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	assert.NoError(t, appConfig.GetThemeLoadError())
	assert.Equal(t,
		ColorPatterns{{Pattern: "*", Color: "#b4befe"}, {Pattern: "^feature/", Color: "green"}},
		appConfig.GetUserConfig().Gui.Theme.BranchColorPatterns,
	)
}

// Converting a theme moves its settings to the end of gui.theme, which can put
// an alias ahead of its anchor, so that the converted theme doesn't parse
func TestReloadUserConfigForRepoExplainsErrorInConvertedTheme(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	themePath := filepath.Join(configDir, "themes", "mocha.yml")
	writeThemeTestFile(t, themePath, "gui: {authorColors: {'*': &lav '#b4befe'}, theme: {activeBorderColor: [*lav, bold]}}\n")
	selectThemeForTest(t, configDir, "mocha")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	assert.Equal(t, "", appConfig.GetAppliedTheme())
	themeLoadError := appConfig.GetThemeLoadError()
	var themeFileError *ThemeFileError
	if assert.ErrorAs(t, themeLoadError, &themeFileError) {
		assert.Equal(t, themePath, themeFileError.Path)
	}
	assert.ErrorContains(t, themeLoadError, "unknown anchor 'lav' referenced")
	assert.ErrorContains(t, themeLoadError, convertedThemeMessage)
	assert.ErrorContains(t, themeLoadError, "the lines and keys above refer to the converted layout")
	assert.ErrorContains(t, themeLoadError, "Moving authorColors and branchColorPatterns into gui.theme yourself")
}

func TestReloadUserConfigForRepoFallsBackWhenThemeHasOldAndNewAuthorColors(t *testing.T) {
	appConfig, configDir := newThemeTestAppConfig(t, "")
	themePath := filepath.Join(configDir, "themes", "mocha.yml")
	writeThemeTestFile(t, themePath, `gui:
  authorColors:
    '*': '#b4befe'
  theme:
    authorColors:
      John: green
`)
	selectThemeForTest(t, configDir, "mocha")

	assert.NoError(t, appConfig.ReloadUserConfigForRepo(nil))

	assert.Equal(t, "", appConfig.GetAppliedTheme())
	assert.Empty(t, appConfig.GetUserConfig().Gui.Theme.AuthorColors)
	var themeFileError *ThemeFileError
	if assert.ErrorAs(t, appConfig.GetThemeLoadError(), &themeFileError) {
		assert.Equal(t, themePath, themeFileError.Path)
	}
	assert.ErrorContains(t, appConfig.GetThemeLoadError(),
		"Couldn't convert the theme to the current layout for key gui.authorColors: new key `authorColors' already exists")
	assert.NotContains(t, appConfig.GetThemeLoadError().Error(), "migrate", "theme files are never migrated")
}

// newThemeTestAppConfig creates an AppConfig whose config dir, which is also
// its state dir, is a fresh temp dir containing the given config.yml.
func newThemeTestAppConfig(t *testing.T, globalConfig string) (*AppConfig, string) {
	t.Helper()

	configDir := t.TempDir()
	t.Setenv("CONFIG_DIR", configDir)
	t.Setenv("LG_CONFIG_FILE", "")
	writeThemeTestFile(t, filepath.Join(configDir, ConfigFilename), globalConfig)

	appConfig, err := NewAppConfig("lazygit", "unversioned", "", "", "", false, t.TempDir())
	assert.NoError(t, err)
	return appConfig, configDir
}

func writeThemeTestFile(t *testing.T, path string, content string) {
	t.Helper()

	assert.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	assert.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

// rewriteThemeTestFile writes a file after an AppConfig has loaded the config,
// like an edit made while lazygit runs. It sets the modification time well
// past the one recorded when loading, so that the change is noticed even where
// timestamps are coarse.
func rewriteThemeTestFile(t *testing.T, path string, content string) {
	t.Helper()

	writeThemeTestFile(t, path, content)
	later := time.Now().Add(time.Hour)
	assert.NoError(t, os.Chtimes(path, later, later))
}

func selectThemeForTest(t *testing.T, configDir string, name string) {
	t.Helper()

	writeThemeTestFile(t, filepath.Join(configDir, selectedThemeFileName), "name: "+name+"\n")
}
