package oscommands

import (
	//"github.com/jesseduffield/lazygit/pkg/common"
	//"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"testing"
)

func TestProcessOutput(t *testing.T) {
	dummyLog := utils.NewDummyLog()
	scenarios := []struct {
		name   string
		runner ICmdObjRunner
	}{
		{
			name:   "hi",
			runner: cmdObjRunner{dummyLog, NewNullGuiIO(utils.NewDummyLog())},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
		})
	}
}
