# About terminfo [![GoDoc][1]][2]

Package `terminfo` provides a pure-Go implementation of reading information
from the terminfo database.

`terminfo` is meant as a replacement for `ncurses` in simple Go programs.

## Installing

Install in the usual Go way:

```sh
$ go get -u github.com/xo/terminfo
```

## Using

Please see the [GoDoc API listing][2] for more information on using `terminfo`.

```go
// _examples/simple/main.go
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/xo/terminfo"
)

func main() {
	//r := rand.New(nil)

	// load terminfo
	ti, err := terminfo.LoadFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// cleanup
	defer func() {
		err := recover()
		termreset(ti)
		if err != nil {
			log.Fatal(err)
		}
	}()

	terminit(ti)
	termtitle(ti, "simple example!")
	termputs(ti, 3, 3, "Ctrl-C to exit")
	maxColors := termcolors(ti)
	if maxColors > 256 {
		maxColors = 256
	}
	for i := 0; i < maxColors; i++ {
		termputs(ti, 5+i/16, 5+i%16, ti.Colorf(i, 0, "█"))
	}

	// wait for signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

// terminit initializes the special CA mode on the terminal, and makes the
// cursor invisible.
func terminit(ti *terminfo.Terminfo) {
	buf := new(bytes.Buffer)
	// set the cursor invisible
	ti.Fprintf(buf, terminfo.CursorInvisible)
	// enter special mode
	ti.Fprintf(buf, terminfo.EnterCaMode)
	// clear the screen
	ti.Fprintf(buf, terminfo.ClearScreen)
	os.Stdout.Write(buf.Bytes())
}

// termreset is the inverse of terminit.
func termreset(ti *terminfo.Terminfo) {
	buf := new(bytes.Buffer)
	ti.Fprintf(buf, terminfo.ExitCaMode)
	ti.Fprintf(buf, terminfo.CursorNormal)
	os.Stdout.Write(buf.Bytes())
}

// termputs puts a string at row, col, interpolating v.
func termputs(ti *terminfo.Terminfo, row, col int, s string, v ...interface{}) {
	buf := new(bytes.Buffer)
	ti.Fprintf(buf, terminfo.CursorAddress, row, col)
	fmt.Fprintf(buf, s, v...)
	os.Stdout.Write(buf.Bytes())
}

// sl is the status line terminfo.
var sl *terminfo.Terminfo

// termtitle sets the window title.
func termtitle(ti *terminfo.Terminfo, s string) {
	var once sync.Once
	once.Do(func() {
		if ti.Has(terminfo.HasStatusLine) {
			return
		}
		// load the sl xterm if terminal is an xterm or has COLORTERM
		if strings.Contains(strings.ToLower(os.Getenv("TERM")), "xterm") || os.Getenv("COLORTERM") == "truecolor" {
			sl, _ = terminfo.Load("xterm+sl")
		}
	})
	if sl != nil {
		ti = sl
	}
	if !ti.Has(terminfo.HasStatusLine) {
		return
	}
	buf := new(bytes.Buffer)
	ti.Fprintf(buf, terminfo.ToStatusLine)
	fmt.Fprint(buf, s)
	ti.Fprintf(buf, terminfo.FromStatusLine)
	os.Stdout.Write(buf.Bytes())
}

// termcolors returns the maximum colors available for the terminal.
func termcolors(ti *terminfo.Terminfo) int {
	if colors := ti.Num(terminfo.MaxColors); colors > 0 {
		return colors
	}
	return int(terminfo.ColorLevelBasic)
}
```

[1]: https://godoc.org/github.com/xo/terminfo?status.svg
[2]: https://godoc.org/github.com/xo/terminfo
