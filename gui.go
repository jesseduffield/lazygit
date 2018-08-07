package main

import (

	// "io"
	// "io/ioutil"

	"log"
	"runtime"
	"strings"
	"time"

	// "strings"

	"github.com/golang-collections/collections/stack"
	"github.com/jesseduffield/gocui"
)

// OverlappingEdges determines if panel edges overlap
var OverlappingEdges = false

type stateType struct {
	GitFiles          []GitFile
	Branches          []Branch
	Commits           []Commit
	StashEntries      []StashEntry
	PreviousView      string
	HasMergeConflicts bool
	ConflictIndex     int
	ConflictTop       bool
	Conflicts         []conflict
	EditHistory       *stack.Stack
}

type conflict struct {
	start  int
	middle int
	end    int
}

var state = stateType{
	GitFiles:      make([]GitFile, 0),
	PreviousView:  "files",
	Commits:       make([]Commit, 0),
	StashEntries:  make([]StashEntry, 0),
	ConflictIndex: 0,
	ConflictTop:   true,
	Conflicts:     make([]conflict, 0),
	EditHistory:   stack.New(),
}

func scrollUpMain(g *gocui.Gui, v *gocui.View) error {
	mainView, _ := g.View("main")
	ox, oy := mainView.Origin()
	if oy >= 1 {
		return mainView.SetOrigin(ox, oy-1)
	}
	return nil
}

func scrollDownMain(g *gocui.Gui, v *gocui.View) error {
	mainView, _ := g.View("main")
	ox, oy := mainView.Origin()
	if oy < len(mainView.BufferLines()) {
		return mainView.SetOrigin(ox, oy+1)
	}
	return nil
}

