package helpers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type IRefsHelper interface {
	CheckoutRef(ref string, options types.CheckoutRefOptions) error
	GetCheckedOutRef() *models.Branch
	CreateGitResetMenu(ref string) error
	CreateCheckoutMenu(commit *models.Commit) error
	ResetToRef(ref string, strength string, envVars []string) error
	NewBranch(from string, fromDescription string, suggestedBranchname string) error
}

type RefsHelper struct {
	c *HelperCommon
}

func NewRefsHelper(
	c *HelperCommon,
) *RefsHelper {
	return &RefsHelper{
		c: c,
	}
}

var _ IRefsHelper = &RefsHelper{}

func (self *RefsHelper) CheckoutRef(ref string, options types.CheckoutRefOptions) error {
	waitingStatus := options.WaitingStatus
	if waitingStatus == "" {
		waitingStatus = self.c.Tr.CheckingOutStatus
	}

	cmdOptions := git_commands.CheckoutOptions{Force: false, EnvVars: options.EnvVars}

	refresh := func() {
		self.c.Contexts().Branches.SetSelection(0)
		self.c.Contexts().ReflogCommits.SetSelection(0)
		self.c.Contexts().LocalCommits.SetSelection(0)
		// loading a heap of commits is slow so we limit them whenever doing a reset
		self.c.Contexts().LocalCommits.SetLimitCommits(true)

		_ = self.c.Refresh(types.RefreshOptions{Mode: types.BLOCK_UI, KeepBranchSelectionIndex: true})
	}

	localBranch, found := lo.Find(self.c.Model().Branches, func(branch *models.Branch) bool {
		return branch.Name == ref
	})

	withCheckoutStatus := func(f func(gocui.Task) error) error {
		if found {
			return self.c.WithInlineStatus(localBranch, types.ItemOperationCheckingOut, context.LOCAL_BRANCHES_CONTEXT_KEY, f)
		} else {
			return self.c.WithWaitingStatus(waitingStatus, f)
		}
	}

	return withCheckoutStatus(func(gocui.Task) error {
		if err := self.c.Git().Branch.Checkout(ref, cmdOptions); err != nil {
			// note, this will only work for english-language git commands. If we force git to use english, and the error isn't this one, then the user will receive an english command they may not understand. I'm not sure what the best solution to this is. Running the command once in english and a second time in the native language is one option

			if options.OnRefNotFound != nil && strings.Contains(err.Error(), "did not match any file(s) known to git") {
				return options.OnRefNotFound(ref)
			}

			if IsSwitchBranchUncommittedChangesError(err) {
				// offer to autostash changes
				self.c.OnUIThread(func() error {
					// (Before showing the prompt, render again to remove the inline status)
					self.c.Contexts().Branches.HandleRender()
					self.c.Confirm(types.ConfirmOpts{
						Title:  self.c.Tr.AutoStashTitle,
						Prompt: self.c.Tr.AutoStashPrompt,
						HandleConfirm: func() error {
							return withCheckoutStatus(func(gocui.Task) error {
								if err := self.c.Git().Stash.Push(self.c.Tr.StashPrefix + ref); err != nil {
									return err
								}
								if err := self.c.Git().Branch.Checkout(ref, cmdOptions); err != nil {
									return err
								}
								err := self.c.Git().Stash.Pop(0)
								// Branch switch successful so re-render the UI even if the pop operation failed (e.g. conflict).
								refresh()
								return err
							})
						},
					})

					return nil
				})
				return nil
			}

			return err
		}

		refresh()
		return nil
	})
}

// Shows a prompt to choose between creating a new branch or checking out a detached head
func (self *RefsHelper) CheckoutRemoteBranch(fullBranchName string, localBranchName string) error {
	checkout := func(branchName string) error {
		// Switch to the branches context _before_ starting to check out the
		// branch, so that we see the inline status
		if self.c.Context().Current() != self.c.Contexts().Branches {
			self.c.Context().Push(self.c.Contexts().Branches)
		}
		return self.CheckoutRef(branchName, types.CheckoutRefOptions{})
	}

	// If a branch with this name already exists locally, just check it out. We
	// don't bother checking whether it actually tracks this remote branch, since
	// it's very unlikely that it doesn't.
	if lo.ContainsBy(self.c.Model().Branches, func(branch *models.Branch) bool {
		return branch.Name == localBranchName
	}) {
		return checkout(localBranchName)
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: utils.ResolvePlaceholderString(self.c.Tr.RemoteBranchCheckoutTitle, map[string]string{
			"branchName": fullBranchName,
		}),
		Prompt: self.c.Tr.RemoteBranchCheckoutPrompt,
		Items: []*types.MenuItem{
			{
				Label:   self.c.Tr.CheckoutTypeNewBranch,
				Tooltip: self.c.Tr.CheckoutTypeNewBranchTooltip,
				OnPress: func() error {
					// First create the local branch with the upstream set, and
					// then check it out. We could do that in one step using
					// "git checkout -b", but we want to benefit from all the
					// nice features of the CheckoutRef function.
					if err := self.c.Git().Branch.CreateWithUpstream(localBranchName, fullBranchName); err != nil {
						return err
					}
					// Do a sync refresh to make sure the new branch is visible,
					// so that we see an inline status when checking it out
					if err := self.c.Refresh(types.RefreshOptions{
						Mode:  types.SYNC,
						Scope: []types.RefreshableView{types.BRANCHES},
					}); err != nil {
						return err
					}
					return checkout(localBranchName)
				},
			},
			{
				Label:   self.c.Tr.CheckoutTypeDetachedHead,
				Tooltip: self.c.Tr.CheckoutTypeDetachedHeadTooltip,
				OnPress: func() error {
					return checkout(fullBranchName)
				},
			},
		},
	})
}

