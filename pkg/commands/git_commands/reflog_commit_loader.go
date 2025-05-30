package git_commands

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type ReflogCommitLoader struct {
	*common.Common
	cmd oscommands.ICmdObjBuilder
}

func NewReflogCommitLoader(common *common.Common, cmd oscommands.ICmdObjBuilder) *ReflogCommitLoader {
	return &ReflogCommitLoader{
		Common: common,
		cmd:    cmd,
	}
}

// GetReflogCommits only returns the new reflog commits since the given lastReflogCommit
// if none is passed (i.e. it's value is nil) then we get all the reflog commits
func (self *ReflogCommitLoader) GetReflogCommits(hashPool *utils.StringPool, lastReflogCommit *models.Commit, filterPath string, filterAuthor string) ([]*models.Commit, bool, error) {
	cmdArgs := NewGitCmd("log").
		Config("log.showSignature=false").
		Arg("-g").
		Arg("--format=+%H%x00%ct%x00%gs%x00%P").
		ArgIf(filterAuthor != "", "--author="+filterAuthor).
		ArgIf(filterPath != "", "--follow", "--name-status", "--", filterPath).
		ToArgv()

	cmdObj := self.cmd.New(cmdArgs).DontLog()

	onlyObtainedNewReflogCommits := false

	commits, err := loadCommits(cmdObj, filterPath, func(line string) (*models.Commit, bool) {
		commit, ok := self.parseLine(hashPool, line)
		if !ok {
			return nil, false
		}

		// note that the unix timestamp here is the timestamp of the COMMIT, not the reflog entry itself,
		// so two consecutive reflog entries may have both the same hash and therefore same timestamp.
		// We use the reflog message to disambiguate, and fingers crossed that we never see the same of those
		// twice in a row. Reason being that it would mean we'd be erroneously exiting early.
		if lastReflogCommit != nil && self.sameReflogCommit(commit, lastReflogCommit) {
			onlyObtainedNewReflogCommits = true
			// after this point we already have these reflogs loaded so we'll simply return the new ones
			return nil, true
		}

		return commit, false
	})
	if err != nil {
		return nil, false, err
	}

	return commits, onlyObtainedNewReflogCommits, nil
}

func (self *ReflogCommitLoader) sameReflogCommit(a *models.Commit, b *models.Commit) bool {
	return a.Hash() == b.Hash() && a.UnixTimestamp == b.UnixTimestamp && a.Name == b.Name
}

func (self *ReflogCommitLoader) parseLine(hashPool *utils.StringPool, line string) (*models.Commit, bool) {
	fields := strings.SplitN(line, "\x00", 4)
	if len(fields) <= 3 {
		return nil, false
	}

	unixTimestamp, _ := strconv.Atoi(fields[1])

	parentHashes := fields[3]
	parents := []string{}
	if len(parentHashes) > 0 {
		parents = strings.Split(parentHashes, " ")
	}

	return models.NewCommit(hashPool, models.NewCommitOpts{
		Hash:          fields[0],
		Name:          fields[2],
		UnixTimestamp: int64(unixTimestamp),
		Status:        models.StatusReflog,
		Parents:       parents,
	}), true
}
