package helpers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
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

// hunk describes the lines in a diff hunk. Used for two distinct cases:
//
// - when the hunk contains some deleted lines. Because we're diffing with a
// context of 0, all deleted lines always come first, and then the added lines
// (if any). In this case, numLines is only the number of deleted lines, we
// ignore whether there are also some added lines in the hunk, as this is not
// relevant for our algorithm.
//
// - when the hunk contains only added lines, in which case (obviously) numLines
// is the number of added lines.
type hunk struct {
	filename     string
	startLineIdx int
	numLines     int
}

func (self *FixupHelper) HandleFindBaseCommitForFixupPress() error {
	diff, hasStagedChanges, err := self.getDiff()
	if err != nil {
		return err
	}

	deletedLineHunks, addedLineHunks := parseDiff(diff)

	commits := self.c.Model().Commits

	var hashes []string
	warnAboutAddedLines := false

	if len(deletedLineHunks) > 0 {
		hashes, err = self.blameDeletedLines(deletedLineHunks)
		warnAboutAddedLines = len(addedLineHunks) > 0
	} else if len(addedLineHunks) > 0 {
		hashes, err = self.blameAddedLines(commits, addedLineHunks)
	} else {
		return errors.New(self.c.Tr.NoChangedFiles)
	}

	if err != nil {
		return err
	}

	if len(hashes) == 0 {
		// This should never happen
		return errors.New(self.c.Tr.NoBaseCommitsFound)
	}

	// If a commit can't be found, and the last known commit is already merged,
	// we know that the commit we're looking for is also merged. Otherwise we
	// can't tell.
	notFoundMeansMerged := len(commits) > 0 && commits[len(commits)-1].Status == models.StatusMerged

	const (
		MERGED int = iota
		NOT_MERGED
		CANNOT_TELL
	)

	// Group the hashes into buckets by merged status
	hashGroups := lo.GroupBy(hashes, func(hash string) int {
		commit, _, ok := self.findCommit(commits, hash)
		if ok {
			return lo.Ternary(commit.Status == models.StatusMerged, MERGED, NOT_MERGED)
		}
		return lo.Ternary(notFoundMeansMerged, MERGED, CANNOT_TELL)
	})

	if len(hashGroups[CANNOT_TELL]) > 0 {
		// If we have any commits that we can't tell if they're merged, just
		// show the generic "not in current view" error. This can only happen if
		// a feature branch has more than 300 commits, or there is no main
		// branch. Both are so unlikely that we don't bother returning a more
		// detailed error message (e.g. we could say something about the commits
		// that *are* in the current branch, but it's not worth it).
		return errors.New(self.c.Tr.BaseCommitIsNotInCurrentView)
	}

	if len(hashGroups[NOT_MERGED]) == 0 {
		// If all the commits are merged, show the "already on main branch"
		// error. It isn't worth doing a detailed report of which commits we
		// found.
		return errors.New(self.c.Tr.BaseCommitIsAlreadyOnMainBranch)
	}

	if len(hashGroups[NOT_MERGED]) > 1 {
		// If there are multiple commits that could be the base commit, list
		// them in the error message. But only the candidates from the current
		// branch, not including any that are already merged.
		subjects, err := self.c.Git().Commit.GetHashesAndCommitMessagesFirstLine(hashGroups[NOT_MERGED])
		if err != nil {
			return err
		}
		message := lo.Ternary(hasStagedChanges,
			self.c.Tr.MultipleBaseCommitsFoundStaged,
			self.c.Tr.MultipleBaseCommitsFoundUnstaged)
		return fmt.Errorf("%s\n\n%s", message, subjects)
	}

	// At this point we know that the NOT_MERGED bucket has exactly one commit,
	// and that's the one we want to select.
	_, index, _ := self.findCommit(commits, hashGroups[NOT_MERGED][0])

	doIt := func() error {
		if !hasStagedChanges {
			if err := self.c.Git().WorkingTree.StageAll(); err != nil {
				return err
			}
			_ = self.c.Refresh(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.FILES}})
		}

		self.c.Contexts().LocalCommits.SetSelection(index)
		self.c.Context().Push(self.c.Contexts().LocalCommits)
		return nil
	}

	if warnAboutAddedLines {
		self.c.Confirm(types.ConfirmOpts{
			Title:  self.c.Tr.FindBaseCommitForFixup,
			Prompt: self.c.Tr.HunksWithOnlyAddedLinesWarning,
			HandleConfirm: func() error {
				return doIt()
			},
		})

		return nil
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

// Parse the diff output into hunks, and return two lists of hunks: the first
// are ones that contain deleted lines, the second are ones that contain only
// added lines.
func parseDiff(diff string) ([]*hunk, []*hunk) {
	lines := strings.Split(strings.TrimSuffix(diff, "\n"), "\n")

	deletedLineHunks := []*hunk{}
	addedLineHunks := []*hunk{}

	hunkHeaderRegexp := regexp.MustCompile(`@@ -(\d+)(?:,\d+)? \+\d+(?:,\d+)? @@`)

	var filename string
	var currentHunk *hunk
	numDeletedLines := 0
	numAddedLines := 0
	finishHunk := func() {
		if currentHunk != nil {
			if numDeletedLines > 0 {
				currentHunk.numLines = numDeletedLines
				deletedLineHunks = append(deletedLineHunks, currentHunk)
			} else if numAddedLines > 0 {
				currentHunk.numLines = numAddedLines
				addedLineHunks = append(addedLineHunks, currentHunk)
			}
		}
		numDeletedLines = 0
		numAddedLines = 0
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			finishHunk()
			currentHunk = nil
		} else if strings.HasPrefix(line, "--- ") {
			// For some reason, the line ends with a tab character if the file
			// name contains spaces
			filename = strings.TrimRight(line[6:], "\t")
		} else if strings.HasPrefix(line, "@@ ") {
			finishHunk()
			match := hunkHeaderRegexp.FindStringSubmatch(line)
			startIdx := utils.MustConvertToInt(match[1])
			currentHunk = &hunk{filename, startIdx, 0}
		} else if currentHunk != nil && line[0] == '-' {
			numDeletedLines++
		} else if currentHunk != nil && line[0] == '+' {
			numAddedLines++
		}
	}
	finishHunk()

	return deletedLineHunks, addedLineHunks
}

