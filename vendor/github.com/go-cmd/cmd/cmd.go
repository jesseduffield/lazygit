// Package cmd runs external commands with concurrent access to output and
// status. It wraps the Go standard library os/exec.Command to correctly handle
// reading output (STDOUT and STDERR) while a command is running and killing a
// command. All operations are safe to call from multiple goroutines.
//
// A basic example that runs env and prints its output:
//
//   import (
//       "fmt"
//       "github.com/go-cmd/cmd"
//   )
//
//   func main() {
//       // Create Cmd, buffered output
//       envCmd := cmd.NewCmd("env")
//
//       // Run and wait for Cmd to return Status
//       status := <-envCmd.Start()
//
//       // Print each line of STDOUT from Cmd
//       for _, line := range status.Stdout {
//           fmt.Println(line)
//       }
//   }
//
// Commands can be ran synchronously (blocking) or asynchronously (non-blocking):
//
//   envCmd := cmd.NewCmd("env") // create
//
//   status := <-envCmd.Start() // run blocking
//
//   statusChan := envCmd.Start() // run non-blocking
//   // Do other work while Cmd is running...
//   status <- statusChan // blocking
//
// Start returns a channel to which the final Status is sent when the command
// finishes for any reason. The first example blocks receiving on the channel.
// The second example is non-blocking because it saves the channel and receives
// on it later. Only one final status is sent to the channel; use Done for
// multiple goroutines to wait for the command to finish, then call Status to
// get the final status.
package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

// Cmd represents an external command, similar to the Go built-in os/exec.Cmd.
// A Cmd cannot be reused after calling Start. Exported fields are read-only and
// should not be modified, except Env which can be set before calling Start.
// To create a new Cmd, call NewCmd or NewCmdOptions.
type Cmd struct {
	Name   string
	Args   []string
	Env    []string
	Dir    string
	Stdout chan string // streaming STDOUT if enabled, else nil (see Options)
	Stderr chan string // streaming STDERR if enabled, else nil (see Options)
	*sync.Mutex
	started    bool          // cmd.Start called, no error
	stopped    bool          // Stop called
	done       bool          // run() done
	final      bool          // status finalized in Status
	startTime  time.Time     // if started true
	stdout     *OutputBuffer // low-level stdout buffering and streaming
	stderr     *OutputBuffer // low-level stderr buffering and streaming
	status     Status
	statusChan chan Status   // nil until Start() called
	doneChan   chan struct{} // closed when done running
	buffered   bool          // buffer STDOUT and STDERR to Status.Stdout and Std
}

// Status represents the running status and consolidated return of a Cmd. It can
// be obtained any time by calling Cmd.Status. If StartTs > 0, the command has
// started. If StopTs > 0, the command has stopped. After the command finishes
// for any reason, this combination of values indicates success (presuming the
// command only exits zero on success):
//
//   Exit     = 0
//   Error    = nil
//   Complete = true
//
// Error is a Go error from the underlying os/exec.Cmd.Start or os/exec.Cmd.Wait.
// If not nil, the command either failed to start (it never ran) or it started
// but was terminated unexpectedly (probably signaled). In either case, the
// command failed. Callers should check Error first. If nil, then check Exit and
// Complete.
type Status struct {
	Cmd      string
	PID      int
	Complete bool     // false if stopped or signaled
	Exit     int      // exit code of process
	Error    error    // Go error
	StartTs  int64    // Unix ts (nanoseconds), zero if Cmd not started
	StopTs   int64    // Unix ts (nanoseconds), zero if Cmd not started or running
	Runtime  float64  // seconds, zero if Cmd not started
	Stdout   []string // buffered STDOUT; see Cmd.Status for more info
	Stderr   []string // buffered STDERR; see Cmd.Status for more info
}

// NewCmd creates a new Cmd for the given command name and arguments. The command
// is not started until Start is called. Output buffering is on, streaming output
// is off. To control output, use NewCmdOptions instead.
func NewCmd(name string, args ...string) *Cmd {
	return &Cmd{
		Name:     name,
		Args:     args,
		buffered: true,
		Mutex:    &sync.Mutex{},
		status: Status{
			Cmd:      name,
			PID:      0,
			Complete: false,
			Exit:     -1,
			Error:    nil,
			Runtime:  0,
		},
		doneChan: make(chan struct{}),
	}
}

