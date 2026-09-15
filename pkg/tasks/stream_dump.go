package tasks

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// streamDumpPathEnvVar names a file to record every byte lazygit reads from a
// command's output into. When it is unset, nothing in this file runs.
//
// The recording exists to settle where a diff renderer's OSC 1717 records sit
// in the stream relative to the rows they describe. A record applies to the
// cells written after it and stops applying at the end of the line (see
// viewBuffer.write in pkg/gocui), so a record that arrives on the far side of a
// row break is attributed to the wrong row, or to no row at all. On Windows the
// stream is ConPTY's re-encoding of the renderer's output rather than the
// renderer's own bytes, so the two can differ there.
//
// ReplayStreamDump reads a recording back through the same reading pipeline,
// which lets a stream captured on one machine be studied on another.
const streamDumpPathEnvVar = "LAZYGIT_DUMP_STREAM"

// A recording is a sequence of blocks, each introduced by a line naming its
// kind:
//
//	NOTE <text>
//	DATA <byte count>
//	<that many raw bytes>
//
// The raw bytes are framed by an explicit count rather than delimited, because
// their line structure is the thing under investigation and must survive the
// recording unaltered. One DATA block holds one Read from the command, so the
// blocks also preserve which bytes arrived together.
const (
	streamDumpNotePrefix = "NOTE "
	streamDumpDataPrefix = "DATA "
)

var (
	streamDumpOnce sync.Once
	streamDumpFile *os.File
	// streamDumpMutex serializes writes, which come from one goroutine per
	// running task.
	streamDumpMutex sync.Mutex
)

// streamDump returns the file to record into, or nil when no recording was
// asked for or the file can't be opened. The result is remembered, so several
// tasks record into one file and their blocks interleave in the order the
// bytes actually arrived.
func streamDump() *os.File {
	streamDumpOnce.Do(func() {
		path := os.Getenv(streamDumpPathEnvVar)
		if path == "" {
			return
		}
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return
		}
		streamDumpFile = f
	})
	return streamDumpFile
}

// DumpStreamNote records a line of context about what is being read, so that
// the bytes that follow can be read in the light of it. Formatted like
// fmt.Printf; no-op unless a recording was asked for.
func DumpStreamNote(format string, args ...any) {
	f := streamDump()
	if f == nil {
		return
	}

	streamDumpMutex.Lock()
	defer streamDumpMutex.Unlock()
	fmt.Fprintf(f, "%s%s\n", streamDumpNotePrefix, fmt.Sprintf(format, args...))
}

// dumpStream returns a reader that records everything read from r, or r itself
// when no recording was asked for.
func dumpStream(r io.Reader, cmdStr string) io.Reader {
	f := streamDump()
	if f == nil {
		return r
	}

	DumpStreamNote("%s reading from: %s", time.Now().Format(time.RFC3339Nano), cmdStr)
	return io.TeeReader(r, streamDumpWriter{f})
}

type streamDumpWriter struct {
	f *os.File
}

// Write records one Read's worth of bytes as a DATA block. It always reports
// success: an io.TeeReader surfaces a failed write to its reader as a read
// error, and a recording that can't be written must not break the render it is
// observing.
func (self streamDumpWriter) Write(p []byte) (int, error) {
	streamDumpMutex.Lock()
	defer streamDumpMutex.Unlock()

	fmt.Fprintf(self.f, "%s%d\n", streamDumpDataPrefix, len(p))
	_, _ = self.f.Write(p)
	_, _ = self.f.WriteString("\n")
	return len(p), nil
}
