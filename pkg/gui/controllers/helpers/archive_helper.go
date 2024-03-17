package helpers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type ArchiveHelper struct {
	c *HelperCommon
}

func NewArchiveHelper(
	c *HelperCommon,
) *ArchiveHelper {
	return &ArchiveHelper{
		c: c,
	}
}

func (self *ArchiveHelper) CreateArchive(refName string) error {
	placeholders := map[string]string{"ref": refName}

	return self.c.Prompt(types.PromptOpts{
		Title: utils.ResolvePlaceholderString(self.c.Tr.ArchiveChooseFileName, placeholders),
		HandleConfirm: func(fileName string) error {
			validArchiveFormats, err := self.c.Git().Archive.GetValidArchiveFormats()
			if err != nil {
				return err
			}

			menuItems := make([]*types.MenuItem, len(validArchiveFormats))

			for i, format := range validArchiveFormats {
				format := format

				menuItems[i] = &types.MenuItem{
					Label: format,
					OnPress: func() error {
						return self.runArchiveCommand(refName, fileName, fileName+"/", format)
					},
				}
			}

			return self.c.Menu(types.CreateMenuOptions{
				Title: self.c.Tr.ArchiveChooseFormatMenuTitle,
				Items: menuItems,
			})
		},
	})
}

func (self *ArchiveHelper) runArchiveCommand(refName string, fileName string, prefix string, suffix string) error {
	return self.c.WithWaitingStatus(self.c.Tr.ArchiveWaitingStatusMessage, func(gocui.Task) error {
		return self.c.Git().Archive.Archive(refName, fileName+suffix, prefix)
	})
}