// Options represents customizations for NewCmdOptions.
type Options struct {
	// If Buffered is true, STDOUT and STDERR are written to Status.Stdout and
	// Status.Stderr. The caller can call Cmd.Status to read output at intervals.
	// See Cmd.Status for more info.
	Buffered bool

	// If Streaming is true, Cmd.Stdout and Cmd.Stderr channels are created and
	// STDOUT and STDERR output lines are written them in real time. This is
	// faster and more efficient than polling Cmd.Status. The caller must read both
	// streaming channels, else lines are dropped silently.
	Streaming bool
}

// NewCmdOptions creates a new Cmd with options. The command is not started
// until Start is called.
func NewCmdOptions(options Options, name string, args ...string) *Cmd {
	out := NewCmd(name, args...)
	out.buffered = options.Buffered
	if options.Streaming {
		out.Stdout = make(chan string, DEFAULT_STREAM_CHAN_SIZE)
		out.Stderr = make(chan string, DEFAULT_STREAM_CHAN_SIZE)
	}
	return out
}

// Start starts the command and immediately returns a channel that the caller
// can use to receive the final Status of the command when it ends. The caller
// can start the command and wait like,
//
//   status := <-myCmd.Start() // blocking
//
// or start the command asynchronously and be notified later when it ends,
//
//   statusChan := myCmd.Start() // non-blocking
//   // Do other work while Cmd is running...
//   status := <-statusChan // blocking
//
// Exactly one Status is sent on the channel when the command ends. The channel
// is not closed. Any Go error is set to Status.Error. Start is idempotent; it
// always returns the same channel.
func (c *Cmd) Start() <-chan Status {
	c.Lock()
	defer c.Unlock()

	if c.statusChan != nil {
		return c.statusChan
	}

	c.statusChan = make(chan Status, 1)
	go c.run()
	return c.statusChan
}

// Stop stops the command by sending its process group a SIGTERM signal.
// Stop is idempotent. An error should only be returned in the rare case that
// Stop is called immediately after the command ends but before Start can
// update its internal state.
func (c *Cmd) Stop() error {
	c.Lock()
	defer c.Unlock()

	// Nothing to stop if Start hasn't been called, or the proc hasn't started,
	// or it's already done.
	if c.statusChan == nil || !c.started || c.done {
		return nil
	}

	// Flag that command was stopped, it didn't complete. This results in
	// status.Complete = false
	c.stopped = true

	// Signal the process group (-pid), not just the process, so that the process
	// and all its children are signaled. Else, child procs can keep running and
	// keep the stdout/stderr fd open and cause cmd.Wait to hang.
	return syscall.Kill(-c.status.PID, syscall.SIGTERM)
}

// Status returns the Status of the command at any time. It is safe to call
// concurrently by multiple goroutines.
//
// With buffered output, Status.Stdout and Status.Stderr contain the full output
// as of the Status call time. For example, if the command counts to 3 and three
// calls are made between counts, Status.Stdout contains:
//
//   "1"
//   "1 2"
//   "1 2 3"
//
// The caller is responsible for tailing the buffered output if needed. Else,
// consider using streaming output. When the command finishes, buffered output
// is complete and final.
//
// Status.Runtime is updated while the command is running and final when it
// finishes.
func (c *Cmd) Status() Status {
	c.Lock()
	defer c.Unlock()

	// Return default status if cmd hasn't been started
	if c.statusChan == nil || !c.started {
		return c.status
	}

	if c.done {
		// No longer running
		if !c.final {
			if c.buffered {
				c.status.Stdout = c.stdout.Lines()
				c.status.Stderr = c.stderr.Lines()
				c.stdout = nil // release buffers
				c.stderr = nil
			}
			c.final = true
		}
	} else {
		// Still running
		c.status.Runtime = time.Now().Sub(c.startTime).Seconds()
		if c.buffered {
			c.status.Stdout = c.stdout.Lines()
			c.status.Stderr = c.stderr.Lines()
		}
	}

	return c.status
}

// Done returns a channel that's closed when the command stops running.
// This method is useful for multiple goroutines to wait for the command
// to finish.Call Status after the command finishes to get its final status.
func (c *Cmd) Done() <-chan struct{} {
	return c.doneChan
}

