package controllers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type StatusController struct {
	baseController
	c *ControllerCommon
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
			Key:         opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:     self.openConfig,
			Description: self.c.Tr.OpenConfig,
			Tooltip:     self.c.Tr.OpenFileTooltip,
		},
		{
			Key:             opts.GetKey(opts.Config.Universal.Edit),
			Handler:         self.editConfig,
			Description:     self.c.Tr.EditConfig,
			Tooltip:         self.c.Tr.EditFileTooltip,
			DisplayOnScreen: true,
		},
		{
			Key:             opts.GetKey(opts.Config.Status.CheckForUpdate),
			Handler:         self.handleCheckForUpdate,
			Description:     self.c.Tr.CheckForUpdate,
			DisplayOnScreen: true,
		},
		{
			Key:             opts.GetKey(opts.Config.Status.RecentRepos),
			Handler:         self.c.Helpers().Repos.CreateRecentReposMenu,
			Description:     self.c.Tr.SwitchRepo,
			DisplayOnScreen: true,
		},
		{
			Key:         opts.GetKey(opts.Config.Status.AllBranchesLogGraph),
			Handler:     func() error { self.showAllBranchLogs(); return nil },
			Description: self.c.Tr.AllBranchesLogGraph,
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
		default:
			self.showDashboard()
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

	self.c.Context().Push(self.Context())

	upstreamStatus := utils.Decolorise(presentation.BranchStatus(currentBranch, types.ItemOperationNone, self.c.Tr, time.Now(), self.c.UserConfig()))
	repoName := self.c.Git().RepoPaths.RepoName()
	workingTreeState := self.c.Git().Status.WorkingTreeState()
	switch workingTreeState {
	case enums.REBASE_MODE_REBASING, enums.REBASE_MODE_MERGING:
		workingTreeStatus := fmt.Sprintf("(%s)", presentation.FormatWorkingTreeStateLower(self.c.Tr, workingTreeState))
		if cursorInSubstring(opts.X, upstreamStatus+" ", workingTreeStatus) {
			return self.c.Helpers().MergeAndRebase.CreateRebaseOptionsMenu()
		}
		if cursorInSubstring(opts.X, upstreamStatus+" "+workingTreeStatus+" ", repoName) {
			return self.c.Helpers().Repos.CreateRecentReposMenu()
		}
	default:
		if cursorInSubstring(opts.X, upstreamStatus+" ", repoName) {
			return self.c.Helpers().Repos.CreateRecentReposMenu()
		}
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

	self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: self.c.Tr.LogTitle,
			Task:  task,
		},
	})
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
			fmt.Sprintf("Keybindings: %s", style.PrintSimpleHyperlink(fmt.Sprintf(constants.Links.Docs.Keybindings, versionStr))),
			fmt.Sprintf("Config Options: %s", style.PrintSimpleHyperlink(fmt.Sprintf(constants.Links.Docs.Config, versionStr))),
			fmt.Sprintf("Tutorial: %s", style.PrintSimpleHyperlink(constants.Links.Docs.Tutorial)),
			fmt.Sprintf("Raise an Issue: %s", style.PrintSimpleHyperlink(constants.Links.Issues)),
			fmt.Sprintf("Release Notes: %s", style.PrintSimpleHyperlink(constants.Links.Releases)),
			style.FgMagenta.Sprintf("Become a sponsor: %s", style.PrintSimpleHyperlink(constants.Links.Donate)), // caffeine ain't free
		}, "\n\n") + "\n"

	self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: self.c.Tr.StatusTitle,
			Task:  types.NewRenderStringTask(dashboardString),
		},
	})
}

func (self *StatusController) handleCheckForUpdate() error {
	return self.c.Helpers().Update.CheckForUpdateInForeground()
}
