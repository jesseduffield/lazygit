package procs_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ionrock/procs"
)

func TestParseEnv(t *testing.T) {
	env := []string{
		"FOO=bar",
		"BAZ=`echo 'hello=world'`",
	}

	m := procs.ParseEnv(env)

	v, ok := m["FOO"]
	if !ok {
		t.Errorf("error missing FOO from env: %#v", m)
	}

	if v != "bar" {
		t.Errorf("error FOO != bar: %s", v)
	}

	v, ok = m["BAZ"]

	if !ok {
		t.Errorf("error missing BAZ from env: %#v", m)
	}

	expectBaz := "`echo 'hello=world'`"
	if v != expectBaz {
		t.Errorf("error BAZ != %s: %s", expectBaz, v)
	}
}

func TestEnvBuilder(t *testing.T) {
	env := procs.Env(map[string]string{
		"FOO": "bar",
		"BAZ": "hello world",
	}, false)

	if len(env) != 2 {
		t.Errorf("error loading env: %s", env)
	}
}

func helperEnvCommand(env map[string]string) *exec.Cmd {
	cmd := exec.Command(os.Args[0], "-test.run=TestEnvBuilderOverrides")
	cmd.Env = procs.Env(env, false)
	return cmd
}

func TestEnvBuilderOverrides(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	for _, envvar := range procs.Env(map[string]string{"FOO": "override"}, true) {
		fmt.Println(envvar)
	}
}

func TestEnvBuilderWithEnv(t *testing.T) {
	cmd := helperEnvCommand(map[string]string{
		"GO_WANT_HELPER_PROCESS": "1",
		"FOO": "default",
	})
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("error running helper: %s", err)
	}

	env := procs.ParseEnv(strings.Split(string(out), "\n"))

	if env["FOO"] != "override" {
		t.Errorf("error overriding envvar: %s", string(out))
	}
}
