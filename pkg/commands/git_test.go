package commands

import (
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"strings"
	"testing"
)

func TestDiff(t *testing.T) {
	dummyLog := logrus.New()
	dummyLog.Out = ioutil.Discard
	osCommand := &OSCommand{
		Log:      dummyLog,
		Platform: getPlatform(),
	}
	gitCommand := &GitCommand{
		Log:       dummyLog,
		OSCommand: osCommand,
	}
	file := File{
		Name:               "asdf",
		DisplayString:      "",
		HasStagedChanges:   false,
		HasUnstagedChanges: true,
		Tracked:            false,
		Deleted:            false,
		HasMergeConflicts:  false,
	}
	osCommand.RunCommand("bash test/repos/ambiguous_ref.sh")
	content := gitCommand.Diff(file)
	if strings.Contains(content, "ambiguous ref") {
		t.Error("Error: ambiguous ref test failed")
	}
}