func (self *RefsHelper) GetCheckedOutRef() *models.Branch {
	if len(self.c.Model().Branches) == 0 {
		return nil
	}

	return self.c.Model().Branches[0]
}

func (self *RefsHelper) ResetToRef(ref string, strength string, envVars []string) error {
	if err := self.c.Git().Commit.ResetToCommit(ref, strength, envVars); err != nil {
		return err
	}

	self.c.Contexts().LocalCommits.SetSelection(0)
	self.c.Contexts().ReflogCommits.SetSelection(0)
	// loading a heap of commits is slow so we limit them whenever doing a reset
	self.c.Contexts().LocalCommits.SetLimitCommits(true)

	if err := self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES, types.BRANCHES, types.REFLOG, types.COMMITS}}); err != nil {
		return err
	}

	return nil
}

func (self *RefsHelper) CreateSortOrderMenu(sortOptionsOrder []string, onSelected func(sortOrder string) error, currentValue string) error {
	type sortMenuOption struct {
		key         types.Key
		label       string
		description string
		sortOrder   string
	}
	availableSortOptions := map[string]sortMenuOption{
		"recency":      {label: self.c.Tr.SortByRecency, description: self.c.Tr.SortBasedOnReflog, key: 'r'},
		"alphabetical": {label: self.c.Tr.SortAlphabetical, description: "--sort=refname", key: 'a'},
		"date":         {label: self.c.Tr.SortByDate, description: "--sort=-committerdate", key: 'd'},
	}
	sortOptions := make([]sortMenuOption, 0, len(sortOptionsOrder))
	for _, key := range sortOptionsOrder {
		sortOption, ok := availableSortOptions[key]
		if !ok {
			panic(fmt.Sprintf("unexpected sort order: %s", key))
		}
		sortOption.sortOrder = key
		sortOptions = append(sortOptions, sortOption)
	}

	menuItems := lo.Map(sortOptions, func(opt sortMenuOption, _ int) *types.MenuItem {
		return &types.MenuItem{
			LabelColumns: []string{
				opt.label,
				style.FgYellow.Sprint(opt.description),
			},
			OnPress: func() error {
				return onSelected(opt.sortOrder)
			},
			Key:    opt.key,
			Widget: types.MakeMenuRadioButton(opt.sortOrder == currentValue),
		}
	})
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.SortOrder,
		Items: menuItems,
	})
}

func (self *RefsHelper) CreateGitResetMenu(ref string) error {
	type strengthWithKey struct {
		strength string
		label    string
		key      types.Key
		tooltip  string
	}
	strengths := []strengthWithKey{
		// not i18'ing because it's git terminology
		{strength: "mixed", label: "Mixed reset", key: 'm', tooltip: self.c.Tr.ResetMixedTooltip},
		{strength: "soft", label: "Soft reset", key: 's', tooltip: self.c.Tr.ResetSoftTooltip},
		{strength: "hard", label: "Hard reset", key: 'h', tooltip: self.c.Tr.ResetHardTooltip},
	}

	menuItems := lo.Map(strengths, func(row strengthWithKey, _ int) *types.MenuItem {
		return &types.MenuItem{
			LabelColumns: []string{
				row.label,
				style.FgRed.Sprintf("reset --%s %s", row.strength, ref),
			},
			OnPress: func() error {
				self.c.LogAction("Reset")
				return self.ResetToRef(ref, row.strength, []string{})
			},
			Key:     row.key,
			Tooltip: row.tooltip,
		}
	})

	return self.c.Menu(types.CreateMenuOptions{
		Title: fmt.Sprintf("%s %s", self.c.Tr.ResetTo, ref),
		Items: menuItems,
	})
}

