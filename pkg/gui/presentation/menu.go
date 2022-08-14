package presentation

import "github.com/jesseduffield/lazygit/pkg/gui/style"

func OpensMenuStyle(str string) string {
	return style.FgMagenta.Sprintf("%s...", str)
}
