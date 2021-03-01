package oscommands

func getPlatform() *Platform {
	return &Platform{
		OS:           "windows",
		CatCmd:       "cmd /c type",
		Shell:        "cmd",
		ShellArg:     "/c",
		EscapedQuote: `\"`,
	}
}
