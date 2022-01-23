package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type RefHelper struct {
	c   *controllers.ControllerCommon
	git *commands.GitCommand

	getState func() *GuiRepoState
}

func NewRefHelper(
	c *controllers.ControllerCommon,
	git *commands.GitCommand,
	getState func() *GuiRepoState,
) *RefHelper {
	return &RefHelper{
		c:        c,
		git:      git,
		getState: getState,
	}
}

var _ controllers.IRefHelper = &RefHelper{}

func (self *RefHelper) CheckoutRef(ref string, options types.CheckoutRefOptions) error {
	waitingStatus := options.WaitingStatus
	if waitingStatus == "" {
		waitingStatus = self.c.Tr.CheckingOutStatus
	}

	cmdOptions := git_commands.CheckoutOptions{Force: false, EnvVars: options.EnvVars}

	onSuccess := func() {
		self.getState().Panels.Branches.SelectedLineIdx = 0
		self.getState().Panels.Commits.SelectedLineIdx = 0
		// loading a heap of commits is slow so we limit them whenever doing a reset
		self.getState().Panels.Commits.LimitCommits = true
	}

	return self.c.WithWaitingStatus(waitingStatus, func() error {
		if err := self.git.Branch.Checkout(ref, cmdOptions); err != nil {
			// note, this will only work for english-language git commands. If we force git to use english, and the error isn't this one, then the user will receive an english command they may not understand. I'm not sure what the best solution to this is. Running the command once in english and a second time in the native language is one option

			if options.OnRefNotFound != nil && strings.Contains(err.Error(), "did not match any file(s) known to git") {
				return options.OnRefNotFound(ref)
			}

			if strings.Contains(err.Error(), "Please commit your changes or stash them before you switch branch") {
				// offer to autostash changes
				return self.c.Ask(popup.AskOpts{

					Title:  self.c.Tr.AutoStashTitle,
					Prompt: self.c.Tr.AutoStashPrompt,
					HandleConfirm: func() error {
						if err := self.git.Stash.Save(self.c.Tr.StashPrefix + ref); err != nil {
							return self.c.Error(err)
						}
						if err := self.git.Branch.Checkout(ref, cmdOptions); err != nil {
							return self.c.Error(err)
						}

						onSuccess()
						if err := self.git.Stash.Pop(0); err != nil {
							if err := self.c.Refresh(types.RefreshOptions{Mode: types.BLOCK_UI}); err != nil {
								return err
							}
							return self.c.Error(err)
						}
						return self.c.Refresh(types.RefreshOptions{Mode: types.BLOCK_UI})
					},
				})
			}

			if err := self.c.Error(err); err != nil {
				return err
			}
		}
		onSuccess()

		return self.c.Refresh(types.RefreshOptions{Mode: types.BLOCK_UI})
	})
}

func (self *RefHelper) ResetToRef(ref string, strength string, envVars []string) error {
	if err := self.git.Commit.ResetToCommit(ref, strength, envVars); err != nil {
		return self.c.Error(err)
	}

	self.getState().Panels.Commits.SelectedLineIdx = 0
	self.getState().Panels.ReflogCommits.SelectedLineIdx = 0
	// loading a heap of commits is slow so we limit them whenever doing a reset
	self.getState().Panels.Commits.LimitCommits = true

	if err := self.c.PushContext(self.getState().Contexts.BranchCommits); err != nil {
		return err
	}

	if err := self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES, types.BRANCHES, types.REFLOG, types.COMMITS}}); err != nil {
		return err
	}

	return nil
}

func (self *RefHelper) CreateGitResetMenu(ref string) error {
	strengths := []string{"soft", "mixed", "hard"}
	menuItems := make([]*popup.MenuItem, len(strengths))
	for i, strength := range strengths {
		strength := strength
		menuItems[i] = &popup.MenuItem{
			DisplayStrings: []string{
				fmt.Sprintf("%s reset", strength),
				style.FgRed.Sprintf("reset --%s %s", strength, ref),
			},
			OnPress: func() error {
				self.c.LogAction("Reset")
				return self.ResetToRef(ref, strength, []string{})
			},
		}
	}

	return self.c.Menu(popup.CreateMenuOptions{
		Title: fmt.Sprintf("%s %s", self.c.Tr.LcResetTo, ref),
		Items: menuItems,
	})
}
