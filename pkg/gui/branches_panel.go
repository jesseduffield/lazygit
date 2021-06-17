package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedBranch() *models.Branch {
	if len(gui.State.Branches) == 0 {
		return nil
	}

	selectedLine := gui.State.Panels.Branches.SelectedLineIdx
	if selectedLine == -1 {
		return nil
	}

	return gui.State.Branches[selectedLine]
}

func (gui *Gui) handleBranchSelect() error {
	var task updateTask
	branch := gui.getSelectedBranch()
	if branch == nil {
		task = NewRenderStringTask(gui.Tr.NoBranchesThisRepo)
	} else {
		task = NewRunPtyTask(gui.Git.Branches().GetBranchGraphCmdObj(branch.Name))
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Log",
			task:  task,
		},
	})
}

// gui.refreshStatus is called at the end of this because that's when we can
// be sure there is a state.Branches array to pick the current branch from
func (gui *Gui) refreshBranches() {
	reflogCommits := gui.State.FilteredReflogCommits
	if gui.State.Modes.Filtering.Active() {
		// in filter mode we filter our reflog commits to just those containing the path
		// however we need all the reflog entries to populate the recencies of our branches
		// which allows us to order them correctly. So if we're filtering we'll just
		// manually load all the reflog commits here
		var err error
		reflogCommits, _, err = gui.Git.Reflog().Load(nil, "")
		if err != nil {
			gui.Log.Error(err)
		}
	}

	gui.State.Branches = gui.Git.Branches().GetBranches(reflogCommits)

	if err := gui.postRefreshUpdate(gui.State.Contexts.Branches); err != nil {
		gui.Log.Error(err)
	}

	gui.refreshStatus()
}

// specific functions

func (gui *Gui) handleBranchPress() error {
	if gui.State.Panels.Branches.SelectedLineIdx == -1 {
		return nil
	}
	if gui.State.Panels.Branches.SelectedLineIdx == 0 {
		return gui.CreateErrorPanel(gui.Tr.AlreadyCheckedOutBranch)
	}
	branch := gui.getSelectedBranch()
	return gui.handleCheckoutRef(branch.Name, handleCheckoutRefOptions{span: gui.Tr.Spans.CheckoutBranch})
}

func (gui *Gui) handleCreatePullRequestPress() error {
	pullRequest := commands.NewPullRequest(gui.Git)

	branch := gui.getSelectedBranch()
	url, err := pullRequest.Create(branch)
	if err != nil {
		return gui.SurfaceError(err)
	}
	gui.OnRunCommand(oscommands.NewCmdLogEntry(fmt.Sprintf("Creating pull request at URL: %s", url), "Create pull request", false))

	return nil
}

func (gui *Gui) handleCopyPullRequestURLPress() error {
	pullRequest := commands.NewPullRequest(gui.Git)

	branch := gui.getSelectedBranch()
	url, err := pullRequest.CopyURL(branch)
	if err != nil {
		return gui.SurfaceError(err)
	}
	gui.OnRunCommand(oscommands.NewCmdLogEntry(fmt.Sprintf("Copying to clipboard: '%s'", url), "Copy URL", false))

	gui.raiseToast(gui.Tr.PullRequestURLCopiedToClipboard)

	return nil
}

func (gui *Gui) handleGitFetch() error {
	return gui.WithPopupWaitingStatus(gui.Tr.FetchWait, func() error {
		err := gui.fetch()
		if err != nil {
			return gui.SurfaceError(err)
		}
		_ = gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC})

		return nil
	})
}

func (gui *Gui) handleForceCheckout() error {
	branch := gui.getSelectedBranch()
	message := gui.Tr.SureForceCheckout
	title := gui.Tr.ForceCheckoutBranch

	return gui.Ask(AskOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			if err := gui.Git.WithSpan(gui.Tr.Spans.ForceCheckoutBranch).Branches().Checkout(branch.Name, commands.CheckoutOpts{Force: true}); err != nil {
				_ = gui.SurfaceError(err)
			}
			return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC})
		},
	})
}

type handleCheckoutRefOptions struct {
	WaitingStatus string
	EnvVars       []string
	onRefNotFound func(ref string) error
	span          string
}

