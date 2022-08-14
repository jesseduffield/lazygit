package helpers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type IRefsHelper interface {
	CheckoutRef(ref string, options types.CheckoutRefOptions) error
	GetCheckedOutRef() *models.Branch
	CreateGitResetMenu(ref string) error
	ResetToRef(ref string, strength string, envVars []string) error
	NewBranch(from string, fromDescription string, suggestedBranchname string) error
}

type RefsHelper struct {
	c        *types.HelperCommon
	git      *commands.GitCommand
	contexts *context.ContextTree
	model    *types.Model
}

func NewRefsHelper(
	c *types.HelperCommon,
	git *commands.GitCommand,
	contexts *context.ContextTree,
	model *types.Model,
) *RefsHelper {
	return &RefsHelper{
		c:        c,
		git:      git,
		contexts: contexts,
		model:    model,
	}
}

var _ IRefsHelper = &RefsHelper{}

func (self *RefsHelper) CheckoutRef(ref string, options types.CheckoutRefOptions) error {
	waitingStatus := options.WaitingStatus
	if waitingStatus == "" {
		waitingStatus = self.c.Tr.CheckingOutStatus
	}

	cmdOptions := git_commands.CheckoutOptions{Force: false, EnvVars: options.EnvVars}

	onSuccess := func() {
		self.contexts.Branches.SetSelectedLineIdx(0)
		self.contexts.ReflogCommits.SetSelectedLineIdx(0)
		self.contexts.LocalCommits.SetSelectedLineIdx(0)
		// loading a heap of commits is slow so we limit them whenever doing a reset
		self.contexts.LocalCommits.SetLimitCommits(true)
	}

	return self.c.WithWaitingStatus(waitingStatus, func() error {
		if err := self.git.Branch.Checkout(ref, cmdOptions); err != nil {
			// note, this will only work for english-language git commands. If we force git to use english, and the error isn't this one, then the user will receive an english command they may not understand. I'm not sure what the best solution to this is. Running the command once in english and a second time in the native language is one option

			if options.OnRefNotFound != nil && strings.Contains(err.Error(), "did not match any file(s) known to git") {
				return options.OnRefNotFound(ref)
			}

			if strings.Contains(err.Error(), "Please commit your changes or stash them before you switch branch") {
				// offer to autostash changes
				return self.c.Confirm(types.ConfirmOpts{
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

func (self *RefsHelper) GetCheckedOutRef() *models.Branch {
	if len(self.model.Branches) == 0 {
		return nil
	}

	return self.model.Branches[0]
}

func (self *RefsHelper) ResetToRef(ref string, strength string, envVars []string) error {
	if err := self.git.Commit.ResetToCommit(ref, strength, envVars); err != nil {
		return self.c.Error(err)
	}

	self.contexts.LocalCommits.SetSelectedLineIdx(0)
	self.contexts.ReflogCommits.SetSelectedLineIdx(0)
	// loading a heap of commits is slow so we limit them whenever doing a reset
	self.contexts.LocalCommits.SetLimitCommits(true)

	if err := self.c.PushContext(self.contexts.LocalCommits); err != nil {
		return err
	}

	if err := self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES, types.BRANCHES, types.REFLOG, types.COMMITS}}); err != nil {
		return err
	}

	return nil
}

func (self *RefsHelper) CreateGitResetMenu(ref string) error {
	type strengthWithKey struct {
		strength string
		key      types.Key
	}
	strengths := []strengthWithKey{
		{strength: "soft", key: 's'},
		{strength: "mixed", key: 'm'},
		{strength: "hard", key: 'h'},
	}

	menuItems := slices.Map(strengths, func(row strengthWithKey) *types.MenuItem {
		return &types.MenuItem{
			LabelColumns: []string{
				fmt.Sprintf("%s reset", row.strength),
				style.FgRed.Sprintf("reset --%s %s", row.strength, ref),
			},
			OnPress: func() error {
				self.c.LogAction("Reset")
				return self.ResetToRef(ref, row.strength, []string{})
			},
			Key: row.key,
		}
	})

	return self.c.Menu(types.CreateMenuOptions{
		Title: fmt.Sprintf("%s %s", self.c.Tr.LcResetTo, ref),
		Items: menuItems,
	})
}

func (self *RefsHelper) NewBranch(from string, fromFormattedName string, suggestedBranchName string) error {
	message := utils.ResolvePlaceholderString(
		self.c.Tr.NewBranchNameBranchOff,
		map[string]string{
			"branchName": fromFormattedName,
		},
	)

	return self.c.Prompt(types.PromptOpts{
		Title:          message,
		InitialContent: suggestedBranchName,
		HandleConfirm: func(response string) error {
			self.c.LogAction(self.c.Tr.Actions.CreateBranch)
			if err := self.git.Branch.New(sanitizedBranchName(response), from); err != nil {
				return err
			}

			if self.c.CurrentContext() != self.contexts.Branches {
				if err := self.c.PushContext(self.contexts.Branches); err != nil {
					return err
				}
			}

			self.contexts.LocalCommits.SetSelectedLineIdx(0)
			self.contexts.Branches.SetSelectedLineIdx(0)

			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		},
	})
}

// sanitizedBranchName will remove all spaces in favor of a dash "-" to meet
// git's branch naming requirement.
func sanitizedBranchName(input string) string {
	return strings.Replace(input, " ", "-", -1)
}
