package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCachedPullRequestTitleYAMLRoundTrip(t *testing.T) {
	t.Parallel()

	state := AppState{
		GithubPullRequests: map[string][]CachedPullRequest{
			"/tmp/repo": {{
				HeadRefName: "example-branch",
				Number:      9,
				Title:       "\nExample pull request title",
				State:       "MERGED",
			}},
		},
	}

	data, err := yaml.Marshal(&state)
	if err != nil {
		t.Fatal(err)
	}

	var decoded AppState
	if err := yaml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("yaml unmarshal failed: %v\n%s", err, data)
	}
}
