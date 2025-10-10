package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/jesseduffield/generics/orderedset"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/jesseduffield/lazygit/pkg/utils/yaml_utils"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

// AppConfig contains the base configuration fields required for lazygit.
type AppConfig struct {
	debug                 bool   `long:"debug" env:"DEBUG" default:"false"`
	version               string `long:"version" env:"VERSION" default:"unversioned"`
	buildDate             string `long:"build-date" env:"BUILD_DATE"`
	name                  string `long:"name" env:"NAME" default:"lazygit"`
	buildSource           string `long:"build-source" env:"BUILD_SOURCE" default:""`
	userConfig            *UserConfig
	globalUserConfigFiles []*ConfigFile
	userConfigFiles       []*ConfigFile
	userConfigDir         string
	tempDir               string
	appState              *AppState
}

type AppConfigurer interface {
	GetDebug() bool

	// build info
	GetVersion() string
	GetName() string
	GetBuildSource() string

	GetUserConfig() *UserConfig
	GetUserConfigPaths() []string
	GetUserConfigDir() string
	ReloadUserConfigForRepo(repoConfigFiles []*ConfigFile) error
	ReloadChangedUserConfigFiles() (error, bool)
	GetTempDir() string

	GetAppState() *AppState
	SaveAppState() error
}

type ConfigFilePolicy int

const (
	ConfigFilePolicyCreateIfMissing ConfigFilePolicy = iota
	ConfigFilePolicyErrorIfMissing
	ConfigFilePolicySkipIfMissing
)

type ConfigFile struct {
	Path    string
	Policy  ConfigFilePolicy
	modDate time.Time
	exists  bool
}

// NewAppConfig makes a new app config
func NewAppConfig(
	name string,
	version,
	commit,
	date string,
	buildSource string,
	debuggingFlag bool,
	tempDir string,
) (*AppConfig, error) {
	configDir, err := findOrCreateConfigDir()
	if err != nil && !os.IsPermission(err) {
		return nil, err
	}

	var configFiles []*ConfigFile
	customConfigFiles := os.Getenv("LG_CONFIG_FILE")
	if customConfigFiles != "" {
		// Load user defined config files
		userConfigPaths := strings.Split(customConfigFiles, ",")
		configFiles = lo.Map(userConfigPaths, func(path string, _ int) *ConfigFile {
			return &ConfigFile{Path: path, Policy: ConfigFilePolicyErrorIfMissing}
		})
	} else {
		// Load default config files
		path := filepath.Join(configDir, ConfigFilename)
		configFile := &ConfigFile{Path: path, Policy: ConfigFilePolicyCreateIfMissing}
		configFiles = []*ConfigFile{configFile}
	}

	userConfig, err := loadUserConfigWithDefaults(configFiles, false)
	if err != nil {
		return nil, err
	}

	appState, err := loadAppState()
	if err != nil {
		return nil, err
	}

	appConfig := &AppConfig{
		name:                  name,
		version:               version,
		buildDate:             date,
		debug:                 debuggingFlag,
		buildSource:           buildSource,
		userConfig:            userConfig,
		globalUserConfigFiles: configFiles,
		userConfigFiles:       configFiles,
		userConfigDir:         configDir,
		tempDir:               tempDir,
		appState:              appState,
	}

	return appConfig, nil
}

func ConfigDir() string {
	_, filePath := findConfigFile(ConfigFilename)

	return filepath.Dir(filePath)
}

func findOrCreateConfigDir() (string, error) {
	folder := ConfigDir()
	return folder, os.MkdirAll(folder, 0o755)
}

func loadUserConfigWithDefaults(configFiles []*ConfigFile, isGuiInitialized bool) (*UserConfig, error) {
	return loadUserConfig(configFiles, GetDefaultConfig(), isGuiInitialized)
}

