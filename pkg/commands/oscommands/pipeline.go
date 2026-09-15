package oscommands

import (
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
// A | B | C. The last command's output is left for the caller to direct.
func wirePipeline(cmdObjs []*CmdObj) ([]*exec.Cmd, error) {
	cmds := lo.Map(cmdObjs, func(cmdObj *CmdObj, _ int) *exec.Cmd {
		return cmdObj.GetCmd()
	})

	for i := range len(cmds) - 1 {
		stdout, err := cmds[i].StdoutPipe()
		if err != nil {
			return nil, err
		}

		cmds[i+1].Stdin = stdout
	}

	return cmds, nil
}

// startPipeline starts every command and reports how many it got going. Every
// one is started before any of them is waited for: waiting closes our end of
// the pipe that feeds the next command, and one that hasn't been started by
// then would inherit a closed stdin.
//
// When a command fails to start, the ones already running are killed, since
// without the rest of the pipeline to drain them they could block forever
// writing to a full pipe. They still have to be reaped, so the count covers
// them too.
func startPipeline(cmds []*exec.Cmd) (int, error) {
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
