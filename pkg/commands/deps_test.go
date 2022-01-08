package commands

import (
	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type commonDeps struct {
	runner     *oscommands.FakeCmdObjRunner
	userConfig *config.UserConfig
	gitConfig  *git_config.FakeGitConfig
	getenv     func(string) string
	removeFile func(string) error
	dotGitDir  string
	common     *common.Common
	cmd        *oscommands.CmdObjBuilder
}

func completeDeps(deps commonDeps) commonDeps {
	if deps.runner == nil {
		deps.runner = oscommands.NewFakeRunner(nil)
	}

	if deps.userConfig == nil {
		deps.userConfig = config.GetDefaultConfig()
	}

	if deps.gitConfig == nil {
		deps.gitConfig = git_config.NewFakeGitConfig(nil)
	}

	if deps.getenv == nil {
		deps.getenv = func(string) string { return "" }
	}

	if deps.removeFile == nil {
		deps.removeFile = func(string) error { return errors.New("unexpected call to removeFile") }
	}

	if deps.dotGitDir == "" {
		deps.dotGitDir = ".git"
	}

	if deps.common == nil {
		deps.common = utils.NewDummyCommonWithUserConfig(deps.userConfig)
	}

	if deps.cmd == nil {
		deps.cmd = oscommands.NewDummyCmdObjBuilder(deps.runner)
	}

	return deps
}

func buildConfigCommands(deps commonDeps) *ConfigCommands {
	deps = completeDeps(deps)
	common := utils.NewDummyCommonWithUserConfig(deps.userConfig)

	return NewConfigCommands(common, deps.gitConfig)
}

func buildOSCommand(deps commonDeps) *oscommands.OSCommand {
	deps = completeDeps(deps)

	return oscommands.NewDummyOSCommandWithDeps(oscommands.OSCommandDeps{
		Common:       deps.common,
		GetenvFn:     deps.getenv,
		Cmd:          deps.cmd,
		RemoveFileFn: deps.removeFile,
	})
}

func buildFileLoader(deps commonDeps) *loaders.FileLoader {
	deps = completeDeps(deps)

	configCommands := buildConfigCommands(deps)

	return loaders.NewFileLoader(deps.common, deps.cmd, configCommands)
}

func buildSubmoduleCommands(deps commonDeps) *SubmoduleCommands {
	deps = completeDeps(deps)

	return NewSubmoduleCommands(deps.common, deps.cmd, deps.dotGitDir)
}

func buildCommitCommands(deps commonDeps) *CommitCommands {
	deps = completeDeps(deps)
	return NewCommitCommands(deps.common, deps.cmd)
}

func buildWorkingTreeCommands(deps commonDeps) *WorkingTreeCommands {
	deps = completeDeps(deps)
	osCommand := buildOSCommand(deps)
	submoduleCommands := buildSubmoduleCommands(deps)
	fileLoader := buildFileLoader(deps)

	return NewWorkingTreeCommands(deps.common, deps.cmd, submoduleCommands, osCommand, fileLoader)
}

func buildStashCommands(deps commonDeps) *StashCommands {
	deps = completeDeps(deps)
	osCommand := buildOSCommand(deps)
	fileLoader := buildFileLoader(deps)
	workingTreeCommands := buildWorkingTreeCommands(deps)

	return NewStashCommands(deps.common, deps.cmd, osCommand, fileLoader, workingTreeCommands)
}

func buildRebaseCommands(deps commonDeps) *RebaseCommands {
	deps = completeDeps(deps)
	configCommands := buildConfigCommands(deps)
	osCommand := buildOSCommand(deps)
	workingTreeCommands := buildWorkingTreeCommands(deps)
	commitCommands := buildCommitCommands(deps)

	return NewRebaseCommands(deps.common, deps.cmd, osCommand, commitCommands, workingTreeCommands, configCommands, deps.dotGitDir)
}

func buildSyncCommands(deps commonDeps) *SyncCommands {
	deps = completeDeps(deps)

	return NewSyncCommands(deps.common, deps.cmd)
}

func buildFileCommands(deps commonDeps) *FileCommands {
	deps = completeDeps(deps)
	configCommands := buildConfigCommands(deps)
	osCommand := buildOSCommand(deps)

	return NewFileCommands(deps.common, deps.cmd, configCommands, osCommand)
}
