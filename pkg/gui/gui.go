package gui

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"fmt"

	"github.com/golang-collections/collections/stack"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/updates"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

// lazygitTitle is the icon that gets display when the user focusses
// on the status view
const (
	lazygitTitle = `
   _                       _ _
  | |                     (_) |
  | | __ _ _____   _  __ _ _| |_
  | |/ _` + "`" + ` |_  / | | |/ _` + "`" + ` | | __|
  | | (_| |/ /| |_| | (_| | | |_
  |_|\__,_/___|\__, |\__, |_|\__|
                __/ | __/ |
               |___/ |___/       `
)

var (
	dashboardString = fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s\n\n%s",
		lazygitTitle,
		"Keybindings: https://github.com/jesseduffield/lazygit/blob/master/docs/Keybindings.md",
		"Config Options: https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md",
		"Tutorial: https://www.youtube.com/watch?v=VDXvbHZYeKY",
		"Raise an Issue: https://github.com/jesseduffield/lazygit/issues",
	)
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

// Teml is short for template used to make the required map[string]interface{} shorter when using gui.Tr.SLocalize and gui.Tr.TemplateLocalize
type Teml i18n.Teml

// Gui wraps the gocui Gui object which handles rendering and events
type Gui struct {
	g             *gocui.Gui
	Log           *logrus.Entry
	GitCommand    *commands.GitCommand
	OSCommand     *commands.OSCommand
	SubProcess    *exec.Cmd
	State         guiState
	Config        config.AppConfigurer
	Tr            *i18n.Localizer
	Errors        SentinelErrors
	Updater       *updates.Updater
	statusManager *statusManager
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
	Updating          bool
	Keys              []Binding
}

// NewGui builds a new gui handler.
func NewGui(log *logrus.Entry, gitCommand *commands.GitCommand, oSCommand *commands.OSCommand, tr *i18n.Localizer, config config.AppConfigurer, updater *updates.Updater) (*Gui, error) {

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
	}

	gui := &Gui{
		Log:           log,
		GitCommand:    gitCommand,
		OSCommand:     oSCommand,
		State:         initialState,
		Config:        config,
		Tr:            tr,
		Updater:       updater,
		statusManager: &statusManager{},
	}

	gui.GenerateSentinelErrors()

	return gui, nil
}

// Run setup the gui with keybindings and start the mainloop.
// returns an error if something goes wrong.
func (gui *Gui) Run() error {

	var err error

	gui.g, err = gocui.NewGui(gocui.OutputNormal, OverlappingEdges)
	if err != nil {
		gui.Log.Errorf("Failed at newgui: %s\n", err)
		return err
	}

	defer gui.g.Close()

	err = gui.SetColorScheme()
	if err != nil {
		gui.Log.Errorf("Failed at setcolorscheme: %s\n", err)
		return err
	}

	gui.g.SetManagerFunc(gui.layout)

	gui.goEvery(time.Second*60, gui.fetch)
	gui.goEvery(time.Second*10, gui.refreshFiles)
	gui.goEvery(time.Millisecond*50, gui.updateLoader)
	gui.goEvery(time.Millisecond*50, gui.renderAppStatus)

	if err = gui.keybindings(gui.g); err != nil {
		gui.Log.Errorf("Failed to set keybindings at Run: %s\n", err)
		return err
	}

	err = gui.g.MainLoop()
	if err != nil {
		gui.Log.Errorf("Failed to run mainLoop at Run: %s\n", err)
		return err
	}

	return nil
}

// RunWithSubprocesses loops, instantiating a new gocui.Gui with each iteration
// if the error returned from a run is a ErrSubProcess, it runs the subprocess
// otherwise it handles the error, possibly by quitting the application
func (gui *Gui) RunWithSubprocesses() {
	for {

		err := gui.Run()
		if err != nil {
			if err == gocui.ErrQuit {
				break
			} else if err == gui.Errors.ErrSubProcess {
				gui.SubProcess.Stdin = os.Stdin
				gui.SubProcess.Stdout = os.Stdout
				gui.SubProcess.Stderr = os.Stderr

				err = gui.SubProcess.Run()
				if err != nil {
					gui.Log.Errorf("Failed to runWithSubProcess: %s\n", err)
					return
				}

				gui.SubProcess.Stdout = ioutil.Discard
				gui.SubProcess.Stderr = ioutil.Discard
				gui.SubProcess.Stdin = nil
				gui.SubProcess = nil
			} else {
				gui.Log.Errorf("Failed to Run at RunWithSubprocesses: %s\n", err)
				panic(err)
			}
		}
	}
}

