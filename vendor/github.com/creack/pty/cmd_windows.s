// disable default -complete option to `go tool compile` on windows

// We are using `go:linkname` to support golang os/exec.Cmd for extended 
// windows process creation (startupInfoEx), and that is not supported by its
// standard library implementation.

// By default, the go compiler will require all functions in *.go files 
// with body defined, if that's not the case, we have to enable CGO to 
// enable symbol lookup. One solution to disable that compile time check
// is to add some go assembly file to your project.

// For this project, we don't use CGO at all, and should not require users
// to set `CGO_ENABLED=1` when compiling their projects using this package.

// By adding this empty assembly file, the go compiler will enable symbol 
// lookup, so that we can have functions with no body defined in *.go files.
