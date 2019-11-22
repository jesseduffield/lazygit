package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Remote : A git remote
type Remote struct {
	Name     string
	Urls     []string
	Selected bool
	Branches []*RemoteBranch
}

// GetDisplayStrings returns the display string of a remote
func (r *Remote) GetDisplayStrings(isFocused bool) []string {

	branchCount := len(r.Branches)

	return []string{r.Name, utils.ColoredString(fmt.Sprintf("%d branches", branchCount), color.FgBlue)}
}
