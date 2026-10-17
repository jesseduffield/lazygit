package oscommands

import (
	"os/exec"
)

// tmux clipboard provider.
//
// On a headless SSH host there is no X11 or Wayland display, so none of xclip,
// xsel, or wl-copy can reach a clipboard. In that situation atotto/clipboard
// reports itself unsupported and lazygit's default copy path fails with
// "No clipboard utilities available". When we are running inside a tmux
// session, though, tmux itself can act as the clipboard: it keeps its own
// paste buffers and, with `set-clipboard on`, relays them to the outer
// terminal's system clipboard via tmux's own OSC 52 emission.
//
// Unlike a raw OSC 52 escape we could emit ourselves, the tmux buffer is
// readable, so paste works too -- `save-buffer` hands the buffer contents back
// on stdout. That gives real bidirectional clipboard support within the
// session, which is why we prefer it over a write-only OSC 52 fallback.

const (
	// tmuxEnvVar is set by tmux for every process running inside a session. Its
	// presence is how we detect that we are attached to a tmux server.
	tmuxEnvVar = "TMUX"

	// tmuxBinaryName is the tmux executable we shell out to. It must be on PATH.
	tmuxBinaryName = "tmux"

	// tmuxLoadBufferSubcommand fills a tmux paste buffer from a file argument.
	tmuxLoadBufferSubcommand = "load-buffer"

	// tmuxSaveBufferSubcommand writes a tmux paste buffer to a file argument.
	tmuxSaveBufferSubcommand = "save-buffer"

	// tmuxSetClipboardFlag asks tmux to also copy the buffer to the outer
	// terminal's system clipboard (via tmux's OSC 52 emission) in addition to
	// storing it in the tmux buffer. Requires `set-clipboard on` in tmux.
	tmuxSetClipboardFlag = "-w"

	// tmuxStdioArg is tmux's conventional "-" file argument meaning "read from
	// stdin" for load-buffer and "write to stdout" for save-buffer.
	tmuxStdioArg = "-"
)

// inTmux reports whether lazygit is running inside a tmux session that we can
// drive as a clipboard. Both conditions must hold: the TMUX environment
// variable is set (we are attached to a session), and the tmux binary is
// resolvable on PATH (we can actually invoke it).
func (c *OSCommand) inTmux() bool {
	if c.getenvFn(tmuxEnvVar) == "" {
		return false
	}

	_, err := exec.LookPath(tmuxBinaryName)
	return err == nil
}

// copyToClipboardTmux stores str in a tmux paste buffer and, via the
// set-clipboard flag, forwards it to the outer terminal's system clipboard.
// The buffer contents are streamed in on stdin using the "-" file argument.
func (c *OSCommand) copyToClipboardTmux(str string) error {
	return c.Cmd.
		New([]string{tmuxBinaryName, tmuxLoadBufferSubcommand, tmuxSetClipboardFlag, tmuxStdioArg}).
		SetStdin(str).
		Run()
}

// pasteFromClipboardTmux returns the contents of the most recent tmux paste
// buffer by streaming it out on stdout using the "-" file argument.
func (c *OSCommand) pasteFromClipboardTmux() (string, error) {
	return c.Cmd.
		New([]string{tmuxBinaryName, tmuxSaveBufferSubcommand, tmuxStdioArg}).
		RunWithOutput()
}
