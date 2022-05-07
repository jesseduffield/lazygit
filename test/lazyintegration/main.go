package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/integration"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
)

// this program lets you manage integration tests in a TUI. See docs/Integration_Tests.md for more info.

type App struct {
	tests   []*integration.Test
	itemIdx int
	testDir string
	editing bool
	g       *gocui.Gui
}

func (app *App) getCurrentTest() *integration.Test {
	if len(app.tests) > 0 {
		return app.tests[app.itemIdx]
	}
	return nil
}

func (app *App) refreshTests() {
	app.loadTests()
	app.g.Update(func(*gocui.Gui) error {
		listView, err := app.g.View("list")
		if err != nil {
			return err
		}

		listView.Clear()
		for _, test := range app.tests {
			fmt.Fprintln(listView, test.Name)
		}

		return nil
	})
}

func (app *App) loadTests() {
	tests, err := integration.LoadTests(app.testDir)
	if err != nil {
		log.Panicln(err)
	}

	app.tests = tests
	if app.itemIdx > len(app.tests)-1 {
		app.itemIdx = len(app.tests) - 1
	}
}

func main() {
	rootDir := integration.GetRootDirectory()
	testDir := filepath.Join(rootDir, "test", "integration")

	app := &App{testDir: testDir}
	app.loadTests()

	g, err := gocui.NewGui(gocui.OutputTrue, false, gocui.NORMAL, false, gui.RuneReplacements)
	if err != nil {
		log.Panicln(err)
	}

	g.Cursor = false

	app.g = g

	g.SetManagerFunc(app.layout)

	if err := g.SetKeybinding("list", nil, gocui.KeyArrowUp, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		if app.itemIdx > 0 {
			app.itemIdx--
		}
		listView, err := g.View("list")
		if err != nil {
			return err
		}
		listView.FocusPoint(0, app.itemIdx)
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", nil, gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", nil, 'q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", nil, 'r', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		cmd := secureexec.Command("sh", "-c", fmt.Sprintf("INCLUDE_SKIPPED=true MODE=record go run test/runner/main.go %s", currentTest.Name))
		app.runSubprocess(cmd)

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", nil, 's', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		cmd := secureexec.Command("sh", "-c", fmt.Sprintf("INCLUDE_SKIPPED=true MODE=sandbox go run test/runner/main.go %s", currentTest.Name))
		app.runSubprocess(cmd)

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", nil, gocui.KeyEnter, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		cmd := secureexec.Command("sh", "-c", fmt.Sprintf("INCLUDE_SKIPPED=true go run test/runner/main.go %s", currentTest.Name))
		app.runSubprocess(cmd)

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", nil, 'u', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		cmd := secureexec.Command("sh", "-c", fmt.Sprintf("INCLUDE_SKIPPED=true MODE=updateSnapshot go run test/runner/main.go %s", currentTest.Name))
		app.runSubprocess(cmd)

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", nil, 't', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		cmd := secureexec.Command("sh", "-c", fmt.Sprintf("INCLUDE_SKIPPED=true SPEED=1 go run test/runner/main.go %s", currentTest.Name))
		app.runSubprocess(cmd)

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", nil, 'o', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		cmd := secureexec.Command("sh", "-c", fmt.Sprintf("code -r %s/%s/test.json", app.testDir, currentTest.Name))
		if err := cmd.Run(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", nil, 'n', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		// need to duplicate that folder and then re-fetch our tests.
		dir := app.testDir + "/" + app.getCurrentTest().Name
		newDir := dir + "_Copy"

		cmd := secureexec.Command("sh", "-c", fmt.Sprintf("cp -r %s %s", dir, newDir))
		if err := cmd.Run(); err != nil {
			return err
		}

		app.loadTests()

		app.refreshTests()
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", nil, 'm', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		app.editing = true
		if _, err := g.SetCurrentView("editor"); err != nil {
			return err
		}
		editorView, err := g.View("editor")
		if err != nil {
			return err
		}
		editorView.Clear()
		fmt.Fprint(editorView, currentTest.Name)

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", nil, 'd', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		dir := app.testDir + "/" + app.getCurrentTest().Name

		cmd := secureexec.Command("sh", "-c", fmt.Sprintf("rm -rf %s", dir))
		if err := cmd.Run(); err != nil {
			return err
		}

		app.refreshTests()

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("editor", nil, gocui.KeyEnter, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		app.editing = false
		if _, err := g.SetCurrentView("list"); err != nil {
			return err
		}

		editorView, err := g.View("editor")
		if err != nil {
			return err
		}

		dir := app.testDir + "/" + app.getCurrentTest().Name
		newDir := app.testDir + "/" + editorView.Buffer()

		cmd := secureexec.Command("sh", "-c", fmt.Sprintf("mv %s %s", dir, newDir))
		if err := cmd.Run(); err != nil {
			return err
		}

		editorView.Clear()

		app.refreshTests()
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("editor", nil, gocui.KeyEsc, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		app.editing = false
		if _, err := g.SetCurrentView("list"); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	err = g.MainLoop()
	g.Close()
	switch err {
	case gocui.ErrQuit:
		return
	default:
		log.Panicln(err)
	}
}

func (app *App) runSubprocess(cmd *exec.Cmd) {
	if err := gocui.Screen.Suspend(); err != nil {
		panic(err)
	}

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Println(err.Error())
	}
	cmd.Stdin = nil
	cmd.Stderr = nil
	cmd.Stdout = nil

	fmt.Fprintf(os.Stdout, "\n%s", style.FgGreen.Sprint("press enter to return"))
	fmt.Scanln() // wait for enter press

	if err := gocui.Screen.Resume(); err != nil {
		panic(err)
	}
}

func (app *App) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	descriptionViewHeight := 7
	keybindingsViewHeight := 3
	editorViewHeight := 3
	if !app.editing {
		editorViewHeight = 0
	} else {
		descriptionViewHeight = 0
		keybindingsViewHeight = 0
	}
	g.Cursor = app.editing
	g.FgColor = gocui.ColorGreen
	listView, err := g.SetView("list", 0, 0, maxX-1, maxY-descriptionViewHeight-keybindingsViewHeight-editorViewHeight-1, 0)
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		listView.Highlight = true
		listView.Clear()
		for _, test := range app.tests {
			fmt.Fprintln(listView, test.Name)
		}
		listView.Title = "Tests"
		listView.FgColor = gocui.ColorDefault
		if _, err := g.SetCurrentView("list"); err != nil {
			return err
		}
	}

	descriptionView, err := g.SetViewBeneath("description", "list", descriptionViewHeight)
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		descriptionView.Title = "Test description"
		descriptionView.Wrap = true
		descriptionView.FgColor = gocui.ColorDefault
	}

	keybindingsView, err := g.SetViewBeneath("keybindings", "description", keybindingsViewHeight)
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		keybindingsView.Title = "Keybindings"
		keybindingsView.Wrap = true
		keybindingsView.FgColor = gocui.ColorDefault
		fmt.Fprintln(keybindingsView, "up/down: navigate, enter: run test, u: run test and update snapshots, r: record test, s: sandbox, o: open test config, n: duplicate test, m: rename test, d: delete test, t: run test at original speed")
	}

	editorView, err := g.SetViewBeneath("editor", "keybindings", editorViewHeight)
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		editorView.Title = "Enter Name"
		editorView.FgColor = gocui.ColorDefault
		editorView.Editable = true
	}

	currentTest := app.getCurrentTest()
	if currentTest == nil {
		return nil
	}

	descriptionView.Clear()
	fmt.Fprintf(descriptionView, "Speed: %f. %s", currentTest.Speed, currentTest.Description)

	if err := g.SetKeybinding("list", nil, gocui.KeyArrowDown, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		if app.itemIdx < len(app.tests)-1 {
			app.itemIdx++
		}

		listView, err := g.View("list")
		if err != nil {
			return err
		}
		listView.FocusPoint(0, app.itemIdx)
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
