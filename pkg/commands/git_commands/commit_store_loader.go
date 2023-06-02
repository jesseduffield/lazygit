package git_commands

import (
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

// CommitStoreLoader populates a commit store with commits from the git log.
type CommitStoreLoader struct {
	*common.Common
	cmd oscommands.ICmdObjBuilder
}

func NewCommitStoreLoader(
	cmn *common.Common,
	cmd oscommands.ICmdObjBuilder,
) *CommitStoreLoader {
	return &CommitStoreLoader{
		Common: cmn,
		cmd:    cmd,
	}
}

// mutates the given commit store to add commits from the git log
func (self *CommitStoreLoader) Load(commitStore *models.CommitStore) error {
	t := time.Now()

	err := self.getLogCmd().RunAndProcessLines(func(line string) (bool, error) {
		commit := self.extractCommitFromLine(line)
		commitStore.Add(commit)
		return false, nil
	})

	self.Log.Warnf("CommitStoreLoader Load took %s", time.Since(t))

	return err
}

// getLog gets the git log.
func (self *CommitStoreLoader) getLogCmd() oscommands.ICmdObj {
	cmdArgs := NewGitCmd("log").
		Arg("--all").
		Arg(`--pretty=format:%H%x00%P`).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog()
}

func (self *CommitStoreLoader) extractCommitFromLine(line string) models.ImmutableCommit {
	split := strings.SplitN(line, "\x00", 2)

	sha := split[0]

	parentsStr := split[1]
	parents := []string{}
	if len(parentsStr) > 0 {
		parents = strings.Split(parentsStr, " ")
	}

	return models.NewImmutableCommit(sha, parents)
}
