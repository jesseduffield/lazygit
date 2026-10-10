package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ReposController struct {
	baseController
	*ListControllerTrait[*models.Repo]
	c *ControllerCommon
}

var _ types.IController = &ReposController{}

func NewReposController(c *ControllerCommon) *ReposController {
	return &ReposController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().Repos,
			c.Contexts().Repos.GetSelected,
			c.Contexts().Repos.GetSelectedItems,
		),
		c: c,
	}
}

func (self *ReposController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Keys:              opts.GetKeys(opts.Config.Universal.Select),
			Handler:           self.withItem(self.enter),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Switch,
			Tooltip:           self.c.Tr.SwitchToRepoTooltip,
			DisplayOnScreen:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.GoInto),
			Handler:           self.withItem(self.enter),
			GetDisabledReason: self.require(self.singleItemSelected()),
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.Edit),
			Handler:         (&EditConfigAction{c: self.c}).Call,
			Description:     self.c.Tr.EditConfig,
			Tooltip:         self.c.Tr.EditFileTooltip,
			DisplayOnScreen: true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Status.CheckForUpdate),
			Handler:     self.c.Helpers().Update.CheckForUpdateInForeground,
			Description: self.c.Tr.CheckForUpdate,
		},
	}
}

func (self *ReposController) GetOnRenderToMain() func() {
	return func() {
		var task types.UpdateTask
		repo := self.context().GetSelected()
		if repo == nil {
			task = types.NewRenderStringTask(self.c.Tr.NoRepos)
		} else {
			cmdObj := self.c.OS().Cmd.New(
				git_commands.NewGitCmd("status").
					Config("color.status=always").
					Dir(repo.Path).
					ToArgv(),
			).DontLog()
			task = types.NewRunCommandTask(git_commands.ForOtherRepo(cmdObj).GetCmd())
		}

		self.c.RenderToMainViews(types.RefreshMainOpts{
			Pair: self.c.MainViewPairs().Normal,
			Main: &types.ViewUpdateOpts{
				Title: self.c.Tr.ReposTitle,
				Task:  task,
			},
		})
	}
}

func (self *ReposController) GetOnDoubleClick() func() error {
	return self.withItemGraceful(self.enter)
}

func (self *ReposController) enter(repo *models.Repo) error {
	if presentation.IsCurrentRepo(repo, self.c.Git().RepoPaths.WorktreePath()) {
		return nil
	}

	if err := self.c.Helpers().Repos.SwitchToTopLevelRepo(repo.Path, context.REPOS_CONTEXT_KEY); err != nil {
		return err
	}

	// The switch focuses the new repo's repos context with its cursor at the
	// top, so put the cursor back on the repo the user picked.
	self.context().SelectRepo(repo.Path)
	self.c.PostRefreshUpdate(self.context())
	return nil
}

func (self *ReposController) context() *context.ReposContext {
	return self.c.Contexts().Repos
}
