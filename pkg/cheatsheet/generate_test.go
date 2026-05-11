package cheatsheet

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

func TestGetBindingSections(t *testing.T) {
	tr := i18n.EnglishTranslationSet()

	tests := []struct {
		testName string
		bindings []*types.Binding
		expected []*bindingSection
	}{
		{
			testName: "no bindings",
			bindings: []*types.Binding{},
			expected: []*bindingSection{},
		},
		{
			testName: "one binding",
			bindings: []*types.Binding{
				{
					ViewName:    "files",
					Description: "stage file",
					Key:         gocui.NewKeyRune('a'),
				},
			},
			expected: []*bindingSection{
				{
					title: "Files",
					bindings: []*types.Binding{
						{
							ViewName:    "files",
							Description: "stage file",
							Key:         gocui.NewKeyRune('a'),
						},
					},
				},
			},
		},
		{
			testName: "global binding",
			bindings: []*types.Binding{
				{
					ViewName:    "",
					Description: "quit",
					Key:         gocui.NewKeyRune('a'),
				},
			},
			expected: []*bindingSection{
				{
					title: "Global keybindings",
					bindings: []*types.Binding{
						{
							ViewName:    "",
							Description: "quit",
							Key:         gocui.NewKeyRune('a'),
						},
					},
				},
			},
		},
		{
			testName: "grouped bindings",
			bindings: []*types.Binding{
				{
					ViewName:    "files",
					Description: "stage file",
					Key:         gocui.NewKeyRune('a'),
				},
				{
					ViewName:    "files",
					Description: "unstage file",
					Key:         gocui.NewKeyRune('a'),
				},
				{
					ViewName:    "submodules",
					Description: "drop submodule",
					Key:         gocui.NewKeyRune('a'),
				},
			},
			expected: []*bindingSection{
				{
					title: "Files",
					bindings: []*types.Binding{
						{
							ViewName:    "files",
							Description: "stage file",
							Key:         gocui.NewKeyRune('a'),
						},
						{
							ViewName:    "files",
							Description: "unstage file",
							Key:         gocui.NewKeyRune('a'),
						},
					},
				},
				{
					title: "Submodules",
					bindings: []*types.Binding{
						{
							ViewName:    "submodules",
							Description: "drop submodule",
							Key:         gocui.NewKeyRune('a'),
						},
					},
				},
			},
		},
		{
			testName: "with navigation bindings",
			bindings: []*types.Binding{
				{
					ViewName:    "files",
					Description: "stage file",
					Key:         gocui.NewKeyRune('a'),
				},
				{
					ViewName:    "files",
					Description: "unstage file",
					Key:         gocui.NewKeyRune('a'),
				},
				{
					ViewName:    "files",
					Description: "scroll",
					Key:         gocui.NewKeyRune('a'),
					Tag:         "navigation",
				},
				{
					ViewName:    "commits",
					Description: "revert commit",
					Key:         gocui.NewKeyRune('a'),
				},
			},
			expected: []*bindingSection{
				{
					title: "List panel navigation",
					bindings: []*types.Binding{
						{
							ViewName:    "files",
							Description: "scroll",
							Key:         gocui.NewKeyRune('a'),
							Tag:         "navigation",
						},
					},
				},
				{
					title: "Commits",
					bindings: []*types.Binding{
						{
							ViewName:    "commits",
							Description: "revert commit",
							Key:         gocui.NewKeyRune('a'),
						},
					},
				},
				{
					title: "Files",
					bindings: []*types.Binding{
						{
							ViewName:    "files",
							Description: "stage file",
							Key:         gocui.NewKeyRune('a'),
						},
						{
							ViewName:    "files",
							Description: "unstage file",
							Key:         gocui.NewKeyRune('a'),
						},
					},
				},
			},
		},
		{
			testName: "with duplicate navigation bindings",
			bindings: []*types.Binding{
				{
					ViewName:    "files",
					Description: "stage file",
					Key:         gocui.NewKeyRune('a'),
				},
				{
					ViewName:    "files",
					Description: "unstage file",
					Key:         gocui.NewKeyRune('a'),
				},
				{
					ViewName:    "files",
					Description: "scroll",
					Key:         gocui.NewKeyRune('a'),
					Tag:         "navigation",
				},
				{
					ViewName:    "commits",
					Description: "revert commit",
					Key:         gocui.NewKeyRune('a'),
				},
				{
					ViewName:    "commits",
					Description: "scroll",
					Key:         gocui.NewKeyRune('a'),
					Tag:         "navigation",
				},
				{
					ViewName:    "commits",
					Description: "page up",
					Key:         gocui.NewKeyRune('a'),
					Tag:         "navigation",
				},
			},
			expected: []*bindingSection{
				{
					title: "List panel navigation",
					bindings: []*types.Binding{
						{
							ViewName:    "files",
							Description: "scroll",
							Key:         gocui.NewKeyRune('a'),
							Tag:         "navigation",
						},
						{
							ViewName:    "commits",
							Description: "page up",
							Key:         gocui.NewKeyRune('a'),
							Tag:         "navigation",
						},
					},
				},
				{
					title: "Commits",
					bindings: []*types.Binding{
						{
							ViewName:    "commits",
							Description: "revert commit",
							Key:         gocui.NewKeyRune('a'),
						},
					},
				},
				{
					title: "Files",
					bindings: []*types.Binding{
						{
							ViewName:    "files",
							Description: "stage file",
							Key:         gocui.NewKeyRune('a'),
						},
						{
							ViewName:    "files",
							Description: "unstage file",
							Key:         gocui.NewKeyRune('a'),
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			actual := getBindingSections(test.bindings, tr)
			assert.EqualValues(t, test.expected, actual)
		})
	}
}
