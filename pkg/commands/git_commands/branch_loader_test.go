package git_commands

// "*|feat/detect-purge|origin/feat/detect-purge|[ahead 1]"
import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestObtainBranch(t *testing.T) {
	type scenario struct {
		testName       string
		input          []string
		expectedBranch *models.Branch
	}

	scenarios := []scenario{
		{
			testName: "TrimHeads",
			input:    []string{"", "heads/a_branch", "", "", "subject", "123"},
			expectedBranch: &models.Branch{
				Name:       "a_branch",
				Pushables:  "?",
				Pullables:  "?",
				Head:       false,
				Subject:    "subject",
				CommitHash: "123",
			},
		},
		{
			testName: "NoUpstream",
			input:    []string{"", "a_branch", "", "", "subject", "123"},
			expectedBranch: &models.Branch{
				Name:       "a_branch",
				Pushables:  "?",
				Pullables:  "?",
				Head:       false,
				Subject:    "subject",
				CommitHash: "123",
			},
		},
		{
			testName: "IsHead",
			input:    []string{"*", "a_branch", "", "", "subject", "123"},
			expectedBranch: &models.Branch{
				Name:       "a_branch",
				Pushables:  "?",
				Pullables:  "?",
				Head:       true,
				Subject:    "subject",
				CommitHash: "123",
			},
		},
		{
			testName: "IsBehindAndAhead",
			input:    []string{"", "a_branch", "a_remote/a_branch", "[behind 2, ahead 3]", "subject", "123"},
			expectedBranch: &models.Branch{
				Name:       "a_branch",
				Pushables:  "3",
				Pullables:  "2",
				Head:       false,
				Subject:    "subject",
				CommitHash: "123",
			},
		},
		{
			testName: "RemoteBranchIsGone",
			input:    []string{"", "a_branch", "a_remote/a_branch", "[gone]", "subject", "123"},
			expectedBranch: &models.Branch{
				Name:         "a_branch",
				UpstreamGone: true,
				Pushables:    "?",
				Pullables:    "?",
				Head:         false,
				Subject:      "subject",
				CommitHash:   "123",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			branch := obtainBranch(s.input)
			assert.EqualValues(t, s.expectedBranch, branch)
		})
	}
}
