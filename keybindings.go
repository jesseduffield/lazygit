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
		Binding{ViewName: "", Key: 'q', Modifier: gocui.ModNone, Handler: quit},
		Binding{ViewName: "", Key: gocui.KeyCtrlC, Modifier: gocui.ModNone, Handler: quit},
		Binding{ViewName: "", Key: gocui.KeyPgup, Modifier: gocui.ModNone, Handler: scrollUpMain},
		Binding{ViewName: "", Key: gocui.KeyPgdn, Modifier: gocui.ModNone, Handler: scrollDownMain},
		Binding{ViewName: "", Key: 'P', Modifier: gocui.ModNone, Handler: pushFiles},
		Binding{ViewName: "", Key: 'p', Modifier: gocui.ModNone, Handler: pullFiles},
		Binding{ViewName: "", Key: 'R', Modifier: gocui.ModNone, Handler: handleRefresh},
		Binding{ViewName: "files", Key: 'c', Modifier: gocui.ModNone, Handler: handleCommitPress},
		Binding{ViewName: "files", Key: gocui.KeySpace, Modifier: gocui.ModNone, Handler: handleFilePress},
		Binding{ViewName: "files", Key: 'd', Modifier: gocui.ModNone, Handler: handleFileRemove},
		Binding{ViewName: "files", Key: 'm', Modifier: gocui.ModNone, Handler: handleSwitchToMerge},
		Binding{ViewName: "files", Key: 'e', Modifier: gocui.ModNone, Handler: handleFileEdit},
		Binding{ViewName: "files", Key: 'o', Modifier: gocui.ModNone, Handler: handleFileOpen},
		Binding{ViewName: "files", Key: 's', Modifier: gocui.ModNone, Handler: handleSublimeFileOpen},
		Binding{ViewName: "files", Key: 'v', Modifier: gocui.ModNone, Handler: handleVsCodeFileOpen},
		Binding{ViewName: "files", Key: 'i', Modifier: gocui.ModNone, Handler: handleIgnoreFile},
		Binding{ViewName: "files", Key: 'r', Modifier: gocui.ModNone, Handler: handleRefreshFiles},
		Binding{ViewName: "files", Key: 'S', Modifier: gocui.ModNone, Handler: handleStashSave},
		Binding{ViewName: "files", Key: 'a', Modifier: gocui.ModNone, Handler: handleAbortMerge},
		Binding{ViewName: "files", Key: 't', Modifier: gocui.ModNone, Handler: handleAddPatch},
		Binding{ViewName: "files", Key: 'D', Modifier: gocui.ModNone, Handler: handleResetHard},
		Binding{ViewName: "main", Key: gocui.KeyEsc, Modifier: gocui.ModNone, Handler: handleEscapeMerge},
		Binding{ViewName: "main", Key: gocui.KeySpace, Modifier: gocui.ModNone, Handler: handlePickHunk},
		Binding{ViewName: "main", Key: 'b', Modifier: gocui.ModNone, Handler: handlePickBothHunks},
		Binding{ViewName: "main", Key: gocui.KeyArrowLeft, Modifier: gocui.ModNone, Handler: handleSelectPrevConflict},
		Binding{ViewName: "main", Key: gocui.KeyArrowRight, Modifier: gocui.ModNone, Handler: handleSelectNextConflict},
		Binding{ViewName: "main", Key: gocui.KeyArrowUp, Modifier: gocui.ModNone, Handler: handleSelectTop},
		Binding{ViewName: "main", Key: gocui.KeyArrowDown, Modifier: gocui.ModNone, Handler: handleSelectBottom},
		Binding{ViewName: "main", Key: 'h', Modifier: gocui.ModNone, Handler: handleSelectPrevConflict},
		Binding{ViewName: "main", Key: 'l', Modifier: gocui.ModNone, Handler: handleSelectNextConflict},
		Binding{ViewName: "main", Key: 'k', Modifier: gocui.ModNone, Handler: handleSelectTop},
		Binding{ViewName: "main", Key: 'j', Modifier: gocui.ModNone, Handler: handleSelectBottom},
		Binding{ViewName: "main", Key: 'z', Modifier: gocui.ModNone, Handler: handlePopFileSnapshot},
		Binding{ViewName: "branches", Key: gocui.KeySpace, Modifier: gocui.ModNone, Handler: handleBranchPress},
		Binding{ViewName: "branches", Key: 'c', Modifier: gocui.ModNone, Handler: handleCheckoutByName},
		Binding{ViewName: "branches", Key: 'F', Modifier: gocui.ModNone, Handler: handleForceCheckout},
		Binding{ViewName: "branches", Key: 'n', Modifier: gocui.ModNone, Handler: handleNewBranch},
		Binding{ViewName: "branches", Key: 'm', Modifier: gocui.ModNone, Handler: handleMerge},
		Binding{ViewName: "commits", Key: 's', Modifier: gocui.ModNone, Handler: handleCommitSquashDown},
		Binding{ViewName: "commits", Key: 'r', Modifier: gocui.ModNone, Handler: handleRenameCommit},
		Binding{ViewName: "commits", Key: 'g', Modifier: gocui.ModNone, Handler: handleResetToCommit},
		Binding{ViewName: "stash", Key: gocui.KeySpace, Modifier: gocui.ModNone, Handler: handleStashApply},
		Binding{ViewName: "stash", Key: 'g', Modifier: gocui.ModNone, Handler: handleStashPop},
		Binding{ViewName: "stash", Key: 'd', Modifier: gocui.ModNone, Handler: handleStashDrop},
	}

	// Would make these keybindings global but that interferes with editing
	// input in the confirmation panel
	for _, viewName := range []string{"files", "branches", "commits", "stash"} {
		bindings = append(bindings, []Binding{
			Binding{ViewName: viewName, Key: gocui.KeyTab, Modifier: gocui.ModNone, Handler: nextView},
			Binding{ViewName: viewName, Key: gocui.KeyArrowLeft, Modifier: gocui.ModNone, Handler: previousView},
			Binding{ViewName: viewName, Key: gocui.KeyArrowRight, Modifier: gocui.ModNone, Handler: nextView},
			Binding{ViewName: viewName, Key: gocui.KeyArrowUp, Modifier: gocui.ModNone, Handler: cursorUp},
			Binding{ViewName: viewName, Key: gocui.KeyArrowDown, Modifier: gocui.ModNone, Handler: cursorDown},
			Binding{ViewName: viewName, Key: 'h', Modifier: gocui.ModNone, Handler: previousView},
			Binding{ViewName: viewName, Key: 'l', Modifier: gocui.ModNone, Handler: nextView},
			Binding{ViewName: viewName, Key: 'k', Modifier: gocui.ModNone, Handler: cursorUp},
			Binding{ViewName: viewName, Key: 'j', Modifier: gocui.ModNone, Handler: cursorDown},
		}...)
	}

	for _, binding := range bindings {
		if err := g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}
	return nil
}
