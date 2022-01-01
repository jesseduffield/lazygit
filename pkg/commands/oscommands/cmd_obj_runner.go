package oscommands

import (
	"bufio"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

type ICmdObjRunner interface {
	Run(cmdObj ICmdObj) error
	RunWithOutput(cmdObj ICmdObj) (string, error)
	RunAndProcessLines(cmdObj ICmdObj, onLine func(line string) (bool, error)) error
}

type cmdObjRunner struct {
	log   *logrus.Entry
	guiIO *guiIO
}

var _ ICmdObjRunner = &cmdObjRunner{}

func (self *cmdObjRunner) runWithCredentialHandling(cmdObj ICmdObj) error {
	switch cmdObj.GetCredentialStrategy() {
	case PROMPT:
		return self.RunCommandWithOutputLive(cmdObj, self.guiIO.promptForCredentialFn)
	case FAIL:
		return self.RunCommandWithOutputLive(cmdObj, func(s string) string { return "\n" })
	}

	// we should never land here
	return errors.New("runWithCredentialHandling called but cmdObj does not have a a credential strategy")
}

func (self *cmdObjRunner) Run(cmdObj ICmdObj) error {
	if cmdObj.GetCredentialStrategy() == NONE {
		_, err := self.RunWithOutput(cmdObj)
		return err
	} else {
		return self.runWithCredentialHandling(cmdObj)
	}
}

func (self *cmdObjRunner) logCmdObj(cmdObj ICmdObj) {
	self.guiIO.logCommandFn(cmdObj.ToString(), true)
}

func (self *cmdObjRunner) RunWithOutput(cmdObj ICmdObj) (string, error) {
	if cmdObj.GetCredentialStrategy() != NONE {
		err := self.runWithCredentialHandling(cmdObj)
		// for now we're not capturing output, just because it would take a little more
		// effort and there's currently no use case for it. Some commands call RunWithOutput
		// but ignore the output, hence why we've got this check here.
		return "", err
	}

	if cmdObj.ShouldLog() {
		self.logCmdObj(cmdObj)
	}

	output, err := sanitisedCommandOutput(cmdObj.GetCmd().CombinedOutput())
	if err != nil {
		self.log.WithField("command", cmdObj.ToString()).Error(output)
	}
	return output, err
}

func (self *cmdObjRunner) RunAndProcessLines(cmdObj ICmdObj, onLine func(line string) (bool, error)) error {
	if cmdObj.GetCredentialStrategy() != NONE {
		return errors.New("cannot call RunAndProcessLines with credential strategy. If you're seeing this then a contributor to Lazygit has accidentally called this method! Please raise an issue")
	}

	if cmdObj.ShouldLog() {
		self.logCmdObj(cmdObj)
	}

	cmd := cmdObj.GetCmd()
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdoutPipe)
	scanner.Split(bufio.ScanLines)
	if err := cmd.Start(); err != nil {
		return err
	}

	for scanner.Scan() {
		line := scanner.Text()
		stop, err := onLine(line)
		if err != nil {
			return err
		}
		if stop {
			_ = cmd.Process.Kill()
			break
		}
	}

	_ = cmd.Wait()

	return nil
}

func sanitisedCommandOutput(output []byte, err error) (string, error) {
	outputString := string(output)
	if err != nil {
		// errors like 'exit status 1' are not very useful so we'll create an error
		// from the combined output
		if outputString == "" {
			return "", utils.WrapError(err)
		}
		return outputString, errors.New(outputString)
	}
	return outputString, nil
}
