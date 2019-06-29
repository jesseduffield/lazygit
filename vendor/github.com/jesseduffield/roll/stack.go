package roll

import (
	"fmt"
	"hash/crc32"
	"os"
	"runtime"
	"strings"
)

var (
	knownFilePathPatterns = []string{
		"github.com/",
		"code.google.com/",
		"bitbucket.org/",
		"launchpad.net/",
		"gopkg.in/",
	}
)

func getCallers(skip int) (pc []uintptr) {
	pc = make([]uintptr, 1000)
	i := runtime.Callers(skip+1, pc)
	return pc[0:i]
}

// -- rollbarFrames

type rollbarFrame struct {
	Filename string `json:"filename"`
	Method   string `json:"method"`
	Line     int    `json:"lineno"`
}

type rollbarFrames []rollbarFrame

// buildRollbarFrames takes a slice of function pointers and returns a Rollbar
// API payload containing the filename, method name, and line number of each
// function.
func buildRollbarFrames(callers []uintptr) (frames rollbarFrames) {
	frames = rollbarFrames{}

	// 2016-08-24 - runtime.CallersFrames was added in Go 1.7, which should
	// replace the following code when roll is able to require Go 1.7+.
	for _, caller := range callers {
		frame := rollbarFrame{
			Filename: "???",
			Method:   "???",
		}
		if fn := runtime.FuncForPC(caller); fn != nil {
			name, line := fn.FileLine(caller)
			frame.Filename = scrubFile(name)
			frame.Line = line
			frame.Method = scrubFunction(fn.Name())
		}
		frames = append(frames, frame)
	}

	return frames
}

// fingerprint returns a checksum that uniquely identifies a stacktrace by the
// filename, method name, and line number of every frame in the stack.
func (f rollbarFrames) fingerprint() string {
	hash := crc32.NewIEEE()
	for _, frame := range f {
		fmt.Fprintf(hash, "%s%s%d", frame.Filename, frame.Method, frame.Line)
	}
	return fmt.Sprintf("%x", hash.Sum32())
}

// -- Helpers

// scrubFile removes unneeded information from the path of a source file. This
// makes them shorter in Rollbar UI as well as making them the same, regardless
// of the machine the code was compiled on.
//
// Example:
//   /home/foo/go/src/github.com/stvp/roll/rollbar.go -> github.com/stvp/roll/rollbar.go
func scrubFile(s string) string {
	var i int
	for _, pattern := range knownFilePathPatterns {
		i = strings.Index(s, pattern)
		if i != -1 {
			return s[i:]
		}
	}
	return s
}

// scrubFunction removes unneeded information from the full name of a function.
//
// Example:
//   github.com/stvp/roll.getCallers -> roll.getCallers
func scrubFunction(name string) string {
	end := strings.LastIndex(name, string(os.PathSeparator))
	return name[end+1 : len(name)]
}
