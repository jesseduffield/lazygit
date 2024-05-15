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

func (self *MergeConflictsHelper) EscapeMerge() error {
	self.resetMergeState()

	// doing this in separate UI thread so that we're not still holding the lock by the time refresh the file
	self.c.OnUIThread(func() error {
		// There is a race condition here: refreshing the files scope can trigger the
		// confirmation context to be pushed if all conflicts are resolved (prompting
		// to continue the merge/rebase. In that case, we don't want to then push the
		// files context over it.
		// So long as both places call OnUIThread, we're fine.
		if self.c.IsCurrentContext(self.c.Contexts().MergeConflicts) {
			return self.c.PushContext(self.c.Contexts().Files)
		}
		return nil
	})
	return nil
}

func (self *MergeConflictsHelper) SetConflictsAndRender(path string) (bool, error) {
	hasConflicts, err := self.setMergeStateWithoutLock(path)
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

	return self.c.PushContext(self.c.Contexts().MergeConflicts)
}

func (self *MergeConflictsHelper) context() *context.MergeConflictsContext {
	return self.c.Contexts().MergeConflicts
}

func (self *MergeConflictsHelper) Render() error {
	content := self.context().GetContentToRender()

	var task types.UpdateTask
	if self.context().IsUserScrolling() {
		task = types.NewRenderStringWithoutScrollTask(content)
	} else {
		originY := self.context().GetOriginY()
		task = types.NewRenderStringWithScrollTask(content, 0, originY)
	}

	return self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().MergeConflicts,
		Main: &types.ViewUpdateOpts{
			Task: task,
		},
	})
}

func (self *MergeConflictsHelper) RefreshMergeState() error {
	self.c.Contexts().MergeConflicts.GetMutex().Lock()
	defer self.c.Contexts().MergeConflicts.GetMutex().Unlock()

	if self.c.CurrentContext().GetKey() != context.MERGE_CONFLICTS_CONTEXT_KEY {
		return nil
	}

	hasConflicts, err := self.SetConflictsAndRender(self.c.Contexts().MergeConflicts.GetState().GetPath())
	if err != nil {
		return err
	}

	if !hasConflicts {
		return self.EscapeMerge()
	}

	return nil
}
