# safeexec

A Go module that provides a stabler alternative to `exec.LookPath()` that:
- Avoids a Windows security risk of executing commands found in the current directory; and
- Allows executing commands found in PATH, even if they come from relative PATH entries.

This is an alternative to [`golang.org/x/sys/execabs`](https://pkg.go.dev/golang.org/x/sys/execabs).

## Usage
```go
import (
    "os/exec"
    "github.com/cli/safeexec"
)

func gitStatus() error {
    gitBin, err := safeexec.LookPath("git")
    if err != nil {
        return err
    }
    cmd := exec.Command(gitBin, "status")
    return cmd.Run()
}
```

## Background
### Windows security vulnerability with Go <= 1.18
Go 1.18 (and older) standard library has a security vulnerability when executing programs:
```go
import "os/exec"

func gitStatus() error {
    // On Windows, this will result in `.\git.exe` or `.\git.bat` being executed
    // if either were found in the current working directory.
    cmd := exec.Command("git", "status")
    return cmd.Run()
}
```

For historic reasons, Go used to implicitly [include the current directory](https://github.com/golang/go/issues/38736) in the PATH resolution on Windows. The `safeexec` package avoids searching the current directory on Windows.

### Relative PATH entries with Go 1.19+

Go 1.19 (and newer) standard library [throws an error](https://github.com/golang/go/issues/43724) if `exec.LookPath("git")` resolved to an executable relative to the current directory. This can happen on other platforms if the PATH environment variable contains relative entries, e.g. `PATH=./bin:$PATH`. The `safeexec` package allows respecting relative PATH entries as it assumes that the responsibility for keeping PATH safe lies outside of the Go program.

## TODO

Ideally, this module would also provide `exec.Command()` and `exec.CommandContext()` equivalents that delegate to the patched version of `LookPath`. However, this doesn't seem possible since `LookPath` may return an error, while `exec.Command/CommandContext()` themselves do not return an error. In the standard library, the resulting `exec.Cmd` struct stores the LookPath error in a private field, but that functionality isn't available to us.
