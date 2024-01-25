package clients

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/jesseduffield/lazygit/pkg/gui"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests"
	"github.com/samber/lo"
)

// This program lets you run integration tests from a TUI. See pkg/integration/README.md for more info.

var SLOW_INPUT_DELAY = 600

func RunTUI(raceDetector bool) {
	rootDir := utils.GetLazyRootDirectory()
	testDir := filepath.Join(rootDir, "test", "integration")

	app := newApp(testDir)
	app.loadTests()

	g, err := gocui.NewGui(gocui.NewGuiOpts{
		OutputMode:       gocui.OutputTrue,
		RuneReplacements: gui.RuneReplacements,
	})
	if err != nil {
		log.Panicln(err)
	}

	g.Cursor = false

	app.g = g

	g.SetManagerFunc(app.layout)

	if err := g.SetKeybinding("list", gocui.KeyArrowUp, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
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

	if err := g.SetKeybinding("list", gocui.KeyArrowDown, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		if app.itemIdx < len(app.filteredTests)-1 {
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

	if err := g.SetKeybinding("list", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", 'q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", 's', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		suspendAndRunTest(currentTest, true, false, raceDetector, 0)

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", gocui.KeyEnter, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		suspendAndRunTest(currentTest, false, false, raceDetector, 0)

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", 't', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		suspendAndRunTest(currentTest, false, false, raceDetector, SLOW_INPUT_DELAY)

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", 'd', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		suspendAndRunTest(currentTest, false, true, raceDetector, 0)

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", 'o', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		cmd := exec.Command("sh", "-c", fmt.Sprintf("code -r pkg/integration/tests/%s.go", currentTest.Name()))
		if err := cmd.Run(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", 'O', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		currentTest := app.getCurrentTest()
		if currentTest == nil {
			return nil
		}

		cmd := exec.Command("sh", "-c", fmt.Sprintf("code test/_results/%s", currentTest.Name()))
		if err := cmd.Run(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("list", '/', gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		app.filtering = true
		if _, err := g.SetCurrentView("editor"); err != nil {
			return err
		}
		editorView, err := g.View("editor")
		if err != nil {
			return err
		}
		editorView.Clear()

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	// not using the editor yet, but will use it to help filter the list
	if err := g.SetKeybinding("editor", gocui.KeyEsc, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		app.filtering = false
		if _, err := g.SetCurrentView("list"); err != nil {
			return err
		}

		app.filteredTests = app.allTests
		app.renderTests()
		app.editorView.TextArea.Clear()
		app.editorView.Clear()
		app.editorView.Reset()

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("editor", gocui.KeyEnter, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		app.filtering = false

		if _, err := g.SetCurrentView("list"); err != nil {
			return err
		}

		app.renderTests()

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

type app struct {
	allTests      []*components.IntegrationTest
	filteredTests []*components.IntegrationTest
	itemIdx       int
	testDir       string
	filtering     bool
	g             *gocui.Gui
	listView      *gocui.View
	editorView    *gocui.View
}

func newApp(testDir string) *app {
	return &app{testDir: testDir, allTests: tests.GetTests(utils.GetLazyRootDirectory())}
}

func (self *app) getCurrentTest() *components.IntegrationTest {
	self.adjustCursor()
	if len(self.filteredTests) > 0 {
		return self.filteredTests[self.itemIdx]
	}
	return nil
}

func (self *app) loadTests() {
	self.filteredTests = self.allTests

	self.adjustCursor()
}

func (self *app) adjustCursor() {
	self.itemIdx = utils.Clamp(self.itemIdx, 0, len(self.filteredTests)-1)
}

func (self *app) filterWithString(needle string) {
	if needle == "" {
		self.filteredTests = self.allTests
	} else {
		self.filteredTests = lo.Filter(self.allTests, func(test *components.IntegrationTest, _ int) bool {
			return strings.Contains(test.Name(), needle)
		})
	}

	self.renderTests()
	self.g.Update(func(g *gocui.Gui) error { return nil })
}

func (self *app) renderTests() {
	self.listView.Clear()
	for _, test := range self.filteredTests {
		fmt.Fprintln(self.listView, test.Name())
	}
}

func (self *app) wrapEditor(f func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool) func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	return func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
		matched := f(v, key, ch, mod)
		if matched {
			self.filterWithString(v.TextArea.GetContent())
		}
		return matched
	}
}

func suspendAndRunTest(test *components.IntegrationTest, sandbox bool, waitForDebugger bool, raceDetector bool, inputDelay int) {
	if err := gocui.Screen.Suspend(); err != nil {
		panic(err)
	}

	runTuiTest(test, sandbox, waitForDebugger, raceDetector, inputDelay)

	fmt.Fprintf(os.Stdout, "\n%s", style.FgGreen.Sprint("press enter to return"))
	fmt.Scanln() // wait for enter press

	if err := gocui.Screen.Resume(); err != nil {
		panic(err)
	}
}

func (self *app) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	descriptionViewHeight := 7
	keybindingsViewHeight := 3
	editorViewHeight := 3
	if !self.filtering {
		editorViewHeight = 0
	} else {
		descriptionViewHeight = 0
		keybindingsViewHeight = 0
	}
	g.Cursor = self.filtering
	g.FgColor = gocui.ColorGreen
	listView, err := g.SetView("list", 0, 0, maxX-1, maxY-descriptionViewHeight-keybindingsViewHeight-editorViewHeight-1, 0)
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		if self.listView == nil {
			self.listView = listView
		}

		listView.Highlight = true
		listView.SelBgColor = gocui.ColorBlue
		self.renderTests()
		listView.Title = "Tests"
		listView.FgColor = gocui.ColorDefault
		if _, err := g.SetCurrentView("list"); err != nil {
			return err
		}
	}

	descriptionView, err := g.SetViewBeneath("description", "list", descriptionViewHeight)
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		descriptionView.Title = "Test description"
		descriptionView.Wrap = true
		descriptionView.FgColor = gocui.ColorDefault
	}

	keybindingsView, err := g.SetViewBeneath("keybindings", "description", keybindingsViewHeight)
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		keybindingsView.Title = "Keybindings"
		keybindingsView.Wrap = true
		keybindingsView.FgColor = gocui.ColorDefault
		fmt.Fprintln(keybindingsView, "up/down: navigate, enter: run test, t: run test slow, s: sandbox, d: debug test, o: open test file, shift+o: open test snapshot directory, forward-slash: filter")
	}

	editorView, err := g.SetViewBeneath("editor", "keybindings", editorViewHeight)
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		if self.editorView == nil {
			self.editorView = editorView
		}

		editorView.Title = "Filter"
		editorView.FgColor = gocui.ColorDefault
		editorView.Editable = true
		editorView.Editor = gocui.EditorFunc(self.wrapEditor(gocui.SimpleEditor))
	}

	currentTest := self.getCurrentTest()
	if currentTest == nil {
		return nil
	}

	descriptionView.Clear()
	fmt.Fprint(descriptionView, currentTest.Description())

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func runTuiTest(test *components.IntegrationTest, sandbox bool, waitForDebugger bool, raceDetector bool, inputDelay int) {
	err := components.RunTests(components.RunTestArgs{
		Tests:           []*components.IntegrationTest{test},
		Logf:            log.Printf,
		RunCmd:          runCmdInTerminal,
		TestWrapper:     runAndPrintError,
		Sandbox:         sandbox,
		WaitForDebugger: waitForDebugger,
		RaceDetector:    raceDetector,
		CodeCoverageDir: "",
		InputDelay:      inputDelay,
		MaxAttempts:     1,
	})
	if err != nil {
		log.Println(err.Error())
	}
}

func runAndPrintError(test *components.IntegrationTest, f func() error) {
	if err := f(); err != nil {
		log.Println(err.Error())
	}
}
