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

	onSuccess := func() {
		self.c.Contexts().Branches.SetSelection(0)
		self.c.Contexts().ReflogCommits.SetSelection(0)
		self.c.Contexts().LocalCommits.SetSelection(0)
		// loading a heap of commits is slow so we limit them whenever doing a reset
		self.c.Contexts().LocalCommits.SetLimitCommits(true)
	}

	refreshOptions := types.RefreshOptions{Mode: types.BLOCK_UI, KeepBranchSelectionIndex: true}

	f := func(gocui.Task) error {
		if err := self.c.Git().Branch.Checkout(ref, cmdOptions); err != nil {
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
						if err := self.c.Git().Stash.Push(self.c.Tr.StashPrefix + ref); err != nil {
							return err
						}
						if err := self.c.Git().Branch.Checkout(ref, cmdOptions); err != nil {
							return err
						}

						onSuccess()
						if err := self.c.Git().Stash.Pop(0); err != nil {
							if err := self.c.Refresh(refreshOptions); err != nil {
								return err
							}
							return err
						}
						return self.c.Refresh(refreshOptions)
					},
				})
			}

			return err
		}
		onSuccess()

		return self.c.Refresh(refreshOptions)
	}

	localBranch, found := lo.Find(self.c.Model().Branches, func(branch *models.Branch) bool {
		return branch.Name == ref
	})
	if found {
		return self.c.WithInlineStatus(localBranch, types.ItemOperationCheckingOut, context.LOCAL_BRANCHES_CONTEXT_KEY, f)
	} else {
		return self.c.WithWaitingStatus(waitingStatus, f)
	}
}

// Shows a prompt to choose between creating a new branch or checking out a detached head
func (self *RefsHelper) CheckoutRemoteBranch(fullBranchName string, localBranchName string) error {
	checkout := func(branchName string) error {
		// Switch to the branches context _before_ starting to check out the
		// branch, so that we see the inline status
		if self.c.CurrentContext() != self.c.Contexts().Branches {
			if err := self.c.PushContext(self.c.Contexts().Branches); err != nil {
				return err
			}
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

func (self *RefsHelper) NewBranch(from string, fromFormattedName string, suggestedBranchName string) error {
	message := utils.ResolvePlaceholderString(
		self.c.Tr.NewBranchNameBranchOff,
		map[string]string{
			"branchName": fromFormattedName,
		},
	)

	if suggestedBranchName == "" {
		suggestedBranchName = self.c.UserConfig.Git.BranchPrefix
	}

	return self.c.Prompt(types.PromptOpts{
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
				return err
			}

			if self.c.CurrentContext() != self.c.Contexts().Branches {
				if err := self.c.PushContext(self.c.Contexts().Branches); err != nil {
					return err
				}
			}

			self.c.Contexts().LocalCommits.SetSelection(0)
			self.c.Contexts().Branches.SetSelection(0)

			return self.c.Refresh(types.RefreshOptions{Mode: types.BLOCK_UI, KeepBranchSelectionIndex: true})
		},
	})
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
