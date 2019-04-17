package getter

import (
	"testing"
)

func TestGitDetector(t *testing.T) {
	cases := []struct {
		Input  string
		Output string
	}{
		{"git@github.com:hashicorp/foo.git", "git::ssh://git@github.com/hashicorp/foo.git"},
		{
			"git@github.com:org/project.git?ref=test-branch",
			"git::ssh://git@github.com/org/project.git?ref=test-branch",
		},
		{
			"git@github.com:hashicorp/foo.git//bar",
			"git::ssh://git@github.com/hashicorp/foo.git//bar",
		},
		{
			"git@github.com:hashicorp/foo.git?foo=bar",
			"git::ssh://git@github.com/hashicorp/foo.git?foo=bar",
		},
		{
			"git@github.xyz.com:org/project.git",
			"git::ssh://git@github.xyz.com/org/project.git",
		},
		{
			"git@github.xyz.com:org/project.git?ref=test-branch",
			"git::ssh://git@github.xyz.com/org/project.git?ref=test-branch",
		},
		{
			"git@github.xyz.com:org/project.git//module/a",
			"git::ssh://git@github.xyz.com/org/project.git//module/a",
		},
		{
			"git@github.xyz.com:org/project.git//module/a?ref=test-branch",
			"git::ssh://git@github.xyz.com/org/project.git//module/a?ref=test-branch",
		},
	}

	pwd := "/pwd"
	f := new(GitDetector)
	for i, tc := range cases {
		output, ok, err := f.Detect(tc.Input, pwd)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if !ok {
			t.Fatal("not ok")
		}

		if output != tc.Output {
			t.Fatalf("%d: bad: %#v", i, output)
		}
	}
}
