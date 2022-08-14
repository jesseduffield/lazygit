package loaders

// "*|feat/detect-purge|origin/feat/detect-purge|[ahead 1]"
import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestObtainBanch(t *testing.T) {
	type scenario struct {
		testName       string
		input          []string
		expectedBranch *models.Branch
	}

	scenarios := []scenario{
		{
			testName:       "TrimHeads",
			input:          []string{"", "heads/a_branch", "", ""},
			expectedBranch: &models.Branch{Name: "a_branch", Pushables: "?", Pullables: "?", Head: false},
		},
		{
			testName:       "NoUpstream",
			input:          []string{"", "a_branch", "", ""},
			expectedBranch: &models.Branch{Name: "a_branch", Pushables: "?", Pullables: "?", Head: false},
		},
		{
			testName:       "IsHead",
			input:          []string{"*", "a_branch", "", ""},
			expectedBranch: &models.Branch{Name: "a_branch", Pushables: "?", Pullables: "?", Head: true},
		},
		{
			testName:       "IsBehindAndAhead",
			input:          []string{"", "a_branch", "a_remote/a_branch", "[behind 2, ahead 3]"},
			expectedBranch: &models.Branch{Name: "a_branch", Pushables: "3", Pullables: "2", Head: false},
		},
		{
			testName:       "RemoteBranchIsGone",
			input:          []string{"", "a_branch", "a_remote/a_branch", "[gone]"},
			expectedBranch: &models.Branch{Name: "a_branch", UpstreamGone: true, Pushables: "?", Pullables: "?", Head: false},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			branch := obtainBranch(s.input)
			assert.EqualValues(t, s.expectedBranch, branch)
		})
	}
}
