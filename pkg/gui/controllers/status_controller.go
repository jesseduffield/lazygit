package controllers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type StatusController struct {
	baseController
	c *ControllerCommon

	// Deployments are fetched from GitHub lazily, the first time the status
	// panel's deployments view is shown for a repo, and then cached for the rest
	// of the session. This avoids a network round-trip (and the git/gh-auth work
	// needed to set it up) on every refresh while the panel is focused. All of
	// these fields are only ever touched on the UI thread, so they need no lock.
	deploymentsRepoPath string // the repo deploymentsContent belongs to
	deploymentsContent  string // rendered content, valid once deploymentsFetched
	deploymentsFetched  bool   // a fetch has completed for deploymentsRepoPath
	deploymentsFetching bool   // a fetch is currently in flight
}

var _ types.IController = &StatusController{}

func NewStatusController(
	c *ControllerCommon,
) *StatusController {
	return &StatusController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *StatusController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Keys:        opts.GetKeys(opts.Config.Universal.OpenFile),
			Handler:     self.openConfig,
			Description: self.c.Tr.OpenConfig,
			Tooltip:     self.c.Tr.OpenFileTooltip,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.Edit),
			Handler:         self.editConfig,
			Description:     self.c.Tr.EditConfig,
			Tooltip:         self.c.Tr.EditFileTooltip,
			DisplayOnScreen: true,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Status.CheckForUpdate),
			Handler:         self.handleCheckForUpdate,
			Description:     self.c.Tr.CheckForUpdate,
			DisplayOnScreen: true,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Status.RecentRepos),
			Handler:         self.c.Helpers().Repos.CreateRecentReposMenu,
			Description:     self.c.Tr.SwitchRepo,
			DisplayOnScreen: true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Status.AllBranchesLogGraph),
			Handler:     func() error { self.switchToOrRotateAllBranchesLogs(); return nil },
			Description: self.c.Tr.AllBranchesLogGraph,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Status.AllBranchesLogGraphReverse),
			Handler:     func() error { self.switchToOrRotateAllBranchesLogsBackward(); return nil },
			Description: self.c.Tr.AllBranchesLogGraphReverse,
		},
	}

	return bindings
}

func (self *StatusController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName: self.Context().GetViewName(),
			Key:      gocui.MouseLeft,
			Handler:  self.onClick,
		},
	}
}

func (self *StatusController) GetOnRenderToMain() func() {
	return func() {
		switch self.c.UserConfig().Gui.StatusPanelView {
		case "dashboard":
			self.showDashboard()
		case "allBranchesLog":
			self.showAllBranchLogs()
		case "deployments":
			self.showDeployments()
		default:
			self.showDeployments()
		}
	}
}

func (self *StatusController) Context() types.Context {
	return self.c.Contexts().Status
}

func (self *StatusController) onClick(opts gocui.ViewMouseBindingOpts) error {
	// TODO: move into some abstraction (status is currently not a listViewContext where a lot of this code lives)
	currentBranch := self.c.Helpers().Refs.GetCheckedOutRef()
	if currentBranch == nil {
		// need to wait for branches to refresh
		return nil
	}

	self.c.Context().Push(self.Context(), types.OnFocusOpts{})

	upstreamStatus := utils.Decolorise(presentation.BranchStatus(currentBranch, types.ItemOperationNone, self.c.Tr, time.Now(), self.c.UserConfig()))
	repoName := self.c.Git().RepoPaths.RepoName()
	workingTreeState := self.c.Git().Status.WorkingTreeState()
	if workingTreeState.Any() {
		workingTreeStatus := fmt.Sprintf("(%s)", workingTreeState.LowerCaseTitle(self.c.Tr))
		if cursorInSubstring(opts.X, upstreamStatus+" ", workingTreeStatus) {
			return self.c.Helpers().MergeAndRebase.CreateRebaseOptionsMenu()
		}
		if cursorInSubstring(opts.X, upstreamStatus+" "+workingTreeStatus+" ", repoName) {
			return self.c.Helpers().Repos.CreateRecentReposMenu()
		}
	} else if cursorInSubstring(opts.X, upstreamStatus+" ", repoName) {
		return self.c.Helpers().Repos.CreateRecentReposMenu()
	}

	return nil
}

