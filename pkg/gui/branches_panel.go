package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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
		task = gui.createRenderStringTask(gui.Tr.NoBranchesThisRepo)
	} else {
		cmd := gui.OSCommand.ExecutableFromString(
			gui.GitCommand.GetBranchGraphCmdStr(branch.Name),
		)

		task = gui.createRunPtyTask(cmd)
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
		reflogCommits, _, err = gui.GitCommand.GetReflogCommits(nil, "")
		if err != nil {
			gui.Log.Error(err)
		}
	}

	builder, err := commands.NewBranchListBuilder(gui.Log, gui.GitCommand, reflogCommits)
	if err != nil {
		_ = gui.surfaceError(err)
	}
	gui.State.Branches = builder.Build()

	if err := gui.postRefreshUpdate(gui.Contexts.Branches.Context); err != nil {
		gui.Log.Error(err)
	}

	gui.refreshStatus()
}

// specific functions

func (gui *Gui) handleBranchPress(g *gocui.Gui, v *gocui.View) error {
	if gui.State.Panels.Branches.SelectedLineIdx == -1 {
		return nil
	}
	if gui.State.Panels.Branches.SelectedLineIdx == 0 {
		return gui.createErrorPanel(gui.Tr.AlreadyCheckedOutBranch)
	}
	branch := gui.getSelectedBranch()
	return gui.handleCheckoutRef(branch.Name, handleCheckoutRefOptions{})
}

func (gui *Gui) handleCreatePullRequestPress(g *gocui.Gui, v *gocui.View) error {
	pullRequest := commands.NewPullRequest(gui.GitCommand)

	branch := gui.getSelectedBranch()
	if err := pullRequest.Create(branch); err != nil {
		return gui.surfaceError(err)
	}

	return nil
}

func (gui *Gui) handleCopyPullRequestURLPress(g *gocui.Gui, v *gocui.View) error {
	pullRequest := commands.NewPullRequest(gui.GitCommand)

	branch := gui.getSelectedBranch()
	if err := pullRequest.CopyURL(branch); err != nil {
		return gui.surfaceError(err)
	}

	gui.raiseToast(gui.Tr.PullRequestURLCopiedToClipboard)

	return nil
}

func (gui *Gui) handleGitFetch(g *gocui.Gui, v *gocui.View) error {
	if err := gui.createLoaderPanel(gui.Tr.FetchWait); err != nil {
		return err
	}
	go utils.Safe(func() {
		err := gui.fetch(true)
		gui.handleCredentialsPopup(err)
		_ = gui.refreshSidePanels(refreshOptions{mode: ASYNC})
	})
	return nil
}

func (gui *Gui) handleForceCheckout(g *gocui.Gui, v *gocui.View) error {
	branch := gui.getSelectedBranch()
	message := gui.Tr.SureForceCheckout
	title := gui.Tr.ForceCheckoutBranch

	return gui.ask(askOpts{
		title:  title,
		prompt: message,
		handleConfirm: func() error {
			if err := gui.GitCommand.Checkout(branch.Name, commands.CheckoutOptions{Force: true}); err != nil {
				_ = gui.surfaceError(err)
			}
			return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
		},
	})
}

type handleCheckoutRefOptions struct {
	WaitingStatus string
	EnvVars       []string
	onRefNotFound func(ref string) error
}

func (gui *Gui) handleCheckoutRef(ref string, options handleCheckoutRefOptions) error {
	waitingStatus := options.WaitingStatus
	if waitingStatus == "" {
		waitingStatus = gui.Tr.CheckingOutStatus
	}

	cmdOptions := commands.CheckoutOptions{Force: false, EnvVars: options.EnvVars}

	onSuccess := func() {
		gui.State.Panels.Branches.SelectedLineIdx = 0
		gui.State.Panels.Commits.SelectedLineIdx = 0
		// loading a heap of commits is slow so we limit them whenever doing a reset
		gui.State.Panels.Commits.LimitCommits = true
	}

	return gui.WithWaitingStatus(waitingStatus, func() error {
		if err := gui.GitCommand.Checkout(ref, cmdOptions); err != nil {
			// note, this will only work for english-language git commands. If we force git to use english, and the error isn't this one, then the user will receive an english command they may not understand. I'm not sure what the best solution to this is. Running the command once in english and a second time in the native language is one option

			if options.onRefNotFound != nil && strings.Contains(err.Error(), "did not match any file(s) known to git") {
				return options.onRefNotFound(ref)
			}

			if strings.Contains(err.Error(), "Please commit your changes or stash them before you switch branch") {
				// offer to autostash changes
				return gui.ask(askOpts{

					title:  gui.Tr.AutoStashTitle,
					prompt: gui.Tr.AutoStashPrompt,
					handleConfirm: func() error {
						if err := gui.GitCommand.StashSave(gui.Tr.StashPrefix + ref); err != nil {
							return gui.surfaceError(err)
						}
						if err := gui.GitCommand.Checkout(ref, cmdOptions); err != nil {
							return gui.surfaceError(err)
						}

						onSuccess()
						if err := gui.GitCommand.StashDo(0, "pop"); err != nil {
							if err := gui.refreshSidePanels(refreshOptions{mode: BLOCK_UI}); err != nil {
								return err
							}
							return gui.surfaceError(err)
						}
						return gui.refreshSidePanels(refreshOptions{mode: BLOCK_UI})
					},
				})
			}

			if err := gui.surfaceError(err); err != nil {
				return err
			}
		}
		onSuccess()

		return gui.refreshSidePanels(refreshOptions{mode: BLOCK_UI})
	})
}

