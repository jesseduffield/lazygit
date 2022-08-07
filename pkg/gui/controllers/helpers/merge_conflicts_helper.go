package helpers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MergeConflictsHelper struct {
	c        *types.HelperCommon
	contexts *context.ContextTree
	git      *commands.GitCommand
}

func NewMergeConflictsHelper(
	c *types.HelperCommon,
	contexts *context.ContextTree,
	git *commands.GitCommand,
) *MergeConflictsHelper {
	return &MergeConflictsHelper{
		c:        c,
		contexts: contexts,
		git:      git,
	}
}

func (self *MergeConflictsHelper) GetMergingOptions() map[string]string {
	keybindingConfig := self.c.UserConfig.Keybinding

	return map[string]string{
		fmt.Sprintf("%s %s", keybindings.Label(keybindingConfig.Universal.PrevItem), keybindings.Label(keybindingConfig.Universal.NextItem)):   self.c.Tr.LcSelectHunk,
		fmt.Sprintf("%s %s", keybindings.Label(keybindingConfig.Universal.PrevBlock), keybindings.Label(keybindingConfig.Universal.NextBlock)): self.c.Tr.LcNavigateConflicts,
		keybindings.Label(keybindingConfig.Universal.Select):   self.c.Tr.LcPickHunk,
		keybindings.Label(keybindingConfig.Main.PickBothHunks): self.c.Tr.LcPickAllHunks,
		keybindings.Label(keybindingConfig.Universal.Undo):     self.c.Tr.LcUndo,
	}
}

func (self *MergeConflictsHelper) SetMergeState(path string) (bool, error) {
	self.context().GetMutex().Lock()
	defer self.context().GetMutex().Unlock()

	return self.setMergeStateWithoutLock(path)
}

func (self *MergeConflictsHelper) setMergeStateWithoutLock(path string) (bool, error) {
	content, err := self.git.File.Cat(path)
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
		return self.c.PushContext(self.contexts.Files)
	})
	return nil
}

func (self *MergeConflictsHelper) SetConflictsAndRender(path string, isFocused bool) (bool, error) {
	hasConflicts, err := self.setMergeStateWithoutLock(path)
	if err != nil {
		return false, err
	}

	if hasConflicts {
		return true, self.context().Render(isFocused)
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

	return self.c.PushContext(self.contexts.MergeConflicts)
}

func (self *MergeConflictsHelper) context() *context.MergeConflictsContext {
	return self.contexts.MergeConflicts
}
