package main

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

func keybindings(g *gocui.Gui) error {
	bindings := []Binding{
		{ViewName: "", Key: 'q', Modifier: gocui.ModNone, Handler: quit},
		{ViewName: "", Key: gocui.KeyCtrlC, Modifier: gocui.ModNone, Handler: quit},
		{ViewName: "", Key: gocui.KeyCtrlU, Modifier: gocui.ModNone, Handler: scrollUpMain},
		{ViewName: "", Key: gocui.KeyCtrlD, Modifier: gocui.ModNone, Handler: scrollDownMain},
		{ViewName: "", Key: gocui.KeyPgup, Modifier: gocui.ModNone, Handler: scrollUpMain},
		{ViewName: "", Key: gocui.KeyPgdn, Modifier: gocui.ModNone, Handler: scrollDownMain},
		{ViewName: "", Key: 'P', Modifier: gocui.ModNone, Handler: pushFiles},
		{ViewName: "", Key: 'p', Modifier: gocui.ModNone, Handler: pullFiles},
		{ViewName: "", Key: 'R', Modifier: gocui.ModNone, Handler: handleRefresh},
		{ViewName: "files", Key: 'c', Modifier: gocui.ModNone, Handler: handleCommitPress},
		{ViewName: "files", Key: 'C', Modifier: gocui.ModNone, Handler: handleCommitEditorPress},
		{ViewName: "files", Key: gocui.KeySpace, Modifier: gocui.ModNone, Handler: handleFilePress},
		{ViewName: "files", Key: 'd', Modifier: gocui.ModNone, Handler: handleFileRemove},
		{ViewName: "files", Key: 'm', Modifier: gocui.ModNone, Handler: handleSwitchToMerge},
		{ViewName: "files", Key: 'e', Modifier: gocui.ModNone, Handler: handleFileEdit},
		{ViewName: "files", Key: 'o', Modifier: gocui.ModNone, Handler: handleFileOpen},
		{ViewName: "files", Key: 's', Modifier: gocui.ModNone, Handler: handleSublimeFileOpen},
		{ViewName: "files", Key: 'v', Modifier: gocui.ModNone, Handler: handleVsCodeFileOpen},
		{ViewName: "files", Key: 'i', Modifier: gocui.ModNone, Handler: handleIgnoreFile},
		{ViewName: "files", Key: 'r', Modifier: gocui.ModNone, Handler: handleRefreshFiles},
		{ViewName: "files", Key: 'S', Modifier: gocui.ModNone, Handler: handleStashSave},
		{ViewName: "files", Key: 'a', Modifier: gocui.ModNone, Handler: handleAbortMerge},
		{ViewName: "files", Key: 't', Modifier: gocui.ModNone, Handler: handleAddPatch},
		{ViewName: "files", Key: 'D', Modifier: gocui.ModNone, Handler: handleResetHard},
		{ViewName: "main", Key: gocui.KeyEsc, Modifier: gocui.ModNone, Handler: handleEscapeMerge},
		{ViewName: "main", Key: gocui.KeySpace, Modifier: gocui.ModNone, Handler: handlePickHunk},
		{ViewName: "main", Key: 'b', Modifier: gocui.ModNone, Handler: handlePickBothHunks},
		{ViewName: "main", Key: gocui.KeyArrowLeft, Modifier: gocui.ModNone, Handler: handleSelectPrevConflict},
		{ViewName: "main", Key: gocui.KeyArrowRight, Modifier: gocui.ModNone, Handler: handleSelectNextConflict},
		{ViewName: "main", Key: gocui.KeyArrowUp, Modifier: gocui.ModNone, Handler: handleSelectTop},
		{ViewName: "main", Key: gocui.KeyArrowDown, Modifier: gocui.ModNone, Handler: handleSelectBottom},
		{ViewName: "main", Key: 'h', Modifier: gocui.ModNone, Handler: handleSelectPrevConflict},
		{ViewName: "main", Key: 'l', Modifier: gocui.ModNone, Handler: handleSelectNextConflict},
		{ViewName: "main", Key: 'k', Modifier: gocui.ModNone, Handler: handleSelectTop},
		{ViewName: "main", Key: 'j', Modifier: gocui.ModNone, Handler: handleSelectBottom},
		{ViewName: "main", Key: 'z', Modifier: gocui.ModNone, Handler: handlePopFileSnapshot},
		{ViewName: "branches", Key: gocui.KeySpace, Modifier: gocui.ModNone, Handler: handleBranchPress},
		{ViewName: "branches", Key: 'c', Modifier: gocui.ModNone, Handler: handleCheckoutByName},
		{ViewName: "branches", Key: 'F', Modifier: gocui.ModNone, Handler: handleForceCheckout},
		{ViewName: "branches", Key: 'n', Modifier: gocui.ModNone, Handler: handleNewBranch},
		{ViewName: "branches", Key: 'd', Modifier: gocui.ModNone, Handler: handleDeleteBranch},
		{ViewName: "branches", Key: 'm', Modifier: gocui.ModNone, Handler: handleMerge},
		{ViewName: "commits", Key: 's', Modifier: gocui.ModNone, Handler: handleCommitSquashDown},
		{ViewName: "commits", Key: 'r', Modifier: gocui.ModNone, Handler: handleRenameCommit},
		{ViewName: "commits", Key: 'g', Modifier: gocui.ModNone, Handler: handleResetToCommit},
		{ViewName: "commits", Key: 'f', Modifier: gocui.ModNone, Handler: handleCommitFixup},
		{ViewName: "stash", Key: gocui.KeySpace, Modifier: gocui.ModNone, Handler: handleStashApply},
		{ViewName: "stash", Key: 'g', Modifier: gocui.ModNone, Handler: handleStashPop},
		{ViewName: "stash", Key: 'd', Modifier: gocui.ModNone, Handler: handleStashDrop},
		{ViewName: "commitMessage", Key: gocui.KeyEnter, Modifier: gocui.ModNone, Handler: handleCommitConfirm},
		{ViewName: "commitMessage", Key: gocui.KeyEsc, Modifier: gocui.ModNone, Handler: handleCommitClose},
		{ViewName: "commitMessage", Key: gocui.KeyTab, Modifier: gocui.ModNone, Handler: handleNewlineCommitMessage},
	}

	// Would make these keybindings global but that interferes with editing
	// input in the confirmation panel
	for _, viewName := range []string{"files", "branches", "commits", "stash"} {
		bindings = append(bindings, []Binding{
			{ViewName: viewName, Key: gocui.KeyTab, Modifier: gocui.ModNone, Handler: nextView},
			{ViewName: viewName, Key: gocui.KeyArrowLeft, Modifier: gocui.ModNone, Handler: previousView},
			{ViewName: viewName, Key: gocui.KeyArrowRight, Modifier: gocui.ModNone, Handler: nextView},
			{ViewName: viewName, Key: gocui.KeyArrowUp, Modifier: gocui.ModNone, Handler: cursorUp},
			{ViewName: viewName, Key: gocui.KeyArrowDown, Modifier: gocui.ModNone, Handler: cursorDown},
			{ViewName: viewName, Key: 'h', Modifier: gocui.ModNone, Handler: previousView},
			{ViewName: viewName, Key: 'l', Modifier: gocui.ModNone, Handler: nextView},
			{ViewName: viewName, Key: 'k', Modifier: gocui.ModNone, Handler: cursorUp},
			{ViewName: viewName, Key: 'j', Modifier: gocui.ModNone, Handler: cursorDown},
		}...)
	}

	for _, binding := range bindings {
		if err := g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}
	return nil
}
