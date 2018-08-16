package gui

import (

	// "io"
	// "io/ioutil"

	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	// "strings"

	"github.com/Sirupsen/logrus"
	"github.com/golang-collections/collections/stack"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/i18n"
)

// OverlappingEdges determines if panel edges overlap
var OverlappingEdges = false

// SentinelErrors are the errors that have special meaning and need to be checked
// by calling functions. The less of these, the better
type SentinelErrors struct {
	ErrSubProcess error
	ErrNoFiles    error
}

// GenerateSentinelErrors makes the sentinel errors for the gui. We're defining it here
// because we can't do package-scoped errors with localization, and also because
// it seems like package-scoped variables are bad in general
// https://dave.cheney.net/2017/06/11/go-without-package-scoped-variables
// In the future it would be good to implement some of the recommendations of
// that article. For now, if we don't need an error to be a sentinel, we will just
// define it inline. This has implications for error messages that pop up everywhere
// in that we'll be duplicating the default values. We may need to look at
// having a default localisation bundle defined, and just using keys-only when
// localising things in the code.
func (gui *Gui) GenerateSentinelErrors() {
	gui.Errors = SentinelErrors{
		ErrSubProcess: errors.New(gui.Tr.SLocalize("RunningSubprocess")),
		ErrNoFiles:    errors.New(gui.Tr.SLocalize("NoChangedFiles")),
	}
}

// Gui wraps the gocui Gui object which handles rendering and events
type Gui struct {
	g          *gocui.Gui
	Log        *logrus.Logger
	GitCommand *commands.GitCommand
	OSCommand  *commands.OSCommand
	Version    string
	SubProcess *exec.Cmd
	State      guiState
	Tr         *i18n.Localizer
	Errors     SentinelErrors
}

type guiState struct {
	Files             []commands.File
	Branches          []commands.Branch
	Commits           []commands.Commit
	StashEntries      []commands.StashEntry
	PreviousView      string
	HasMergeConflicts bool
	ConflictIndex     int
	ConflictTop       bool
	Conflicts         []commands.Conflict
	EditHistory       *stack.Stack
	Platform          commands.Platform
	Version           string
}

// NewGui builds a new gui handler
func NewGui(log *logrus.Logger, gitCommand *commands.GitCommand, oSCommand *commands.OSCommand, tr *i18n.Localizer, version string) (*Gui, error) {
	initialState := guiState{
		Files:         make([]commands.File, 0),
		PreviousView:  "files",
		Commits:       make([]commands.Commit, 0),
		StashEntries:  make([]commands.StashEntry, 0),
		ConflictIndex: 0,
		ConflictTop:   true,
		Conflicts:     make([]commands.Conflict, 0),
		EditHistory:   stack.New(),
		Platform:      *oSCommand.Platform,
		Version:       version,
	}

	gui := &Gui{
		Log:        log,
		GitCommand: gitCommand,
		OSCommand:  oSCommand,
		Version:    version,
		State:      initialState,
		Tr:         tr,
	}

	gui.GenerateSentinelErrors()

	return gui, nil
}

func (gui *Gui) scrollUpMain(g *gocui.Gui, v *gocui.View) error {
	mainView, _ := g.View("main")
	ox, oy := mainView.Origin()
	if oy >= 1 {
		return mainView.SetOrigin(ox, oy-1)
	}
	return nil
}

func (gui *Gui) scrollDownMain(g *gocui.Gui, v *gocui.View) error {
	mainView, _ := g.View("main")
	ox, oy := mainView.Origin()
	if oy < len(mainView.BufferLines()) {
		return mainView.SetOrigin(ox, oy+1)
	}
	return nil
}