// layout is called for every screen re-render e.g. when the screen is resized
func (gui *Gui) layout(g *gocui.Gui) error {
	gui.g.Highlight = true
	width, height := gui.g.Size()
	version := gui.Config.GetVersion()
	leftSideWidth := width / 3
	statusFilesBoundary := 2
	filesBranchesBoundary := 2 * height / 5   // height - 20
	commitsBranchesBoundary := 3 * height / 5 // height - 10
	commitsStashBoundary := height - 5        // height - 5
	optionsVersionBoundary := width - utils.Max(len(version), 1)
	minimumHeight := 16
	minimumWidth := 10

	appStatus := gui.statusManager.getStatusString()
	appStatusOptionsBoundary := 0
	if appStatus != "" {
		appStatusOptionsBoundary = len(appStatus) + 2
	}

	panelSpacing := 1
	if OverlappingEdges {
		panelSpacing = 0
	}

	if height < minimumHeight || width < minimumWidth {
		v, err := gui.g.SetView("limit", 0, 0, utils.Max(width-1, 2), utils.Max(height-1, 2), 0)
		if err != nil {
			if err != gocui.ErrUnknownView {
				gui.Log.Errorf("Failed to create limit view at layout: %s\n", err)
				return err
			}

			v.Title = gui.Tr.SLocalize("NotEnoughSpace")
			v.Wrap = true
		}

		return nil
	}

	_ = gui.g.DeleteView("limit")

	optionsTop := height - 2

	// hiding options if there's not enough space
	if height < 30 {
		optionsTop = height - 1
	}

	v, err := gui.g.SetView("main", leftSideWidth+panelSpacing, 0, width-1, optionsTop, gocui.LEFT)
	if err != nil {

		if err != gocui.ErrUnknownView {
			gui.Log.Errorf("Failed to create files view in main-layout: %s\n", err)
			return err
		}

		v.Title = gui.Tr.SLocalize("DiffTitle")
		v.Wrap = true
		v.FgColor = gocui.ColorWhite

	}

	v, err = gui.g.SetView("status", 0, 0, leftSideWidth, statusFilesBoundary, gocui.BOTTOM|gocui.RIGHT)
	if err != nil {

		if err != gocui.ErrUnknownView {
			gui.Log.Errorf("Failed to create status view in status-layout: %s\n", err)
			return err
		}

		v.Title = gui.Tr.SLocalize("StatusTitle")
		v.FgColor = gocui.ColorWhite
	}

	filesView, err := gui.g.SetView("files", 0, statusFilesBoundary+panelSpacing, leftSideWidth, filesBranchesBoundary, gocui.TOP|gocui.BOTTOM)
	if err != nil {

		if err != gocui.ErrUnknownView {
			gui.Log.Errorf("Failed to create files view in files-layout: %s\n", err)
			return err
		}

		filesView.Highlight = true
		filesView.Title = gui.Tr.SLocalize("FilesTitle")
		filesView.FgColor = gocui.ColorWhite

		err = gui.registerRefresher("files", gui.refreshFiles)
		if err != nil {
			gui.Log.Errorf("Failed to register refresher at files-layout: %s\n", err)
			return err
		}

	}

	v, err = gui.g.SetView("branches", 0, filesBranchesBoundary+panelSpacing, leftSideWidth, commitsBranchesBoundary, gocui.TOP|gocui.BOTTOM)
	if err != nil {

		if err != gocui.ErrUnknownView {
			gui.Log.Errorf("Failed to create branches view in branches-layout: %s\n", err)
			return err
		}

		v.Title = gui.Tr.SLocalize("BranchesTitle")
		v.FgColor = gocui.ColorWhite

		err = gui.registerRefresher("branches", gui.refreshBranches)
		if err != nil {
			gui.Log.Errorf("Failed to create files view in branches-layout: %s\n", err)
			return err
		}

	}

	v, err = gui.g.SetView("commits", 0, commitsBranchesBoundary+panelSpacing, leftSideWidth, commitsStashBoundary, gocui.TOP|gocui.BOTTOM)
	if err != nil {

		if err != gocui.ErrUnknownView {
			gui.Log.Errorf("Failed to create commits view in commits-layout: %s\n", err)
			return err
		}

		v.Title = gui.Tr.SLocalize("CommitsTitle")
		v.FgColor = gocui.ColorWhite

		err = gui.registerRefresher("commits", gui.refreshCommits)
		if err != nil {
			gui.Log.Errorf("Failed to register refresher at commits-layout: %s\n", err)
			return err
		}
	}

	v, err = gui.g.SetView("stash", 0, commitsStashBoundary+panelSpacing, leftSideWidth, optionsTop, gocui.TOP|gocui.RIGHT)
	if err != nil {

		if err != gocui.ErrUnknownView {
			gui.Log.Errorf("Failed to create stash view in stash-layout: %s\n", err)
			return err
		}

		v.Title = gui.Tr.SLocalize("StashTitle")
		v.FgColor = gocui.ColorWhite
	}

	v, err = gui.g.SetView("options", appStatusOptionsBoundary-1, optionsTop, optionsVersionBoundary-1, optionsTop+2, 0)
	if err != nil {

		if err != gocui.ErrUnknownView {
			gui.Log.Errorf("Failed to create options view in options-layout: %s\n", err)
			return err
		}

		v.Frame = false

		v.FgColor, err = gui.GetOptionsPanelTextColor()
		if err != nil {
			gui.Log.Errorf("Failed to get color in options-layout: %s\n", err)
			return err
		}
	}

	v, _ = gui.g.View("commitMessage")
	if v == nil {

		// doesn't matter where this view starts because it will be hidden
		v, err = gui.g.SetView("commitMessage", 0, 0, width/2, height/2, 0)
		if err != nil {

			if err != gocui.ErrUnknownView {
				gui.Log.Errorf("Failed to create commitMessage view in commitmessage-layout: %s\n", err)
				return err
			}

			_, err = gui.g.SetViewOnBottom("commitMessage")
			if err != nil {
				gui.Log.Errorf("Failed to set commitmessage view to bottom in commitmessage-layout: %s\n", err)
				return err
			}

			v.Title = gui.Tr.SLocalize("CommitMessage")
			v.FgColor = gocui.ColorWhite
			v.Editable = true
			v.Editor = gocui.EditorFunc(gui.simpleEditor)
		}
	}

	v, err = gui.g.SetView("appStatus", -1, optionsTop, width, optionsTop+2, 0)
	if err != nil {

		if err != gocui.ErrUnknownView {
			gui.Log.Errorf("Failed to create appstatus view in appstatus-layout: %s\n", err)
			return err
		}

		v.BgColor = gocui.ColorDefault
		v.FgColor = gocui.ColorCyan
		v.Frame = false

		_, err = gui.g.SetViewOnBottom("appStatus")
		if err != nil {
			gui.Log.Errorf("Failed to set appstatus view to bottom in appstatus-layout: %s\n", err)
			return err
		}
	}

	v, err = gui.g.SetView("version", optionsVersionBoundary-1, optionsTop, width, optionsTop+2, 0)
	if err != nil {

		if err != gocui.ErrUnknownView {
			gui.Log.Errorf("Failed to create version view in version-layout: %s\n", err)
			return err
		}

		v.BgColor = gocui.ColorDefault
		v.FgColor = gocui.ColorGreen
		v.Frame = false

		err = gui.renderString(gui.g, "version", version)
		if err != nil {
			gui.Log.Errorf("Failed to render string version in version-layout: %s\n", err)
			return err
		}

		// these are only called once (it's a place to put all the things you want
		// to happen on startup after the screen is first rendered)
		gui.Updater.CheckForNewUpdate(gui.onBackgroundUpdateCheckFinish, false)

		err = gui.handleFileSelect()
		if err != nil {
			gui.Log.Errorf("Failed to handleFileSelect at layout: %s\n", err)
			return err
		}

		err = gui.refreshFiles()
		if err != nil {
			gui.Log.Errorf("Failed to refreshFiles at layout: %s\n", err)
			return err
		}

		err = gui.refreshBranches()
		if err != nil {
			gui.Log.Errorf("Failed to refreshBranches at layout: %s\n", err)
			return err
		}

		err = gui.refreshCommits()
		if err != nil {
			gui.Log.Errorf("Failed to refreshCommits at layout: %s\n", err)
			return err
		}

		err = gui.refreshStashEntries()
		if err != nil {
			gui.Log.Errorf("Failed to refreshStashEntries at layout: %s\n", err)
			return err
		}

		err = gui.switchFocus(gui.g, nil, filesView)
		if err != nil {
			gui.Log.Errorf("Failed to create switchFocus in appstatus-layout: %s\n", err)
			return err
		}

		if gui.Config.GetUserConfig().GetString("reporting") == "undetermined" {
			err = gui.promptAnonymousReporting()
			if err != nil {
				gui.Log.Errorf("Failed to promptAnonReporting in appstatus-layout: %s\n", err)
				return err
			}
		}
	}

	err = gui.resizeCurrentPopupPanel(gui.g)
	if err != nil {
		gui.Log.Errorf("Failed to resizeCurrentPopupPanel at layout: %s\n", err)
		return err
	}

	return nil
}

