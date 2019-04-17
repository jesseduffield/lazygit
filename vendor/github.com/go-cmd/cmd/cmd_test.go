package cmd_test

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/go-cmd/cmd"
	"github.com/go-test/deep"
)

func TestCmdOK(t *testing.T) {
	now := time.Now().Unix()

	p := cmd.NewCmd("echo", "foo")
	gotStatus := <-p.Start()
	expectStatus := cmd.Status{
		Cmd:      "echo",
		PID:      gotStatus.PID, // nondeterministic
		Complete: true,
		Exit:     0,
		Error:    nil,
		Runtime:  gotStatus.Runtime, // nondeterministic
		Stdout:   []string{"foo"},
		Stderr:   []string{},
	}
	if gotStatus.StartTs < now {
		t.Error("StartTs < now")
	}
	if gotStatus.StopTs < gotStatus.StartTs {
		t.Error("StopTs < StartTs")
	}
	gotStatus.StartTs = 0
	gotStatus.StopTs = 0
	if diffs := deep.Equal(gotStatus, expectStatus); diffs != nil {
		t.Error(diffs)
	}
	if gotStatus.PID < 0 {
		t.Errorf("got PID %d, expected non-zero", gotStatus.PID)
	}
	if gotStatus.Runtime < 0 {
		t.Errorf("got runtime %f, expected non-zero", gotStatus.Runtime)
	}
}

func TestCmdNonzeroExit(t *testing.T) {
	p := cmd.NewCmd("false")
	gotStatus := <-p.Start()
	expectStatus := cmd.Status{
		Cmd:      "false",
		PID:      gotStatus.PID, // nondeterministic
		Complete: true,
		Exit:     1,
		Error:    nil,
		Runtime:  gotStatus.Runtime, // nondeterministic
		Stdout:   []string{},
		Stderr:   []string{},
	}
	gotStatus.StartTs = 0
	gotStatus.StopTs = 0
	if diffs := deep.Equal(gotStatus, expectStatus); diffs != nil {
		t.Error(diffs)
	}
	if gotStatus.PID < 0 {
		t.Errorf("got PID %d, expected non-zero", gotStatus.PID)
	}
	if gotStatus.Runtime < 0 {
		t.Errorf("got runtime %f, expected non-zero", gotStatus.Runtime)
	}
}

func TestCmdStop(t *testing.T) {
	// Count to 3 sleeping 5s between counts. The long sleep is because we want
	// to kill the proc right after count "1" to ensure Stdout only contains "1"
	// and also to ensure that the proc is really killed instantly because if
	// it's not then timeout below will trigger.
	p := cmd.NewCmd("./test/count-and-sleep", "3", "5")

	// Start process in bg and get chan to receive final Status when done
	statusChan := p.Start()

	// Give it a second
	time.Sleep(1 * time.Second)

	// Kill the process
	err := p.Stop()
	if err != nil {
		t.Error(err)
	}

	// The final status should be returned instantly
	timeout := time.After(1 * time.Second)
	var gotStatus cmd.Status
	select {
	case gotStatus = <-statusChan:
	case <-timeout:
		t.Fatal("timeout waiting for statusChan")
	}

	start := time.Unix(0, gotStatus.StartTs)
	stop := time.Unix(0, gotStatus.StopTs)
	d := stop.Sub(start).Seconds()
	if d < 0.90 || d > 2 {
		t.Errorf("stop - start time not between 0.9s and 2.0s: %s - %s = %f", stop, start, d)
	}
	gotStatus.StartTs = 0
	gotStatus.StopTs = 0

	expectStatus := cmd.Status{
		Cmd:      "./test/count-and-sleep",
		PID:      gotStatus.PID,                    // nondeterministic
		Complete: false,                            // signaled by Stop
		Exit:     -1,                               // signaled by Stop
		Error:    errors.New("signal: terminated"), // signaled by Stop
		Runtime:  gotStatus.Runtime,                // nondeterministic
		Stdout:   []string{"1"},
		Stderr:   []string{},
	}
	if diffs := deep.Equal(gotStatus, expectStatus); diffs != nil {
		t.Error(diffs)
	}
	if gotStatus.PID < 0 {
		t.Errorf("got PID %d, expected non-zero", gotStatus.PID)
	}
	if gotStatus.Runtime < 0 {
		t.Errorf("got runtime %f, expected non-zero", gotStatus.Runtime)
	}

	// Stop should be idempotent
	err = p.Stop()
	if err != nil {
		t.Error(err)
	}

	// Start should be idempotent, too. It just returns the same statusChan again.
	c2 := p.Start()
	if diffs := deep.Equal(statusChan, c2); diffs != nil {
		t.Error(diffs)
	}
}

