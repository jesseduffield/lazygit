package oscommands

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/samber/lo"
)

// pipelineString names a chain of commands the way a shell would write it.
func pipelineString(cmdObjs []*CmdObj) string {
	return strings.Join(
		lo.Map(cmdObjs, func(cmdObj *CmdObj, _ int) string {
			return cmdObj.ToString()
		}),
		" | ",
	)
}

// wirePipeline connects each command's output to the next one's input, like
// A | B | C, and returns the commands along with the parent's ends of those
// pipes. The last command's output is left for the caller to direct.
//
// The parent's ends have to be closed once the commands are running.
// startPipeline does that; see there for why it matters.
func wirePipeline(cmdObjs []*CmdObj) ([]*exec.Cmd, []io.Closer, error) {
	cmds := lo.Map(cmdObjs, func(cmdObj *CmdObj, _ int) *exec.Cmd {
		return cmdObj.GetCmd()
	})

	parentEnds := []io.Closer{}
	for i := range len(cmds) - 1 {
		reader, writer, err := os.Pipe()
		if err != nil {
			closeAll(parentEnds)
			return nil, nil, err
		}

		cmds[i].Stdout = writer
		cmds[i+1].Stdin = reader
		parentEnds = append(parentEnds, reader, writer)
	}

	return cmds, parentEnds, nil
}

// startPipeline starts every command and reports how many it got going. Every
// one is started before any of them is waited for: waiting closes our end of
// the pipe that feeds the next command, and one that hasn't been started by
// then would inherit a closed stdin.
//
// Once they are all running, each of them holds its own ends of the pipes it
// reads and writes, and the parent lets go of its copies. Both directions
// matter. While the parent holds the read end of a link, a command writing
// into it never learns that the command meant to read it is gone, and keeps
// running after the pipeline has been brought down. While the parent holds the
// write end, the command reading it never reaches the end of its input.
//
// When a command fails to start, the ones already running are killed, since
// without the rest of the pipeline to drain them they could block forever
// writing to a full pipe. They still have to be reaped, so the count covers
// them too.
func startPipeline(cmds []*exec.Cmd, parentEnds []io.Closer) (int, error) {
	defer closeAll(parentEnds)

	for i, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			for _, started := range cmds[:i] {
				_ = started.Process.Kill()
			}

			return i, err
		}
	}

	return len(cmds), nil
}

func closeAll(closers []io.Closer) {
	for _, closer := range closers {
		_ = closer.Close()
	}
}
