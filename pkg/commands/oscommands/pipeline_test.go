package oscommands

import (
	"fmt"
	"io"
	"os"
	"strings"
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
		os.Exit(3)
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

	os.Exit(0)
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