func TestCmdNotStarted(t *testing.T) {
	// Call everything _but_ Start.
	p := cmd.NewCmd("echo", "foo")

	gotStatus := p.Status()
	expectStatus := cmd.Status{
		Cmd:      "echo",
		PID:      0,
		Complete: false,
		Exit:     -1,
		Error:    nil,
		Runtime:  0,
		Stdout:   nil,
		Stderr:   nil,
	}
	if diffs := deep.Equal(gotStatus, expectStatus); diffs != nil {
		t.Error(diffs)
	}

	err := p.Stop()
	if err != nil {
		t.Error(err)
	}
}

func TestCmdOutput(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "cmd.TestCmdOutput")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	t.Logf("temp file: %s", tmpfile.Name())
	os.Remove(tmpfile.Name())

	p := cmd.NewCmd("./test/touch-file-count", tmpfile.Name())

	p.Start()

	touchFile := func(file string) {
		if err := exec.Command("touch", file).Run(); err != nil {
			t.Fatal(err)
		}
		time.Sleep(600 * time.Millisecond)
	}
	var s cmd.Status
	var stdout []string

	touchFile(tmpfile.Name())
	s = p.Status()
	stdout = []string{"1"}
	if diffs := deep.Equal(s.Stdout, stdout); diffs != nil {
		t.Log(s.Stdout)
		t.Error(diffs)
	}

	touchFile(tmpfile.Name())
	s = p.Status()
	stdout = []string{"1", "2"}
	if diffs := deep.Equal(s.Stdout, stdout); diffs != nil {
		t.Log(s.Stdout)
		t.Error(diffs)
	}

	// No more output yet
	s = p.Status()
	stdout = []string{"1", "2"}
	if diffs := deep.Equal(s.Stdout, stdout); diffs != nil {
		t.Log(s.Stdout)
		t.Error(diffs)
	}

	// +2 lines
	touchFile(tmpfile.Name())
	touchFile(tmpfile.Name())
	s = p.Status()
	stdout = []string{"1", "2", "3", "4"}
	if diffs := deep.Equal(s.Stdout, stdout); diffs != nil {
		t.Log(s.Stdout)
		t.Error(diffs)
	}

	// Kill the process
	if err := p.Stop(); err != nil {
		t.Error(err)
	}
}

func TestCmdNotFound(t *testing.T) {
	p := cmd.NewCmd("cmd-does-not-exist")
	gotStatus := <-p.Start()
	gotStatus.StartTs = 0
	gotStatus.StopTs = 0
	expectStatus := cmd.Status{
		Cmd:      "cmd-does-not-exist",
		PID:      0,
		Complete: false,
		Exit:     -1,
		Error:    errors.New(`exec: "cmd-does-not-exist": executable file not found in $PATH`),
		Runtime:  0,
		Stdout:   nil,
		Stderr:   nil,
	}
	if diffs := deep.Equal(gotStatus, expectStatus); diffs != nil {
		t.Logf("%+v", gotStatus)
		t.Error(diffs)
	}
}