func (gui *Gui) handleCheckoutByName(g *gocui.Gui, v *gocui.View) error {
	return gui.prompt(promptOpts{
		title:               gui.Tr.BranchName + ":",
		findSuggestionsFunc: gui.findBranchNameSuggestions,
		handleConfirm: func(response string) error {
			return gui.handleCheckoutRef(response, handleCheckoutRefOptions{
				onRefNotFound: func(ref string) error {

					return gui.ask(askOpts{
						title:  gui.Tr.BranchNotFoundTitle,
						prompt: fmt.Sprintf("%s %s%s", gui.Tr.BranchNotFoundPrompt, ref, "?"),
						handleConfirm: func() error {
							return gui.createNewBranchWithName(ref)
						},
					})
				},
			})
		}},
	)
}

func (gui *Gui) getCheckedOutBranch() *models.Branch {
	if len(gui.State.Branches) == 0 {
		return nil
	}

	return gui.State.Branches[0]
}

func (gui *Gui) createNewBranchWithName(newBranchName string) error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}

	if err := gui.GitCommand.NewBranch(newBranchName, branch.Name); err != nil {
		return gui.surfaceError(err)
	}

	gui.State.Panels.Branches.SelectedLineIdx = 0
	return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
}

func (gui *Gui) handleDeleteBranch(g *gocui.Gui, v *gocui.View) error {
	return gui.deleteBranch(false)
}

func (gui *Gui) deleteBranch(force bool) error {
	selectedBranch := gui.getSelectedBranch()
	if selectedBranch == nil {
		return nil
	}
	checkedOutBranch := gui.getCheckedOutBranch()
	if checkedOutBranch.Name == selectedBranch.Name {
		return gui.createErrorPanel(gui.Tr.CantDeleteCheckOutBranch)
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

	return gui.ask(askOpts{
		title:  title,
		prompt: message,
		handleConfirm: func() error {
			if err := gui.GitCommand.DeleteBranch(selectedBranch.Name, force); err != nil {
				errMessage := err.Error()
				if !force && strings.Contains(errMessage, "is not fully merged") {
					return gui.deleteNamedBranch(selectedBranch, true)
				}
				return gui.createErrorPanel(errMessage)
			}
			return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{BRANCHES}})
		},
	})
}

func (gui *Gui) mergeBranchIntoCheckedOutBranch(branchName string) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	if gui.GitCommand.IsHeadDetached() {
		return gui.createErrorPanel("Cannot merge branch in detached head state. You might have checked out a commit directly or a remote branch, in which case you should checkout the local branch you want to be on")
	}
	checkedOutBranchName := gui.getCheckedOutBranch().Name
	if checkedOutBranchName == branchName {
		return gui.createErrorPanel(gui.Tr.CantMergeBranchIntoItself)
	}
	prompt := utils.ResolvePlaceholderString(
		gui.Tr.ConfirmMerge,
		map[string]string{
			"checkedOutBranch": checkedOutBranchName,
			"selectedBranch":   branchName,
		},
	)

	return gui.ask(askOpts{
		title:  gui.Tr.MergingTitle,
		prompt: prompt,
		handleConfirm: func() error {
			err := gui.GitCommand.Merge(branchName, commands.MergeOpts{})
			return gui.handleGenericMergeCommandResult(err)
		},
	})
}

func (gui *Gui) handleMerge(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	selectedBranchName := gui.getSelectedBranch().Name
	return gui.mergeBranchIntoCheckedOutBranch(selectedBranchName)
}

func (gui *Gui) handleRebaseOntoLocalBranch(g *gocui.Gui, v *gocui.View) error {
	selectedBranchName := gui.getSelectedBranch().Name
	return gui.handleRebaseOntoBranch(selectedBranchName)
}

