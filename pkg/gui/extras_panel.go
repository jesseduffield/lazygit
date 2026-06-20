package gui

import (
	"errors"
	"io"
	"path/filepath"
	"time"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) handleCreateExtrasMenuPanel() error {
	noGitOutputDisabledReason := func() *types.DisabledReason {
		if gui.hasGitOutput() {
			return nil
		}
		return &types.DisabledReason{Text: gui.c.Tr.NoGitOutputToCopy}
	}

	return gui.c.Menu(types.CreateMenuOptions{
		Title: gui.c.Tr.CommandLog,
		Items: []*types.MenuItem{
			{
				Label: gui.c.Tr.ToggleShowCommandLog,
				Keys:  []gocui.Key{gocui.NewKeyRune('t')},
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
				Keys:    []gocui.Key{gocui.NewKeyRune('f')},
				OnPress: gui.handleFocusCommandLog,
			},
			{
				Label:          gui.c.Tr.CopyGitOutputToClipboard,
				Keys:           []gocui.Key{gocui.NewKeyRune('c')},
				OnPress:        gui.handleCopyLastGitOutputToClipboard,
				DisabledReason: noGitOutputDisabledReason(),
			},
			{
				Label:          gui.c.Tr.CopyAllGitOutputToClipboard,
				Keys:           []gocui.Key{gocui.NewKeyRune('a')},
				OnPress:        gui.handleCopyAllGitOutputToClipboard,
				DisabledReason: noGitOutputDisabledReason(),
			},
			{
				Label:   gui.c.Tr.OpenCommandLogInEditor,
				Keys:    []gocui.Key{gocui.NewKeyRune('o')},
				OnPress: gui.handleOpenCommandLogInEditor,
			},
		},
	})
}

func (gui *Gui) handleCopyLastGitOutputToClipboard() error {
	output := gui.lastGitOutput()
	if output == "" {
		return errors.New(gui.c.Tr.NoGitOutputToCopy)
	}

	if err := gui.os.CopyToClipboardQuiet(output); err != nil {
		return err
	}

	gui.c.Toast(gui.c.Tr.GitOutputCopiedToClipboard)
	return nil
}

func (gui *Gui) handleCopyAllGitOutputToClipboard() error {
	output := gui.allGitOutput()
	if output == "" {
		return errors.New(gui.c.Tr.NoGitOutputToCopy)
	}

	if err := gui.os.CopyToClipboardQuiet(output); err != nil {
		return err
	}

	gui.c.Toast(gui.c.Tr.GitOutputCopiedToClipboard)
	return nil
}

func (gui *Gui) handleOpenCommandLogInEditor() error {
	content := gui.commandLogContent()
	if content == "" {
		return errors.New(gui.c.Tr.NoCommandLogToOpenInEditor)
	}

	filepath := filepath.Join(
		gui.os.GetTempDir(),
		gui.c.Git().RepoPaths.RepoName(),
		time.Now().Format("Jan _2 15.04.05.000000000")+"-command-log.txt",
	)
	if err := gui.os.CreateFileWithContent(filepath, content); err != nil {
		return err
	}

	return gui.Helpers().Files.EditFiles([]string{filepath})
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
	return &prefixWriter{
		writer: gui.Views.Extras,
		prefix: style.FgMagenta.Sprintf("\n\n%s\n", gui.c.Tr.GitOutput),
	}
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