func TestCmdLost(t *testing.T) {
	// Test something like the kernel OOM killing the proc. So the proc is
	// stopped outside our control.
	p := cmd.NewCmd("./test/count-and-sleep", "3", "5")

	statusChan := p.Start()

	// Give it a second
	time.Sleep(1 * time.Second)

	// Get the PID and kill it
	s := p.Status()
	if s.PID <= 0 {
		t.Fatalf("got PID %d, expected PID > 0", s.PID)
	}
	pgid, err := syscall.Getpgid(s.PID)
	if err != nil {
		t.Fatal(err)
	}
	syscall.Kill(-pgid, syscall.SIGKILL) // -pid = process group of pid

	// Even though killed externally, our wait should return instantly
	timeout := time.After(1 * time.Second)
	var gotStatus cmd.Status
	select {
	case gotStatus = <-statusChan:
	case <-timeout:
		t.Fatal("timeout waiting for statusChan")
	}
	gotStatus.Runtime = 0 // nondeterministic
	gotStatus.StartTs = 0
	gotStatus.StopTs = 0

	expectStatus := cmd.Status{
		Cmd:      "./test/count-and-sleep",
		PID:      s.PID,
		Complete: false,
		Exit:     -1,
		Error:    errors.New("signal: killed"),
		Runtime:  0,
		Stdout:   []string{"1"},
		Stderr:   []string{},
	}
	if diffs := deep.Equal(gotStatus, expectStatus); diffs != nil {
		t.Logf("%+v\n", gotStatus)
		t.Error(diffs)
	}
}