func (gui *Gui) handleRebaseOntoBranch(selectedBranchName string) error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	checkedOutBranch := gui.getCheckedOutBranch().Name
	if selectedBranchName == checkedOutBranch {
		return gui.createErrorPanel(gui.Tr.CantRebaseOntoSelf)
	}
	prompt := utils.ResolvePlaceholderString(
		gui.Tr.ConfirmRebase,
		map[string]string{
			"checkedOutBranch": checkedOutBranch,
			"selectedBranch":   selectedBranchName,
		},
	)

	return gui.ask(askOpts{
		title:  gui.Tr.RebasingTitle,
		prompt: prompt,
		handleConfirm: func() error {
			err := gui.GitCommand.RebaseBranch(selectedBranchName)
			return gui.handleGenericMergeCommandResult(err)
		},
	})
}

func (gui *Gui) handleFastForward(g *gocui.Gui, v *gocui.View) error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}
	if branch.Pushables == "" {
		return nil
	}
	if branch.Pushables == "?" {
		return gui.createErrorPanel(gui.Tr.FwdNoUpstream)
	}
	if branch.Pushables != "0" {
		return gui.createErrorPanel(gui.Tr.FwdCommitsToPush)
	}

	upstream, err := gui.GitCommand.GetUpstreamForBranch(branch.Name)
	if err != nil {
		return gui.surfaceError(err)
	}

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
	go utils.Safe(func() {
		_ = gui.createLoaderPanel(message)

		if gui.State.Panels.Branches.SelectedLineIdx == 0 {
			_ = gui.pullWithMode("ff-only", PullFilesOptions{})
		} else {
			err := gui.GitCommand.FastForward(branch.Name, remoteName, remoteBranchName, gui.promptUserForCredential)
			gui.handleCredentialsPopup(err)
			_ = gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []int{BRANCHES}})
		}
	})
	return nil
}

func (gui *Gui) handleCreateResetToBranchMenu(g *gocui.Gui, v *gocui.View) error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}

	return gui.createResetMenu(branch.Name)
}

func (gui *Gui) handleRenameBranch(g *gocui.Gui, v *gocui.View) error {
	branch := gui.getSelectedBranch()
	if branch == nil {
		return nil
	}

	// TODO: find a way to not checkout the branch here if it's not the current branch (i.e. find some
	// way to get it to show up in the reflog)

	promptForNewName := func() error {
		return gui.prompt(promptOpts{
			title:          gui.Tr.NewBranchNamePrompt + " " + branch.Name + ":",
			initialContent: branch.Name,
			handleConfirm: func(newBranchName string) error {
				if err := gui.GitCommand.RenameBranch(branch.Name, newBranchName); err != nil {
					return gui.surfaceError(err)
				}
				// need to checkout so that the branch shows up in our reflog and therefore
				// doesn't get lost among all the other branches when we switch to something else
				if err := gui.GitCommand.Checkout(newBranchName, commands.CheckoutOptions{Force: false}); err != nil {
					return gui.surfaceError(err)
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
			},
		})
	}

	// I could do an explicit check here for whether the branch is tracking a remote branch
	// but if we've selected it we'll already know that via Pullables and Pullables.
	// Bit of a hack but I'm lazy.
	notTrackingRemote := branch.Pullables == "?"
	if notTrackingRemote {
		return promptForNewName()
	}

	return gui.ask(askOpts{
		title:         gui.Tr.LcRenameBranch,
		prompt:        gui.Tr.RenameBranchWarning,
		handleConfirm: promptForNewName,
	})
}

func (gui *Gui) currentBranch() *models.Branch {
	if len(gui.State.Branches) == 0 {
		return nil
	}
	return gui.State.Branches[0]
}

func (gui *Gui) handleNewBranchOffCurrentItem() error {
	context := gui.currentSideContext()

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
		// will set to the remote's existing name
		prefilledName = item.ID()
	}

	return gui.prompt(promptOpts{
		title:          message,
		initialContent: prefilledName,
		handleConfirm: func(response string) error {
			if err := gui.GitCommand.NewBranch(sanitizedBranchName(response), item.ID()); err != nil {
				return err
			}

			// if we're currently in the branch commits context then the selected commit
			// is about to go to the top of the list
			if context.GetKey() == BRANCH_COMMITS_CONTEXT_KEY {
				context.GetPanelState().SetSelectedLineIdx(0)
			}

			if context.GetKey() != gui.Contexts.Branches.Context.GetKey() {
				if err := gui.pushContext(gui.Contexts.Branches.Context); err != nil {
					return err
				}
			}

			gui.State.Panels.Branches.SelectedLineIdx = 0

			return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
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

func (gui *Gui) findBranchNameSuggestions(input string) []*types.Suggestion {
	branchNames := gui.getBranchNames()

	matchingBranchNames := utils.FuzzySearch(sanitizedBranchName(input), branchNames)

	suggestions := make([]*types.Suggestion, len(matchingBranchNames))
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
	return strings.ReplaceAll(input, " ", "-")
}
