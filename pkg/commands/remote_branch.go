package commands

import (
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Remote Branch : A git remote branch
type RemoteBranch struct {
	Name       string
	Selected   bool
	RemoteName string
}

// GetDisplayStrings returns the display string of branch
func (b *RemoteBranch) GetDisplayStrings(isFocused bool) []string {
	displayName := utils.ColoredString(b.Name, GetBranchColor(b.Name))

	return []string{displayName}
}
