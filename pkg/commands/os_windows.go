package commands

func getPlatform() *Platform {
	return &Platform{
		os:                   "windows",
		catCmd:               "type",
		shell:                "cmd",
		shellArg:             "/c",
		escapedQuote:         `\"`,
		fallbackEscapedQuote: "\\'",
	}
}
