package controllers

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/commands/models"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
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
}

func NewSwitchToDiffFilesController(
	c *ControllerCommon,
	context CanSwitchToDiffFiles,
) *SwitchToDiffFilesController {
	return &SwitchToDiffFilesController{
		baseController: baseController{},
		c:              c,
		context:        context,
	}
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

func (self *SwitchToDiffFilesController) GetFocusedMainViewActions() types.FocusedMainViewActions {
	return self
}

func (self *SwitchToDiffFilesController) OnClick(mainViewName string, clickedLineIdx int) error {
	info, ok := self.c.Helpers().Staging.GetDiffLineInfo(mainViewName, clickedLineIdx)
	if !ok {
		return nil
	}
	line, isDeletion := info.PatchSelectLine()

	// Capture before self.enter() pushes the commit files panel, which
	// re-renders the main view. We escape "all the way out" to this side
	// panel (skipping the commit files panel), then focus the main view.
	snapshot := focusedMainViewSnapshot(self.c, mainViewName, self.context)

	if err := self.enter(); err != nil {
		return err
	}

	context := self.c.Contexts().CommitFiles
	var node *filetree.CommitFileNode

	relativePath, err := filepath.Rel(self.c.Git().RepoPaths.WorktreePath(), info.Path)
	if err != nil {
		return err
	}
	relativePath = "./" + relativePath
	context.CommitFileTreeViewModel.ExpandToPath(relativePath)
	self.c.PostRefreshUpdate(context)

	idx, ok := context.CommitFileTreeViewModel.GetIndexForPath(relativePath)
	if !ok {
		return nil
	}

	context.SetSelectedLineIdx(idx)
	context.GetViewTrait().FocusPoint(
		context.ModelIndexToViewIndex(idx), false)
	node = context.GetSelected()
	return self.c.Helpers().CommitFiles.EnterCommitFile(node, snapshot, types.OnFocusOpts{ClickedWindowName: "main", ClickedViewLineIdx: line, ClickedViewRealLineIdx: line, ClickedViewRealLineIsDeletion: isDeletion, SelectLineInDefaultMode: true})
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

func (self *SwitchToDiffFilesController) enter() error {
	ref := self.context.GetSelectedRef()
	refsRange := self.context.GetSelectedRefRangeForDiffFiles()
	commitFilesContext := self.c.Contexts().CommitFiles

	canRebase := self.canRebase(ref, refsRange)

	commitFilesContext.ClearFilter()
	commitFilesContext.ReInit(ref, refsRange)
	commitFilesContext.SetSelection(0)
	commitFilesContext.SetCanRebase(canRebase)
	commitFilesContext.SetParentContext(self.context)
	commitFilesContext.SetWindowName(self.context.GetWindowName())
	commitFilesContext.GetView().TitlePrefix = self.context.GetView().TitlePrefix

	self.c.Refresh(types.RefreshOptions{
		Scope: []types.RefreshableView{types.COMMIT_FILES},
	})

	if filterPath := self.c.Modes().Filtering.GetPath(); filterPath != "" {
		path, err := filepath.Rel(self.c.Git().RepoPaths.RepoPath(), filterPath)
		if err != nil {
			path = filterPath
		}
		commitFilesContext.CommitFileTreeViewModel.SelectPath(
			filepath.ToSlash(path), self.c.UserConfig().Gui.ShowRootItemInFileTree)
	}

	self.c.Context().Push(commitFilesContext, types.OnFocusOpts{})
	return nil
}

// canRebase reports whether patches built from the selected ref may modify commits —
// true only for commits of the currently checked-out branch, and not while diffing a
// different ref or over a range. Shared by entering the commit files panel and toggling
// patch lines straight from the main view.
func (self *SwitchToDiffFilesController) canRebase(ref models.Ref, refsRange *types.RefRange) bool {
	canRebase := self.context.CanRebase()
	if canRebase {
		if self.c.Modes().Diffing.Active() {
			if self.c.Modes().Diffing.Ref != ref.RefName() {
				canRebase = false
			}
		} else if refsRange != nil {
			canRebase = false
		}
	}
	return canRebase
}

// PrimaryAction toggles the selected line(s) of the whole-commit diff into or out of the
// custom patch when space is pressed in the focused main view of the commits / sub-commits
// / stash panels. The patch target is the panel's selected ref (or range), matching the
// diff the main view shows. Unlike the commit files panel there are no per-file patch
// indicators to update, so the toggle refreshes cheaply: it re-renders just this panel's
// main + secondary views (leaving the commit list untouched, which a list refresh would
// needlessly reload on every keystroke), re-running the same diff command (scroll
// preserved) and repainting the inclusion gutter.
func (self *SwitchToDiffFilesController) PrimaryAction(mainViewName string, firstLineIdx int, lastLineIdx int) error {
	ref := self.context.GetSelectedRef()
	if ref == nil {
		return nil
	}
	refsRange := self.context.GetSelectedRefRangeForDiffFiles()

	from, to := context.FromAndToForDiff(ref, refsRange)
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(from)
	canRebase := self.canRebase(ref, refsRange)

	return togglePatchFromFocusedMainView(self.c, mainViewName, firstLineIdx, lastLineIdx,
		from, to, reverse, canRebase,
		func() {
			self.c.OnUIThread(func() error {
				self.c.PostRefreshUpdate(self.context)
				return nil
			})
		})
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