func loadUserConfig(configFiles []*ConfigFile, base *UserConfig, isGuiInitialized bool) (*UserConfig, error) {
	for _, configFile := range configFiles {
		path := configFile.Path
		statInfo, err := os.Stat(path)
		if err == nil {
			configFile.exists = true
			configFile.modDate = statInfo.ModTime()
		} else {
			if !os.IsNotExist(err) {
				return nil, err
			}

			switch configFile.Policy {
			case ConfigFilePolicyErrorIfMissing:
				return nil, err

			case ConfigFilePolicySkipIfMissing:
				configFile.exists = false
				continue

			case ConfigFilePolicyCreateIfMissing:
				file, err := os.Create(path)
				if err != nil {
					if os.IsPermission(err) {
						// apparently when people have read-only permissions they prefer us to fail silently
						continue
					}
					return nil, err
				}
				file.Close()

				configFile.exists = true
				statInfo, err := os.Stat(configFile.Path)
				if err != nil {
					return nil, err
				}
				configFile.modDate = statInfo.ModTime()
			}
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		content, err = migrateUserConfig(path, content, isGuiInitialized)
		if err != nil {
			return nil, err
		}

		existingCustomCommands := base.CustomCommands

		if err := yaml.Unmarshal(content, base); err != nil {
			return nil, fmt.Errorf("The config at `%s` couldn't be parsed, please inspect it before opening up an issue.\n%w", path, err)
		}

		base.CustomCommands = append(base.CustomCommands, existingCustomCommands...)

		if err := base.Validate(); err != nil {
			return nil, fmt.Errorf("The config at `%s` has a validation error.\n%w", path, err)
		}
	}

	return base, nil
}

type ChangesSet = orderedset.OrderedSet[string]

func NewChangesSet() *ChangesSet {
	return orderedset.New[string]()
}

// Do any backward-compatibility migrations of things that have changed in the
// config over time; examples are renaming a key to a better name, moving a key
// from one container to another, or changing the type of a key (e.g. from bool
// to an enum).
func migrateUserConfig(path string, content []byte, isGuiInitialized bool) ([]byte, error) {
	changes := NewChangesSet()

	changedContent, didChange, err := computeMigratedConfig(path, content, changes)
	if err != nil {
		return nil, err
	}

	// Nothing to do if config didn't change
	if !didChange {
		return content, nil
	}

	changesText := "The following changes were made:\n\n"
	changesText += strings.Join(lo.Map(changes.ToSliceFromOldest(), func(change string, _ int) string {
		return fmt.Sprintf("- %s\n", change)
	}), "")

	// Write config back
	if !isGuiInitialized {
		fmt.Printf("The user config file %s must be migrated. Attempting to do this automatically.\n", path)
		fmt.Println(changesText)
	}
	if err := os.WriteFile(path, changedContent, 0o644); err != nil {
		errorMsg := fmt.Sprintf("While attempting to write back migrated user config to %s, an error occurred: %s", path, err)
		if isGuiInitialized {
			errorMsg += "\n\n" + changesText
		}
		return nil, errors.New(errorMsg)
	}
	if !isGuiInitialized {
		fmt.Printf("Config file saved successfully to %s\n", path)
	}
	return changedContent, nil
}

// A pure function helper for testing purposes
func computeMigratedConfig(path string, content []byte, changes *ChangesSet) ([]byte, bool, error) {
	var err error
	var rootNode yaml.Node
	err = yaml.Unmarshal(content, &rootNode)
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse YAML: %w", err)
	}
	var originalCopy yaml.Node
	err = yaml.Unmarshal(content, &originalCopy)
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse YAML, but only the second time!?!? How did that happen: %w", err)
	}

	pathsToReplace := []struct {
		oldPath []string
		newName string
	}{
		{[]string{"gui", "skipUnstageLineWarning"}, "skipDiscardChangeWarning"},
		{[]string{"keybinding", "universal", "executeCustomCommand"}, "executeShellCommand"},
		{[]string{"gui", "windowSize"}, "screenMode"},
		{[]string{"keybinding", "files", "openMergeTool"}, "openMergeOptions"},
	}

	for _, pathToReplace := range pathsToReplace {
		err, didReplace := yaml_utils.RenameYamlKey(&rootNode, pathToReplace.oldPath, pathToReplace.newName)
		if err != nil {
			return nil, false, fmt.Errorf("Couldn't migrate config file at `%s` for key %s: %w", path, strings.Join(pathToReplace.oldPath, "."), err)
		}
		if didReplace {
			changes.Add(fmt.Sprintf("Renamed '%s' to '%s'", strings.Join(pathToReplace.oldPath, "."), pathToReplace.newName))
		}
	}

	err = changeNullKeybindingsToDisabled(&rootNode, changes)
	if err != nil {
		return nil, false, fmt.Errorf("Couldn't migrate config file at `%s`: %w", path, err)
	}

	err = changeElementToSequence(&rootNode, []string{"git", "commitPrefix"}, changes)
	if err != nil {
		return nil, false, fmt.Errorf("Couldn't migrate config file at `%s`: %w", path, err)
	}

	err = changeCommitPrefixesMap(&rootNode, changes)
	if err != nil {
		return nil, false, fmt.Errorf("Couldn't migrate config file at `%s`: %w", path, err)
	}

	err = changeCustomCommandStreamAndOutputToOutputEnum(&rootNode, changes)
	if err != nil {
		return nil, false, fmt.Errorf("Couldn't migrate config file at `%s`: %w", path, err)
	}

	err = migrateAllBranchesLogCmd(&rootNode, changes)
	if err != nil {
		return nil, false, fmt.Errorf("Couldn't migrate config file at `%s`: %w", path, err)
	}

	err = migratePagers(&rootNode, changes)
	if err != nil {
		return nil, false, fmt.Errorf("Couldn't migrate config file at `%s`: %w", path, err)
	}

	// Add more migrations here...

	if reflect.DeepEqual(rootNode, originalCopy) {
		return nil, false, nil
	}

	newContent, err := yaml_utils.YamlMarshal(&rootNode)
	if err != nil {
		return nil, false, fmt.Errorf("Failed to remarsal!\n %w", err)
	}
	return newContent, true, nil
}

