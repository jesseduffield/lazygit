package git_commands

type CustomCommands struct {
	*GitCommon
}

func NewCustomCommands(gitCommon *GitCommon) *CustomCommands {
	return &CustomCommands{
		GitCommon: gitCommon,
	}
}

// Only to be used for the sake of running custom commands specified by the user.
// If you want to run a new command, try finding a place for it in one of the neighbouring
// files, or creating a new BlahCommands struct to hold it.
func (self *CustomCommands) RunWithOutput(cmdStr string) (string, error) {
	return self.cmd.New(cmdStr).RunWithOutput()
}
