package loaders

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
					Name:         "Myfile",
					ChangeStatus: "MM",
				},
			},
		},
		{
			testName: "two files",
			input:    "MM\x00Myfile\x00M \x00MyOtherFile\x00",
			output: []*models.CommitFile{
				{
					Name:         "Myfile",
					ChangeStatus: "MM",
				},
				{
					Name:         "MyOtherFile",
					ChangeStatus: "M ",
				},
			},
		},
		{
			testName: "three files",
			input:    "MM\x00Myfile\x00M \x00MyOtherFile\x00 M\x00YetAnother\x00",
			output: []*models.CommitFile{
				{
					Name:         "Myfile",
					ChangeStatus: "MM",
				},
				{
					Name:         "MyOtherFile",
					ChangeStatus: "M ",
				},
				{
					Name:         "YetAnother",
					ChangeStatus: " M",
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
