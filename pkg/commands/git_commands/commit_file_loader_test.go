package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestGetCommitFilesFromFilenames(t *testing.T) {
	tests := []struct {
		testName string
		input    string
		output   []*models.CommitFile
	}{
		{
			testName: "no files",
			input:    "",
			output:   []*models.CommitFile{},
		},
		{
			testName: "one file",
			input:    "MM\x00Myfile\x00",
			output: []*models.CommitFile{
				{
					Path:         "Myfile",
					ChangeStatus: "MM",
				},
			},
		},
		{
			testName: "two files",
			input:    "MM\x00Myfile\x00M \x00MyOtherFile\x00",
			output: []*models.CommitFile{
				{
					Path:         "Myfile",
					ChangeStatus: "MM",
				},
				{
					Path:         "MyOtherFile",
					ChangeStatus: "M ",
				},
			},
		},
		{
			testName: "three files",
			input:    "MM\x00Myfile\x00M \x00MyOtherFile\x00 M\x00YetAnother\x00",
			output: []*models.CommitFile{
				{
					Path:         "Myfile",
					ChangeStatus: "MM",
				},
				{
					Path:         "MyOtherFile",
					ChangeStatus: "M ",
				},
				{
					Path:         "YetAnother",
					ChangeStatus: " M",
				},
			},
		},
		{
			testName: "a rename among regular files",
			input:    "M\x00Myfile\x00R100\x00before\x00after\x00A\x00Added\x00",
			output: []*models.CommitFile{
				{
					Path:         "Myfile",
					ChangeStatus: "M",
				},
				{
					Path:         "after",
					PreviousPath: "before",
					ChangeStatus: "R",
				},
				{
					Path:         "Added",
					ChangeStatus: "A",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			result := getCommitFilesFromFilenames(test.input)
			assert.Equal(t, test.output, result)
		})
	}
}
