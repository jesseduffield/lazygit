package cheatsheet

import (
	"testing"

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
					Contexts:    []string{"files"},
					Description: "stage file",
				},
			},
			expected: []*bindingSection{
				{
					title: "Files",
					bindings: []*types.Binding{
						{
							ViewName:    "files",
							Contexts:    []string{"files"},
							Description: "stage file",
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
				},
			},
			expected: []*bindingSection{
				{
					title: "Global Keybindings",
					bindings: []*types.Binding{
						{
							ViewName:    "",
							Description: "quit",
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
					Contexts:    []string{"files"},
				},
				{
					ViewName:    "files",
					Description: "unstage file",
					Contexts:    []string{"files"},
				},
				{
					ViewName:    "files",
					Description: "drop submodule",
					Contexts:    []string{"submodules"},
				},
			},
			expected: []*bindingSection{
				{
					title: "Files",
					bindings: []*types.Binding{
						{
							ViewName:    "files",
							Description: "stage file",
							Contexts:    []string{"files"},
						},
						{
							ViewName:    "files",
							Description: "unstage file",
							Contexts:    []string{"files"},
						},
					},
				},
				{
					title: "Submodules",
					bindings: []*types.Binding{
						{
							ViewName:    "files",
							Description: "drop submodule",
							Contexts:    []string{"submodules"},
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
					Contexts:    []string{"files"},
					Description: "stage file",
				},
				{
					ViewName:    "files",
					Contexts:    []string{"files"},
					Description: "unstage file",
				},
				{
					ViewName:    "files",
					Contexts:    []string{"files"},
					Description: "scroll",
					Tag:         "navigation",
				},
				{
					ViewName:    "commits",
					Contexts:    []string{"commits"},
					Description: "revert commit",
				},
			},
			expected: []*bindingSection{
				{
					title: "List Panel Navigation",
					bindings: []*types.Binding{
						{
							ViewName:    "files",
							Contexts:    []string{"files"},
							Description: "scroll",
							Tag:         "navigation",
						},
					},
				},
				{
					title: "Commits",
					bindings: []*types.Binding{
						{
							ViewName:    "commits",
							Contexts:    []string{"commits"},
							Description: "revert commit",
						},
					},
				},
				{
					title: "Files",
					bindings: []*types.Binding{
						{
							ViewName:    "files",
							Contexts:    []string{"files"},
							Description: "stage file",
						},
						{
							ViewName:    "files",
							Contexts:    []string{"files"},
							Description: "unstage file",
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
					Contexts:    []string{"files"},
					Description: "stage file",
				},
				{
					ViewName:    "files",
					Contexts:    []string{"files"},
					Description: "unstage file",
				},
				{
					ViewName:    "files",
					Contexts:    []string{"files"},
					Description: "scroll",
					Tag:         "navigation",
				},
				{
					ViewName:    "commits",
					Contexts:    []string{"commits"},
					Description: "revert commit",
				},
				{
					ViewName:    "commits",
					Contexts:    []string{"commits"},
					Description: "scroll",
					Tag:         "navigation",
				},
				{
					ViewName:    "commits",
					Contexts:    []string{"commits"},
					Description: "page up",
					Tag:         "navigation",
				},
			},
			expected: []*bindingSection{
				{
					title: "List Panel Navigation",
					bindings: []*types.Binding{
						{
							ViewName:    "files",
							Contexts:    []string{"files"},
							Description: "scroll",
							Tag:         "navigation",
						},
						{
							ViewName:    "commits",
							Contexts:    []string{"commits"},
							Description: "page up",
							Tag:         "navigation",
						},
					},
				},
				{
					title: "Commits",
					bindings: []*types.Binding{
						{
							ViewName:    "commits",
							Contexts:    []string{"commits"},
							Description: "revert commit",
						},
					},
				},
				{
					title: "Files",
					bindings: []*types.Binding{
						{
							ViewName:    "files",
							Contexts:    []string{"files"},
							Description: "stage file",
						},
						{
							ViewName:    "files",
							Contexts:    []string{"files"},
							Description: "unstage file",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			actual := getBindingSections(test.bindings, &tr)
			assert.EqualValues(t, test.expected, actual)
		})
	}
}
