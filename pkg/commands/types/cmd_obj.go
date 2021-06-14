package types

import "os/exec"

// A command object is a general way to represent a command to be run on the
// command line. If you want to log the command you'll use .ToString() and
// if you want to run it you'll use .GetCmd()
type ICmdObj interface {
	GetCmd() *exec.Cmd
	ToString() string
	AddEnvVars(...string) ICmdObj
}
