package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
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
	return self.gui.helpers.Refresh.Refresh(opts)
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
	return self.gui.State.ContextMgr.Push(context, opts...)
}

func (self *guiCommon) PopContext() error {
	return self.gui.State.ContextMgr.Pop()
}

func (self *guiCommon) ReplaceContext(context types.Context) error {
	return self.gui.State.ContextMgr.Replace(context)
}

func (self *guiCommon) RemoveContexts(contexts []types.Context) error {
	return self.gui.State.ContextMgr.RemoveContexts(contexts)
}

func (self *guiCommon) CurrentContext() types.Context {
	return self.gui.State.ContextMgr.Current()
}

func (self *guiCommon) CurrentStaticContext() types.Context {
	return self.gui.State.ContextMgr.CurrentStatic()
}

func (self *guiCommon) CurrentSideContext() types.Context {
	return self.gui.State.ContextMgr.CurrentSide()
}

func (self *guiCommon) IsCurrentContext(c types.Context) bool {
	return self.gui.State.ContextMgr.IsCurrent(c)
}

func (self *guiCommon) Context() types.IContextMgr {
	return self.gui.State.ContextMgr
}

func (self *guiCommon) ContextForKey(key types.ContextKey) types.Context {
	return self.gui.State.ContextMgr.ContextForKey(key)
}

func (self *guiCommon) ActivateContext(context types.Context) error {
	return self.gui.State.ContextMgr.ActivateContext(context, types.OnFocusOpts{})
}

func (self *guiCommon) GetAppState() *config.AppState {
	return self.gui.Config.GetAppState()
}

func (self *guiCommon) SaveAppState() error {
	return self.gui.Config.SaveAppState()
}

func (self *guiCommon) GetConfig() config.AppConfigurer {
	return self.gui.Config
}

func (self *guiCommon) ResetViewOrigin(view *gocui.View) {
	self.gui.resetViewOrigin(view)
}

func (self *guiCommon) SetViewContent(view *gocui.View, content string) {
	self.gui.setViewContent(view, content)
}

func (self *guiCommon) Render() {
	self.gui.render()
}

func (self *guiCommon) Views() types.Views {
	return self.gui.Views
}

func (self *guiCommon) Git() *commands.GitCommand {
	return self.gui.git
}

func (self *guiCommon) OS() *oscommands.OSCommand {
	return self.gui.os
}

func (self *guiCommon) Modes() *types.Modes {
	return self.gui.State.Modes
}

func (self *guiCommon) Model() *types.Model {
	return self.gui.State.Model
}

func (self *guiCommon) Mutexes() types.Mutexes {
	return self.gui.Mutexes
}

func (self *guiCommon) GocuiGui() *gocui.Gui {
	return self.gui.g
}

func (self *guiCommon) OnUIThread(f func() error) {
	self.gui.onUIThread(f)
}

func (self *guiCommon) OnWorker(f func(gocui.Task)) {
	self.gui.onWorker(f)
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

func (self *guiCommon) State() types.IStateAccessor {
	return self.gui.stateAccessor
}

func (self *guiCommon) KeybindingsOpts() types.KeybindingsOpts {
	return self.gui.keybindingOpts()
}

func (self *guiCommon) IsAnyModeActive() bool {
	return self.gui.helpers.Mode.IsAnyModeActive()
}

func (self *guiCommon) GetInitialKeybindingsWithCustomCommands() ([]*types.Binding, []*gocui.ViewMouseBinding) {
	return self.gui.GetInitialKeybindingsWithCustomCommands()
}

func (self *guiCommon) AfterLayout(f func() error) {
	select {
	case self.gui.afterLayoutFuncs <- f:
	default:
		// hopefully this never happens
		self.gui.c.Log.Error("afterLayoutFuncs channel is full, skipping function")
	}
}

func (self *guiCommon) InDemo() bool {
	return self.gui.integrationTest != nil && self.gui.integrationTest.IsDemo()
}
