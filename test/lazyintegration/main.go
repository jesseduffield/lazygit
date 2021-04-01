package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
)

var errSubProcess = errors.New("subprocess")

type App struct {
	tests      []*IntegrationTest
	itemIdx    int
	subProcess *exec.Cmd
	testDir    string
	editing    bool
	g          *gocui.Gui
}

func (app *App) getCurrentTest() *IntegrationTest {
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
	tests, err := loadTests(app.testDir)
	if err != nil {
		log.Panicln(err)
	}

	app.tests = tests
	if app.itemIdx > len(app.tests)-1 {
		app.itemIdx = len(app.tests) - 1
	}
}

func main() {
	rootDir := getRootDirectory()
	testDir := filepath.Join(rootDir, "test", "integration")

	app := &App{testDir: testDir}
	app.loadTests()

Loop:
	for {
		g, err := gocui.NewGui(gocui.OutputNormal, false, false)
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

			cmd := secureexec.Command("sh", "-c", fmt.Sprintf("RECORD_EVENTS=true go test pkg/gui/gui_test.go -run /%s", currentTest.Name))
			app.subProcess = cmd

			return errSubProcess
		}); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding("list", nil, gocui.KeyEnter, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
			currentTest := app.getCurrentTest()
			if currentTest == nil {
				return nil
			}

			cmd := secureexec.Command("sh", "-c", fmt.Sprintf("go test pkg/gui/gui_test.go -run /%s", currentTest.Name))
			app.subProcess = cmd

			return errSubProcess
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
		if err != nil {
			switch err {
			case gocui.ErrQuit:
				break Loop

			case errSubProcess:
				cmd := app.subProcess
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				if err := cmd.Run(); err != nil {
					log.Println(err.Error())
				}
				cmd.Stdin = nil
				cmd.Stderr = nil
				cmd.Stdout = nil

				fmt.Fprintf(os.Stdout, "\n%s", coloredString("press enter to return", color.FgGreen))
				fmt.Scanln() // wait for enter press

			default:
				log.Panicln(err)
			}
		}
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
	listView.SetCursor(0, app.itemIdx)

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
		fmt.Fprintln(keybindingsView, "up/down: navigate, enter: run test, r: record test, o: open test config, n: duplicate test, m: rename test, d: delete test")
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
	fmt.Fprintf(descriptionView, "Speed: %d. %s", currentTest.Speed, currentTest.Description)

	if err := g.SetKeybinding("list", nil, gocui.KeyArrowDown, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		if app.itemIdx < len(app.tests)-1 {
			app.itemIdx++
		}
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func getRootDirectory() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		_, err := os.Stat(filepath.Join(path, ".git"))

		if err == nil {
			return path
		}

		if !os.IsNotExist(err) {
			panic(err)
		}

		path = filepath.Dir(path)

		if path == "/" {
			panic("must run in lazygit folder or child folder")
		}
	}
}

type IntegrationTest struct {
	Name        string `json:"name"`
	Speed       int    `json:"speed"`
	Description string `json:"description"`
}

func loadTests(testDir string) ([]*IntegrationTest, error) {
	paths, err := filepath.Glob(filepath.Join(testDir, "/*/test.json"))
	if err != nil {
		return nil, err
	}
	tests := make([]*IntegrationTest, len(paths))

	for i, path := range paths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		test := &IntegrationTest{}

		err = json.Unmarshal(data, test)
		if err != nil {
			return nil, err
		}
		test.Name = strings.TrimPrefix(filepath.Dir(path), testDir+"/")

		tests[i] = test
	}

	return tests, nil
}

func coloredString(str string, colorAttributes ...color.Attribute) string {
	colour := color.New(colorAttributes...)
	return coloredStringDirect(str, colour)
}

// coloredStringDirect used for aggregating a few color attributes rather than
// just sending a single one
func coloredStringDirect(str string, colour *color.Color) string {
	return colour.SprintFunc()(fmt.Sprint(str))
}
