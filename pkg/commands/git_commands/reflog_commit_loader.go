package git_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
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
func (self *ReflogCommitLoader) GetReflogCommits(lastReflogCommit *models.Commit, filterPath string) ([]*models.Commit, bool, error) {
	commits := make([]*models.Commit, 0)

	filterPathArg := ""
	if filterPath != "" {
		filterPathArg = fmt.Sprintf(" --follow -- %s", self.cmd.Quote(filterPath))
	}

	cmdObj := self.cmd.New(fmt.Sprintf(`git -c log.showSignature=false log -g --abbrev=40 --format="%s"%s`, "%h%x00%ct%x00%gs%x00%p", filterPathArg)).DontLog()
	onlyObtainedNewReflogCommits := false
	err := cmdObj.RunAndProcessLines(func(line string) (bool, error) {
		fields := strings.SplitN(line, "\x00", 4)
		if len(fields) <= 3 {
			return false, nil
		}

		unixTimestamp, _ := strconv.Atoi(fields[1])

		parentHashes := fields[3]
		parents := []string{}
		if len(parentHashes) > 0 {
			parents = strings.Split(parentHashes, " ")
		}

		commit := &models.Commit{
			Sha:           fields[0],
			Name:          fields[2],
			UnixTimestamp: int64(unixTimestamp),
			Status:        "reflog",
			Parents:       parents,
		}

		// note that the unix timestamp here is the timestamp of the COMMIT, not the reflog entry itself,
		// so two consecutive reflog entries may have both the same SHA and therefore same timestamp.
		// We use the reflog message to disambiguate, and fingers crossed that we never see the same of those
		// twice in a row. Reason being that it would mean we'd be erroneously exiting early.
		if lastReflogCommit != nil && commit.Sha == lastReflogCommit.Sha && commit.UnixTimestamp == lastReflogCommit.UnixTimestamp && commit.Name == lastReflogCommit.Name {
			onlyObtainedNewReflogCommits = true
			// after this point we already have these reflogs loaded so we'll simply return the new ones
			return true, nil
		}

		commits = append(commits, commit)
		return false, nil
	})
	if err != nil {
		return nil, false, err
	}

	return commits, onlyObtainedNewReflogCommits, nil
}