// --------------------------------------------------------------------------

func (c *Cmd) run() {
	defer func() {
		c.statusChan <- c.Status() // unblocks Start if caller is waiting
		close(c.doneChan)
	}()

	// //////////////////////////////////////////////////////////////////////
	// Setup command
	// //////////////////////////////////////////////////////////////////////
	cmd := exec.Command(c.Name, c.Args...)

	// Set process group ID so the cmd and all its children become a new
	// process group. This allows Stop to SIGTERM the cmd's process group
	// without killing this process (i.e. this code here).
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Write stdout and stderr to buffers that are safe to read while writing
	// and don't cause a race condition.
	if c.buffered && c.Stdout != nil {
		// Buffered and streaming, create both and combine with io.MultiWriter
		c.stdout = NewOutputBuffer()
		c.stderr = NewOutputBuffer()
		cmd.Stdout = io.MultiWriter(NewOutputStream(c.Stdout), c.stdout)
		cmd.Stderr = io.MultiWriter(NewOutputStream(c.Stderr), c.stderr)
	} else if c.buffered {
		// Buffered only
		c.stdout = NewOutputBuffer()
		c.stderr = NewOutputBuffer()
		cmd.Stdout = c.stdout
		cmd.Stderr = c.stderr
	} else if c.Stdout != nil {
		// Streaming only
		cmd.Stdout = NewOutputStream(c.Stdout)
		cmd.Stderr = NewOutputStream(c.Stderr)
	} else {
		// No output (effectively >/dev/null 2>&1)
		cmd.Stdout = nil
		cmd.Stderr = nil
	}

	// Set the runtime environment for the command as per os/exec.Cmd.  If Env
	// is nil, use the current process' environment.
	cmd.Env = c.Env
	cmd.Dir = c.Dir

	// //////////////////////////////////////////////////////////////////////
	// Start command
	// //////////////////////////////////////////////////////////////////////
	now := time.Now()
	if err := cmd.Start(); err != nil {
		c.Lock()
		c.status.Error = err
		c.status.StartTs = now.UnixNano()
		c.status.StopTs = time.Now().UnixNano()
		c.done = true
		c.Unlock()
		return
	}

	// Set initial status
	c.Lock()
	c.startTime = now              // command is running
	c.status.PID = cmd.Process.Pid // command is running
	c.status.StartTs = now.UnixNano()
	c.started = true
	c.Unlock()

	// //////////////////////////////////////////////////////////////////////
	// Wait for command to finish or be killed
	// //////////////////////////////////////////////////////////////////////
	err := cmd.Wait()
	now = time.Now()

	// Get exit code of the command. According to the manual, Wait() returns:
	// "If the command fails to run or doesn't complete successfully, the error
	// is of type *ExitError. Other error types may be returned for I/O problems."
	exitCode := 0
	signaled := false
	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			// This is the normal case which is not really an error. It's string
			// representation is only "*exec.ExitError". It only means the cmd
			// did not exit zero and caller should see ExitError.Stderr, which
			// we already have. So first we'll have this as the real/underlying
			// type, then discard err so status.Error doesn't contain a useless
			// "*exec.ExitError". With the real type we can get the non-zero
			// exit code and determine if the process was signaled, which yields
			// a more specific error message, so we set err again in that case.
			exiterr := err.(*exec.ExitError)
			err = nil
			if waitStatus, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = waitStatus.ExitStatus() // -1 if signaled
				if waitStatus.Signaled() {
					signaled = true
					err = errors.New(exiterr.Error()) // "signal: terminated"
				}
			}
		default:
			// I/O problem according to the manual ^. Don't change err.
		}
	}

	// Set final status
	c.Lock()
	if !c.stopped && !signaled {
		c.status.Complete = true
	}
	c.status.Runtime = now.Sub(c.startTime).Seconds()
	c.status.StopTs = now.UnixNano()
	c.status.Exit = exitCode
	c.status.Error = err
	c.done = true
	c.Unlock()
}

// //////////////////////////////////////////////////////////////////////////
// Output
// //////////////////////////////////////////////////////////////////////////

