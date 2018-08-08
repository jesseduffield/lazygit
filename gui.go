package main

import (

	// "io"
	// "io/ioutil"

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

func layout(g *gocui.Gui) error {
	g.Highlight = true
	g.SelFgColor = gocui.ColorWhite | gocui.AttrBold
	if runtime.GOOS != "windows" && runtime.GOOS != "linux" {
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

	if v, err := g.SetView("options", -1, optionsTop, width-len(version)-2, optionsTop+2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.BgColor = gocui.ColorDefault
		v.FgColor = gocui.ColorBlue
		v.Frame = false
	}

	if v, err := g.SetView("version", width-len(version)-1, optionsTop, width, optionsTop+2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.BgColor = gocui.ColorDefault
		v.FgColor = gocui.ColorGreen
		v.Frame = false
		renderString(g, "version", version)

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
