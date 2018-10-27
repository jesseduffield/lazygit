package procs_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/ionrock/procs"
)

func newProcess() *procs.Process {
	return &procs.Process{
		Cmds: []*exec.Cmd{
			exec.Command("echo", "foo"),
			exec.Command("grep", "foo"),
		},
	}
}

func TestProcess(t *testing.T) {
	p := newProcess()

	err := p.Run()
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}

	out, _ := p.Output()
	if !bytes.Equal(bytes.TrimSpace(out), []byte("foo")) {
		t.Errorf("wrong output: expected foo but got %s", out)
	}
}

func TestProcessWithOutput(t *testing.T) {
	p := newProcess()

	p.OutputHandler = func(line string) string {
		return fmt.Sprintf("x | %s", line)
	}

	err := p.Run()

	if err != nil {
		t.Fatalf("error running program: %s", err)
	}
	expected := []byte("x | foo")
	out, _ := p.Output()
	if !bytes.Equal(bytes.TrimSpace(out), expected) {
		t.Errorf("wrong output: expected %q but got %q", expected, out)
	}
}

func TestProcessStartAndWait(t *testing.T) {
	p := newProcess()

	p.Start()
	p.Wait()

	out, _ := p.Output()
	expected := []byte("foo")
	if !bytes.Equal(bytes.TrimSpace(out), expected) {
		t.Errorf("wrong output: expected %q but got %q", expected, out)
	}
}

func TestProcessStartAndWaitWithOutput(t *testing.T) {
	p := newProcess()
	p.OutputHandler = func(line string) string {
		return fmt.Sprintf("x | %s", line)
	}

	p.Start()
	p.Wait()

	out, _ := p.Output()
	expected := []byte("x | foo")
	if !bytes.Equal(bytes.TrimSpace(out), expected) {
		t.Errorf("wrong output: expected %q but got %q", expected, out)
	}
}

func TestProcessFromString(t *testing.T) {
	p := procs.NewProcess("echo 'foo'")
	err := p.Run()
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}
	out, _ := p.Output()
	if !bytes.Equal(bytes.TrimSpace(out), []byte("foo")) {
		t.Errorf("wrong output: expected foo but got %s", out)
	}
}

func TestProcessFromStringWithPipe(t *testing.T) {
	p := procs.NewProcess("echo 'foo' | grep foo")
	err := p.Run()
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}

	out, _ := p.Output()
	if !bytes.Equal(bytes.TrimSpace(out), []byte("foo")) {
		t.Errorf("wrong output: expected foo but got %s", out)
	}
}

func TestStderrOutput(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	fmt.Fprintln(os.Stdout, "stdout output")
	fmt.Fprintln(os.Stderr, "stderr output")
	os.Exit(1)
}

func TestProcessPipeWithFailures(t *testing.T) {

	// This will run a piped command with a failure part way
	// through. We want to be sure we get output on stderr.
	p := procs.NewProcess(fmt.Sprintf("echo 'foo' | %s -test.run=TestStderrOutput | grep foo", os.Args[0]))
	p.Env = map[string]string{"GO_WANT_HELPER_PROCESS": "1"}

	err := p.Run()
	if err == nil {
		t.Fatal("expected error running program")

	}

	out, _ := p.Output()
	expected := []byte("") // expecting no output b/c the grep foo won't run
	if !bytes.Equal(out, expected) {
		t.Errorf("wrong stdout output: expected '%s' but got '%s'", expected, out)
	}

	errOut, _ := p.ErrOutput()
	expected = []byte("stderr output")
	if !bytes.Equal(bytes.TrimSpace(errOut), expected) {
		t.Errorf("wrong stderr output: expected '%s' but got '%s'", expected, out)
	}
}