func TestCmdBothOutput(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "cmd.TestStreamingOutput")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(tmpfile.Name()); err != nil {
		t.Fatal(err)
	}

	touchFile := func(file string) {
		if err := exec.Command("touch", file).Run(); err != nil {
			t.Fatal(err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Streams a count to stdout and stderr until given file exists
	// Output like:
	//   stdout 1
	//   stderr 1
	//   stdout 2
	//   stderr 2
	// Where each is printed on stdout and stderr as indicated.
	p := cmd.NewCmdOptions(cmd.Options{Buffered: true, Streaming: true}, "./test/stream", tmpfile.Name())
	p.Start()
	time.Sleep(250 * time.Millisecond) // give test/stream a moment to print something

	timeout := time.After(10 * time.Second) // test timeout

	// test/stream is spewing output, so we should be able to read it while
	// the cmd is running. Try and fetch 3 lines from stdout and stderr.
	i := 0
	stdoutPrevLine := ""
	stderrPrevLine := ""
	readLines := 3
	lines := 0
	for i < readLines {
		i++

		// STDOUT
		select {
		case curLine := <-p.Stdout:
			t.Logf("got line: '%s'", curLine)
			if curLine == "" {
				// Shouldn't happen because test/stream doesn't print empty lines.
				// This indicates a bug in the stream buffer handling.
				t.Fatal("got empty line")
			}
			if stdoutPrevLine != "" && curLine == stdoutPrevLine {
				t.Fatalf("current line == previous line, expected new output:\ncprev: %s\ncur: %s\n", stdoutPrevLine, curLine)
			}
			stdoutPrevLine = curLine
			lines++
		case <-timeout:
			t.Fatal("timeout reading streaming output")
		default:
		}

		// STDERR
		select {
		case curLine := <-p.Stderr:
			t.Logf("got line: '%s'", curLine)
			if curLine == "" {
				// Shouldn't happen because test/stream doesn't print empty lines.
				// This indicates a bug in the stream buffer handling.
				t.Fatal("got empty line")
			}
			if stderrPrevLine != "" && curLine == stderrPrevLine {
				t.Fatalf("current line == previous line, expected new output:\ncprev: %s\ncur: %s\n", stderrPrevLine, curLine)
			}
			stderrPrevLine = curLine
			lines++
		case <-timeout:
			t.Fatal("timeout reading streaming output")
		default:
		}

		time.Sleep(200 * time.Millisecond)
	}

	// readLines * 2 (stdout and stderr)
	if lines != readLines*2 {
		t.Fatalf("read %d lines from streaming output, expected 6", lines)
	}

	s := p.Status()
	if len(s.Stdout) < readLines {
		t.Fatalf("read %d lines from buffered STDOUT, expected %d", len(s.Stdout), readLines)
	}
	if len(s.Stderr) < readLines {
		t.Fatalf("read %d lines from buffered STDERR, expected %d", len(s.Stderr), readLines)
	}

	// Stop test/stream
	touchFile(tmpfile.Name())

	s = p.Status()
	if s.Exit != 0 {
		t.Errorf("got exit %d, expected 0", s.Exit)
	}

	// Kill the process
	if err := p.Stop(); err != nil {
		t.Error(err)
	}
}

func TestCmdOnlyStreamingOutput(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "cmd.TestStreamingOutput")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(tmpfile.Name()); err != nil {
		t.Fatal(err)
	}

	touchFile := func(file string) {
		if err := exec.Command("touch", file).Run(); err != nil {
			t.Fatal(err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Streams a count to stdout and stderr until given file exists
	// Output like:
	//   stdout 1
	//   stderr 1
	//   stdout 2
	//   stderr 2
	// Where each is printed on stdout and stderr as indicated.
	p := cmd.NewCmdOptions(cmd.Options{Buffered: false, Streaming: true}, "./test/stream", tmpfile.Name())
	p.Start()
	time.Sleep(250 * time.Millisecond) // give test/stream a moment to print something

	timeout := time.After(10 * time.Second) // test timeout

	// test/stream is spewing output, so we should be able to read it while
	// the cmd is running. Try and fetch 3 lines from stdout and stderr.
	i := 0
	stdoutPrevLine := ""
	stderrPrevLine := ""
	readLines := 3
	lines := 0
	for i < readLines {
		i++
		t.Log(i)

		// STDOUT
		select {
		case curLine := <-p.Stdout:
			t.Logf("got line: '%s'", curLine)
			if curLine == "" {
				// Shouldn't happen because test/stream doesn't print empty lines.
				// This indicates a bug in the stream buffer handling.
				t.Fatal("got empty line")
			}
			if stdoutPrevLine != "" && curLine == stdoutPrevLine {
				t.Fatalf("current line == previous line, expected new output:\ncprev: %s\ncur: %s\n", stdoutPrevLine, curLine)
			}
			stdoutPrevLine = curLine
			lines++
		case <-timeout:
			t.Fatal("timeout reading streaming output")
		default:
		}

		// STDERR
		select {
		case curLine := <-p.Stderr:
			t.Logf("got line: '%s'", curLine)
			if curLine == "" {
				// Shouldn't happen because test/stream doesn't print empty lines.
				// This indicates a bug in the stream buffer handling.
				t.Fatal("got empty line")
			}
			if stderrPrevLine != "" && curLine == stderrPrevLine {
				t.Fatalf("current line == previous line, expected new output:\ncprev: %s\ncur: %s\n", stderrPrevLine, curLine)
			}
			stderrPrevLine = curLine
			lines++
		case <-timeout:
			t.Fatal("timeout reading streaming output")
		default:
		}

		time.Sleep(200 * time.Millisecond)
	}

	// readLines * 2 (stdout and stderr)
	if lines != readLines*2 {
		t.Fatalf("read %d lines from streaming output, expected 6", lines)
	}

	s := p.Status()
	if len(s.Stdout) != 0 {
		t.Fatalf("read %d lines from buffered STDOUT, expected 0", len(s.Stdout))
	}
	if len(s.Stderr) != 0 {
		t.Fatalf("read %d lines from buffered STDERR, expected 0", len(s.Stderr))
	}

	// Stop test/stream
	touchFile(tmpfile.Name())

	s = p.Status()
	if s.Exit != 0 {
		t.Errorf("got exit %d, expected 0", s.Exit)
	}

	// Kill the process
	if err := p.Stop(); err != nil {
		t.Error(err)
	}
}

func TestStreamingMultipleLines(t *testing.T) {
	lines := make(chan string, 5)
	out := cmd.NewOutputStream(lines)

	// Quick side test: Lines() chan string should be the same chan string
	// we created the object with
	if out.Lines() != lines {
		t.Errorf("Lines() does not return the given string chan")
	}

	// Write two short lines
	input := "foo\nbar\n"
	n, err := out.Write([]byte(input))
	if n != len(input) {
		t.Errorf("Write n = %d, expected %d", n, len(input))
	}
	if err != nil {
		t.Errorf("got err '%v', expected nil", err)
	}

	// Get one line
	var gotLine string
	select {
	case gotLine = <-lines:
	default:
		t.Fatal("blocked on <-lines")
	}

	// "foo" should be sent before "bar" because that was the input
	if gotLine != "foo" {
		t.Errorf("got line: '%s', expected 'foo'", gotLine)
	}

	// Get next line
	select {
	case gotLine = <-lines:
	default:
		t.Fatal("blocked on <-lines")
	}

	if gotLine != "bar" {
		t.Errorf("got line: '%s', expected 'bar'", gotLine)
	}
}

func TestStreamingBlankLines(t *testing.T) {
	lines := make(chan string, 5)
	out := cmd.NewOutputStream(lines)

	// Blank line in the middle
	input := "foo\n\nbar\n"
	expectLines := []string{"foo", "", "bar"}
	gotLines := []string{}
	n, err := out.Write([]byte(input))
	if n != len(input) {
		t.Errorf("Write n = %d, expected %d", n, len(input))
	}
	if err != nil {
		t.Errorf("got err '%v', expected nil", err)
	}
LINES1:
	for {
		select {
		case line := <-lines:
			gotLines = append(gotLines, line)
		default:
			break LINES1
		}
	}
	if diffs := deep.Equal(gotLines, expectLines); diffs != nil {
		t.Error(diffs)
	}

	// All blank lines
	input = "\n\n\n"
	expectLines = []string{"", "", ""}
	gotLines = []string{}
	n, err = out.Write([]byte(input))
	if n != len(input) {
		t.Errorf("Write n = %d, expected %d", n, len(input))
	}
	if err != nil {
		t.Errorf("got err '%v', expected nil", err)
	}
LINES2:
	for {
		select {
		case line := <-lines:
			gotLines = append(gotLines, line)
		default:
			break LINES2
		}
	}
	if diffs := deep.Equal(gotLines, expectLines); diffs != nil {
		t.Error(diffs)
	}

	// Blank lines at end
	input = "foo\n\n\n"
	expectLines = []string{"foo", "", ""}
	gotLines = []string{}
	n, err = out.Write([]byte(input))
	if n != len(input) {
		t.Errorf("Write n = %d, expected %d", n, len(input))
	}
	if err != nil {
		t.Errorf("got err '%v', expected nil", err)
	}
LINES3:
	for {
		select {
		case line := <-lines:
			gotLines = append(gotLines, line)
		default:
			break LINES3
		}
	}
	if diffs := deep.Equal(gotLines, expectLines); diffs != nil {
		t.Error(diffs)
	}
}

func TestStreamingCarriageReturn(t *testing.T) {
	// Carriage return should be stripped
	lines := make(chan string, 5)
	out := cmd.NewOutputStream(lines)

	input := "foo\r\nbar\r\n"
	expectLines := []string{"foo", "bar"}
	gotLines := []string{}
	n, err := out.Write([]byte(input))
	if n != len(input) {
		t.Errorf("Write n = %d, expected %d", n, len(input))
	}
	if err != nil {
		t.Errorf("got err '%v', expected nil", err)
	}
LINES1:
	for {
		select {
		case line := <-lines:
			gotLines = append(gotLines, line)
		default:
			break LINES1
		}
	}
	if diffs := deep.Equal(gotLines, expectLines); diffs != nil {
		t.Error(diffs)
	}
}

func TestStreamingLineBuffering(t *testing.T) {
	// Lines not terminated with newline are held in the line buffer until next
	// write. When line is later terminated with newline, we prepend the buffered
	// line and send the complete line.
	lines := make(chan string, 1)
	out := cmd.NewOutputStream(lines)

	// Write 3 unterminated lines. Without a newline, they'll be buffered until...
	for i := 0; i < 3; i++ {
		input := "foo"
		n, err := out.Write([]byte(input))
		if err != nil {
			t.Errorf("got err '%v', expected nil", err)
		}
		if n != len(input) {
			t.Errorf("Write n = %d, expected %d", n, len(input))
		}

		// Should not get a line yet because it's not newline terminated
		var gotLine string
		select {
		case gotLine = <-lines:
			t.Errorf("got line '%s', expected no line yet", gotLine)
		default:
		}
	}

	// Write a line with newline that terminate the previous input
	input := "bar\n"
	n, err := out.Write([]byte(input))
	if err != nil {
		t.Errorf("got err '%v', expected nil", err)
	}
	if n != len(input) {
		t.Errorf("Write n = %d, expected %d", n, len(input))
	}

	// Now we get the previously buffered part of the line "foofoofoo" plus
	// the newline terminated part "bar"
	var gotLine string
	select {
	case gotLine = <-lines:
	default:
		t.Fatal("blocked receiving line")
	}

	expectLine := "foofoofoobar"
	if gotLine != expectLine {
		t.Errorf("got line '%s', expected '%s'", gotLine, expectLine)
	}
}

func TestStreamingErrLineBufferOverflow1(t *testing.T) {
	// Overflow the line buffer in 1 write. The first line "bc" is sent,
	// but the remaining line can't be buffered because it's +2 bytes larger
	// than the line buffer.
	longLine := make([]byte, 3+cmd.DEFAULT_LINE_BUFFER_SIZE+2) // "bc\nAAA...zz"
	longLine[0] = 'b'
	longLine[1] = 'c'
	longLine[2] = '\n'
	for i := 3; i < cmd.DEFAULT_LINE_BUFFER_SIZE; i++ {
		longLine[i] = 'A'
	}
	// These 2 chars cause ErrLineBufferOverflow:
	longLine[cmd.DEFAULT_LINE_BUFFER_SIZE] = 'z'
	longLine[cmd.DEFAULT_LINE_BUFFER_SIZE+1] = 'z'

	lines := make(chan string, 5)
	out := cmd.NewOutputStream(lines)

	// Write the long line, it should only write (n) 3 bytes for "bc\n"
	n, err := out.Write(longLine)
	if n != 3 { // "bc\n"
		t.Errorf("Write n = %d, expected 3", n)
	}
	switch err.(type) {
	case cmd.ErrLineBufferOverflow:
		lbo := err.(cmd.ErrLineBufferOverflow)
		if lbo.BufferSize != cmd.DEFAULT_LINE_BUFFER_SIZE {
			t.Errorf("ErrLineBufferOverflow.BufferSize = %d, expected %d", lbo.BufferSize, cmd.DEFAULT_LINE_BUFFER_SIZE)
		}
		if lbo.BufferFree != cmd.DEFAULT_LINE_BUFFER_SIZE {
			t.Errorf("ErrLineBufferOverflow.BufferFree = %d, expected %d", lbo.BufferFree, cmd.DEFAULT_LINE_BUFFER_SIZE)
		}
		if lbo.Line != string(longLine[3:]) {
			t.Errorf("ErrLineBufferOverflow.Line = '%s', expected '%s'", lbo.Line, string(longLine[3:]))
		}
		if lbo.Error() == "" {
			t.Errorf("ErrLineBufferOverflow.Error() string is empty, expected something")
		}
	default:
		t.Errorf("got err '%v', expected cmd.ErrLineBufferOverflow", err)
	}

	// "bc" should be sent before the overflow error
	var gotLine string
	select {
	case gotLine = <-lines:
	default:
		t.Fatal("blocked on <-lines")
	}
	if gotLine != "bc" {
		t.Errorf("got line '%s', expected 'bc'", gotLine)
	}

	// Streaming should still work after an overflow. However, Go is going to
	// stop any time Write() returns an error.
	n, err = out.Write([]byte("foo\n"))
	if n != 4 {
		t.Errorf("got n %d, expected 4", n)
	}
	if err != nil {
		t.Errorf("got err '%v', expected nil", err)
	}

	select {
	case gotLine = <-lines:
	default:
		t.Fatal("blocked on <-lines")
	}
	if gotLine != "foo" {
		t.Errorf("got line: '%s', expected 'foo'", gotLine)
	}
}

func TestStreamingErrLineBufferOverflow2(t *testing.T) {
	// Overflow line buffer on 2nd write. So first write puts something in the
	// buffer, and then 2nd overflows it instead of completing the line.
	lines := make(chan string, 1)
	out := cmd.NewOutputStream(lines)

	// Get "bar" into the buffer by omitting its newline
	input := "foo\nbar"
	n, err := out.Write([]byte(input))
	if err != nil {
		t.Errorf("got err '%v', expected nil", err)
	}
	if n != len(input) {
		t.Errorf("Write n = %d, expected %d", n, len(input))
	}

	// Only "foo" sent, not "bar" yet
	var gotLine string
	select {
	case gotLine = <-lines:
	default:
		t.Fatal("blocked on <-lines")
	}
	if gotLine != "foo" {
		t.Errorf("got line '%s', expected 'foo'", gotLine)
	}

	// Buffer contains "bar", now wverflow it on 2nd write
	longLine := make([]byte, cmd.DEFAULT_LINE_BUFFER_SIZE)
	for i := 0; i < cmd.DEFAULT_LINE_BUFFER_SIZE; i++ {
		longLine[i] = 'X'
	}
	n, err = out.Write(longLine)
	if n != 0 {
		t.Errorf("Write n = %d, expected 0", n)
	}
	switch err.(type) {
	case cmd.ErrLineBufferOverflow:
		lbo := err.(cmd.ErrLineBufferOverflow)
		// Buffer has "bar" so it's free is total - 3
		if lbo.BufferFree != cmd.DEFAULT_LINE_BUFFER_SIZE-3 {
			t.Errorf("ErrLineBufferOverflow.BufferFree = %d, expected %d", lbo.BufferFree, cmd.DEFAULT_LINE_BUFFER_SIZE)
		}
		// Up to but not include "bc\n" because it should have been truncated
		expectLine := "bar" + string(longLine)
		if lbo.Line != expectLine {
			t.Errorf("ErrLineBufferOverflow.Line = '%s', expected '%s'", lbo.Line, expectLine)
		}
	default:
		t.Errorf("got err '%v', expected cmd.ErrLineBufferOverflow", err)
	}
}

func TestStreamingSetLineBufferSize(t *testing.T) {
	// Same overflow as TestStreamingErrLineBufferOverflow1 but before we use
	// stream output, we'll increase buffer size by calling SetLineBufferSize
	// which should prevent the overflow
	longLine := make([]byte, 3+cmd.DEFAULT_LINE_BUFFER_SIZE+2) // "bc\nAAA...z\n"
	longLine[0] = 'b'
	longLine[1] = 'c'
	longLine[2] = '\n'
	for i := 3; i < cmd.DEFAULT_LINE_BUFFER_SIZE; i++ {
		longLine[i] = 'A'
	}
	longLine[cmd.DEFAULT_LINE_BUFFER_SIZE] = 'z'
	longLine[cmd.DEFAULT_LINE_BUFFER_SIZE+1] = '\n'

	lines := make(chan string, 5)
	out := cmd.NewOutputStream(lines)
	out.SetLineBufferSize(cmd.DEFAULT_LINE_BUFFER_SIZE * 2)

	n, err := out.Write(longLine)
	if err != nil {
		t.Errorf("error '%v', expected nil", err)
	}
	if n != len(longLine) {
		t.Errorf("Write n = %d, expected %d", n, len(longLine))
	}

	// First we get "bc"
	var gotLine string
	select {
	case gotLine = <-lines:
	default:
		t.Fatal("blocked on <-lines")
	}
	if gotLine != "bc" {
		t.Errorf("got line '%s', expected 'bc'", gotLine)
	}

	// Then we get the long line because the buffer was large enough to hold it
	select {
	case gotLine = <-lines:
	default:
		t.Fatal("blocked on <-lines")
	}
	expectLine := string(longLine[3 : cmd.DEFAULT_LINE_BUFFER_SIZE+1]) // not newline
	if gotLine != expectLine {
		t.Errorf("got line: '%s', expected '%s'", gotLine, expectLine)
	}
}

func TestDone(t *testing.T) {
	// Count to 3 sleeping 1s between counts
	p := cmd.NewCmd("./test/count-and-sleep", "3", "1")
	statusChan := p.Start()

	// For 2s while cmd is running, Done() chan should block, which means
	// it's still running
	runningTimer := time.After(2 * time.Second)
TIMER:
	for {
		select {
		case <-runningTimer:
			break TIMER
		default:
		}
		select {
		case <-p.Done():
			t.Fatal("Done chan is closed before runningTime finished")
		default:
			// Done chan blocked, cmd is still running
		}
		time.Sleep(400 * time.Millisecond)
	}

	// Wait for cmd to complete
	var s1 cmd.Status
	select {
	case s1 = <-statusChan:
		t.Logf("got status: %+v", s1)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for cmd to complete")
	}

	// After cmd completes, Done chan should be closed and not block
	select {
	case <-p.Done():
	default:
		t.Fatal("Done chan did not block after cmd completed")
	}

	// After command completes, we should be able to get exact same
	// Status that's returned on the Start() chan
	s2 := p.Status()
	if diff := deep.Equal(s1, s2); diff != nil {
		t.Error(diff)
	}
}

func TestCmdEnvOK(t *testing.T) {
	now := time.Now().Unix()

	p := cmd.NewCmd("env")
	p.Env = []string{"FOO=foo"}
	gotStatus := <-p.Start()
	expectStatus := cmd.Status{
		Cmd:      "env",
		PID:      gotStatus.PID, // nondeterministic
		Complete: true,
		Exit:     0,
		Error:    nil,
		Runtime:  gotStatus.Runtime, // nondeterministic
		Stdout:   []string{"FOO=foo"},
		Stderr:   []string{},
	}
	if gotStatus.StartTs < now {
		t.Error("StartTs < now")
	}
	if gotStatus.StopTs < gotStatus.StartTs {
		t.Error("StopTs < StartTs")
	}
	gotStatus.StartTs = 0
	gotStatus.StopTs = 0
	if diffs := deep.Equal(gotStatus, expectStatus); diffs != nil {
		t.Error(diffs)
	}
	if gotStatus.PID < 0 {
		t.Errorf("got PID %d, expected non-zero", gotStatus.PID)
	}
	if gotStatus.Runtime < 0 {
		t.Errorf("got runtime %f, expected non-zero", gotStatus.Runtime)
	}
}

func TestCmdNoOutput(t *testing.T) {
	// Set both output options to false to discard all output
	p := cmd.NewCmdOptions(
		cmd.Options{
			Buffered:  false,
			Streaming: false,
		},
		"echo", "hell-world")
	s := <-p.Start()
	if s.Exit != 0 {
		t.Errorf("got exit %d, expected 0", s.Exit)
	}
	if len(s.Stdout) != 0 {
		t.Errorf("got stdout, expected no output: %v", s.Stdout)
	}
	if len(s.Stderr) != 0 {
		t.Errorf("got stderr, expected no output: %v", s.Stderr)
	}
}
