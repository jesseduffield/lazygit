package commands

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
)

// Platform stores the os state
type platform struct {
	os           string
	shell        string
	shellArg     string
	escapedQuote string
}

// OSCommand holds all the os commands
type OSCommand struct {
	Log      *logrus.Logger
	Platform platform
}

// NewOSCommand os command runner
func NewOSCommand(log *logrus.Logger) (*OSCommand, error) {
	osCommand := &OSCommand{
		Log:      log,
		Platform: getPlatform(),
	}
	return osCommand, nil
}

// RunCommand wrapper around commands
func (c *OSCommand) RunCommand(command string) (string, error) {
	c.Log.WithField("command", command).Info("RunCommand")
	splitCmd := strings.Split(command, " ")
	cmdOut, err := exec.Command(splitCmd[0], splitCmd[1:]...).CombinedOutput()
	return sanitisedCommandOutput(cmdOut, err)
}

// RunDirectCommand wrapper around direct commands
func (c *OSCommand) RunDirectCommand(command string) (string, error) {
	c.Log.WithField("command", command).Info("RunDirectCommand")

	cmdOut, err := exec.
		Command(c.Platform.shell, c.Platform.shellArg, command).
		CombinedOutput()
	return sanitisedCommandOutput(cmdOut, err)
}

func sanitisedCommandOutput(output []byte, err error) (string, error) {
	outputString := string(output)
	if outputString == "" && err != nil {
		return err.Error(), err
	}
	return outputString, err
}

func getPlatform() platform {
	switch runtime.GOOS {
	case "windows":
		return platform{
			os:           "windows",
			shell:        "cmd",
			shellArg:     "/c",
			escapedQuote: "\\\"",
		}
	default:
		return platform{
			os:           runtime.GOOS,
			shell:        "bash",
			shellArg:     "-c",
			escapedQuote: "\"",
		}
	}
}
