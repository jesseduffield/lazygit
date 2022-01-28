package gui

import (
	"io"

	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

func (gui *Gui) handleCreateExtrasMenuPanel() error {
	return gui.PopupHandler.Menu(popup.CreateMenuOptions{
		Title: gui.Tr.CommandLog,
		Items: []*popup.MenuItem{
			{
				DisplayString: gui.Tr.ToggleShowCommandLog,
				OnPress: func() error {
					currentContext := gui.currentStaticContext()
					if gui.ShowExtrasWindow && currentContext.GetKey() == COMMAND_LOG_CONTEXT_KEY {
						if err := gui.returnFromContext(); err != nil {
							return err
						}
					}
					show := !gui.ShowExtrasWindow
					gui.ShowExtrasWindow = show
					gui.Config.GetAppState().HideCommandLog = !show
					_ = gui.Config.SaveAppState()
					return nil
				},
			},
			{
				DisplayString: gui.Tr.FocusCommandLog,
				OnPress:       gui.handleFocusCommandLog,
			},
		},
	})
}

func (gui *Gui) handleFocusCommandLog() error {
	gui.ShowExtrasWindow = true
	gui.State.Contexts.CommandLog.SetParentContext(gui.currentSideContext())
	return gui.pushContext(gui.State.Contexts.CommandLog)
}

func (gui *Gui) scrollUpExtra() error {
	gui.Views.Extras.Autoscroll = false

	return gui.scrollUpView(gui.Views.Extras)
}

func (gui *Gui) scrollDownExtra() error {
	gui.Views.Extras.Autoscroll = false

	if err := gui.scrollDownView(gui.Views.Extras); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) getCmdWriter() io.Writer {
	return &prefixWriter{writer: gui.Views.Extras, prefix: style.FgMagenta.Sprintf("\n\n%s\n", gui.Tr.GitOutput)}
}

// Ensures that the first write is preceded by writing a prefix.
// This allows us to say 'Git output:' before writing the actual git output.
// We could just write directly to the view in this package before running the command but we already have code in the commands package that writes to the same view beforehand (with the command it's about to run) so things would be out of order.
type prefixWriter struct {
	prefix        string
	prefixWritten bool
	writer        io.Writer
}

func (self *prefixWriter) Write(p []byte) (n int, err error) {
	if !self.prefixWritten {
		self.prefixWritten = true
		// assuming we can write this prefix in one go
		_, err = self.writer.Write([]byte(self.prefix))
		if err != nil {
			return
		}
	}
	return self.writer.Write(p)
}
