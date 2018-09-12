package commands

func getPlatform() *Platform {
	return &Platform{
		os:                   "windows",
		shell:                "cmd",
		shellArg:             "/c",
		escapedQuote:         `\"`,
		fallbackEscapedQuote: "\\'",
	}
}
