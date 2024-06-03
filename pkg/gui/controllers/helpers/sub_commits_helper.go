package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type SubCommitsHelper struct {
	c *HelperCommon

	refreshHelper *RefreshHelper
	setSubCommits func([]*models.Commit)
}

func NewSubCommitsHelper(
	c *HelperCommon,
	refreshHelper *RefreshHelper,
	setSubCommits func([]*models.Commit),
) *SubCommitsHelper {
	return &SubCommitsHelper{
		c:             c,
		refreshHelper: refreshHelper,
		setSubCommits: setSubCommits,
	}
}

type ViewSubCommitsOpts struct {
	Ref                     types.Ref
	RefToShowDivergenceFrom string
	TitleRef                string
	Context                 types.Context
	ShowBranchHeads         bool
}

func (self *SubCommitsHelper) ViewSubCommits(opts ViewSubCommitsOpts) error {
	commits, err := self.c.Git().Loaders.CommitLoader.GetCommits(
		git_commands.GetCommitsOptions{
			Limit:                   true,
			FilterPath:              self.c.Modes().Filtering.GetPath(),
			FilterAuthor:            self.c.Modes().Filtering.GetAuthor(),
			IncludeRebaseCommits:    false,
			RefName:                 opts.Ref.FullRefName(),
			RefForPushedStatus:      opts.Ref.FullRefName(),
			RefToShowDivergenceFrom: opts.RefToShowDivergenceFrom,
			MainBranches:            self.c.Model().MainBranches,
		},
	)
	if err != nil {
		return err
	}

	self.setSubCommits(commits)
	self.refreshHelper.RefreshAuthors(commits)

	subCommitsContext := self.c.Contexts().SubCommits
	subCommitsContext.SetSelection(0)
	subCommitsContext.SetParentContext(opts.Context)
	subCommitsContext.SetWindowName(opts.Context.GetWindowName())
	subCommitsContext.SetTitleRef(utils.TruncateWithEllipsis(opts.TitleRef, 50))
	subCommitsContext.SetRef(opts.Ref)
	subCommitsContext.SetRefToShowDivergenceFrom(opts.RefToShowDivergenceFrom)
	subCommitsContext.SetLimitCommits(true)
	subCommitsContext.SetShowBranchHeads(opts.ShowBranchHeads)
	subCommitsContext.ClearSearchString()
	subCommitsContext.GetView().ClearSearch()
	subCommitsContext.GetView().TitlePrefix = opts.Context.GetView().TitlePrefix

	err = self.c.PostRefreshUpdate(self.c.Contexts().SubCommits)
	if err != nil {
		return err
	}

	return self.c.PushContext(self.c.Contexts().SubCommits)
}
