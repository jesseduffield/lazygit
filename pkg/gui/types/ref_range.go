package types

import "github.com/jesseduffield/lazygit/pkg/commands/models"

type RefRange struct {
	From models.Ref
	To   models.Ref
}
