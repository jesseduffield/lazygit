package helpers

import (
	"testing"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

func TestCommitSignatureSubTitle(t *testing.T) {
	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelNone)
	defer color.ForceSetColorLevel(oldColorLevel)

	icons.SetNerdFontsVersion("")
	defer icons.SetNerdFontsVersion("")

	helper := &DiffHelper{
		c: &HelperCommon{Common: common.NewDummyCommon()},
	}

	assert.Equal(t, helper.c.Tr.CommitSignatureVerifiedSubTitle, helper.CommitSignatureSubTitle(&models.Commit{
		SignatureStatus: models.GitSignatureStatusGood,
	}))
	assert.Equal(t, helper.c.Tr.CommitSignatureVerifiedSubTitle, helper.CommitSignatureSubTitle(&models.Commit{
		SignatureStatus: models.GitSignatureStatusBad,
	}))
	assert.Equal(t, "", helper.CommitSignatureSubTitle(&models.Commit{
		SignatureStatus: models.GitSignatureStatusNone,
	}))
}

func TestCommitSignatureSubTitleWithIcons(t *testing.T) {
	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelNone)
	defer color.ForceSetColorLevel(oldColorLevel)

	icons.SetNerdFontsVersion("3")
	defer icons.SetNerdFontsVersion("")

	helper := &DiffHelper{
		c: &HelperCommon{Common: common.NewDummyCommon()},
	}

	assert.Equal(t, icons.SIGNED_COMMIT_ICON+" "+helper.c.Tr.CommitSignatureVerifiedSubTitle, helper.CommitSignatureSubTitle(&models.Commit{
		SignatureStatus: models.GitSignatureStatusGood,
	}))
	assert.Equal(t, icons.SIGNED_COMMIT_ICON+" "+helper.c.Tr.CommitSignatureVerifiedSubTitle, helper.CommitSignatureSubTitle(&models.Commit{
		SignatureStatus: models.GitSignatureStatusCannotCheck,
	}))
}

func TestCombineSubTitles(t *testing.T) {
	helper := &DiffHelper{}

	assert.Equal(t, "Verified | Ignoring whitespace", helper.CombineSubTitles("Verified", "", "Ignoring whitespace"))
	assert.Equal(t, "", helper.CombineSubTitles("", ""))
}