// returns the list of commit hashes that introduced the lines which have now been deleted
func (self *FixupHelper) blameDeletedLines(deletedLineHunks []*hunk) ([]string, error) {
	errg := errgroup.Group{}
	hashChan := make(chan string)

	for _, h := range deletedLineHunks {
		errg.Go(func() error {
			blameOutput, err := self.c.Git().Blame.BlameLineRange(h.filename, "HEAD", h.startLineIdx, h.numLines)
			if err != nil {
				return err
			}
			blameLines := strings.Split(strings.TrimSuffix(blameOutput, "\n"), "\n")
			for _, line := range blameLines {
				hashChan <- strings.Split(line, " ")[0]
			}
			return nil
		})
	}

	go func() {
		// We don't care about the error here, we'll check it later (in the
		// return statement below). Here we only wait for all the goroutines to
		// finish so that we can close the channel.
		_ = errg.Wait()
		close(hashChan)
	}()

	result := set.New[string]()
	for hash := range hashChan {
		result.Add(hash)
	}

	return result.ToSlice(), errg.Wait()
}

func (self *FixupHelper) blameAddedLines(commits []*models.Commit, addedLineHunks []*hunk) ([]string, error) {
	errg := errgroup.Group{}
	hashesChan := make(chan []string)

	for _, h := range addedLineHunks {
		errg.Go(func() error {
			result := make([]string, 0, 2)

			appendBlamedLine := func(blameOutput string) {
				blameLines := strings.Split(strings.TrimSuffix(blameOutput, "\n"), "\n")
				if len(blameLines) == 1 {
					result = append(result, strings.Split(blameLines[0], " ")[0])
				}
			}

			// Blame the line before this hunk, if there is one
			if h.startLineIdx > 0 {
				blameOutput, err := self.c.Git().Blame.BlameLineRange(h.filename, "HEAD", h.startLineIdx, 1)
				if err != nil {
					return err
				}
				appendBlamedLine(blameOutput)
			}

			// Blame the line after this hunk. We don't know how many lines the
			// file has, so we can't check if there is a line after the hunk;
			// let the error tell us.
			blameOutput, err := self.c.Git().Blame.BlameLineRange(h.filename, "HEAD", h.startLineIdx+1, 1)
			if err != nil {
				// If this fails, we're probably at the end of the file (we
				// could have checked this beforehand, but it's expensive). If
				// there was a line before this hunk, this is fine, we'll just
				// return that one; if not, the hunk encompasses the entire
				// file, and we can't blame the lines before and after the hunk.
				// This is an error.
				if h.startLineIdx == 0 {
					return errors.New("Entire file") // TODO i18n
				}
			} else {
				appendBlamedLine(blameOutput)
			}

			hashesChan <- result
			return nil
		})
	}

	go func() {
		// We don't care about the error here, we'll check it later (in the
		// return statement below). Here we only wait for all the goroutines to
		// finish so that we can close the channel.
		_ = errg.Wait()
		close(hashesChan)
	}()

	result := set.New[string]()
	for hashes := range hashesChan {
		if len(hashes) == 1 {
			result.Add(hashes[0])
		} else if len(hashes) > 1 {
			if hashes[0] == hashes[1] {
				result.Add(hashes[0])
			} else {
				_, index1, ok1 := self.findCommit(commits, hashes[0])
				_, index2, ok2 := self.findCommit(commits, hashes[1])
				if ok1 && ok2 {
					result.Add(lo.Ternary(index1 < index2, hashes[0], hashes[1]))
				} else if ok1 {
					result.Add(hashes[0])
				} else if ok2 {
					result.Add(hashes[1])
				} else {
					return nil, errors.New(self.c.Tr.NoBaseCommitsFound)
				}
			}
		}
	}

	return result.ToSlice(), errg.Wait()
}

func (self *FixupHelper) findCommit(commits []*models.Commit, hash string) (*models.Commit, int, bool) {
	return lo.FindIndexOf(commits, func(commit *models.Commit) bool {
		return commit.Hash == hash
	})
}
