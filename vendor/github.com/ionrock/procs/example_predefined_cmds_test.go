package procs_test

import (
	"fmt"
	"os/exec"

	"github.com/ionrock/procs"
)

func Example_predefinedCmds() {
	p := procs.Process{
		Cmds: []*exec.Cmd{
			exec.Command("echo", "foo"),
			exec.Command("grep", "foo"),
		},
	}

	p.Run()
	out, _ := p.Output()
	fmt.Println(string(out))
	// Output:
	// foo
}
