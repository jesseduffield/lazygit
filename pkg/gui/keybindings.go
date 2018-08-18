package gui

import "github.com/jesseduffield/gocui"

// Binding - a keybinding mapping a key and modifier to a handler. The keypress
// is only handled if the given view has focus, or handled globally if the view
// is ""
type Binding struct {
	ViewName string
	Handler  func(*gocui.Gui, *gocui.View) error
	Key      interface{} // FIXME: find out how to get `gocui.Key | rune`
	Modifier gocui.Modifier
}

func (gui *Gui) keybindings(g *gocui.Gui) error {
	bindings := []Binding{
		{ViewName: "", Key: 'q', Modifier: gocui.ModNone, Handler: gui.quit},
		{ViewName: "", Key: gocui.KeyCtrlC, Modifier: gocui.ModNone, Handler: gui.quit},
		{ViewName: "", Key: gocui.KeyPgup, Modifier: gocui.ModNone, Handler: gui.scrollUpMain},
		{ViewName: "", Key: gocui.KeyPgdn, Modifier: gocui.ModNone, Handler: gui.scrollDownMain},
		{ViewName: "", Key: gocui.KeyCtrlU, Modifier: gocui.ModNone, Handler: gui.scrollUpMain},
		{ViewName: "", Key: gocui.KeyCtrlD, Modifier: gocui.ModNone, Handler: gui.scrollDownMain},
		{ViewName: "", Key: 'P', Modifier: gocui.ModNone, Handler: gui.pushFiles},
		{ViewName: "", Key: 'p', Modifier: gocui.ModNone, Handler: gui.pullFiles},
		{ViewName: "", Key: 'R', Modifier: gocui.ModNone, Handler: gui.handleRefresh},
		{ViewName: "status", Key: 'e', Modifier: gocui.ModNone, Handler: gui.handleEditConfig},
		{ViewName: "status", Key: 'o', Modifier: gocui.ModNone, Handler: gui.handleOpenConfig},
		{ViewName: "files", Key: 'c', Modifier: gocui.ModNone, Handler: gui.handleCommitPress},
		{ViewName: "files", Key: 'C', Modifier: gocui.ModNone, Handler: gui.handleCommitEditorPress},
		{ViewName: "files", Key: gocui.KeySpace, Modifier: gocui.ModNone, Handler: gui.handleFilePress},
		{ViewName: "files", Key: 'd', Modifier: gocui.ModNone, Handler: gui.handleFileRemove},
		{ViewName: "files", Key: 'm', Modifier: gocui.ModNone, Handler: gui.handleSwitchToMerge},
		{ViewName: "files", Key: 'e', Modifier: gocui.ModNone, Handler: gui.handleFileEdit},
		{ViewName: "files", Key: 'o', Modifier: gocui.ModNone, Handler: gui.handleFileOpen},
		{ViewName: "files", Key: 's', Modifier: gocui.ModNone, Handler: gui.handleSublimeFileOpen},
		{ViewName: "files", Key: 'v', Modifier: gocui.ModNone, Handler: gui.handleVsCodeFileOpen},
		{ViewName: "files", Key: 'i', Modifier: gocui.ModNone, Handler: gui.handleIgnoreFile},
		{ViewName: "files", Key: 'r', Modifier: gocui.ModNone, Handler: gui.handleRefreshFiles},
		{ViewName: "files", Key: 'S', Modifier: gocui.ModNone, Handler: gui.handleStashSave},
		{ViewName: "files", Key: 'a', Modifier: gocui.ModNone, Handler: gui.handleAbortMerge},
		{ViewName: "files", Key: 't', Modifier: gocui.ModNone, Handler: gui.handleAddPatch},
		{ViewName: "files", Key: 'D', Modifier: gocui.ModNone, Handler: gui.handleResetHard},
		{ViewName: "main", Key: gocui.KeyEsc, Modifier: gocui.ModNone, Handler: gui.handleEscapeMerge},
		{ViewName: "main", Key: gocui.KeySpace, Modifier: gocui.ModNone, Handler: gui.handlePickHunk},
		{ViewName: "main", Key: 'b', Modifier: gocui.ModNone, Handler: gui.handlePickBothHunks},
		{ViewName: "main", Key: gocui.KeyArrowLeft, Modifier: gocui.ModNone, Handler: gui.handleSelectPrevConflict},
		{ViewName: "main", Key: gocui.KeyArrowRight, Modifier: gocui.ModNone, Handler: gui.handleSelectNextConflict},
		{ViewName: "main", Key: gocui.KeyArrowUp, Modifier: gocui.ModNone, Handler: gui.handleSelectTop},
		{ViewName: "main", Key: gocui.KeyArrowDown, Modifier: gocui.ModNone, Handler: gui.handleSelectBottom},
		{ViewName: "main", Key: 'h', Modifier: gocui.ModNone, Handler: gui.handleSelectPrevConflict},
		{ViewName: "main", Key: 'l', Modifier: gocui.ModNone, Handler: gui.handleSelectNextConflict},
		{ViewName: "main", Key: 'k', Modifier: gocui.ModNone, Handler: gui.handleSelectTop},
		{ViewName: "main", Key: 'j', Modifier: gocui.ModNone, Handler: gui.handleSelectBottom},
		{ViewName: "main", Key: 'z', Modifier: gocui.ModNone, Handler: gui.handlePopFileSnapshot},
		{ViewName: "branches", Key: gocui.KeySpace, Modifier: gocui.ModNone, Handler: gui.handleBranchPress},
		{ViewName: "branches", Key: 'c', Modifier: gocui.ModNone, Handler: gui.handleCheckoutByName},
		{ViewName: "branches", Key: 'F', Modifier: gocui.ModNone, Handler: gui.handleForceCheckout},
		{ViewName: "branches", Key: 'n', Modifier: gocui.ModNone, Handler: gui.handleNewBranch},
		{ViewName: "branches", Key: 'd', Modifier: gocui.ModNone, Handler: gui.handleDeleteBranch},
		{ViewName: "branches", Key: 'm', Modifier: gocui.ModNone, Handler: gui.handleMerge},
		{ViewName: "commits", Key: 's', Modifier: gocui.ModNone, Handler: gui.handleCommitSquashDown},
		{ViewName: "commits", Key: 'r', Modifier: gocui.ModNone, Handler: gui.handleRenameCommit},
		{ViewName: "commits", Key: 'g', Modifier: gocui.ModNone, Handler: gui.handleResetToCommit},
		{ViewName: "commits", Key: 'f', Modifier: gocui.ModNone, Handler: gui.handleCommitFixup},
		{ViewName: "stash", Key: gocui.KeySpace, Modifier: gocui.ModNone, Handler: gui.handleStashApply},
		{ViewName: "stash", Key: 'g', Modifier: gocui.ModNone, Handler: gui.handleStashPop},
		{ViewName: "stash", Key: 'd', Modifier: gocui.ModNone, Handler: gui.handleStashDrop},
		{ViewName: "commitMessage", Key: gocui.KeyEnter, Modifier: gocui.ModNone, Handler: gui.handleCommitConfirm},
		{ViewName: "commitMessage", Key: gocui.KeyEsc, Modifier: gocui.ModNone, Handler: gui.handleCommitClose},
		{ViewName: "commitMessage", Key: gocui.KeyTab, Modifier: gocui.ModNone, Handler: gui.handleNewlineCommitMessage},
	}

	// Would make these keybindings global but that interferes with editing
	// input in the confirmation panel
	for _, viewName := range []string{"status", "files", "branches", "commits", "stash"} {
		bindings = append(bindings, []Binding{
			{ViewName: viewName, Key: gocui.KeyTab, Modifier: gocui.ModNone, Handler: gui.nextView},
			{ViewName: viewName, Key: gocui.KeyArrowLeft, Modifier: gocui.ModNone, Handler: gui.previousView},
			{ViewName: viewName, Key: gocui.KeyArrowRight, Modifier: gocui.ModNone, Handler: gui.nextView},
			{ViewName: viewName, Key: gocui.KeyArrowUp, Modifier: gocui.ModNone, Handler: gui.cursorUp},
			{ViewName: viewName, Key: gocui.KeyArrowDown, Modifier: gocui.ModNone, Handler: gui.cursorDown},
			{ViewName: viewName, Key: 'h', Modifier: gocui.ModNone, Handler: gui.previousView},
			{ViewName: viewName, Key: 'l', Modifier: gocui.ModNone, Handler: gui.nextView},
			{ViewName: viewName, Key: 'k', Modifier: gocui.ModNone, Handler: gui.cursorUp},
			{ViewName: viewName, Key: 'j', Modifier: gocui.ModNone, Handler: gui.cursorDown},
		}...)
	}

	for _, binding := range bindings {
		if err := g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}
	return nil
}
