package helpers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type FixupHelper struct {
	c *HelperCommon
}

func NewFixupHelper(
	c *HelperCommon,
) *FixupHelper {
	return &FixupHelper{
		c: c,
	}
}

type deletedLineInfo struct {
	filename     string
	startLineIdx int
	numLines     int
}

func (self *FixupHelper) HandleFindBaseCommitForFixupPress() error {
	diff, hasStagedChanges, err := self.getDiff()
	if err != nil {
		return err
	}
	if diff == "" {
		return errors.New(self.c.Tr.NoChangedFiles)
	}

	deletedLineInfos, hasHunksWithOnlyAddedLines := self.parseDiff(diff)
	if len(deletedLineInfos) == 0 {
		return errors.New(self.c.Tr.NoDeletedLinesInDiff)
	}

	hashes := self.blameDeletedLines(deletedLineInfos)

	if len(hashes) == 0 {
		// This should never happen
		return errors.New(self.c.Tr.NoBaseCommitsFound)
	}
	if len(hashes) > 1 {
		subjects, err := self.c.Git().Commit.GetHashesAndCommitMessagesFirstLine(hashes)
		if err != nil {
			return err
		}
		message := lo.Ternary(hasStagedChanges,
			self.c.Tr.MultipleBaseCommitsFoundStaged,
			self.c.Tr.MultipleBaseCommitsFoundUnstaged)
		return fmt.Errorf("%s\n\n%s", message, subjects)
	}

	commit, index, ok := lo.FindIndexOf(self.c.Model().Commits, func(commit *models.Commit) bool {
		return commit.Hash == hashes[0]
	})
	if !ok {
		commits := self.c.Model().Commits
		if commits[len(commits)-1].Status == models.StatusMerged {
			// If the commit is not found, it's most likely because it's already
			// merged, and more than 300 commits away. Check if the last known
			// commit is already merged; if so, show the "already merged" error.
			return errors.New(self.c.Tr.BaseCommitIsAlreadyOnMainBranch)
		}
		// If we get here, the current branch must have more then 300 commits. Unlikely...
		return errors.New(self.c.Tr.BaseCommitIsNotInCurrentView)
	}
	if commit.Status == models.StatusMerged {
		return errors.New(self.c.Tr.BaseCommitIsAlreadyOnMainBranch)
	}

	doIt := func() error {
		if !hasStagedChanges {
			if err := self.c.Git().WorkingTree.StageAll(); err != nil {
				return err
			}
			_ = self.c.Refresh(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.FILES}})
		}

		self.c.Contexts().LocalCommits.SetSelection(index)
		return self.c.PushContext(self.c.Contexts().LocalCommits)
	}

	if hasHunksWithOnlyAddedLines {
		return self.c.Confirm(types.ConfirmOpts{
			Title:  self.c.Tr.FindBaseCommitForFixup,
			Prompt: self.c.Tr.HunksWithOnlyAddedLinesWarning,
			HandleConfirm: func() error {
				return doIt()
			},
		})
	}

	return doIt()
}

func (self *FixupHelper) getDiff() (string, bool, error) {
	args := []string{"-U0", "--ignore-submodules=all", "HEAD", "--"}

	// Try staged changes first
	hasStagedChanges := true
	diff, err := self.c.Git().Diff.DiffIndexCmdObj(append([]string{"--cached"}, args...)...).RunWithOutput()

	if err == nil && diff == "" {
		hasStagedChanges = false
		// If there are no staged changes, try unstaged changes
		diff, err = self.c.Git().Diff.DiffIndexCmdObj(args...).RunWithOutput()
	}

	return diff, hasStagedChanges, err
}

func (self *FixupHelper) parseDiff(diff string) ([]*deletedLineInfo, bool) {
	lines := strings.Split(strings.TrimSuffix(diff, "\n"), "\n")

	deletedLineInfos := []*deletedLineInfo{}
	hasHunksWithOnlyAddedLines := false

	hunkHeaderRegexp := regexp.MustCompile(`@@ -(\d+)(?:,\d+)? \+\d+(?:,\d+)? @@`)

	var filename string
	var currentLineInfo *deletedLineInfo
	finishHunk := func() {
		if currentLineInfo != nil {
			if currentLineInfo.numLines > 0 {
				deletedLineInfos = append(deletedLineInfos, currentLineInfo)
			} else {
				hasHunksWithOnlyAddedLines = true
			}
		}
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			finishHunk()
			currentLineInfo = nil
		} else if strings.HasPrefix(line, "--- ") {
			// For some reason, the line ends with a tab character if the file
			// name contains spaces
			filename = strings.TrimRight(line[6:], "\t")
		} else if strings.HasPrefix(line, "@@ ") {
			finishHunk()
			match := hunkHeaderRegexp.FindStringSubmatch(line)
			startIdx := utils.MustConvertToInt(match[1])
			currentLineInfo = &deletedLineInfo{filename, startIdx, 0}
		} else if currentLineInfo != nil && line[0] == '-' {
			currentLineInfo.numLines++
		}
	}
	finishHunk()

	return deletedLineInfos, hasHunksWithOnlyAddedLines
}

// returns the list of commit hashes that introduced the lines which have now been deleted
func (self *FixupHelper) blameDeletedLines(deletedLineInfos []*deletedLineInfo) []string {
	var wg sync.WaitGroup
	hashChan := make(chan string)

	for _, info := range deletedLineInfos {
		wg.Add(1)
		go func(info *deletedLineInfo) {
			defer wg.Done()

			blameOutput, err := self.c.Git().Blame.BlameLineRange(info.filename, "HEAD", info.startLineIdx, info.numLines)
			if err != nil {
				self.c.Log.Errorf("Error blaming file '%s': %v", info.filename, err)
				return
			}
			blameLines := strings.Split(strings.TrimSuffix(blameOutput, "\n"), "\n")
			for _, line := range blameLines {
				hashChan <- strings.Split(line, " ")[0]
			}
		}(info)
	}

	go func() {
		wg.Wait()
		close(hashChan)
	}()

	result := set.New[string]()
	for hash := range hashChan {
		result.Add(hash)
	}

	return result.ToSlice()
}
