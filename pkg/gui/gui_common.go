package gui

import (
	"errors"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// hacking this by including the gui struct for now until we split more things out
type guiCommon struct {
	gui *Gui
	types.IPopupHandler
}

var _ types.IGuiCommon = &guiCommon{}

func (self *guiCommon) LogAction(msg string) {
	self.gui.LogAction(msg)
}

func (self *guiCommon) LogCommand(cmdStr string, isCommandLine bool) {
	self.gui.LogCommand(cmdStr, isCommandLine)
}

func (self *guiCommon) Refresh(opts types.RefreshOptions) error {
	return self.gui.Refresh(opts)
}

func (self *guiCommon) PostRefreshUpdate(context types.Context) error {
	return self.gui.postRefreshUpdate(context)
}

func (self *guiCommon) RunSubprocessAndRefresh(cmdObj oscommands.ICmdObj) error {
	return self.gui.runSubprocessWithSuspenseAndRefresh(cmdObj)
}

func (self *guiCommon) RunSubprocess(cmdObj oscommands.ICmdObj) (bool, error) {
	return self.gui.runSubprocessWithSuspense(cmdObj)
}

func (self *guiCommon) PushContext(context types.Context, opts ...types.OnFocusOpts) error {
	singleOpts := types.OnFocusOpts{}
	if len(opts) > 0 {
		// using triple dot but you should only ever pass one of these opt structs
		if len(opts) > 1 {
			return errors.New("cannot pass multiple opts to pushContext")
		}

		singleOpts = opts[0]
	}

	return self.gui.pushContext(context, singleOpts)
}

func (self *guiCommon) PopContext() error {
	return self.gui.popContext()
}

func (self *guiCommon) RemoveContexts(contexts []types.Context) error {
	return self.gui.removeContexts(contexts)
}

func (self *guiCommon) CurrentContext() types.Context {
	return self.gui.currentContext()
}

func (self *guiCommon) CurrentStaticContext() types.Context {
	return self.gui.currentStaticContext()
}

func (self *guiCommon) IsCurrentContext(c types.Context) bool {
	return self.CurrentContext().GetKey() == c.GetKey()
}

func (self *guiCommon) ActivateContext(context types.Context) error {
	return self.gui.activateContext(context, types.OnFocusOpts{})
}

func (self *guiCommon) GetAppState() *config.AppState {
	return self.gui.Config.GetAppState()
}

func (self *guiCommon) SaveAppState() error {
	return self.gui.Config.SaveAppState()
}

func (self *guiCommon) Render() {
	self.gui.render()
}

func (self *guiCommon) OpenSearch() {
	_ = self.gui.handleOpenSearch(self.gui.currentViewName())
}

func (self *guiCommon) OnUIThread(f func() error) {
	self.gui.onUIThread(f)
}

func (self *guiCommon) RenderToMainViews(opts types.RefreshMainOpts) error {
	return self.gui.refreshMainViews(opts)
}

func (self *guiCommon) MainViewPairs() types.MainViewPairs {
	return types.MainViewPairs{
		Normal:         self.gui.normalMainContextPair(),
		Staging:        self.gui.stagingMainContextPair(),
		PatchBuilding:  self.gui.patchBuildingMainContextPair(),
		MergeConflicts: self.gui.mergingMainContextPair(),
	}
}
