# pty

Pty is a Go package for using unix pseudo-terminals.

## Install

    go get github.com/kr/pty

## Example

### Command

```go
package main

import (
	"github.com/kr/pty"
	"io"
	"os"
	"os/exec"
)

func main() {
	c := exec.Command("grep", "--color=auto", "bar")
	f, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

	go func() {
		f.Write([]byte("foo\n"))
		f.Write([]byte("bar\n"))
		f.Write([]byte("baz\n"))
		f.Write([]byte{4}) // EOT
	}()
	io.Copy(os.Stdout, f)
}
```

### Shell

```go
package main

import (
        "io"
        "log"
        "os"
        "os/exec"
        "os/signal"
        "syscall"

        "github.com/kr/pty"
        "golang.org/x/crypto/ssh/terminal"
)

func test() error {
        // Create arbitrary command.
        c := exec.Command("bash")

        // Start the command with a pty.
        ptmx, err := pty.Start(c)
        if err != nil {
                return err
        }
        // Make sure to close the pty at the end.
        defer func() { _ = ptmx.Close() }() // Best effort.

        // Handle pty size.
        ch := make(chan os.Signal, 1)
        signal.Notify(ch, syscall.SIGWINCH)
        go func() {
                for range ch {
                        if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
                                log.Printf("error resizing pty: %s", err)
                        }
                }
        }()
        ch <- syscall.SIGWINCH // Initial resize.

        // Set stdin in raw mode.
        oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
        if err != nil {
                panic(err)
        }
        defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

        // Copy stdin to the pty and the pty to stdout.
        go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
        _, _ = io.Copy(os.Stdout, ptmx)

        return nil
}

func main() {
        if err := test(); err != nil {
                log.Fatal(err)
        }
}
```