func (gui *Gui) handleCheckoutRef(ref string, options handleCheckoutRefOptions) error {
	waitingStatus := options.WaitingStatus
	if waitingStatus == "" {
		waitingStatus = gui.Tr.CheckingOutStatus
	}

	cmdOptions := commands.CheckoutOpts{Force: false, EnvVars: options.EnvVars}

	onSuccess := func() {
		gui.State.Panels.Branches.SelectedLineIdx = 0
		gui.State.Panels.Commits.SelectedLineIdx = 0
		// loading a heap of commits is slow so we limit them whenever doing a reset
		gui.State.Panels.Commits.LimitCommits = true
	}

	gitCommand := gui.Git.WithSpan(options.span)

	return gui.WithWaitingStatus(waitingStatus, func() error {
		if err := gitCommand.Branches().Checkout(ref, cmdOptions); err != nil {
			// note, this will only work for english-language git commands. If we force git to use english, and the error isn't this one, then the user will receive an english command they may not understand. I'm not sure what the best solution to this is. Running the command once in english and a second time in the native language is one option

			if options.onRefNotFound != nil && strings.Contains(err.Error(), "did not match any file(s) known to git") {
				return options.onRefNotFound(ref)
			}

			if strings.Contains(err.Error(), "Please commit your changes or stash them before you switch branch") {
				// offer to autostash changes
				return gui.Ask(AskOpts{

					Title:  gui.Tr.AutoStashTitle,
					Prompt: gui.Tr.AutoStashPrompt,
					HandleConfirm: func() error {
						if err := gitCommand.Stash().Save(gui.Tr.StashPrefix + ref); err != nil {
							return gui.SurfaceError(err)
						}
						if err := gitCommand.Branches().Checkout(ref, cmdOptions); err != nil {
							return gui.SurfaceError(err)
						}

						onSuccess()
						if err := gitCommand.Stash().Do(0, "pop"); err != nil {
							if err := gui.RefreshSidePanels(RefreshOptions{Mode: BLOCK_UI}); err != nil {
								return err
							}
							return gui.SurfaceError(err)
						}
						return gui.RefreshSidePanels(RefreshOptions{Mode: BLOCK_UI})
					},
				})
			}

			if err := gui.SurfaceError(err); err != nil {
				return err
			}
		}
		onSuccess()

		return gui.RefreshSidePanels(RefreshOptions{Mode: BLOCK_UI})
	})
}

func (gui *Gui) handleCheckoutByName() error {
	return gui.Prompt(PromptOpts{
		Title:               gui.Tr.BranchName + ":",
		FindSuggestionsFunc: gui.findBranchNameSuggestions,
		HandleConfirm: func(response string) error {
			return gui.handleCheckoutRef(response, handleCheckoutRefOptions{
				span: "Checkout branch",
				onRefNotFound: func(ref string) error {

					return gui.Ask(AskOpts{
						Title:  gui.Tr.BranchNotFoundTitle,
						Prompt: fmt.Sprintf("%s %s%s", gui.Tr.BranchNotFoundPrompt, ref, "?"),
						HandleConfirm: func() error {
							return gui.createNewBranchWithName(ref)
						},
					})
				},
			})
		}},
	)
}

func (gui *Gui) createNewBranchWithName(newBranchName string) error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}

	if err := gui.Git.Branches().NewBranch(newBranchName, branch.Name); err != nil {
		return gui.SurfaceError(err)
	}

	gui.State.Panels.Branches.SelectedLineIdx = 0
	return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC})
}

func (gui *Gui) handleDeleteBranch() error {
	return gui.deleteBranch(false)
}

func (gui *Gui) deleteBranch(force bool) error {
	selectedBranch := gui.getSelectedBranch()
	if selectedBranch == nil {
		return nil
	}
	checkedOutBranch := gui.CurrentBranch()
	if checkedOutBranch.Name == selectedBranch.Name {
		return gui.CreateErrorPanel(gui.Tr.CantDeleteCheckOutBranch)
	}
	return gui.deleteNamedBranch(selectedBranch, force)
}

