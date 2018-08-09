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
	Platform          platform
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
	Platform:      getPlatform(),
}

type platform struct {
	os           string
	shell        string
	shellArg     string
	escapedQuote string
}

func getPlatform() platform {
	switch runtime.GOOS {
	case "windows":
		return platform{
			os:           "windows",
			shell:        "cmd",
			shellArg:     "/c",
			escapedQuote: "\\\"",
		}
	default:
		return platform{
			os:           runtime.GOOS,
			shell:        "bash",
			shellArg:     "-c",
			escapedQuote: "\"",
		}
	}
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// layout is called for every screen re-render e.g. when the screen is resized
func layout(g *gocui.Gui) error {
	g.Highlight = true
	g.SelFgColor = gocui.ColorWhite | gocui.AttrBold
	width, height := g.Size()
	leftSideWidth := width / 3
	statusFilesBoundary := 2
	filesBranchesBoundary := 2 * height / 5   // height - 20
	commitsBranchesBoundary := 3 * height / 5 // height - 10
	commitsStashBoundary := height - 5        // height - 5
	minimumHeight := 16
	minimumWidth := 10

	panelSpacing := 1
	if OverlappingEdges {
		panelSpacing = 0
	}

	if height < minimumHeight || width < minimumWidth {
		v, err := g.SetView("limit", 0, 0, max(width-1, 2), max(height-1, 2), 0)
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

	// If the confirmation panel is already displayed, just resize the width,
	// otherwise continue
	if v, err := g.View("confirmation"); err == nil {
		_, y := v.Size()
		x0, y0, x1, _ := getConfirmationPanelDimensions(g, "")
		g.SetView("confirmation", x0, y0, x1, y0+y+1, 0)
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

func fetch(g *gocui.Gui) error {
	gitFetch()
	refreshStatus(g)
	return nil
}

func updateLoader(g *gocui.Gui) error {
	if confirmationView, _ := g.View("confirmation"); confirmationView != nil {
		content := trimmedContent(confirmationView)
		if strings.Contains(content, "...") {
			staticContent := strings.Split(content, "...")[0] + "..."
			renderString(g, "confirmation", staticContent+" "+loader())
		}
	}
	return nil
}

func goEvery(g *gocui.Gui, interval time.Duration, function func(*gocui.Gui) error) {
	go func() {
		for range time.Tick(interval) {
			function(g)
		}
	}()
}

func run() (err error) {
	g, err := gocui.NewGui(gocui.OutputNormal, OverlappingEdges)
	if err != nil {
		return
	}
	defer g.Close()

	goEvery(g, time.Second*60, fetch)
	goEvery(g, time.Second*10, refreshFiles)
	goEvery(g, time.Millisecond*10, updateLoader)

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
