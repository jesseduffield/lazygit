package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) branchesRenderToMain() error {
	var task types.UpdateTask
	branch := gui.State.Contexts.Branches.GetSelected()
	if branch == nil {
		task = types.NewRenderStringTask(gui.c.Tr.NoBranchesThisRepo)
	} else {
		cmdObj := gui.git.Branch.GetGraphCmdObj(branch.FullRefName())

		task = types.NewRunPtyTask(cmdObj.GetCmd())
	}

	return gui.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: gui.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: gui.c.Tr.LogTitle,
			Task:  task,
		},
	})
}

func (gui *Gui) refreshGithubPullRequests() error {
	if err := gui.git.Gh.BaseRepo(); err == nil {
		err := gui.setGithubPullRequests()
		if err != nil {
			return err
		}

		gui.refreshBranches()

		return nil
	}

	// when config not exits
	err := gui.refreshRemotes()
	if err != nil {
		return err
	}

	_ = gui.c.Prompt(types.PromptOpts{
		Title:               gui.c.Tr.SelectRemoteRepository,
		InitialContent:      "",
		FindSuggestionsFunc: gui.helpers.Suggestions.GetRemoteRepoSuggestionsFunc(),
		HandleConfirm: func(repository string) error {
			return gui.c.WithWaitingStatus(gui.c.Tr.LcSelectingRemote, func() error {
				_, err := gui.git.Gh.SetBaseRepo(repository)
				if err != nil {
					return err
				}

				err = gui.setGithubPullRequests()
				if err != nil {
					return err
				}
				gui.refreshBranches()
				return nil
			})
		},
	})

	return nil
}

func (gui *Gui) setGithubPullRequests() error {
	prs, err := gui.git.Gh.GithubMostRecentPRs()
	if err != nil {
		return gui.c.Error(err)
	}

	gui.State.Model.PullRequests = prs
	return nil
}