// prompAnonymouseReporting ask the user to help by sending logs.
// returns an error if something goes wrong.
func (gui *Gui) promptAnonymousReporting() error {
	return gui.createConfirmationPanel(nil, gui.Tr.SLocalize("AnonymousReportingTitle"), gui.Tr.SLocalize("AnonymousReportingPrompt"),
		func(g *gocui.Gui, v *gocui.View) error {
			return gui.Config.WriteToUserConfig("reporting", "on")
		}, func(g *gocui.Gui, v *gocui.View) error {
			return gui.Config.WriteToUserConfig("reporting", "off")
		})
}

// Fetch fetches the commits.
// returns an error if something goes wrong.
func (gui *Gui) fetch() error {

	err := gui.GitCommand.Fetch()
	if err != nil {
		gui.Log.Errorf("Failed to fetch at fetch: %s\n", err)
		return err
	}

	err = gui.refreshStatus()
	if err != nil {
		gui.Log.Errorf("Failed to refreshStatus at fetch: %s\n", err)
		return err
	}

	return nil
}

// updateloader shows a little loader
func (gui *Gui) updateLoader() error {

	view, _ := gui.g.View("confirmation")
	if view != nil {

		content := gui.trimmedContent(view)
		if strings.Contains(content, "...") {
			staticContent := strings.Split(content, "...")[0] + "..."

			err := gui.renderString(gui.g, "confirmation", staticContent+" "+utils.Loader())
			if err != nil {
				gui.Log.Errorf("Failed to render string at updateLoader: %s\n", err)
				return err
			}
		}
	}
	return nil
}

