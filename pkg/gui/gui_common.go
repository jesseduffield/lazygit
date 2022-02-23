package gui

import (
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
	return self.gui.pushContext(context, opts...)
}

func (self *guiCommon) PopContext() error {
	return self.gui.returnFromContext()
}

func (self *guiCommon) CurrentContext() types.Context {
	return self.gui.currentContext()
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
	self.gui.OnUIThread(f)
}
