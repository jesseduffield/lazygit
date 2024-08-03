package oscommands

func GetPlatform() *Platform {
	return &Platform{
		OS:                  "windows",
		Shell:               "cmd",
		InteractiveShell:    "cmd",
		ShellArg:            "/c",
		InteractiveShellArg: "",
	}
}