// renderAppStatus renders the app status
// returns an error if something goes wrong.
func (gui *Gui) renderAppStatus() error {
	appStatus := gui.statusManager.getStatusString()
	if appStatus != "" {
		return gui.renderString(gui.g, "appStatus", appStatus)
	}

	return nil
}

// renderGlobalOptions renders the global options.
// returns an error if something goes wrong.
func (gui *Gui) renderGlobalOptions() error {
	return gui.renderOptionsMap(gui.g, map[string]string{
		"PgUp/PgDn": gui.Tr.SLocalize("scroll"),
		"← → ↑ ↓":   gui.Tr.SLocalize("navigate"),
		"esc/q":     gui.Tr.SLocalize("close"),
		"x":         gui.Tr.SLocalize("menu"),
	})
}

// goEvery is a little goroutine that executes the function every
// interval duration.
func (gui *Gui) goEvery(interval time.Duration, function func() error) {
	go func() {
		for range time.Tick(interval) {
			err := function()
			if err != nil {
				gui.Log.Errorf("Failed to exectute function in goevery: %s\n", err)
			}
		}
	}()
}

// handleRefresh is a macro for refresh
func (gui *Gui) handleRefresh(g *gocui.Gui, v *gocui.View) error {
	gui.refresh()
	return nil
}

// quit handles the quit keys
func (gui *Gui) quit(g *gocui.Gui, v *gocui.View) error {

	if gui.State.Updating {

		err := gui.createUpdateQuitConfirmation(v)
		if err != nil {
			gui.Log.Errorf("Failed to create update quit confirmation at quit: %s\n", err)
		}

		return nil
	}

	if gui.Config.GetUserConfig().GetBool("confirmOnQuit") {
		err := gui.createConfirmationPanel(v, "", gui.Tr.SLocalize("ConfirmQuit"),
			func(g *gocui.Gui, v *gocui.View) error {
				return gocui.ErrQuit
			}, nil)
		if err != nil {
			gui.Log.Errorf("Failed to create confirmation panel at quit: %s\n", err)
			return err
		}

		return nil
	}

	return gocui.ErrQuit
}
