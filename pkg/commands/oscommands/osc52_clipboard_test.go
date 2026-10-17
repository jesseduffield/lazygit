package oscommands

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBuildOSC52Sequence verifies the escape sequence is framed correctly:
// the `\x1b]52;c;` prefix, the base64 of the payload, and the ST terminator.
func TestBuildOSC52Sequence(t *testing.T) {
	const payload = "commit-hash-abc123"

	sequence, err := buildOSC52Sequence(payload)
	assert.NoError(t, err)

	expected := osc52SequencePrefix + base64.StdEncoding.EncodeToString([]byte(payload)) + osc52SequenceSuffix
	assert.Equal(t, expected, sequence)

	// Guard against accidental terminator/prefix changes that would silently
	// break terminals which only accept the exact framing.
	assert.Equal(t, "\x1b]52;c;", osc52SequencePrefix)
	assert.Equal(t, "\x1b\\", osc52SequenceSuffix)
}

// TestBuildOSC52SequenceRejectsEmpty verifies we never assemble an empty
// payload, which OSC 52 interprets as a request to clear the clipboard.
func TestBuildOSC52SequenceRejectsEmpty(t *testing.T) {
	sequence, err := buildOSC52Sequence("")
	assert.Error(t, err)
	assert.Equal(t, "", sequence)
}

// TestOSC52ClipboardCacheStartsEmpty verifies paste returns the empty string
// before anything has been copied, rather than panicking or erroring.
func TestOSC52ClipboardCacheStartsEmpty(t *testing.T) {
	osCommand := NewDummyOSCommand()

	assert.Equal(t, "", osCommand.pasteFromClipboardOSC52())
}

// TestOSC52ClipboardCacheRoundTrip verifies the central fallback contract:
// because OSC 52 read-back is unreliable, paste must return exactly the value
// most recently copied.
func TestOSC52ClipboardCacheRoundTrip(t *testing.T) {
	osCommand := NewDummyOSCommand()

	osCommand.rememberOSC52Clipboard("first-copied-value")
	assert.Equal(t, "first-copied-value", osCommand.pasteFromClipboardOSC52())

	// A subsequent copy must overwrite, not append or stick to the old value.
	osCommand.rememberOSC52Clipboard("second-copied-value")
	assert.Equal(t, "second-copied-value", osCommand.pasteFromClipboardOSC52())
}