// os/exec.Cmd.StdoutPipe is usually used incorrectly. The docs are clear:
// "it is incorrect to call Wait before all reads from the pipe have completed."
// Therefore, we can't read from the pipe in another goroutine because it
// causes a race condition: we'll read in one goroutine and the original
// goroutine that calls Wait will write on close which is what Wait does.
// The proper solution is using an io.Writer for cmd.Stdout. I couldn't find
// an io.Writer that's also safe for concurrent reads (as lines in a []string
// no less), so I created one:

// OutputBuffer represents command output that is saved, line by line, in an
// unbounded buffer. It is safe for multiple goroutines to read while the command
// is running and after it has finished. If output is small (a few megabytes)
// and not read frequently, an output buffer is a good solution.
//
// A Cmd in this package uses an OutputBuffer for both STDOUT and STDERR by
// default when created by calling NewCmd. To use OutputBuffer directly with
// a Go standard library os/exec.Command:
//
//   import "os/exec"
//   import "github.com/go-cmd/cmd"
//   runnableCmd := exec.Command(...)
//   stdout := cmd.NewOutputBuffer()
//   runnableCmd.Stdout = stdout
//
// While runnableCmd is running, call stdout.Lines() to read all output
// currently written.
type OutputBuffer struct {
	buf   *bytes.Buffer
	lines []string
	*sync.Mutex
}

// NewOutputBuffer creates a new output buffer. The buffer is unbounded and safe
// for multiple goroutines to read while the command is running by calling Lines.
func NewOutputBuffer() *OutputBuffer {
	out := &OutputBuffer{
		buf:   &bytes.Buffer{},
		lines: []string{},
		Mutex: &sync.Mutex{},
	}
	return out
}

// Write makes OutputBuffer implement the io.Writer interface. Do not call
// this function directly.
func (rw *OutputBuffer) Write(p []byte) (n int, err error) {
	rw.Lock()
	n, err = rw.buf.Write(p) // and bytes.Buffer implements io.Writer
	rw.Unlock()
	return // implicit
}

// Lines returns lines of output written by the Cmd. It is safe to call while
// the Cmd is running and after it has finished. Subsequent calls returns more
// lines, if more lines were written. "\r\n" are stripped from the lines.
func (rw *OutputBuffer) Lines() []string {
	rw.Lock()
	// Scanners are io.Readers which effectively destroy the buffer by reading
	// to EOF. So once we scan the buf to lines, the buf is empty again.
	s := bufio.NewScanner(rw.buf)
	for s.Scan() {
		rw.lines = append(rw.lines, s.Text())
	}
	rw.Unlock()
	return rw.lines
}

// --------------------------------------------------------------------------

const (
	// DEFAULT_LINE_BUFFER_SIZE is the default size of the OutputStream line buffer.
	// The default value is usually sufficient, but if ErrLineBufferOverflow errors
	// occur, try increasing the size by calling OutputBuffer.SetLineBufferSize.
	DEFAULT_LINE_BUFFER_SIZE = 16384

	// DEFAULT_STREAM_CHAN_SIZE is the default string channel size for a Cmd when
	// Options.Streaming is true. The string channel size can have a minor
	// performance impact if too small by causing OutputStream.Write to block
	// excessively.
	DEFAULT_STREAM_CHAN_SIZE = 1000
)

// ErrLineBufferOverflow is returned by OutputStream.Write when the internal
// line buffer is filled before a newline character is written to terminate a
// line. Increasing the line buffer size by calling OutputStream.SetLineBufferSize
// can help prevent this error.
type ErrLineBufferOverflow struct {
	Line       string // Unterminated line that caused the error
	BufferSize int    // Internal line buffer size
	BufferFree int    // Free bytes in line buffer
}

func (e ErrLineBufferOverflow) Error() string {
	return fmt.Sprintf("line does not contain newline and is %d bytes too long to buffer (buffer size: %d)",
		len(e.Line)-e.BufferSize, e.BufferSize)
}

