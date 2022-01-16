package gui

import (
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// after selecting the 200th commit, we'll load in all the rest
const COMMIT_THRESHOLD = 200

// list panel functions

func (gui *Gui) getSelectedLocalCommit() *models.Commit {
	selectedLine := gui.State.Panels.Commits.SelectedLineIdx
	if selectedLine == -1 || selectedLine > len(gui.State.Commits)-1 {
		return nil
	}

	return gui.State.Commits[selectedLine]
}

func (gui *Gui) onCommitFocus() error {
	state := gui.State.Panels.Commits
	if state.SelectedLineIdx > COMMIT_THRESHOLD && state.LimitCommits {
		state.LimitCommits = false
		go utils.Safe(func() {
			if err := gui.refreshCommitsWithLimit(); err != nil {
				_ = gui.c.Error(err)
			}
		})
	}

	gui.escapeLineByLinePanel()

	return nil
}

func (gui *Gui) branchCommitsRenderToMain() error {
	var task updateTask
	commit := gui.getSelectedLocalCommit()
	if commit == nil {
		task = NewRenderStringTask(gui.c.Tr.NoCommitsThisBranch)
	} else {
		cmdObj := gui.git.Commit.ShowCmdObj(commit.Sha, gui.State.Modes.Filtering.GetPath())
		task = NewRunPtyTask(cmdObj.GetCmd())
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Patch",
			task:  task,
		},
		secondary: gui.secondaryPatchPanelUpdateOpts(),
	})
}

// during startup, the bottleneck is fetching the reflog entries. We need these
// on startup to sort the branches by recency. So we have two phases: INITIAL, and COMPLETE.
// In the initial phase we don't get any reflog commits, but we asynchronously get them
// and refresh the branches after that
func (gui *Gui) refreshReflogCommitsConsideringStartup() {
	switch gui.State.StartupStage {
	case INITIAL:
		go utils.Safe(func() {
			_ = gui.refreshReflogCommits()
			gui.refreshBranches()
			gui.State.StartupStage = COMPLETE
		})

	case COMPLETE:
		_ = gui.refreshReflogCommits()
	}
}

// whenever we change commits, we should update branches because the upstream/downstream
// counts can change. Whenever we change branches we should probably also change commits
// e.g. in the case of switching branches.
func (gui *Gui) refreshCommits() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go utils.Safe(func() {
		gui.refreshReflogCommitsConsideringStartup()

		gui.refreshBranches()
		wg.Done()
	})

	go utils.Safe(func() {
		_ = gui.refreshCommitsWithLimit()
		context, ok := gui.State.Contexts.CommitFiles.GetParentContext()
		if ok && context.GetKey() == BRANCH_COMMITS_CONTEXT_KEY {
			// This makes sense when we've e.g. just amended a commit, meaning we get a new commit SHA at the same position.
			// However if we've just added a brand new commit, it pushes the list down by one and so we would end up
			// showing the contents of a different commit than the one we initially entered.
			// Ideally we would know when to refresh the commit files context and when not to,
			// or perhaps we could just pop that context off the stack whenever cycling windows.
			// For now the awkwardness remains.
			commit := gui.getSelectedLocalCommit()
			if commit != nil {
				gui.State.Panels.CommitFiles.refName = commit.RefName()
				_ = gui.refreshCommitFilesView()
			}
		}
		wg.Done()
	})

	wg.Wait()
}

func (gui *Gui) refreshCommitsWithLimit() error {
	gui.Mutexes.BranchCommitsMutex.Lock()
	defer gui.Mutexes.BranchCommitsMutex.Unlock()

	commits, err := gui.git.Loaders.Commits.GetCommits(
		loaders.GetCommitsOptions{
			Limit:                gui.State.Panels.Commits.LimitCommits,
			FilterPath:           gui.State.Modes.Filtering.GetPath(),
			IncludeRebaseCommits: true,
			RefName:              gui.refForLog(),
			All:                  gui.ShowWholeGitGraph,
		},
	)
	if err != nil {
		return err
	}
	gui.State.Commits = commits

	return gui.c.PostRefreshUpdate(gui.State.Contexts.BranchCommits)
}

func (gui *Gui) refForLog() string {
	bisectInfo := gui.git.Bisect.GetInfo()
	gui.State.BisectInfo = bisectInfo

	if !bisectInfo.Started() {
		return "HEAD"
	}

	// need to see if our bisect's current commit is reachable from our 'new' ref.
	if bisectInfo.Bisecting() && !gui.git.Bisect.ReachableFromStart(bisectInfo) {
		return bisectInfo.GetNewSha()
	}

	return bisectInfo.GetStartSha()
}

func (gui *Gui) refreshRebaseCommits() error {
	gui.Mutexes.BranchCommitsMutex.Lock()
	defer gui.Mutexes.BranchCommitsMutex.Unlock()

	updatedCommits, err := gui.git.Loaders.Commits.MergeRebasingCommits(gui.State.Commits)
	if err != nil {
		return err
	}
	gui.State.Commits = updatedCommits

	return gui.c.PostRefreshUpdate(gui.State.Contexts.BranchCommits)
}
