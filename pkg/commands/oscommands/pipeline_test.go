package oscommands

import (
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// The pipeline tests need programs to run, and the test binary is the one
// program every platform we test on is sure to have. pipelineMember builds a
// command that re-runs this binary in the role a member of the pipeline is to
// play; the roles are in TestPipelineMember.
const pipelineRoleEnvVar = "LAZYGIT_TEST_PIPELINE_ROLE"

func pipelineMember(role string) *CmdObj {
	return NewDummyOSCommand().Cmd.
		New([]string{os.Args[0], "-test.run=^TestPipelineMember$"}).
		AddEnvVars(pipelineRoleEnvVar + "=" + role)
}

// TestPipelineMember is the program the pipeline tests run, not a test of its
// own. It exits before the testing package reports anything, so that its output
// is what the role wrote and nothing else.
//
// For the same reason it exits with syscall.Exit, which skips the exit hooks
// that os.Exit runs. In a binary built with -cover, one of these hooks writes
// coverage data to $GOCOVERDIR and prints an error to stderr if that fails. On
// Windows it fails now and then if two members of a pipeline exit at the same
// time, because both of them replace the same file in that directory.
func TestPipelineMember(t *testing.T) {
	switch os.Getenv(pipelineRoleEnvVar) {
	case "":
		t.Skip("not a test; the pipeline tests run this binary in a role")
	case "count":
		for i := 1; i <= 3; i++ {
			fmt.Printf("line %d\n", i)
		}
	case "upcase":
		input, _ := io.ReadAll(os.Stdin)
		fmt.Print(strings.ToUpper(string(input)))
	case "copy":
		_, _ = io.Copy(os.Stdout, os.Stdin)
	case "complain":
		fmt.Fprintln(os.Stderr, "something went wrong")
		syscall.Exit(3)
	case "flood":
		// A failed write means the reader is gone, and there is no point
		// writing to nobody. On platforms that raise a signal for it instead,
		// this process is already dead by the time the write returns.
		for i := 1; ; i++ {
			if _, err := fmt.Printf("line %d\n", i); err != nil {
				break
			}
		}
	}

	syscall.Exit(0)
}

func TestStartPipelineStreamsTheOutputOfTheLastCommand(t *testing.T) {
	pipeline, reader, err := NewDummyOSCommand().StartPipeline(
		pipelineMember("count"),
		pipelineMember("upcase"),
	)
	assert.NoError(t, err)

	output, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, "LINE 1\nLINE 2\nLINE 3\n", string(output))

	assert.NoError(t, pipeline.Wait())
	assert.NoError(t, reader.Close())
}

func TestStartPipelineReadsWhatTheCommandsComplainAbout(t *testing.T) {
	pipeline, reader, err := NewDummyOSCommand().StartPipeline(
		pipelineMember("count"),
		pipelineMember("complain"),
	)
	assert.NoError(t, err)

	output, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, "something went wrong\n", string(output))

	// The failure of the command nearest the output is the one reported, even
	// though the one feeding it was left writing into a pipe nobody reads.
	assert.ErrorContains(t, pipeline.Wait(), "exit status 3")

	assert.NoError(t, reader.Close())
}

func TestClosingAPipelinesOutputBringsItDown(t *testing.T) {
	pipeline, reader, err := NewDummyOSCommand().StartPipeline(
		pipelineMember("flood"),
		pipelineMember("copy"),
	)
	assert.NoError(t, err)

	// Read some output first, so that both commands are past their startup and
	// really running when the reader goes.
	buf := make([]byte, len("line 1\n"))
	_, err = io.ReadFull(reader, buf)
	assert.NoError(t, err)
	assert.Equal(t, "line 1\n", string(buf))

	assert.NoError(t, reader.Close())

	done := make(chan error, 1)
	go func() { done <- pipeline.Wait() }()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("the pipeline was still running long after its output was closed")
	}
}

func TestStartPipelineReportsACommandItCannotStart(t *testing.T) {
	osCommand := NewDummyOSCommand()

	_, _, err := osCommand.StartPipeline(
		pipelineMember("count"),
		osCommand.Cmd.New([]string{"lazygit-no-such-command"}),
	)
	assert.Error(t, err)
}

func TestPipeCommandsReturnsWhenALaterCommandDiesEarly(t *testing.T) {
	done := make(chan error, 1)
	go func() {
		done <- NewDummyOSCommand().PipeCommands(
			pipelineMember("flood"),
			pipelineMember("complain"),
		)
	}()

	select {
	case err := <-done:
		assert.ErrorContains(t, err, "something went wrong")
	case <-time.After(10 * time.Second):
		t.Fatal("PipeCommands was still waiting for a command whose output nothing reads")
	}
}
