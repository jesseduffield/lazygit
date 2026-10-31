package oscommands

import (
	"encoding/base64"
	"os"

	"github.com/go-errors/errors"
)

// OSC 52 clipboard fallback.
//
// OSC 52 is a terminal escape sequence that lets a program set the system
// clipboard purely by writing to the terminal, with no external helper binary
// and no tmux. The terminal emulator decodes the base64 payload and updates the
// clipboard on the machine the terminal is running on. This is the last-resort
// clipboard path for a headless SSH session that is *not* inside tmux but whose
// terminal understands OSC 52. Neovim's and Helix's clipboard providers use
// the same mechanism as their final fallback.
//
// The catch is read-back: most terminals refuse OSC 52 clipboard *reads* for
// security reasons, so there is no dependable way to paste from the terminal
// clipboard. We therefore remember whatever we last copied and hand that back
// on paste, matching how Neovim's and Helix's clipboard providers behave.
//
// CRITICAL PATH: terminal escape emission -- human review recommended. These
// bytes are written verbatim to the controlling terminal; a malformed sequence
// corrupts the display for the rest of the session.

const (
	// osc52SequencePrefix opens an OSC 52 "set clipboard" request targeting the
	// `c` (clipboard) selection. The base64-encoded payload follows immediately.
	osc52SequencePrefix = "\x1b]52;c;"

	// osc52SequenceSuffix is the String Terminator (ESC backslash). ST is the
	// standard OSC terminator; terminals that accept OSC 52 accept it, whereas
	// BEL is a legacy xterm alternative.
	osc52SequenceSuffix = "\x1b\\"

	// controllingTerminalPath is the process's controlling terminal. We target
	// it rather than os.Stdout so the sequence reaches the real terminal even
	// while the gocui TUI owns stdout, and so it survives any stdout
	// redirection.
	controllingTerminalPath = "/dev/tty"

	// terminalOpenMode is the permission argument to os.OpenFile. It is ignored
	// because we open without O_CREATE, but the call still requires a value.
	terminalOpenMode os.FileMode = 0
)

// copyToClipboardOSC52 emits an OSC 52 escape sequence to the controlling
// terminal to set the system clipboard, then remembers the value so a later
// paste can return it. It returns an error if the terminal cannot be opened or
// the sequence cannot be fully written; on error the remembered value is left
// unchanged so paste never reports a copy that did not happen.
//
// Note: OSC 52 provides no acknowledgement, so a terminal that does not
// implement it silently drops the payload. A successful write here means the
// sequence was emitted, not that the clipboard was definitely updated. Very
// large payloads may also be truncated or rejected by the terminal; that limit
// is inherent to OSC 52 and is left to the terminal.
func (c *OSCommand) copyToClipboardOSC52(str string) error {
	// An empty payload is a defined OSC 52 request to *clear* the clipboard,
	// which a "copy" action never intends. Treat it as a no-op so we never wipe
	// the user's clipboard by accident, but still remember the empty value.
	if str == "" {
		c.rememberOSC52Clipboard(str)
		return nil
	}

	sequence, err := buildOSC52Sequence(str)
	if err != nil {
		return err
	}

	terminal, err := os.OpenFile(controllingTerminalPath, os.O_WRONLY, terminalOpenMode)
	if err != nil {
		return errors.Errorf("opening controlling terminal %s for OSC52 clipboard copy: %v", controllingTerminalPath, err)
	}
	defer terminal.Close()

	written, err := terminal.WriteString(sequence)
	if err != nil {
		return errors.Errorf("writing OSC52 clipboard sequence: %v", err)
	}

	// Postcondition: a short write means the terminal received a truncated,
	// malformed escape sequence. Surface it rather than pretend the copy worked.
	if written != len(sequence) {
		return errors.Errorf("short write emitting OSC52 clipboard sequence: wrote %d of %d bytes", written, len(sequence))
	}

	c.rememberOSC52Clipboard(str)
	return nil
}

// buildOSC52Sequence frames str into a complete OSC 52 set-clipboard escape
// sequence: prefix, base64-encoded payload, terminator. It is pure and
// side-effect free so the framing can be verified in isolation from the
// terminal write. str must be non-empty; an empty payload would assemble a
// clipboard-clearing request, which buildOSC52Sequence refuses.
func buildOSC52Sequence(str string) (string, error) {
	if str == "" {
		return "", errors.Errorf("refusing to build empty OSC52 clipboard payload")
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	sequence := osc52SequencePrefix + encoded + osc52SequenceSuffix

	// Postcondition: the framed sequence must be strictly longer than its
	// prefix and suffix combined, otherwise the encoding step produced nothing.
	if len(sequence) <= len(osc52SequencePrefix)+len(osc52SequenceSuffix) {
		return "", errors.Errorf("assembled OSC52 sequence is missing its payload")
	}

	return sequence, nil
}

// pasteFromClipboardOSC52 returns the value most recently copied through the
// OSC 52 fallback. It never queries the terminal because OSC 52 read-back is
// not dependably supported.
func (c *OSCommand) pasteFromClipboardOSC52() string {
	return c.recallOSC52Clipboard()
}

// rememberOSC52Clipboard stores the last copied value under the mutex so copy
// and paste can be called from different goroutines safely.
func (c *OSCommand) rememberOSC52Clipboard(str string) {
	c.osc52ClipboardMutex.Lock()
	defer c.osc52ClipboardMutex.Unlock()
	c.osc52Clipboard = str
}

// recallOSC52Clipboard reads the last copied value under the mutex.
func (c *OSCommand) recallOSC52Clipboard() string {
	c.osc52ClipboardMutex.Lock()
	defer c.osc52ClipboardMutex.Unlock()
	return c.osc52Clipboard
}
