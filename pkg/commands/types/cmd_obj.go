package types

import "os/exec"

type ICmdObj interface {
	ToCmd() *exec.Cmd
	ToString() string
	AddEnvVars(...string)
}