func (gui *Gui) deleteNamedBranch(selectedBranch *models.Branch, force bool) error {
	title := gui.Tr.DeleteBranch
	var templateStr string
	if force {
		templateStr = gui.Tr.ForceDeleteBranchMessage
	} else {
		templateStr = gui.Tr.DeleteBranchMessage
	}
	message := utils.ResolvePlaceholderString(
		templateStr,
		map[string]string{
			"selectedBranchName": selectedBranch.Name,
		},
	)

	return gui.Ask(AskOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			if err := gui.Git.WithSpan(gui.Tr.Spans.DeleteBranch).Branches().Delete(selectedBranch.Name, force); err != nil {
				errMessage := err.Error()
				if !force && strings.Contains(errMessage, "is not fully merged") {
					return gui.deleteNamedBranch(selectedBranch, true)
				}
				return gui.CreateErrorPanel(errMessage)
			}
			return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC, Scope: []RefreshableView{BRANCHES}})
		},
	})
}

func (gui *Gui) mergeBranchIntoCheckedOutBranch(branchName string) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	if gui.Git.Status().IsHeadDetached() {
		return gui.CreateErrorPanel("Cannot merge branch in detached head state. You might have checked out a commit directly or a remote branch, in which case you should checkout the local branch you want to be on")
	}
	checkedOutBranchName := gui.CurrentBranch().Name
	if checkedOutBranchName == branchName {
		return gui.CreateErrorPanel(gui.Tr.CantMergeBranchIntoItself)
	}
	prompt := utils.ResolvePlaceholderString(
		gui.Tr.ConfirmMerge,
		map[string]string{
			"checkedOutBranch": checkedOutBranchName,
			"selectedBranch":   branchName,
		},
	)

	return gui.Ask(AskOpts{
		Title:  gui.Tr.MergingTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			err := gui.Git.WithSpan(gui.Tr.Spans.Merge).Branches().Merge(branchName, commands.MergeOpts{})
			return gui.handleGenericMergeCommandResult(err)
		},
	})
}

func (gui *Gui) handleMerge() error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	selectedBranchName := gui.getSelectedBranch().Name
	return gui.mergeBranchIntoCheckedOutBranch(selectedBranchName)
}

func (gui *Gui) handleRebaseOntoLocalBranch() error {
	selectedBranchName := gui.getSelectedBranch().Name
	return gui.handleRebaseOntoBranch(selectedBranchName)
}

func (gui *Gui) handleRebaseOntoBranch(selectedBranchName string) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	checkedOutBranch := gui.CurrentBranch().Name
	if selectedBranchName == checkedOutBranch {
		return gui.CreateErrorPanel(gui.Tr.CantRebaseOntoSelf)
	}
	prompt := utils.ResolvePlaceholderString(
		gui.Tr.ConfirmRebase,
		map[string]string{
			"checkedOutBranch": checkedOutBranch,
			"selectedBranch":   selectedBranchName,
		},
	)

	return gui.Ask(AskOpts{
		Title:  gui.Tr.RebasingTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			err := gui.Git.WithSpan(gui.Tr.Spans.RebaseBranch).Rebasing().RebaseBranch(selectedBranchName)
			return gui.handleGenericMergeCommandResult(err)
		},
	})
}

func (gui *Gui) handleFastForward() error {
	branch := gui.getSelectedBranch()
	if branch == nil || !branch.IsRealBranch() {
		return nil
	}

	if !branch.IsTrackingRemote() {
		return gui.CreateErrorPanel(gui.Tr.FwdNoUpstream)
	}
	if branch.HasCommitsToPush() {
		return gui.CreateErrorPanel(gui.Tr.FwdCommitsToPush)
	}

	upstream, err := gui.Git.Branches().GetUpstream(branch.Name)
	if err != nil {
		return gui.SurfaceError(err)
	}

	span := gui.Tr.Spans.FastForwardBranch

	split := strings.Split(upstream, "/")
	remoteName := split[0]
	remoteBranchName := strings.Join(split[1:], "/")

	message := utils.ResolvePlaceholderString(
		gui.Tr.Fetching,
		map[string]string{
			"from": fmt.Sprintf("%s/%s", remoteName, remoteBranchName),
			"to":   branch.Name,
		},
	)

	return gui.WithPopupWaitingStatus(message, func() error {
		if gui.State.Panels.Branches.SelectedLineIdx == 0 {
			_ = gui.pullWithMode("ff-only", PullFilesOptions{span: span})
		} else {
			err := gui.Git.WithSpan(span).Sync().FastForward(branch.Name, remoteName, remoteBranchName)
			if err != nil {
				return gui.SurfaceError(err)
			}
			_ = gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC, Scope: []RefreshableView{BRANCHES}})
		}

		return nil
	})
}

