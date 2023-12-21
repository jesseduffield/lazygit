package git_commands

// "*|feat/detect-purge|origin/feat/detect-purge|[ahead 1]"
import (
	"strconv"
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestObtainBranch(t *testing.T) {
	type scenario struct {
		testName                 string
		input                    []string
		storeCommitDateAsRecency bool
		expectedBranch           *models.Branch
	}

	// Use a time stamp of 2 1/2 hours ago, resulting in a recency string of "2h"
	now := time.Now().Unix()
	timeStamp := strconv.Itoa(int(now - 2.5*60*60))

	scenarios := []scenario{
		{
			testName:                 "TrimHeads",
			input:                    []string{"", "heads/a_branch", "", "", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
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
			testName:                 "NoUpstream",
			input:                    []string{"", "a_branch", "", "", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
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
			testName:                 "IsHead",
			input:                    []string{"*", "a_branch", "", "", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
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
			testName:                 "IsBehindAndAhead",
			input:                    []string{"", "a_branch", "a_remote/a_branch", "[behind 2, ahead 3]", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
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
			testName:                 "RemoteBranchIsGone",
			input:                    []string{"", "a_branch", "a_remote/a_branch", "[gone]", "subject", "123", timeStamp},
			storeCommitDateAsRecency: false,
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
		{
			testName:                 "WithCommitDateAsRecency",
			input:                    []string{"", "a_branch", "", "", "subject", "123", timeStamp},
			storeCommitDateAsRecency: true,
			expectedBranch: &models.Branch{
				Name:       "a_branch",
				Recency:    "2h",
				Pushables:  "?",
				Pullables:  "?",
				Head:       false,
				Subject:    "subject",
				CommitHash: "123",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			branch := obtainBranch(s.input, s.storeCommitDateAsRecency)
			assert.EqualValues(t, s.expectedBranch, branch)
		})
	}
}
