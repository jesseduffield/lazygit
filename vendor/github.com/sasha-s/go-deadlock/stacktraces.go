package deadlock

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const stackBufSize = 50

var stackBufPool = sync.Pool{
	New: func() interface{} {
		return new([stackBufSize]uintptr)
	},
}

// callers returns a stack trace backed by a pooled buffer. The caller must
// eventually return buf via releaseStackBuf — typically through the
// postLock/postUnlock path which stores it in stackGID.buf.
func callers(skip int) ([]uintptr, *[stackBufSize]uintptr) {
	buf := stackBufPool.Get().(*[stackBufSize]uintptr)
	n := runtime.Callers(2+skip, buf[:])
	return buf[:n], buf
}

// releaseStackBuf returns a pooled stack buffer obtained from callers(). Safe to
// call with nil (e.g. when the buffer was already handed off via stackGID.buf).
func releaseStackBuf(buf *[stackBufSize]uintptr) {
	if buf != nil {
		stackBufPool.Put(buf)
	}
}

// copyStack creates an independent copy of a stack trace. Required when storing
// stacks in long-lived structures (e.g. l.order) because the originals are backed
// by pooled buffers that will be recycled in postUnlock.
func copyStack(s []uintptr) []uintptr {
	c := make([]uintptr, len(s))
	copy(c, s)
	return c
}

func printStack(w io.Writer, stack []uintptr) {
	home := os.Getenv("HOME")
	usr, err := user.Current()
	if err == nil {
		home = usr.HomeDir
	}
	cwd, _ := os.Getwd()

	for i, pc := range stack {
		f := runtime.FuncForPC(pc)
		name := f.Name()
		pkg := ""
		if pos := strings.LastIndex(name, "/"); pos >= 0 {
			name = name[pos+1:]
		}
		if pos := strings.Index(name, "."); pos >= 0 {
			pkg = name[:pos]
			name = name[pos+1:]
		}
		file, line := f.FileLine(pc)
		if (pkg == "runtime" && name == "goexit") || (pkg == "testing" && name == "tRunner") {
			fmt.Fprintln(w)
			return
		}
		tail := ""
		if i == 0 {
			tail = " <<<<<" // Make the line performing a lock prominent.
		}
		// Shorten the file name.
		clean := file
		if cwd != "" {
			cl, err := filepath.Rel(cwd, file)
			if err == nil {
				clean = cl
			}
		}
		if home != "" {
			s2 := strings.Replace(file, home, "~", 1)
			if len(clean) > len(s2) {
				clean = s2
			}
		}
		fmt.Fprintf(w, "%s:%d %s.%s %s%s\n", clean, line-1, pkg, name, code(file, line), tail)
	}
	fmt.Fprintln(w)
}

var fileSources struct {
	sync.Mutex
	lines map[string][][]byte
}

// Reads souce file lines from disk if not cached already.
func getSourceLines(file string) [][]byte {
	fileSources.Lock()
	defer fileSources.Unlock()
	if fileSources.lines == nil {
		fileSources.lines = map[string][][]byte{}
	}
	if lines, ok := fileSources.lines[file]; ok {
		return lines
	}
	text, _ := ioutil.ReadFile(file)
	fileSources.lines[file] = bytes.Split(text, []byte{'\n'})
	return fileSources.lines[file]
}

func code(file string, line int) string {
	lines := getSourceLines(file)
	line -= 2
	if line >= len(lines) || line < 0 {
		return "???"
	}
	return "{ " + string(bytes.TrimSpace(lines[line])) + " }"
}

// Stacktraces for all goroutines.
func stacks() []byte {
	buf := make([]byte, 1024*16)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}
