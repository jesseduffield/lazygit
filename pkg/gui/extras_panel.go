package gui

import (
	"io"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) handleCreateExtrasMenuPanel() error {
	return gui.c.Menu(types.CreateMenuOptions{
		Title: gui.c.Tr.CommandLog,
		Items: []*types.MenuItem{
			{
				Label: gui.c.Tr.ToggleShowCommandLog,
				OnPress: func() error {
					currentContext := gui.c.Context().CurrentStatic()
					if gui.c.State().GetShowExtrasWindow() && currentContext.GetKey() == context.COMMAND_LOG_CONTEXT_KEY {
						gui.c.Context().Pop()
					}
					show := !gui.c.State().GetShowExtrasWindow()
					gui.c.State().SetShowExtrasWindow(show)
					gui.c.GetAppState().HideCommandLog = !show
					gui.c.SaveAppStateAndLogError()
					return nil
				},
			},
			{
				Label:   gui.c.Tr.FocusCommandLog,
				OnPress: gui.handleFocusCommandLog,
			},
		},
	})
}

func (gui *Gui) handleFocusCommandLog() error {
	gui.c.State().SetShowExtrasWindow(true)
	// TODO: is this necessary? Can't I just call 'return from context'?
	gui.State.Contexts.CommandLog.SetParentContext(gui.c.Context().CurrentSide())
	gui.c.Context().Push(gui.State.Contexts.CommandLog, types.OnFocusOpts{})
	return nil
}

func (gui *Gui) scrollUpExtra() error {
	gui.Views.Extras.Autoscroll = false

	gui.scrollUpView(gui.Views.Extras)

	return nil
}

func (gui *Gui) scrollDownExtra() error {
	gui.Views.Extras.Autoscroll = false

	gui.scrollDownView(gui.Views.Extras)

	return nil
}

func (gui *Gui) pageUpExtrasPanel() error {
	gui.Views.Extras.Autoscroll = false

	gui.Views.Extras.ScrollUp(gui.Contexts().CommandLog.GetViewTrait().PageDelta())

	return nil
}

func (gui *Gui) pageDownExtrasPanel() error {
	gui.Views.Extras.Autoscroll = false

	gui.Views.Extras.ScrollDown(gui.Contexts().CommandLog.GetViewTrait().PageDelta())

	return nil
}

func (gui *Gui) goToExtrasPanelTop() error {
	gui.Views.Extras.Autoscroll = false

	gui.Views.Extras.ScrollUp(gui.Views.Extras.ViewLinesHeight())

	return nil
}

func (gui *Gui) goToExtrasPanelBottom() error {
	gui.Views.Extras.Autoscroll = true

	gui.Views.Extras.ScrollDown(gui.Views.Extras.ViewLinesHeight())

	return nil
}

func (gui *Gui) getCmdWriter() io.Writer {
	return &prefixWriter{writer: gui.Views.Extras, prefix: style.FgMagenta.Sprintf("\n\n%s\n", gui.c.Tr.GitOutput)}
}

// Ensures that the first write is preceded by writing a prefix.
// This allows us to say 'Git output:' before writing the actual git output.
// We could just write directly to the view in this package before running the command but we already have code in the commands package that writes to the same view beforehand (with the command it's about to run) so things would be out of order.
type prefixWriter struct {
	prefix        string
	prefixWritten bool
	writer        io.Writer
}

func (self *prefixWriter) Write(p []byte) (int, error) {
	if !self.prefixWritten {
		self.prefixWritten = true
		// assuming we can write this prefix in one go
		n, err := self.writer.Write([]byte(self.prefix))
		if err != nil {
			return n, err
		}
	}
	return self.writer.Write(p)
}