func handleRefresh(g *gocui.Gui, v *gocui.View) error {
	return refreshSidePanels(g)
}

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
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, previousView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, scrollUpMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, scrollDownMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'P', gocui.ModNone, pushFiles); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'p', gocui.ModNone, pullFiles); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'R', gocui.ModNone, handleRefresh); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", 'c', gocui.ModNone, handleCommitPress); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", gocui.KeySpace, gocui.ModNone, handleFilePress); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", 'd', gocui.ModNone, handleFileRemove); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", 'm', gocui.ModNone, handleSwitchToMerge); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", 'o', gocui.ModNone, handleFileOpen); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", 's', gocui.ModNone, handleSublimeFileOpen); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", 'e', gocui.ModNone, handleFileEdit); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", 'i', gocui.ModNone, handleIgnoreFile); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", 'r', gocui.ModNone, handleRefreshFiles); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", 'S', gocui.ModNone, handleStashSave); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", 'a', gocui.ModNone, handleAbortMerge); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, handleSelectTop); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyEsc, gocui.ModNone, handleEscapeMerge); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, handleSelectBottom); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeySpace, gocui.ModNone, handlePickHunk); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'b', gocui.ModNone, handlePickBothHunks); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowLeft, gocui.ModNone, handleSelectPrevConflict); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowRight, gocui.ModNone, handleSelectNextConflict); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'z', gocui.ModNone, handlePopFileSnapshot); err != nil {
		return err
	}
	if err := g.SetKeybinding("branches", gocui.KeySpace, gocui.ModNone, handleBranchPress); err != nil {
		return err
	}
	if err := g.SetKeybinding("branches", 'c', gocui.ModNone, handleCheckoutByName); err != nil {
		return err
	}
	if err := g.SetKeybinding("branches", 'F', gocui.ModNone, handleForceCheckout); err != nil {
		return err
	}
	if err := g.SetKeybinding("branches", 'n', gocui.ModNone, handleNewBranch); err != nil {
		return err
	}
	if err := g.SetKeybinding("branches", 'm', gocui.ModNone, handleMerge); err != nil {
		return err
	}
	if err := g.SetKeybinding("commits", 's', gocui.ModNone, handleCommitSquashDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("commits", 'r', gocui.ModNone, handleRenameCommit); err != nil {
		return err
	}
	if err := g.SetKeybinding("commits", 'g', gocui.ModNone, handleResetToCommit); err != nil {
		return err
	}
	if err := g.SetKeybinding("stash", gocui.KeySpace, gocui.ModNone, handleStashApply); err != nil {
		return err
	}
	// TODO: come up with a better keybinding (p/P used for pushing/pulling which
	// I'd like to be global. Perhaps all global keybindings should use a modifier
	// like command? But then there's gonna be hotkey conflicts with the terminal
	if err := g.SetKeybinding("stash", 'k', gocui.ModNone, handleStashPop); err != nil {
		return err
	}
	if err := g.SetKeybinding("stash", 'd', gocui.ModNone, handleStashDrop); err != nil {
		return err
	bindings := []Binding{
		Binding{ViewName: "", Key: gocui.KeyArrowLeft, Modifier: gocui.ModNone, Handler: previousView},
		Binding{ViewName: "", Key: gocui.KeyArrowRight, Modifier: gocui.ModNone, Handler: nextView},
		Binding{ViewName: "", Key: gocui.KeyTab, Modifier: gocui.ModNone, Handler: nextView},
		Binding{ViewName: "", Key: 'q', Modifier: gocui.ModNone, Handler: quit},
		Binding{ViewName: "", Key: gocui.KeyCtrlC, Modifier: gocui.ModNone, Handler: quit},
		Binding{ViewName: "", Key: gocui.KeyArrowDown, Modifier: gocui.ModNone, Handler: cursorDown},
		Binding{ViewName: "", Key: gocui.KeyArrowUp, Modifier: gocui.ModNone, Handler: cursorUp},
		Binding{ViewName: "", Key: gocui.KeyPgup, Modifier: gocui.ModNone, Handler: scrollUpMain},
		Binding{ViewName: "", Key: gocui.KeyPgdn, Modifier: gocui.ModNone, Handler: scrollDownMain},
		Binding{ViewName: "", Key: 'P', Modifier: gocui.ModNone, Handler: pushFiles},
		Binding{ViewName: "", Key: 'p', Modifier: gocui.ModNone, Handler: pullFiles},
		Binding{ViewName: "", Key: 'R', Modifier: gocui.ModNone, Handler: handleRefresh},
		Binding{ViewName: "files", Key: 'c', Modifier: gocui.ModNone, Handler: handleCommitPress},
		Binding{ViewName: "files", Key: gocui.KeySpace, Modifier: gocui.ModNone, Handler: handleFilePress},
		Binding{ViewName: "files", Key: 'd', Modifier: gocui.ModNone, Handler: handleFileRemove},
		Binding{ViewName: "files", Key: 'm', Modifier: gocui.ModNone, Handler: handleSwitchToMerge},
		Binding{ViewName: "files", Key: 'o', Modifier: gocui.ModNone, Handler: handleFileOpen},
		Binding{ViewName: "files", Key: 's', Modifier: gocui.ModNone, Handler: handleSublimeFileOpen},
		Binding{ViewName: "files", Key: 'v', Modifier: gocui.ModNone, Handler: handleVsCodeFileOpen},
		Binding{ViewName: "files", Key: 'i', Modifier: gocui.ModNone, Handler: handleIgnoreFile},
		Binding{ViewName: "files", Key: 'r', Modifier: gocui.ModNone, Handler: handleRefreshFiles},
		Binding{ViewName: "files", Key: 'S', Modifier: gocui.ModNone, Handler: handleStashSave},
		Binding{ViewName: "files", Key: 'a', Modifier: gocui.ModNone, Handler: handleAbortMerge},
		Binding{ViewName: "main", Key: gocui.KeyArrowUp, Modifier: gocui.ModNone, Handler: handleSelectTop},
		Binding{ViewName: "main", Key: gocui.KeyArrowDown, Modifier: gocui.ModNone, Handler: handleSelectBottom},
		Binding{ViewName: "main", Key: gocui.KeyEsc, Modifier: gocui.ModNone, Handler: handleEscapeMerge},
		Binding{ViewName: "main", Key: gocui.KeySpace, Modifier: gocui.ModNone, Handler: handlePickHunk},
		Binding{ViewName: "main", Key: 'b', Modifier: gocui.ModNone, Handler: handlePickBothHunks},
		Binding{ViewName: "main", Key: gocui.KeyArrowLeft, Modifier: gocui.ModNone, Handler: handleSelectPrevConflict},
		Binding{ViewName: "main", Key: gocui.KeyArrowRight, Modifier: gocui.ModNone, Handler: handleSelectNextConflict},
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
		Binding{ViewName: "stash", Key: 'k', Modifier: gocui.ModNone, Handler: handleStashPop},
		Binding{ViewName: "stash", Key: 'd', Modifier: gocui.ModNone, Handler: handleStashDrop},
	}
	for _, binding := range bindings {
		if err := g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}
	return nil
}

func layout(g *gocui.Gui) error {
	g.Highlight = true
	g.SelFgColor = gocui.ColorWhite | gocui.AttrBold
	if runtime.GOOS != "windows" {
		g.FgColor = gocui.ColorBlack
	}
	width, height := g.Size()
	leftSideWidth := width / 3
	statusFilesBoundary := 2
	filesBranchesBoundary := 2 * height / 5   // height - 20
	commitsBranchesBoundary := 3 * height / 5 // height - 10
	commitsStashBoundary := height - 5        // height - 5
	minimumHeight := 16

	panelSpacing := 1
	if OverlappingEdges {
		panelSpacing = 0
	}

	if height < minimumHeight {
		v, err := g.SetView("limit", 0, 0, width-1, height-1, 0)
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Not enough space to render panels"
			v.Wrap = true
		}
		return nil
	}

	g.DeleteView("limit")

	optionsTop := height - 2
	// hiding options if there's not enough space
	if height < 30 {
		optionsTop = height - 1
	}

	v, err := g.SetView("main", leftSideWidth+panelSpacing, 0, width-1, optionsTop, gocui.LEFT)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Diff"
		v.Wrap = true
		v.FgColor = gocui.ColorWhite
	}

	if v, err := g.SetView("status", 0, 0, leftSideWidth, statusFilesBoundary, gocui.BOTTOM|gocui.RIGHT); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Status"
		v.FgColor = gocui.ColorWhite
	}

	filesView, err := g.SetView("files", 0, statusFilesBoundary+panelSpacing, leftSideWidth, filesBranchesBoundary, gocui.TOP|gocui.BOTTOM)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		filesView.Highlight = true
		filesView.Title = "Files"
		v.FgColor = gocui.ColorWhite
	}

	if v, err := g.SetView("branches", 0, filesBranchesBoundary+panelSpacing, leftSideWidth, commitsBranchesBoundary, gocui.TOP|gocui.BOTTOM); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Branches"
		v.FgColor = gocui.ColorWhite
	}

	if v, err := g.SetView("commits", 0, commitsBranchesBoundary+panelSpacing, leftSideWidth, commitsStashBoundary, gocui.TOP|gocui.BOTTOM); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Commits"
		v.FgColor = gocui.ColorWhite
	}

	if v, err := g.SetView("stash", 0, commitsStashBoundary+panelSpacing, leftSideWidth, optionsTop, gocui.TOP|gocui.RIGHT); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Stash"
		v.FgColor = gocui.ColorWhite
	}

	if v, err := g.SetView("options", -1, optionsTop, width, optionsTop+2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.BgColor = gocui.ColorDefault
		v.FgColor = gocui.ColorBlue
		v.Frame = false
		v.Title = "Options"

		// these are only called once
		handleFileSelect(g, filesView)
		refreshFiles(g)
		refreshBranches(g)
		refreshCommits(g)
		refreshStashEntries(g)
		nextView(g, nil)
	}

	return nil
}

func fetch(g *gocui.Gui) {
	gitFetch()
	refreshStatus(g)
}

func updateLoader(g *gocui.Gui) {
	if confirmationView, _ := g.View("confirmation"); confirmationView != nil {
		content := trimmedContent(confirmationView)
		if strings.Contains(content, "...") {
			staticContent := strings.Split(content, "...")[0] + "..."
			renderString(g, "confirmation", staticContent+" "+loader())
		}
	}
}

func run() (err error) {
	g, err := gocui.NewGui(gocui.OutputNormal, OverlappingEdges)
	if err != nil {
		return
	}
	defer g.Close()

	// periodically fetching to check for upstream differences
	go func() {
		for range time.Tick(time.Second * 60) {
			fetch(g)
		}
	}()

	go func() {
		for range time.Tick(time.Millisecond * 10) {
			updateLoader(g)
		}
	}()

	g.SetManagerFunc(layout)

	if err = keybindings(g); err != nil {
		return
	}

	err = g.MainLoop()
	return
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
