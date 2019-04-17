package roll

import (
	"strings"
	"testing"
)

func TestBuildRollbarFrames(t *testing.T) {
	frames := buildRollbarFrames(getCallers(0))

	if len(frames) != 4 {
		t.Fatalf("expected 4 frames, got %d", len(frames))
	}

	if !strings.Contains(frames[0].Filename, "github.com/stvp/roll/stack.go") {
		t.Errorf("expected %#v, got %#v", "github.com/stvp/roll/stack.go", frames[0].Filename)
	}

	if !strings.Contains(frames[0].Method, "roll.getCallers") {
		t.Errorf("expected %#v, got %#v", "roll.getCallers", frames[0].Method)
	}
}

func TestRollbarFramesFingerprint(t *testing.T) {
	tests := []struct {
		Fingerprint string
		Title       string
		Frames      rollbarFrames
	}{
		{
			"9344290d",
			"broken",
			rollbarFrames{
				{"foo.go", "Oops", 1},
			},
		},
		{
			"9344290d",
			"very broken",
			rollbarFrames{
				{"foo.go", "Oops", 1},
			},
		},
		{
			"a4d78b7",
			"broken",
			rollbarFrames{
				{"foo.go", "Oops", 2},
			},
		},
		{
			"50e0fcb3",
			"broken",
			rollbarFrames{
				{"foo.go", "Oops", 1},
				{"foo.go", "Oops", 2},
			},
		},
	}

	for i, test := range tests {
		fp := test.Frames.fingerprint()
		if fp != test.Fingerprint {
			t.Errorf("tests[%d]: got %s", i, fp)
		}
	}
}

func TestScrubFile(t *testing.T) {
	tests := []struct {
		Given    string
		Expected string
	}{
		{"", ""},
		{"foo.go", "foo.go"},
		{"/home/foo/go/src/github.com/stvp/rollbar.go", "github.com/stvp/rollbar.go"},
		{"/home/foo/go/src/gopkg.in/yaml.v1/encode.go", "gopkg.in/yaml.v1/encode.go"},
	}
	for i, test := range tests {
		got := scrubFile(test.Given)
		if got != test.Expected {
			t.Errorf("tests[%d]: got %s", i, got)
		}
	}
}

func TestScrubFunction(t *testing.T) {
	tests := []struct {
		Given    string
		Expected string
	}{
		{"", ""},
		{"roll.getCallers", "roll.getCallers"},
		{"github.com/stvp/roll.getCallers", "roll.getCallers"},
	}
	for i, test := range tests {
		got := scrubFunction(test.Given)
		if got != test.Expected {
			t.Errorf("tests[%d]: got %s", i, got)
		}
	}
}
