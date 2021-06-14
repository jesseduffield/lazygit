package commands

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/stretchr/testify/assert"
)

// TestGitCommandGetStashEntries is a function.
func TestGitCommandGetStashEntries(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func([]*models.StashEntry)
	}

	scenarios := []scenario{
		{
			"No stash entries found",
			func(string, ...string) *exec.Cmd {
				return secureexec.Command("echo")
			},
			func(entries []*models.StashEntry) {
				assert.Len(t, entries, 0)
			},
		},
		{
			"Several stash entries found",
			func(string, ...string) *exec.Cmd {
				return secureexec.Command("echo", "WIP on add-pkg-commands-test: 55c6af2 increase parallel build\nWIP on master: bb86a3f update github template")
			},
			func(entries []*models.StashEntry) {
				expected := []*models.StashEntry{
					{
						Index: 0,
						Name:  "WIP on add-pkg-commands-test: 55c6af2 increase parallel build",
					},
					{
						Index: 1,
						Name:  "WIP on master: bb86a3f update github template",
					},
				}

				assert.Len(t, entries, 2)
				assert.EqualValues(t, expected, entries)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGit()
			gitCmd.GetOSCommand().Command = s.command

			s.test(gitCmd.GetStashEntries(""))
		})
	}
}