func changeNullKeybindingsToDisabled(rootNode *yaml.Node, changes *ChangesSet) error {
	return yaml_utils.Walk(rootNode, func(node *yaml.Node, path string) {
		if strings.HasPrefix(path, "keybinding.") && node.Kind == yaml.ScalarNode && node.Tag == "!!null" {
			node.Value = "<disabled>"
			node.Tag = "!!str"
			changes.Add(fmt.Sprintf("Changed 'null' to '<disabled>' for keybinding '%s'", path))
		}
	})
}

func changeElementToSequence(rootNode *yaml.Node, path []string, changes *ChangesSet) error {
	return yaml_utils.TransformNode(rootNode, path, func(node *yaml.Node) error {
		if node.Kind == yaml.MappingNode {
			nodeContentCopy := node.Content
			node.Kind = yaml.SequenceNode
			node.Value = ""
			node.Tag = "!!seq"
			node.Content = []*yaml.Node{{
				Kind:    yaml.MappingNode,
				Content: nodeContentCopy,
			}}

			changes.Add(fmt.Sprintf("Changed '%s' to an array of strings", strings.Join(path, ".")))

			return nil
		}
		return nil
	})
}

func changeCommitPrefixesMap(rootNode *yaml.Node, changes *ChangesSet) error {
	return yaml_utils.TransformNode(rootNode, []string{"git", "commitPrefixes"}, func(prefixesNode *yaml.Node) error {
		if prefixesNode.Kind == yaml.MappingNode {
			for _, contentNode := range prefixesNode.Content {
				if contentNode.Kind == yaml.MappingNode {
					nodeContentCopy := contentNode.Content
					contentNode.Kind = yaml.SequenceNode
					contentNode.Value = ""
					contentNode.Tag = "!!seq"
					contentNode.Content = []*yaml.Node{{
						Kind:    yaml.MappingNode,
						Content: nodeContentCopy,
					}}
					changes.Add("Changed 'git.commitPrefixes' elements to arrays of strings")
				}
			}
		}
		return nil
	})
}

