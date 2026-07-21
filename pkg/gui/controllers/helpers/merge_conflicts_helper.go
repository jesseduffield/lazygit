package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MergeConflictsHelper struct {
	c *HelperCommon
}

func NewMergeConflictsHelper(
	c *HelperCommon,
) *MergeConflictsHelper {
	return &MergeConflictsHelper{
		c: c,
	}
}

func (self *MergeConflictsHelper) SetMergeState(path string) (bool, error) {
	self.context().GetMutex().Lock()
	defer self.context().GetMutex().Unlock()

	return self.setMergeStateWithoutLock(path)
}

func (self *MergeConflictsHelper) setMergeStateWithoutLock(path string) (bool, error) {
	content, err := self.c.Git().File.Cat(path)
	if err != nil {
		return false, err
	}

	if path != self.context().GetState().GetPath() {
		self.context().SetUserScrolling(false)
	}

	self.context().GetState().SetContent(content, path)

	return !self.context().GetState().NoConflicts(), nil
}

func (self *MergeConflictsHelper) ResetMergeState() {
	self.context().GetMutex().Lock()
	defer self.context().GetMutex().Unlock()

	self.resetMergeState()
}

func (self *MergeConflictsHelper) resetMergeState() {
	self.context().SetUserScrolling(false)
	self.context().GetState().Reset()
}

// EscapeMerge returns from the merge conflicts view to the files context. It
// must be called on the UI thread, without the merge-conflicts mutex held:
// pushing the files context renders the newly focused file to the main view,
// which can take the mutex again (via SetMergeState).
func (self *MergeConflictsHelper) EscapeMerge() {
	self.ResetMergeState()

	// The files refresh may already have opened the prompt to continue the
	// rebase/merge on top of us (if all conflicts are resolved); in that case
	// don't push the files context over it.
	if self.c.Context().IsCurrent(self.c.Contexts().MergeConflicts) {
		self.c.Context().Push(self.c.Contexts().Files, types.OnFocusOpts{})
	}
}

// SetConflictsAndRender re-reads the file being merged and re-renders the
// merge conflicts view. Returns whether the file still has conflicts.
func (self *MergeConflictsHelper) SetConflictsAndRender() (bool, error) {
	self.context().GetMutex().Lock()
	defer self.context().GetMutex().Unlock()

	hasConflicts, err := self.setMergeStateWithoutLock(self.context().GetState().GetPath())
	if err != nil {
		return false, err
	}

	if hasConflicts {
		return true, self.context().Render()
	}

	return false, nil
}

func (self *MergeConflictsHelper) SwitchToMerge(path string) error {
	if self.context().GetState().GetPath() != path {
		hasConflicts, err := self.SetMergeState(path)
		if err != nil {
			return err
		}
		if !hasConflicts {
			return nil
		}
	}

	self.c.Context().Push(self.c.Contexts().MergeConflicts, types.OnFocusOpts{})
	return nil
}

func (self *MergeConflictsHelper) context() *context.MergeConflictsContext {
	return self.c.Contexts().MergeConflicts
}

func (self *MergeConflictsHelper) Render() {
	content := self.context().GetContentToRender()

	var task types.UpdateTask
	if self.context().IsUserScrolling() {
		task = types.NewRenderStringWithoutScrollTask(content)
	} else {
		originY := self.context().GetOriginY()
		task = types.NewRenderStringWithScrollTask(content, 0, originY)
	}

	self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().MergeConflicts,
		Main: &types.ViewUpdateOpts{
			Task: task,
		},
	})
}

func (self *MergeConflictsHelper) RefreshMergeState() error {
	if self.c.Context().Current().GetKey() != context.MERGE_CONFLICTS_CONTEXT_KEY {
		return nil
	}

	hasConflicts, err := self.SetConflictsAndRender()
	if err != nil {
		return err
	}

	if !hasConflicts {
		self.EscapeMerge()
	}

	return nil
}
