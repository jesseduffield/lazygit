package commands

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func newDummyBranchListBuilder() (*BranchListBuilder, *oscommands.OSCommand) {
	osCommand := oscommands.NewDummyOSCommand()

	return &BranchListBuilder{
		Log:           utils.NewDummyLog(),
		GitCommand:    NewDummyGitCommandWithOSCommand(osCommand),
		ReflogCommits: []*models.Commit{},
	}, osCommand
}

func TestBranchListBuilderLoadsMergedStatus(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func([]*models.Branch)
	}

	scenarios := []scenario{
		{
			"exactly the branches listed by branch --merged are marked as merged",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)

				switch args[0] {
				case "for-each-ref":
					assert.EqualValues(t, []string{"for-each-ref", "--sort=-committerdate", "--format=%(HEAD)|%(refname:short)|%(upstream:short)|%(upstream:track)", "refs/heads"}, args)
					return exec.Command("echo", "*|merged1||\n |unmerged1||\n |merged2||\n |unmerged2||")
				case "branch":
					assert.EqualValues(t, []string{"branch", "--format=%(refname:short)", "--merged"}, args)
					return exec.Command("echo", "merged1\nmerged2")
				}
				return nil
			},
			func(output []*models.Branch) {
				merged := []string{}
				unmerged := []string{}
				for _, b := range output {
					if b.Merged {
						merged = append(merged, b.Name)
					} else {
						unmerged = append(unmerged, b.Name)
					}
				}
				assert.EqualValues(t, []string{"merged1", "merged2"}, merged)
				assert.EqualValues(t, []string{"unmerged1", "unmerged2"}, unmerged)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			b, os := newDummyBranchListBuilder()
			os.SetCommand(s.command)
			s.test(b.Build())
		})
	}
}
