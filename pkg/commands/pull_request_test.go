package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetRepoInfoFromURL is a function.
func TestGetRepoInfoFromURL(t *testing.T) {
	type scenario struct {
		testName string
		repoURL  string
		test     func(*RepoInformation)
	}

	scenarios := []scenario{
		{
			"Returns repository information for git remote url",
			"git@github.com:petersmith/super_calculator",
			func(repoInfo *RepoInformation) {
				assert.EqualValues(t, repoInfo.Owner, "petersmith")
				assert.EqualValues(t, repoInfo.Repository, "super_calculator")
			},
		},
		{
			"Returns repository information for http remote url",
			"https://my_username@bitbucket.org/johndoe/social_network.git",
			func(repoInfo *RepoInformation) {
				assert.EqualValues(t, repoInfo.Owner, "johndoe")
				assert.EqualValues(t, repoInfo.Repository, "social_network")
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			s.test(GetRepoInfoFromURL(s.repoURL))
		})
	}
}
