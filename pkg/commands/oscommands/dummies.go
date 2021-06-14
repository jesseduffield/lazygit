package oscommands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// NewDummyOS creates a new dummy OSCommand for testing
func NewDummyOS() *OS {
	return NewOS(utils.NewDummyLog(), config.NewDummyAppConfig())
}