func runeCount(str string) int {
	return len([]rune(str))
}

func cursorInSubstring(cx int, prefix string, substring string) bool {
	return cx >= runeCount(prefix) && cx < runeCount(prefix+substring)
}

func lazygitTitle() string {
	return `
   _                       _ _
  | |                     (_) |
  | | __ _ _____   _  __ _ _| |_
  | |/ _` + "`" + ` |_  / | | |/ _` + "`" + ` | | __|
  | | (_| |/ /| |_| | (_| | | |_
  |_|\__,_/___|\__, |\__, |_|\__|
                __/ | __/ |
               |___/ |___/       `
}

func (self *StatusController) askForConfigFile(action func(file string) error) error {
	confPaths := self.c.GetConfig().GetUserConfigPaths()
	switch len(confPaths) {
	case 0:
		return errors.New(self.c.Tr.NoConfigFileFoundErr)
	case 1:
		return action(confPaths[0])
	default:
		menuItems := lo.Map(confPaths, func(path string, _ int) *types.MenuItem {
			return &types.MenuItem{
				Label: path,
				OnPress: func() error {
					return action(path)
				},
			}
		})

		return self.c.Menu(types.CreateMenuOptions{
			Title: self.c.Tr.SelectConfigFile,
			Items: menuItems,
		})
	}
}

func (self *StatusController) openConfig() error {
	return self.askForConfigFile(self.c.Helpers().Files.OpenFile)
}

func (self *StatusController) editConfig() error {
	return self.askForConfigFile(func(file string) error {
		return self.c.Helpers().Files.EditFiles([]string{file})
	})
}

func (self *StatusController) showAllBranchLogs() {
	cmdObj := self.c.Git().Branch.AllBranchesLogCmdObj()
	task := types.NewRunPtyTask(cmdObj.GetCmd())

	title := self.c.Tr.LogTitle
	if i, n := self.c.Git().Branch.GetAllBranchesLogIdxAndCount(); n > 1 {
		title = fmt.Sprintf(self.c.Tr.LogXOfYTitle, i+1, n)
	}
	self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: title,
			Task:  task,
		},
	})
}

// Switches to the all branches view, or, if already on that view,
// rotates to the next command in the list, and then renders it.
func (self *StatusController) switchToOrRotateAllBranchesLogs() {
	// A bit of a hack to ensure we only rotate to the next branch log command
	// if we currently are looking at a branch log. Otherwise, we should just show
	// the current index (if we are coming from the dashboard).
	if self.c.Views().Main.Title != self.c.Tr.StatusTitle {
		self.c.Git().Branch.RotateAllBranchesLogIdx()
	}
	self.showAllBranchLogs()
}

// Switches to the all branches view, or, if already on that view,
// rotates to the previous command in the list, and then renders it.
func (self *StatusController) switchToOrRotateAllBranchesLogsBackward() {
	// A bit of a hack to ensure we only rotate to the previous branch log command
	// if we currently are looking at a branch log. Otherwise, we should just show
	// the current index (if we are coming from the dashboard).
	if self.c.Views().Main.Title != self.c.Tr.StatusTitle {
		self.c.Git().Branch.RotateAllBranchesLogIdxBackward()
	}
	self.showAllBranchLogs()
}

func (self *StatusController) showDashboard() {
	versionStr := "master"
	version, err := types.ParseVersionNumber(self.c.GetConfig().GetVersion())
	if err == nil {
		// Don't just take the version string as is, but format it again. This
		// way it will be correct even if a distribution omits the "v", or the
		// ".0" at the end.
		versionStr = fmt.Sprintf("v%d.%d.%d", version.Major, version.Minor, version.Patch)
	}

	dashboardString := strings.Join(
		[]string{
			lazygitTitle(),
			fmt.Sprintf("Copyright %d Jesse Duffield", time.Now().Year()),
			fmt.Sprintf("Keybindings: %s", fmt.Sprintf(constants.Links.Docs.Keybindings, versionStr)),
			fmt.Sprintf("Config Options: %s", fmt.Sprintf(constants.Links.Docs.Config, versionStr)),
			fmt.Sprintf("Tutorial: %s", constants.Links.Docs.Tutorial),
			fmt.Sprintf("Raise an Issue: %s", constants.Links.Issues),
			fmt.Sprintf("Release Notes: %s", constants.Links.Releases),
			style.FgMagenta.Sprintf("Become a sponsor: %s", constants.Links.Donate), // caffeine ain't free
		}, "\n\n",
	) + "\n"

	self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: self.c.Tr.StatusTitle,
			Task:  types.NewRenderStringTask(dashboardString),
		},
	})
}

