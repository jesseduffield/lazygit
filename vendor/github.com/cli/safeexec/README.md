# safeexec

A Go module that provides a safer alternative to `exec.LookPath()` on Windows.

The following, relatively common approach to running external commands has a subtle vulnerability on Windows:
```go
import "os/exec"

func gitStatus() error {
    // On Windows, this will result in `.\git.exe` or `.\git.bat` being executed
    // if either were found in the current working directory.
    cmd := exec.Command("git", "status")
    return cmd.Run()
}
```

Searching the current directory (surprising behavior) before searching folders listed in the PATH environment variable (expected behavior) seems to be intended in Go and unlikely to be changed: https://github.com/golang/go/issues/38736

Since Go does not provide a version of [`exec.LookPath()`](https://golang.org/pkg/os/exec/#LookPath) that only searches PATH and does not search the current working directory, this module provides a `LookPath` function that works consistently across platforms.

Example use:
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

## TODO

Ideally, this module would also provide `exec.Command()` and `exec.CommandContext()` equivalents that delegate to the patched version of `LookPath`. However, this doesn't seem possible since `LookPath` may return an error, while `exec.Command/CommandContext()` themselves do not return an error. In the standard library, the resulting `exec.Cmd` struct stores the LookPath error in a private field, but that functionality isn't available to us.
