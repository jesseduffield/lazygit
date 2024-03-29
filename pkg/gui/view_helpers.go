package gui

import (
	"regexp"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/spkg/bom"
)

func (gui *Gui) resetViewOrigin(v *gocui.View) {
	if err := v.SetCursor(0, 0); err != nil {
		gui.Log.Error(err)
	}

	if err := v.SetOrigin(0, 0); err != nil {
		gui.Log.Error(err)
	}
}

// Returns the number of lines that we should read initially from a cmd task so
// that the scrollbar has the correct size, along with the number of lines after
// which the view is filled and we can do a first refresh.
func (gui *Gui) linesToReadFromCmdTask(v *gocui.View) tasks.LinesToRead {
	_, height := v.Size()
	_, oy := v.Origin()

	linesForFirstRefresh := height + oy + 10

	// We want to read as many lines initially as necessary to let the
	// scrollbar go to its minimum height, so that the scrollbar thumb doesn't
	// change size as you scroll down.
	minScrollbarHeight := 1
	linesToReadForAccurateScrollbar := height*(height-1)/minScrollbarHeight + oy

	// However, cap it at some arbitrary max limit, so that we don't get
	// performance problems for huge monitors or tiny font sizes
	if linesToReadForAccurateScrollbar > 5000 {
		linesToReadForAccurateScrollbar = 5000
	}

	return tasks.LinesToRead{
		Total:               linesToReadForAccurateScrollbar,
		InitialRefreshAfter: linesForFirstRefresh,
	}
}

func (gui *Gui) cleanString(s string) string {
	output := string(bom.Clean([]byte(s)))
	return utils.NormalizeLinefeeds(output)
}

func (gui *Gui) setViewContent(v *gocui.View, s string) {
	v.SetContent(gui.cleanString(s))
}

func (gui *Gui) currentViewName() string {
	currentView := gui.g.CurrentView()
	if currentView == nil {
		return ""
	}
	return currentView.Name()
}

func (gui *Gui) onViewTabClick(windowName string, tabIndex int) error {
	tabs := gui.viewTabMap()[windowName]
	if len(tabs) == 0 {
		return nil
	}

	viewName := tabs[tabIndex].ViewName

	context, ok := gui.helpers.View.ContextForView(viewName)
	if !ok {
		return nil
	}

	return gui.c.PushContext(context)
}

func (gui *Gui) handleNextTab() error {
	view := getTabbedView(gui)
	if view == nil {
		return nil
	}

	for _, context := range gui.State.Contexts.Flatten() {
		if context.GetViewName() == view.Name() {
			return gui.onViewTabClick(
				context.GetWindowName(),
				utils.ModuloWithWrap(view.TabIndex+1, len(view.Tabs)),
			)
		}
	}

	return nil
}

func (gui *Gui) handlePrevTab() error {
	view := getTabbedView(gui)
	if view == nil {
		return nil
	}

	for _, context := range gui.State.Contexts.Flatten() {
		if context.GetViewName() == view.Name() {
			return gui.onViewTabClick(
				context.GetWindowName(),
				utils.ModuloWithWrap(view.TabIndex-1, len(view.Tabs)),
			)
		}
	}

	return nil
}

func getTabbedView(gui *Gui) *gocui.View {
	// It safe assumption that only static contexts have tabs
	context := gui.c.CurrentStaticContext()
	view, _ := gui.g.View(context.GetViewName())
	return view
}

func (gui *Gui) render() {
	gui.c.OnUIThread(func() error { return nil })
}

// postRefreshUpdate is to be called on a context after the state that it depends on has been refreshed
// if the context's view is set to another context we do nothing.
// if the context's view is the current view we trigger a focus; re-selecting the current item.
func (gui *Gui) postRefreshUpdate(c types.Context) error {
	t := time.Now()
	defer func() {
		gui.Log.Infof("postRefreshUpdate for %s took %s", c.GetKey(), time.Since(t))
	}()

	if err := c.HandleRender(); err != nil {
		return err
	}

	if gui.currentViewName() == c.GetViewName() {
		if err := c.HandleFocus(types.OnFocusOpts{}); err != nil {
			return err
		}
	}

	return nil
}

// handleGenericClick is a generic click handler that can be used for any view.
// It handles opening URLs in the browser when the user clicks on one.
func (gui *Gui) handleGenericClick(view *gocui.View) error {
	cx, cy := view.Cursor()
	word, err := view.Word(cx, cy)
	if err != nil {
		return nil
	}

	// Allow URLs to be wrapped in angle brackets, and the closing bracket to
	// be followed by punctuation:
	re := regexp.MustCompile(`^<?(https://.+?)(>[,.;!]*)?$`)
	matches := re.FindStringSubmatch(word)
	if matches == nil {
		return nil
	}

	// Ignore errors (opening the link via the OS can fail if the
	// `os.openLink` config key references a command that doesn't exist, or
	// that errors when called.)
	_ = gui.c.OS().OpenLink(matches[1])

	return nil
}