func changeCustomCommandStreamAndOutputToOutputEnum(rootNode *yaml.Node, changes *ChangesSet) error {
	return yaml_utils.Walk(rootNode, func(node *yaml.Node, path string) {
		// We are being lazy here and rely on the fact that the only mapping
		// nodes in the tree under customCommands are actual custom commands. If
		// this ever changes (e.g. because we add a struct field to
		// customCommand), then we need to change this to iterate properly.
		if strings.HasPrefix(path, "customCommands[") && node.Kind == yaml.MappingNode {
			output := ""
			if streamKey, streamValue := yaml_utils.RemoveKey(node, "subprocess"); streamKey != nil {
				if streamValue.Kind == yaml.ScalarNode && streamValue.Value == "true" {
					output = "terminal"
					changes.Add("Changed 'subprocess: true' to 'output: terminal' in custom command")
				} else {
					changes.Add("Deleted redundant 'subprocess: false' in custom command")
				}
			}
			if streamKey, streamValue := yaml_utils.RemoveKey(node, "stream"); streamKey != nil {
				if streamValue.Kind == yaml.ScalarNode && streamValue.Value == "true" && output == "" {
					output = "log"
					changes.Add("Changed 'stream: true' to 'output: log' in custom command")
				} else {
					changes.Add(fmt.Sprintf("Deleted redundant 'stream: %v' property in custom command", streamValue.Value))
				}
			}
			if streamKey, streamValue := yaml_utils.RemoveKey(node, "showOutput"); streamKey != nil {
				if streamValue.Kind == yaml.ScalarNode && streamValue.Value == "true" && output == "" {
					changes.Add("Changed 'showOutput: true' to 'output: popup' in custom command")
					output = "popup"
				} else {
					changes.Add(fmt.Sprintf("Deleted redundant 'showOutput: %v' property in custom command", streamValue.Value))
				}
			}
			if output != "" {
				outputKeyNode := &yaml.Node{
					Kind:  yaml.ScalarNode,
					Value: "output",
					Tag:   "!!str",
				}
				outputValueNode := &yaml.Node{
					Kind:  yaml.ScalarNode,
					Value: output,
					Tag:   "!!str",
				}
				node.Content = append(node.Content, outputKeyNode, outputValueNode)
			}
		}
	})
}

// This migration is special because users have already defined
// a single element at `allBranchesLogCmd` and the sequence at `allBranchesLogCmds`.
// Some users have explicitly set `allBranchesLogCmd` to be an empty string in order
// to remove it, so in that case we just delete the element, and add nothing to the list
func migrateAllBranchesLogCmd(rootNode *yaml.Node, changes *ChangesSet) error {
	return yaml_utils.TransformNode(rootNode, []string{"git"}, func(gitNode *yaml.Node) error {
		cmdKeyNode, cmdValueNode := yaml_utils.LookupKey(gitNode, "allBranchesLogCmd")
		// Nothing to do if they do not have the deprecated item
		if cmdKeyNode == nil {
			return nil
		}

		cmdsKeyNode, cmdsValueNode := yaml_utils.LookupKey(gitNode, "allBranchesLogCmds")
		var change string
		if cmdsKeyNode == nil {
			// Create empty sequence node and attach it onto the root git node
			// We will later populate it with the individual allBranchesLogCmd record
			cmdsKeyNode = &yaml.Node{Kind: yaml.ScalarNode, Value: "allBranchesLogCmds"}
			cmdsValueNode = &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{}}
			gitNode.Content = append(gitNode.Content,
				cmdsKeyNode,
				cmdsValueNode,
			)
			change = "Created git.allBranchesLogCmds array containing value of git.allBranchesLogCmd"
		} else {
			if cmdsValueNode.Kind != yaml.SequenceNode {
				return errors.New("You should have an allBranchesLogCmds defined as a sequence!")
			}

			change = "Prepended git.allBranchesLogCmd value to git.allBranchesLogCmds array"
		}

		if cmdValueNode.Value != "" {
			// Prepending the individual element to make it show up first in the list, which was prior behavior
			cmdsValueNode.Content = utils.Prepend(cmdsValueNode.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: cmdValueNode.Value})
			changes.Add(change)
		}

		// Clear out the existing allBranchesLogCmd, now that we have migrated it into the list
		_, _ = yaml_utils.RemoveKey(gitNode, "allBranchesLogCmd")
		changes.Add("Removed obsolete git.allBranchesLogCmd")

		return nil
	})
}

