package oscommands

func getPlatform() *Platform {
	return &Platform{
		OS:       "windows",
		Shell:    "cmd",
		ShellArg: "/c",
	}
}
