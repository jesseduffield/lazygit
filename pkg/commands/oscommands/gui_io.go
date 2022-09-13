package oscommands

import (
	"io"

	"github.com/sirupsen/logrus"
)

// this struct captures some IO stuff
type guiIO struct {
	// this is for logging anything we want. It'll be written to a log file for the sake
	// of debugging.
	log *logrus.Entry

	// this is for us to log the command we're about to run e.g. 'git push'. The GUI
	// will write this to a log panel so that the user can see which commands are being
	// run.
	// The isCommandLineCommand arg is there so that we can style the log differently
	// depending on whether we're directly outputting a command we're about to run that
	// will be run on the command line, or if we're using something from Go's standard lib.
	logCommandFn func(str string, isCommandLineCommand bool)
	// this is for us to directly write the output of a command. We will do this for
	// certain commands like 'git push'. The GUI will write this to a command output panel.
	// We need a new cmd writer per command, hence it being a function.
	newCmdWriterFn func() io.Writer
	// this allows us to request info from the user like username/password, in the event
	// that a command requests it.
	// the 'credential' arg is something like 'username' or 'password'
	promptForCredentialFn func(credential CredentialType) string
}

func NewGuiIO(log *logrus.Entry, logCommandFn func(string, bool), newCmdWriterFn func() io.Writer, promptForCredentialFn func(CredentialType) string) *guiIO {
	return &guiIO{
		log:                   log,
		logCommandFn:          logCommandFn,
		newCmdWriterFn:        newCmdWriterFn,
		promptForCredentialFn: promptForCredentialFn,
	}
}

// we use this function when we want to access the functionality of our OS struct but we
// don't have anywhere to log things, or request input from the user.
func NewNullGuiIO(log *logrus.Entry) *guiIO {
	return &guiIO{
		log:                   log,
		logCommandFn:          func(string, bool) {},
		newCmdWriterFn:        func() io.Writer { return io.Discard },
		promptForCredentialFn: failPromptFn,
	}
}
