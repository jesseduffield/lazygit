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
	RunLineOutputCmd(cmdObj ICmdObj, onLine func(line string) (bool, error)) error
}

type RunExpectation func(ICmdObj) (string, error)

type RealRunner struct {
	log       *logrus.Entry
	logCmdObj func(ICmdObj)
}

func (self *RealRunner) Run(cmdObj ICmdObj) error {
	_, err := self.RunWithOutput(cmdObj)
	return err
}

func (self *RealRunner) RunWithOutput(cmdObj ICmdObj) (string, error) {
	self.logCmdObj(cmdObj)
	output, err := sanitisedCommandOutput(cmdObj.GetCmd().CombinedOutput())
	if err != nil {
		self.log.WithField("command", cmdObj.ToString()).Error(output)
	}
	return output, err
}

func (self *RealRunner) RunLineOutputCmd(cmdObj ICmdObj, onLine func(line string) (bool, error)) error {
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
