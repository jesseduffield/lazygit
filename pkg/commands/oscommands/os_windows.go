package oscommands

func getPlatform() *Platform {
	return &Platform{
		OS:                   "windows",
		CatCmd:               "type",
		Shell:                "cmd",
		ShellArg:             "/c",
		EscapedQuote:         `\"`,
		FallbackEscapedQuote: "\\'",
	}
}