func migratePagers(rootNode *yaml.Node, changes *ChangesSet) error {
	return yaml_utils.TransformNode(rootNode, []string{"git"}, func(gitNode *yaml.Node) error {
		pagingKeyNode, pagingValueNode := yaml_utils.LookupKey(gitNode, "paging")
		if pagingKeyNode == nil || pagingValueNode.Kind != yaml.MappingNode {
			// If there's no "paging" section (or it's not an object), there's nothing to do
			return nil
		}

		pagersKeyNode, _ := yaml_utils.LookupKey(gitNode, "pagers")
		if pagersKeyNode != nil {
			// Conversely, if there *is* already a "pagers" array, we also have nothing to do.
			// This covers the case where the user keeps both the "paging" section and the "pagers"
			// array for the sake of easier testing of old versions.
			return nil
		}

		pagingKeyNode.Value = "pagers"
		pagingContentCopy := pagingValueNode.Content
		pagingValueNode.Kind = yaml.SequenceNode
		pagingValueNode.Tag = "!!seq"
		pagingValueNode.Content = []*yaml.Node{{
			Kind:    yaml.MappingNode,
			Content: pagingContentCopy,
		}}

		changes.Add("Moved git.paging object to git.pagers array")

		return nil
	})
}

func (c *AppConfig) GetDebug() bool {
	return c.debug
}

func (c *AppConfig) GetVersion() string {
	return c.version
}

func (c *AppConfig) GetName() string {
	return c.name
}

// GetBuildSource returns the source of the build. For builds from goreleaser
// this will be binaryBuild
func (c *AppConfig) GetBuildSource() string {
	return c.buildSource
}

// GetUserConfig returns the user config
func (c *AppConfig) GetUserConfig() *UserConfig {
	return c.userConfig
}

// GetAppState returns the app state
func (c *AppConfig) GetAppState() *AppState {
	return c.appState
}

func (c *AppConfig) GetUserConfigPaths() []string {
	return lo.FilterMap(c.userConfigFiles, func(f *ConfigFile, _ int) (string, bool) {
		return f.Path, f.exists
	})
}

func (c *AppConfig) GetUserConfigDir() string {
	return c.userConfigDir
}

func (c *AppConfig) ReloadUserConfigForRepo(repoConfigFiles []*ConfigFile) error {
	configFiles := append(c.globalUserConfigFiles, repoConfigFiles...)
	userConfig, err := loadUserConfigWithDefaults(configFiles, true)
	if err != nil {
		return err
	}

	c.userConfig = userConfig
	c.userConfigFiles = configFiles
	return nil
}

func (c *AppConfig) ReloadChangedUserConfigFiles() (error, bool) {
	fileHasChanged := func(f *ConfigFile) bool {
		info, err := os.Stat(f.Path)
		if err != nil && !os.IsNotExist(err) {
			// If we can't stat the file, assume it hasn't changed
			return false
		}
		exists := err == nil
		return exists != f.exists || (exists && info.ModTime() != f.modDate)
	}

	if lo.NoneBy(c.userConfigFiles, fileHasChanged) {
		return nil, false
	}

	userConfig, err := loadUserConfigWithDefaults(c.userConfigFiles, true)
	if err != nil {
		return err, false
	}

	c.userConfig = userConfig
	return nil, true
}

func (c *AppConfig) GetTempDir() string {
	return c.tempDir
}

