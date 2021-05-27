package oscommands

func getPlatform() *Platform {
	return &Platform{
		OS:           "windows",
		CatCmd:       []string{"cmd", "/c", "type"},
		Shell:        "cmd",
		ShellArg:     "/c",
		EscapedQuote: `\"`,
	}
}