// showDeployments shows the environment deployment statuses for the repo,
// fetched from GitHub. This runs on the UI thread on every render-to-main (i.e.
// on every refresh while the status panel is focused), so it must be cheap: the
// only synchronous work it does is a subprocess-free check for a GitHub remote
// (to fall back to the dashboard for non-GitHub repos). The expensive work
// (resolving the base remote + auth token, and the network request) happens
// once per repo on a worker, and the result is cached and re-rendered from
// there on subsequent calls.
func (self *StatusController) showDeployments() {
	repoPath := self.c.Git().RepoPaths.RepoPath()
	if repoPath != self.deploymentsRepoPath {
		// Switched repos: drop the previous repo's cache.
		self.deploymentsRepoPath = repoPath
		self.deploymentsContent = ""
		self.deploymentsFetched = false
		self.deploymentsFetching = false
	}

	if !self.c.Helpers().Host.HasGithubRemote() {
		self.showDashboard()
		return
	}

	if self.deploymentsFetched {
		self.renderToStatusMain(self.deploymentsContent)
		return
	}

	self.renderToStatusMain(self.c.Tr.FetchingDeploymentsStatus)

	if self.deploymentsFetching {
		return
	}
	self.deploymentsFetching = true

	fetchRepoPath := repoPath
	self.c.OnWorker(func(_ gocui.Task) error {
		content := self.fetchDeploymentsContent()

		self.c.OnUIThread(func() error {
			// Discard the result if the user switched repos while we were fetching.
			if self.deploymentsRepoPath != fetchRepoPath {
				return nil
			}
			self.deploymentsFetching = false
			self.deploymentsFetched = true
			self.deploymentsContent = content

			// Only render if the status panel is still showing its deployments view;
			// the user may have navigated away or switched to the all-branches log
			// (which renders into the same main view) while we were fetching.
			if self.statusMainShowsDeployments() {
				self.renderToStatusMain(content)
			}
			return nil
		})
		return nil
	})
}

// fetchDeploymentsContent runs on a worker. It resolves the GitHub base remote
// and auth token (which shells out to git and reads gh's config) and performs
// the network request, then turns the outcome into a renderable string.
func (self *StatusController) fetchDeploymentsContent() string {
	serviceInfo, token, ok := self.c.Helpers().Host.GithubBaseRemote()
	if !ok {
		return self.c.Tr.DeploymentsNotAuthenticated
	}

	deployments, err := self.c.Git().GitHub.FetchDeployments(&serviceInfo, token)
	switch {
	case err != nil:
		self.c.Log.Error("error fetching deployments from GitHub: " + err.Error())
		return fmt.Sprintf(self.c.Tr.FetchingDeploymentsError, err.Error())
	case len(deployments) == 0:
		return self.c.Tr.NoDeploymentsFound
	default:
		return presentation.GetDeploymentsContent(deployments, self.c.Tr)
	}
}

// statusMainShowsDeployments reports whether the status panel is focused and its
// main view is currently showing the deployments/status content (as opposed to
// the all-branches log, which the user can switch to with a keybinding and which
// renders into the same main view with a different title).
func (self *StatusController) statusMainShowsDeployments() bool {
	return self.c.Context().IsCurrent(self.Context()) &&
		self.c.Views().Main.Title == self.c.Tr.StatusTitle
}

func (self *StatusController) renderToStatusMain(str string) {
	self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: self.c.Tr.StatusTitle,
			Task:  types.NewRenderStringTask(str),
		},
	})
}

func (self *StatusController) handleCheckForUpdate() error {
	return self.c.Helpers().Update.CheckForUpdateInForeground()
}
