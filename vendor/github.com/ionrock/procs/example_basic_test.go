package procs_test

import (
	"fmt"

	"github.com/ionrock/procs"
)

func Example() {

	b := procs.Builder{
		Context: map[string]string{
			"NAME": "eric",
		},
		Templates: []string{
			"echo $NAME |",
			"grep $NAME",
		},
	}

	cmd := b.Command()

	fmt.Println(cmd)

	p := procs.NewProcess(cmd)

	p.Run()
	out, _ := p.Output()
	fmt.Println(string(out))
	// Output:
	// echo eric | grep eric
	// eric
}
