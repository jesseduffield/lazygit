package types

import "os/exec"

type ICmdObj interface {
	GetCmd() *exec.Cmd
	ToString() string
	AddEnvVars(...string)
}
