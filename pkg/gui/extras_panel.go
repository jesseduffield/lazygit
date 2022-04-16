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
				DisplayString: gui.c.Tr.ToggleShowCommandLog,
				OnPress: func() error {
					currentContext := gui.currentStaticContext()
					if gui.ShowExtrasWindow && currentContext.GetKey() == context.COMMAND_LOG_CONTEXT_KEY {
						if err := gui.c.PopContext(); err != nil {
							return err
						}
					}
					show := !gui.ShowExtrasWindow
					gui.ShowExtrasWindow = show
					gui.c.GetAppState().HideCommandLog = !show
					_ = gui.c.SaveAppState()
					return nil
				},
			},
			{
				DisplayString: gui.c.Tr.FocusCommandLog,
				OnPress:       gui.handleFocusCommandLog,
			},
		},
	})
}

func (gui *Gui) handleFocusCommandLog() error {
	gui.ShowExtrasWindow = true
	// TODO: is this necessary? Can't I just call 'return from context'?
	gui.State.Contexts.CommandLog.SetParentContext(gui.currentSideContext())
	return gui.c.PushContext(gui.State.Contexts.CommandLog)
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
