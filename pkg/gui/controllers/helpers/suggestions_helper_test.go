package helpers

import (
	"testing"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

// fakeGuiCommon gives a helper under test the model. Calling any other
// IGuiCommon method panics, because the embedded interface is nil.
type fakeGuiCommon struct {
	types.IGuiCommon
	model *types.Model
}

func (self *fakeGuiCommon) Model() *types.Model {
	return self.model
}

func TestGetBranchNameSuggestionsFunc(t *testing.T) {
	previousColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelMillions)
	t.Cleanup(func() { color.ForceSetColorLevel(previousColorLevel) })
	previousDefaultTextColor := theme.DefaultTextColor
	t.Cleanup(func() { theme.DefaultTextColor = previousDefaultTextColor })
	// The branch colors can't be read back from outside their package, so
	// reset them to what the default config sets.
	t.Cleanup(func() { presentation.SetCustomBranches(nil) })

	theme.DefaultTextColor = style.FgYellow
	presentation.SetCustomBranches(config.ColorPatterns{{Pattern: "^feature/", Color: "green"}})
	branches := lo.Map([]string{"main", "feature/foo", "feature/bar"}, func(name string, _ int) *models.Branch {
		return &models.Branch{Name: name}
	})
	suggestionsHelper := NewSuggestionsHelper(&HelperCommon{
		Common:     common.NewDummyCommon(),
		IGuiCommon: &fakeGuiCommon{model: &types.Model{Branches: branches}},
	})
	findSuggestions := suggestionsHelper.GetBranchNameSuggestionsFunc()

	// The func runs on the suggestions worker, while the UI thread may set
	// new colors. Its labels keep the colors from when it was created.
	presentation.SetCustomBranches(config.ColorPatterns{{Pattern: "^feature/", Color: "red"}})
	theme.DefaultTextColor = style.FgBlue

	assert.Equal(t, []*types.Suggestion{
		{Value: "main", Label: "\x1b[33mmain\x1b[0m"},
		{Value: "feature/foo", Label: "\x1b[32mfeature/foo\x1b[0m"},
		{Value: "feature/bar", Label: "\x1b[32mfeature/bar\x1b[0m"},
	}, findSuggestions(""))
	assert.Equal(t, []*types.Suggestion{
		{Value: "feature/foo", Label: "\x1b[32mfeature/foo\x1b[0m"},
	}, findSuggestions("foo"))
}
