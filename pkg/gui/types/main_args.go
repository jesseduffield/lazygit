package types

// StartArgs is the struct that represents some things we want to do on program start
type StartArgs struct {
	// FilterPath determines which path we're going to filter on so that we only see commits from that file.
	FilterPath string
	// GitArg determines what context we open in
	GitArg GitArg
}

type GitArg string

const (
	GitArgNone   GitArg = ""
	GitArgStatus GitArg = "status"
	GitArgBranch GitArg = "branch"
	GitArgLog    GitArg = "log"
	GitArgStash  GitArg = "stash"
)

func NewStartArgs(filterPath string, gitArg GitArg) StartArgs {
	return StartArgs{
		FilterPath: filterPath,
		GitArg:     gitArg,
	}
}
