package procs_test

import (
	"fmt"

	"github.com/ionrock/procs"
)

func ExampleSplitCommand() {
	parts := procs.SplitCommand("echo 'hello world'")
	for i, p := range parts {
		fmt.Printf("%d %s\n", i+1, p)
	}

	// Output:
	// 1 echo
	// 2 hello world
}

func ExampleSplitCommandEnv() {
	env := map[string]string{
		"GREETING": "hello",
		"NAME":     "world!",
		"PASSWORD": "secret",
	}

	getenv := func(key string) string {
		if v, ok := env[key]; ok && key != "PASSWORD" {
			return v
		}
		return ""
	}

	parts := procs.SplitCommandEnv("echo '$GREETING $NAME $PASSWORD'", getenv)

	for i, p := range parts {
		fmt.Printf("%d %s\n", i+1, p)
	}

	// Output:
	// 1 echo
	// 2 hello world!
}
