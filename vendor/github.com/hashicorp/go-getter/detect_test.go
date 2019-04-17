package getter

import (
	"fmt"
	"testing"
)

func TestDetect(t *testing.T) {
	cases := []struct {
		Input  string
		Pwd    string
		Output string
		Err    bool
	}{
		{"./foo", "/foo", "file:///foo/foo", false},
		{"git::./foo", "/foo", "git::file:///foo/foo", false},
		{
			"git::github.com/hashicorp/foo",
			"",
			"git::https://github.com/hashicorp/foo.git",
			false,
		},
		{
			"./foo//bar",
			"/foo",
			"file:///foo/foo//bar",
			false,
		},
		{
			"git::github.com/hashicorp/foo//bar",
			"",
			"git::https://github.com/hashicorp/foo.git//bar",
			false,
		},
		{
			"git::https://github.com/hashicorp/consul.git",
			"",
			"git::https://github.com/hashicorp/consul.git",
			false,
		},
		{
			"git::https://person@someothergit.com/foo/bar",
			"",
			"git::https://person@someothergit.com/foo/bar",
			false,
		},
		{
			"git::https://person@someothergit.com/foo/bar",
			"/bar",
			"git::https://person@someothergit.com/foo/bar",
			false,
		},
		{
			"./foo/archive//*",
			"/bar",
			"file:///bar/foo/archive//*",
			false,
		},

		// https://github.com/hashicorp/go-getter/pull/124
		{
			"git::ssh://git@my.custom.git/dir1/dir2",
			"",
			"git::ssh://git@my.custom.git/dir1/dir2",
			false,
		},
		{
			"git::git@my.custom.git:dir1/dir2",
			"/foo",
			"git::ssh://git@my.custom.git/dir1/dir2",
			false,
		},
		{
			"git::git@my.custom.git:dir1/dir2",
			"",
			"git::ssh://git@my.custom.git/dir1/dir2",
			false,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d %s", i, tc.Input), func(t *testing.T) {
			output, err := Detect(tc.Input, tc.Pwd, Detectors)
			if err != nil != tc.Err {
				t.Fatalf("%d: bad err: %s", i, err)
			}
			if output != tc.Output {
				t.Fatalf("%d: bad output: %s\nexpected: %s", i, output, tc.Output)
			}
		})
	}
}
