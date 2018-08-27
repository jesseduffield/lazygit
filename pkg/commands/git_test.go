package commands

import (
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func newDummyLog() *logrus.Entry {
	log := logrus.New()
	log.Out = ioutil.Discard
	return log.WithField("test", "test")
}

func newDummyGitCommand() *GitCommand {
	return &GitCommand{
		Log:       newDummyLog(),
		OSCommand: newDummyOSCommand(),
	}
}

func TestGitCommandGetStashEntries(t *testing.T) {
	type scenario struct {
		command func(string, ...string) *exec.Cmd
		test    func([]StashEntry)
	}

	scenarios := []scenario{
		{
			func(string, ...string) *exec.Cmd {
				return exec.Command("echo")
			},
			func(entries []StashEntry) {
				assert.Len(t, entries, 0)
			},
		},
		{
			func(string, ...string) *exec.Cmd {
				return exec.Command("echo", "WIP on add-pkg-commands-test: 55c6af2 increase parallel build\nWIP on master: bb86a3f update github template")
			},
			func(entries []StashEntry) {
				expected := []StashEntry{
					{
						0,
						"WIP on add-pkg-commands-test: 55c6af2 increase parallel build",
						"WIP on add-pkg-commands-test: 55c6af2 increase parallel build",
					},
					{
						1,
						"WIP on master: bb86a3f update github template",
						"WIP on master: bb86a3f update github template",
					},
				}

				assert.Len(t, entries, 2)
				assert.EqualValues(t, expected, entries)
			},
		},
	}

	for _, s := range scenarios {
		gitCmd := newDummyGitCommand()
		gitCmd.OSCommand.command = s.command

		s.test(gitCmd.GetStashEntries())
	}
}

func TestGetStashEntryDiff(t *testing.T) {
	gitCmd := newDummyGitCommand()
	gitCmd.OSCommand.command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"stash", "show", "-p", "--color", "stash@{1}"}, args)

		return exec.Command("echo")
	}

	gitCmd.GetStashEntryDiff(1)
}

func TestGitCommandDiff(t *testing.T) {
	gitCommand := newDummyGitCommand()
	assert.NoError(t, test.GenerateRepo("lots_of_diffs.sh"))

	files := []File{
		{
			Name:               "deleted_staged",
			HasStagedChanges:   false,
			HasUnstagedChanges: true,
			Tracked:            true,
			Deleted:            true,
			HasMergeConflicts:  false,
			DisplayString:      " D deleted_staged",
		},
		{
			Name:               "file with space staged",
			HasStagedChanges:   true,
			HasUnstagedChanges: false,
			Tracked:            false,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      "A  \"file with space staged\"",
		},
		{
			Name:               "file with space unstaged",
			HasStagedChanges:   false,
			HasUnstagedChanges: true,
			Tracked:            false,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      "?? file with space unstaged",
		},
		{
			Name:               "modified_unstaged",
			HasStagedChanges:   true,
			HasUnstagedChanges: false,
			Tracked:            true,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      "M  modified_unstaged",
		},
		{
			Name:               "modified_staged",
			HasStagedChanges:   false,
			HasUnstagedChanges: true,
			Tracked:            true,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      " M modified_staged",
		},
		{
			Name:               "renamed_before -> renamed_after",
			HasStagedChanges:   true,
			HasUnstagedChanges: false,
			Tracked:            true,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      "R  renamed_before -> renamed_after",
		},
		{
			Name:               "untracked_unstaged",
			HasStagedChanges:   false,
			HasUnstagedChanges: true,
			Tracked:            false,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      "?? untracked_unstaged",
		},
		{
			Name:               "untracked_staged",
			HasStagedChanges:   true,
			HasUnstagedChanges: false,
			Tracked:            false,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      "A  untracked_staged",
		},
		{
			Name:               "master",
			HasStagedChanges:   false,
			HasUnstagedChanges: true,
			Tracked:            false,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      "?? master",
		},
	}

	for _, file := range files {
		assert.NotContains(t, gitCommand.Diff(file), "error")
	}
}
