package controllers

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// This controller is for all contexts that contain commit files.

var _ types.IController = &SwitchToDiffFilesController{}

type CanSwitchToDiffFiles interface {
	types.IListContext
	CanRebase() bool
	GetSelectedRef() models.Ref
	GetSelectedRefRangeForDiffFiles() *types.RefRange
}

// Not using our ListControllerTrait because we have our own way of working with
// range selections that's different from ListControllerTrait's
type SwitchToDiffFilesController struct {
	baseController
	c       *ControllerCommon
	context CanSwitchToDiffFiles

	// what this panel offers on the diff it shows in the focused main view
	diffActions *CommitDiffActions
}

func NewSwitchToDiffFilesController(
	c *ControllerCommon,
	context CanSwitchToDiffFiles,
) *SwitchToDiffFilesController {
	controller := &SwitchToDiffFilesController{
		baseController: baseController{},
		c:              c,
		context:        context,
	}
	controller.diffActions = NewCommitDiffActions(c, context, controller.diffTarget)
	return controller
}

// diffTarget is the commit — or stash entry, or range of commits — the panel has
// selected, whose whole diff its main view shows.
func (self *SwitchToDiffFilesController) diffTarget() *commitDiffTarget {
	ref := self.context.GetSelectedRef()
	if ref == nil {
		return nil
	}
	refRange := self.context.GetSelectedRefRangeForDiffFiles()
	from, to := context.FromAndToForDiff(ref, refRange)
	return &commitDiffTarget{from: from, to: to, canRebase: self.canRebase(ref, refRange)}
}

// canRebase reports whether the given selection is one lazygit may rewrite: the panel
// has to allow it in the first place, a range of commits can't be rewritten as one,
// and in diffing mode what the main view shows is a diff against another ref rather
// than the commit itself, unless that other ref is the selected commit.
func (self *SwitchToDiffFilesController) canRebase(ref models.Ref, refRange *types.RefRange) bool {
	if !self.context.CanRebase() {
		return false
	}
	if self.c.Modes().Diffing.Active() {
		return self.c.Modes().Diffing.Ref == ref.RefName()
	}
	return refRange == nil
}

func (self *SwitchToDiffFilesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Keys:              opts.GetKeys(opts.Config.Universal.GoInto),
			Handler:           self.enter,
			GetDisabledReason: self.canEnter,
			Description:       self.c.Tr.ViewItemFiles,
		},
	}

	return bindings
}

func (self *SwitchToDiffFilesController) Context() types.Context {
	return self.context
}

func (self *SwitchToDiffFilesController) GetOnDoubleClick() func() error {
	return func() error {
		if self.canEnter() == nil {
			return self.enter()
		}

		return nil
	}
}

func (self *SwitchToDiffFilesController) GetFocusedMainViewDiffSource() types.FocusedMainViewDiffSource {
	return self.diffActions
}

func (self *SwitchToDiffFilesController) enter() error {
	ref := self.context.GetSelectedRef()
	refsRange := self.context.GetSelectedRefRangeForDiffFiles()
	return viewCommitFiles(self.c, ref, refsRange, self.canRebase(ref, refsRange), self.context)
}

// viewCommitFiles opens the commit files panel on top of parent, listing the files of
// the diff of ref, or of refRange where that isn't nil.
func viewCommitFiles(
	c *ControllerCommon, ref models.Ref, refRange *types.RefRange, canRebase bool, parent types.Context,
) error {
	commitFilesContext := c.Contexts().CommitFiles

	commitFilesContext.ClearFilter()
	commitFilesContext.ReInit(ref, refRange)
	commitFilesContext.SetSelection(0)
	commitFilesContext.SetCanRebase(canRebase)
	commitFilesContext.SetParentContext(parent)
	commitFilesContext.SetWindowName(parent.GetWindowName())
	commitFilesContext.GetView().TitlePrefix = parent.GetView().TitlePrefix

	c.Refresh(types.RefreshOptions{
		Scope: []types.RefreshableView{types.COMMIT_FILES},
		Then: func() error {
			if filterPath := c.Modes().Filtering.GetPath(); filterPath != "" {
				path, err := filepath.Rel(c.Git().RepoPaths.RepoPath(), filterPath)
				if err != nil {
					path = filterPath
				}
				commitFilesContext.CommitFileTreeViewModel.SelectPath(
					filepath.ToSlash(path), c.UserConfig().Gui.ShowRootItemInFileTree)
			}
			c.Context().Push(commitFilesContext, types.OnFocusOpts{})
			return nil
		},
	})
	return nil
}

func (self *SwitchToDiffFilesController) canEnter() *types.DisabledReason {
	refRange := self.context.GetSelectedRefRangeForDiffFiles()
	if refRange != nil {
		return nil
	}
	ref := self.context.GetSelectedRef()
	if ref == nil {
		return &types.DisabledReason{Text: self.c.Tr.NoItemSelected}
	}
	if ref.RefName() == "" {
		return &types.DisabledReason{Text: self.c.Tr.SelectedItemDoesNotHaveFiles}
	}

	return nil
}
