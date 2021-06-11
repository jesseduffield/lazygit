type ICmdObj interface {
	ToCmd() *exec.Cmd
	ToString() string
	AddEnvVars(...string)
}
