package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
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

func (self *guiCommon) PostRefreshUpdate(context types.Context) {
	self.gui.postRefreshUpdate(context)
}

func (self *guiCommon) RunSubprocessAndRefresh(cmdObj oscommands.ICmdObj) error {
	return self.gui.runSubprocessWithSuspenseAndRefresh(cmdObj)
}

func (self *guiCommon) RunSubprocess(cmdObj oscommands.ICmdObj) (bool, error) {
	return self.gui.runSubprocessWithSuspense(cmdObj)
}

func (self *guiCommon) Context() types.IContextMgr {
	return self.gui.State.ContextMgr
}

func (self *guiCommon) ContextForKey(key types.ContextKey) types.Context {
	return self.gui.State.ContextMgr.ContextForKey(key)
}

func (self *guiCommon) GetAppState() *config.AppState {
	return self.gui.Config.GetAppState()
}

func (self *guiCommon) SaveAppState() error {
	return self.gui.Config.SaveAppState()
}

func (self *guiCommon) SaveAppStateAndLogError() {
	if err := self.gui.Config.SaveAppState(); err != nil {
		self.gui.Log.Errorf("error when saving app state: %v", err)
	}
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

func (self *guiCommon) OnWorker(f func(gocui.Task) error) {
	self.gui.onWorker(f)
}

func (self *guiCommon) RenderToMainViews(opts types.RefreshMainOpts) {
	self.gui.refreshMainViews(opts)
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

func (self *guiCommon) CallKeybindingHandler(binding *types.Binding) error {
	return self.gui.callKeybindingHandler(binding)
}

func (self *guiCommon) ResetKeybindings() error {
	return self.gui.resetKeybindings()
}

func (self *guiCommon) IsAnyModeActive() bool {
	return self.gui.helpers.Mode.IsAnyModeActive()
}

func (self *guiCommon) GetInitialKeybindingsWithCustomCommands() ([]*types.Binding, []*gocui.ViewMouseBinding) {
	return self.gui.GetInitialKeybindingsWithCustomCommands()
}

func (self *guiCommon) AfterLayout(f func() error) {
	self.gui.afterLayout(f)
}

func (self *guiCommon) RunningIntegrationTest() bool {
	return self.gui.integrationTest != nil
}

func (self *guiCommon) InDemo() bool {
	return self.gui.integrationTest != nil && self.gui.integrationTest.IsDemo()
}

func (self *guiCommon) WithInlineStatus(item types.HasUrn, operation types.ItemOperation, contextKey types.ContextKey, f func(gocui.Task) error) error {
	self.gui.helpers.InlineStatus.WithInlineStatus(helpers.InlineStatusOpts{Item: item, Operation: operation, ContextKey: contextKey}, f)
	return nil
}
