package commands

import (
	"io/ioutil"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

// this file defines constructors for loaders, passing in all the dependencies required based on a smaller set of arguments passed in by the client.

func NewCommitLoader(
	cmn *common.Common,
	gitCommand *GitCommand,
	osCommand *oscommands.OSCommand,
) *loaders.CommitLoader {
	return loaders.NewCommitLoader(
		cmn,
		gitCommand.Cmd,
		gitCommand.CurrentBranchName,
		gitCommand.RebaseMode,
		ioutil.ReadFile,
		filepath.Walk,
		gitCommand.DotGitDir,
	)
}