func (self *RefsHelper) CreateCheckoutMenu(commit *models.Commit) error {
	branches := lo.Filter(self.c.Model().Branches, func(branch *models.Branch, _ int) bool {
		return commit.Hash == branch.CommitHash && branch.Name != self.c.Model().CheckedOutBranch
	})

	hash := commit.Hash
	var menuItems []*types.MenuItem

	if len(branches) > 0 {
		menuItems = append(menuItems, lo.Map(branches, func(branch *models.Branch, index int) *types.MenuItem {
			var key types.Key
			if index < 9 {
				key = rune(index + 1 + '0') // Convert 1-based index to key
			}
			return &types.MenuItem{
				LabelColumns: []string{fmt.Sprintf(self.c.Tr.Actions.CheckoutBranchAtCommit, branch.Name)},
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.CheckoutBranch)
					return self.CheckoutRef(branch.RefName(), types.CheckoutRefOptions{})
				},
				Key: key,
			}
		})...)
	} else {
		menuItems = append(menuItems, &types.MenuItem{
			LabelColumns:   []string{self.c.Tr.Actions.CheckoutBranch},
			OnPress:        func() error { return nil },
			DisabledReason: &types.DisabledReason{Text: self.c.Tr.NoBranchesFoundAtCommitTooltip},
			Key:            '1',
		})
	}

	menuItems = append(menuItems, &types.MenuItem{
		LabelColumns: []string{fmt.Sprintf(self.c.Tr.Actions.CheckoutCommitAsDetachedHead, utils.ShortHash(hash))},
		OnPress: func() error {
			self.c.LogAction(self.c.Tr.Actions.CheckoutCommit)
			return self.CheckoutRef(hash, types.CheckoutRefOptions{})
		},
		Key: 'd',
	})

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Actions.CheckoutBranchOrCommit,
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

	if suggestedBranchName == "" {
		suggestedBranchName = self.c.UserConfig().Git.BranchPrefix
	}

	refresh := func() error {
		if self.c.Context().Current() != self.c.Contexts().Branches {
			self.c.Context().Push(self.c.Contexts().Branches)
		}

		self.c.Contexts().LocalCommits.SetSelection(0)
		self.c.Contexts().Branches.SetSelection(0)

		return self.c.Refresh(types.RefreshOptions{Mode: types.BLOCK_UI, KeepBranchSelectionIndex: true})
	}

	self.c.Prompt(types.PromptOpts{
		Title:          message,
		InitialContent: suggestedBranchName,
		HandleConfirm: func(response string) error {
			self.c.LogAction(self.c.Tr.Actions.CreateBranch)
			newBranchName := SanitizedBranchName(response)
			newBranchFunc := self.c.Git().Branch.New
			if newBranchName != suggestedBranchName {
				newBranchFunc = self.c.Git().Branch.NewWithoutTracking
			}
			if err := newBranchFunc(newBranchName, from); err != nil {
				if IsSwitchBranchUncommittedChangesError(err) {
					// offer to autostash changes
					self.c.Confirm(types.ConfirmOpts{
						Title:  self.c.Tr.AutoStashTitle,
						Prompt: self.c.Tr.AutoStashPrompt,
						HandleConfirm: func() error {
							if err := self.c.Git().Stash.Push(self.c.Tr.StashPrefix + newBranchName); err != nil {
								return err
							}
							if err := newBranchFunc(newBranchName, from); err != nil {
								return err
							}
							popErr := self.c.Git().Stash.Pop(0)
							// Branch switch successful so re-render the UI even if the pop operation failed (e.g. conflict).
							refreshError := refresh()
							if popErr != nil {
								// An error from pop is the more important one to report to the user
								return popErr
							}
							return refreshError
						},
					})

					return nil
				}

				return err
			}

			return refresh()
		},
	})

	return nil
}

// SanitizedBranchName will remove all spaces in favor of a dash "-" to meet
// git's branch naming requirement.
func SanitizedBranchName(input string) string {
	return strings.Replace(input, " ", "-", -1)
}

// Checks if the given branch name is a remote branch, and returns the name of
// the remote and the bare branch name if it is.
func (self *RefsHelper) ParseRemoteBranchName(fullBranchName string) (string, string, bool) {
	remoteName, branchName, found := strings.Cut(fullBranchName, "/")
	if !found {
		return "", "", false
	}

	// See if the part before the first slash is actually one of our remotes
	if !lo.ContainsBy(self.c.Model().Remotes, func(remote *models.Remote) bool {
		return remote.Name == remoteName
	}) {
		return "", "", false
	}

	return remoteName, branchName, true
}

func IsSwitchBranchUncommittedChangesError(err error) bool {
	return strings.Contains(err.Error(), "Please commit your changes or stash them before you switch branch")
}