func (gui *Gui) handleRefresh(g *gocui.Gui, v *gocui.View) error {
	return gui.refreshSidePanels(g)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// layout is called for every screen re-render e.g. when the screen is resized
func (gui *Gui) layout(g *gocui.Gui) error {
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
			v.Title = gui.Tr.SLocalize("NotEnoughSpace")
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
		v.Title = gui.Tr.SLocalize("DiffTitle")
		v.Wrap = true
		v.FgColor = gocui.ColorWhite
	}

	if v, err := g.SetView("status", 0, 0, leftSideWidth, statusFilesBoundary, gocui.BOTTOM|gocui.RIGHT); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = gui.Tr.SLocalize("StatusTitle")
		v.FgColor = gocui.ColorWhite
	}

	filesView, err := g.SetView("files", 0, statusFilesBoundary+panelSpacing, leftSideWidth, filesBranchesBoundary, gocui.TOP|gocui.BOTTOM)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		filesView.Highlight = true
		filesView.Title = gui.Tr.SLocalize("FilesTitle")
		v.FgColor = gocui.ColorWhite
	}

	if v, err := g.SetView("branches", 0, filesBranchesBoundary+panelSpacing, leftSideWidth, commitsBranchesBoundary, gocui.TOP|gocui.BOTTOM); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = gui.Tr.SLocalize("BranchesTitle")
		v.FgColor = gocui.ColorWhite
	}

	if v, err := g.SetView("commits", 0, commitsBranchesBoundary+panelSpacing, leftSideWidth, commitsStashBoundary, gocui.TOP|gocui.BOTTOM); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = gui.Tr.SLocalize("CommitsTitle")
		v.FgColor = gocui.ColorWhite
	}

	if v, err := g.SetView("stash", 0, commitsStashBoundary+panelSpacing, leftSideWidth, optionsTop, gocui.TOP|gocui.RIGHT); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = gui.Tr.SLocalize("StashTitle")
		v.FgColor = gocui.ColorWhite
	}

	if v, err := g.SetView("options", -1, optionsTop, width-len(gui.Version)-2, optionsTop+2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.FgColor = gocui.ColorBlue
		v.Frame = false
	}

	if gui.getCommitMessageView(g) == nil {
		// doesn't matter where this view starts because it will be hidden
		if commitMessageView, err := g.SetView("commitMessage", 0, 0, width, height, 0); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			g.SetViewOnBottom("commitMessage")
			commitMessageView.Title = gui.Tr.SLocalize("CommitMessage")
			commitMessageView.FgColor = gocui.ColorWhite
			commitMessageView.Editable = true
		}
	}

	if v, err := g.SetView("version", width-len(gui.Version)-1, optionsTop, width, optionsTop+2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.BgColor = gocui.ColorDefault
		v.FgColor = gocui.ColorGreen
		v.Frame = false
		gui.renderString(g, "version", gui.Version)

		// these are only called once
		gui.handleFileSelect(g, filesView)
		gui.refreshFiles(g)
		gui.refreshBranches(g)
		gui.refreshCommits(g)
		gui.refreshStashEntries(g)
		gui.nextView(g, nil)
	}

	gui.resizePopupPanels(g)

	return nil
}

func (gui *Gui) fetch(g *gocui.Gui) error {
	gui.GitCommand.Fetch()
	gui.refreshStatus(g)
	return nil
}

func (gui *Gui) updateLoader(g *gocui.Gui) error {
	if confirmationView, _ := g.View("confirmation"); confirmationView != nil {
		content := gui.trimmedContent(confirmationView)
		if strings.Contains(content, "...") {
			staticContent := strings.Split(content, "...")[0] + "..."
			gui.renderString(g, "confirmation", staticContent+" "+gui.loader())
		}
	}
	return nil
}

func (gui *Gui) goEvery(g *gocui.Gui, interval time.Duration, function func(*gocui.Gui) error) {
	go func() {
		for range time.Tick(interval) {
			function(g)
		}
	}()
}

func (gui *Gui) resizePopupPanels(g *gocui.Gui) error {
	v := g.CurrentView()
	if v.Name() == "commitMessage" || v.Name() == "confirmation" {
		return gui.resizePopupPanel(g, v)
	}
	return nil
}

// Run setup the gui with keybindings and start the mainloop
func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal, OverlappingEdges)
	if err != nil {
		return err
	}
	defer g.Close()

	gui.g = g // TODO: always use gui.g rather than passing g around everywhere

	g.FgColor = gocui.ColorDefault

	gui.goEvery(g, time.Second*60, gui.fetch)
	gui.goEvery(g, time.Second*10, gui.refreshFiles)
	gui.goEvery(g, time.Millisecond*10, gui.updateLoader)

	g.SetManagerFunc(gui.layout)

	if err = gui.keybindings(g); err != nil {
		return err
	}

	err = g.MainLoop()
	return err
}

// RunWithSubprocesses loops, instantiating a new gocui.Gui with each iteration
// if the error returned from a run is a ErrSubProcess, it runs the subprocess
// otherwise it handles the error, possibly by quitting the application
func (gui *Gui) RunWithSubprocesses() {
	for {
		if err := gui.Run(); err != nil {
			if err == gocui.ErrQuit {
				break
			} else if err == gui.Errors.ErrSubProcess {
				gui.SubProcess.Stdin = os.Stdin
				gui.SubProcess.Stdout = os.Stdout
				gui.SubProcess.Stderr = os.Stderr
				gui.SubProcess.Run()
				gui.SubProcess.Stdout = ioutil.Discard
				gui.SubProcess.Stderr = ioutil.Discard
				gui.SubProcess.Stdin = nil
				gui.SubProcess = nil
			} else {
				log.Panicln(err)
			}
		}
	}
}

func (gui *Gui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