// OutputStream represents real time, line by line output from a running Cmd.
// Lines are terminated by a single newline preceded by an optional carriage
// return. Both newline and carriage return are stripped from the line when
// sent to a caller-provided channel.
//
// The caller must begin receiving before starting the Cmd. Write blocks on the
// channel; the caller must always read the channel. The channel is not closed
// by the OutputStream.
//
// A Cmd in this package uses an OutputStream for both STDOUT and STDERR when
// created by calling NewCmdOptions and Options.Streaming is true. To use
// OutputStream directly with a Go standard library os/exec.Command:
//
//   import "os/exec"
//   import "github.com/go-cmd/cmd"
//
//   stdoutChan := make(chan string, 100)
//   go func() {
//       for line := range stdoutChan {
//           // Do something with the line
//       }
//   }()
//
//   runnableCmd := exec.Command(...)
//   stdout := cmd.NewOutputStream(stdoutChan)
//   runnableCmd.Stdout = stdout
//
//
// While runnableCmd is running, lines are sent to the channel as soon as they
// are written and newline-terminated by the command. After the command finishes,
// the caller should wait for the last lines to be sent:
//
//   for len(stdoutChan) > 0 {
//       time.Sleep(10 * time.Millisecond)
//   }
//
// Since the channel is not closed by the OutputStream, the two indications that
// all lines have been sent and received are the command finishing and the
// channel size being zero.
type OutputStream struct {
	streamChan chan string
	bufSize    int
	buf        []byte
	lastChar   int
}

// NewOutputStream creates a new streaming output on the given channel. The
// caller must begin receiving on the channel before the command is started.
// The OutputStream never closes the channel.
func NewOutputStream(streamChan chan string) *OutputStream {
	out := &OutputStream{
		streamChan: streamChan,
		// --
		bufSize:  DEFAULT_LINE_BUFFER_SIZE,
		buf:      make([]byte, DEFAULT_LINE_BUFFER_SIZE),
		lastChar: 0,
	}
	return out
}

// Write makes OutputStream implement the io.Writer interface. Do not call
// this function directly.
func (rw *OutputStream) Write(p []byte) (n int, err error) {
	n = len(p) // end of buffer
	firstChar := 0

LINES:
	for {
		// Find next newline in stream buffer. nextLine starts at 0, but buff
		// can contain multiple lines, like "foo\nbar". So in that case nextLine
		// will be 0 ("foo\nbar\n") then 4 ("bar\n") on next iteration. And i
		// will be 3 and 7, respectively. So lines are [0:3] are [4:7].
		newlineOffset := bytes.IndexByte(p[firstChar:], '\n')
		if newlineOffset < 0 {
			break LINES // no newline in stream, next line incomplete
		}

		// End of line offset is start (nextLine) + newline offset. Like bufio.Scanner,
		// we allow \r\n but strip the \r too by decrementing the offset for that byte.
		lastChar := firstChar + newlineOffset // "line\n"
		if newlineOffset > 0 && p[newlineOffset-1] == '\r' {
			lastChar -= 1 // "line\r\n"
		}

		// Send the line, prepend line buffer if set
		var line string
		if rw.lastChar > 0 {
			line = string(rw.buf[0:rw.lastChar])
			rw.lastChar = 0 // reset buffer
		}
		line += string(p[firstChar:lastChar])
		rw.streamChan <- line // blocks if chan full

		// Next line offset is the first byte (+1) after the newline (i)
		firstChar += newlineOffset + 1
	}

	if firstChar < n {
		remain := len(p[firstChar:])
		bufFree := len(rw.buf[rw.lastChar:])
		if remain > bufFree {
			var line string
			if rw.lastChar > 0 {
				line = string(rw.buf[0:rw.lastChar])
			}
			line += string(p[firstChar:])
			err = ErrLineBufferOverflow{
				Line:       line,
				BufferSize: rw.bufSize,
				BufferFree: bufFree,
			}
			n = firstChar
			return // implicit
		}
		copy(rw.buf[rw.lastChar:], p[firstChar:])
		rw.lastChar += remain
	}

	return // implicit
}

// Lines returns the channel to which lines are sent. This is the same channel
// passed to NewOutputStream.
func (rw *OutputStream) Lines() <-chan string {
	return rw.streamChan
}

// SetLineBufferSize sets the internal line buffer size. The default is DEFAULT_LINE_BUFFER_SIZE.
// This function must be called immediately after NewOutputStream, and it is not
// safe to call by multiple goroutines.
//
// Increasing the line buffer size can help reduce ErrLineBufferOverflow errors.
func (rw *OutputStream) SetLineBufferSize(n int) {
	rw.bufSize = n
	rw.buf = make([]byte, rw.bufSize)
}