// findConfigFile looks for a possibly existing config file.
// This function does NOT create any folders or files.
func findConfigFile(filename string) (exists bool, path string) {
	if envConfigDir := os.Getenv("CONFIG_DIR"); envConfigDir != "" {
		return true, filepath.Join(envConfigDir, filename)
	}

	// look for jesseduffield/lazygit/filename in XDG_CONFIG_HOME and XDG_CONFIG_DIRS
	legacyConfigPath, err := xdg.SearchConfigFile(filepath.Join("jesseduffield", "lazygit", filename))
	if err == nil {
		return true, legacyConfigPath
	}

	// look for lazygit/filename in XDG_CONFIG_HOME and XDG_CONFIG_DIRS
	configFilepath, err := xdg.SearchConfigFile(filepath.Join("lazygit", filename))
	if err == nil {
		return true, configFilepath
	}

	return false, filepath.Join(xdg.ConfigHome, "lazygit", filename)
}

var ConfigFilename = "config.yml"

// stateFilePath looks for a possibly existing state file.
// if none exist, the default path is returned and all parent directories are created.
func stateFilePath(filename string) (string, error) {
	exists, legacyStateFile := findConfigFile(filename)
	if exists {
		return legacyStateFile, nil
	}

	// looks for XDG_STATE_HOME/lazygit/filename
	return xdg.StateFile(filepath.Join("lazygit", filename))
}

// SaveAppState marshalls the AppState struct and writes it to the disk
func (c *AppConfig) SaveAppState() error {
	marshalledAppState, err := yaml.Marshal(c.appState)
	if err != nil {
		return err
	}

	filepath, err := stateFilePath(stateFileName)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, marshalledAppState, 0o644)
	if err != nil && os.IsPermission(err) {
		// apparently when people have read-only permissions they prefer us to fail silently
		return nil
	}

	return err
}

var stateFileName = "state.yml"

// loadAppState loads recorded AppState from file
func loadAppState() (*AppState, error) {
	appState := getDefaultAppState()

	filepath, err := stateFilePath(stateFileName)
	if err != nil {
		if os.IsPermission(err) {
			// apparently when people have read-only permissions they prefer us to fail silently
			return appState, nil
		}
		return nil, err
	}

	appStateBytes, err := os.ReadFile(filepath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if len(appStateBytes) == 0 {
		return appState, nil
	}

	err = yaml.Unmarshal(appStateBytes, appState)
	if err != nil {
		return nil, err
	}

	return appState, nil
}

// SaveGlobalUserConfig saves the UserConfig back to disk. This is only used in
// integration tests, so we are a bit sloppy with error handling.
func (c *AppConfig) SaveGlobalUserConfig() {
	if len(c.globalUserConfigFiles) != 1 {
		panic("expected exactly one global user config file")
	}

	yamlContent, err := yaml.Marshal(c.userConfig)
	if err != nil {
		log.Fatalf("error marshalling user config: %v", err)
	}

	err = os.WriteFile(c.globalUserConfigFiles[0].Path, yamlContent, 0o644)
	if err != nil {
		log.Fatalf("error saving user config: %v", err)
	}
}

// AppState stores data between runs of the app like when the last update check
// was performed and which other repos have been checked out
type AppState struct {
	LastUpdateCheck        int64
	RecentRepos            []string
	StartupPopupVersion    int
	DidShowHunkStagingHint bool
	LastVersion            string // this is the last version the user was using, for the purpose of showing release notes

	// these are for shell commands typed in directly, not for custom commands in the lazygit config.
	// For backwards compatibility we keep the old name in yaml files.
	ShellCommandsHistory []string `yaml:"customcommandshistory"`

	HideCommandLog bool
}

func getDefaultAppState() *AppState {
	return &AppState{}
}

func LogPath() (string, error) {
	if os.Getenv("LAZYGIT_LOG_PATH") != "" {
		return os.Getenv("LAZYGIT_LOG_PATH"), nil
	}

	return stateFilePath("development.log")
}