func (gui *Gui) handleCreateResetToBranchMenu() error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}

	return gui.createResetMenu(branch.Name)
}

func (gui *Gui) handleRenameBranch() error {
	branch := gui.getSelectedBranch()
	if branch == nil || !branch.IsRealBranch() {
		return nil
	}

	promptForNewName := func() error {
		return gui.Prompt(PromptOpts{
			Title:          gui.Tr.NewBranchNamePrompt + " " + branch.Name + ":",
			InitialContent: branch.Name,
			HandleConfirm: func(newBranchName string) error {
				if err := gui.Git.WithSpan(gui.Tr.Spans.RenameBranch).Branches().RenameBranch(branch.Name, newBranchName); err != nil {
					return gui.SurfaceError(err)
				}

				// need to find where the branch is now so that we can re-select it. That means we need to refetch the branches synchronously and then find our branch
				gui.refreshBranches()

				// now that we've got our stuff again we need to find that branch and reselect it.
				for i, newBranch := range gui.State.Branches {
					if newBranch.Name == newBranchName {
						gui.State.Panels.Branches.SetSelectedLineIdx(i)
						if err := gui.State.Contexts.Branches.HandleRender(); err != nil {
							return err
						}
					}
				}

				return nil
			},
		})
	}

	// I could do an explicit check here for whether the branch is tracking a remote branch
	// but if we've selected it we'll already know that via Pullables and Pullables.
	// Bit of a hack but I'm lazy.
	if !branch.IsTrackingRemote() {
		return promptForNewName()
	}

	return gui.Ask(AskOpts{
		Title:         gui.Tr.LcRenameBranch,
		Prompt:        gui.Tr.RenameBranchWarning,
		HandleConfirm: promptForNewName,
	})
}

func (gui *Gui) CurrentBranch() *models.Branch {
	if len(gui.State.Branches) == 0 {
		return nil
	}
	return gui.State.Branches[0]
}

func (gui *Gui) handleNewBranchOffCurrentItem() error {
	context := gui.currentSideListContext()

	item, ok := context.GetSelectedItem()
	if !ok {
		return nil
	}

	message := utils.ResolvePlaceholderString(
		gui.Tr.NewBranchNameBranchOff,
		map[string]string{
			"branchName": item.Description(),
		},
	)

	prefilledName := ""
	if context.GetKey() == REMOTE_BRANCHES_CONTEXT_KEY {
		// will set to the remote's branch name without the remote name
		prefilledName = strings.SplitAfterN(item.ID(), "/", 2)[1]
	}

	return gui.Prompt(PromptOpts{
		Title:          message,
		InitialContent: prefilledName,
		HandleConfirm: func(response string) error {
			if err := gui.Git.WithSpan(gui.Tr.Spans.CreateBranch).Branches().NewBranch(sanitizedBranchName(response), item.ID()); err != nil {
				return err
			}

			// if we're currently in the branch commits context then the selected commit
			// is about to go to the top of the list
			if context.GetKey() == BRANCH_COMMITS_CONTEXT_KEY {
				context.GetPanelState().SetSelectedLineIdx(0)
			}

			if context.GetKey() != gui.State.Contexts.Branches.GetKey() {
				if err := gui.pushContext(gui.State.Contexts.Branches); err != nil {
					return err
				}
			}

			gui.State.Panels.Branches.SelectedLineIdx = 0

			return gui.RefreshSidePanels(RefreshOptions{Mode: ASYNC})
		},
	})
}

func (gui *Gui) getBranchNames() []string {
	result := make([]string, len(gui.State.Branches))

	for i, branch := range gui.State.Branches {
		result[i] = branch.Name
	}

	return result
}

func (gui *Gui) findBranchNameSuggestions(input string) []*Suggestion {
	branchNames := gui.getBranchNames()

	matchingBranchNames := utils.FuzzySearch(sanitizedBranchName(input), branchNames)

	suggestions := make([]*Suggestion, len(matchingBranchNames))
	for i, branchName := range matchingBranchNames {
		suggestions[i] = &types.Suggestion{
			Value: branchName,
			Label: utils.ColoredString(branchName, presentation.GetBranchColor(branchName)),
		}
	}

	return suggestions
}

// sanitizedBranchName will remove all spaces in favor of a dash "-" to meet
// git's branch naming requirement.
func sanitizedBranchName(input string) string {
	return strings.Replace(input, " ", "-", -1)
}
