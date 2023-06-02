package models

import (
	"testing"
)

func TestCommitStore(t *testing.T) {
	scenarios := []struct {
		testName     string
		commits      []ImmutableCommit
		hash         string
		ancestorHash string
		expected     IsAncestorResponse
	}{
		{
			testName:     "Empty commit store, same hash",
			commits:      []ImmutableCommit{},
			hash:         "a",
			ancestorHash: "a",
			expected:     IsAncestorResponseYes,
		},
		{
			testName:     "Empty commit store, different hash",
			commits:      []ImmutableCommit{},
			hash:         "a",
			ancestorHash: "b",
			expected:     IsAncestorResponseUnknown,
		},
		{
			testName: "Hash found, ancestor not found",
			commits: []ImmutableCommit{
				NewImmutableCommit("a", []string{"c"}),
				NewImmutableCommit("c", []string{}),
			},
			hash:         "a",
			ancestorHash: "b",
			expected:     IsAncestorResponseNo,
		},
		{
			testName: "Hash found, ancestor is parent",
			commits: []ImmutableCommit{
				NewImmutableCommit("a", []string{"b"}),
				NewImmutableCommit("b", []string{"c"}),
			},
			hash:         "a",
			ancestorHash: "b",
			expected:     IsAncestorResponseYes,
		},
		{
			testName: "Hash found, ancestor is grandparent",
			commits: []ImmutableCommit{
				NewImmutableCommit("a", []string{"b"}),
				NewImmutableCommit("b", []string{"c"}),
			},
			hash:         "a",
			ancestorHash: "c",
			expected:     IsAncestorResponseYes,
		},
		{
			testName: "Hash found, not an ancestor",
			commits: []ImmutableCommit{
				NewImmutableCommit("a", []string{"b"}),
				NewImmutableCommit("b", []string{"d"}),
				NewImmutableCommit("c", []string{"d"}),
				NewImmutableCommit("d", []string{}),
			},
			hash:         "a",
			ancestorHash: "c",
			expected:     IsAncestorResponseNo,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.testName, func(t *testing.T) {
			commitStore := NewCommitStore()
			commitStore.AddSlice(scenario.commits)

			response := commitStore.IsAncestor(scenario.hash, scenario.ancestorHash)

			if response != scenario.expected {
				responseStr := IsAncestorResponseStrings[response]
				expectedStr := IsAncestorResponseStrings[scenario.expected]

				t.Errorf("Expected %s, got %s", expectedStr, responseStr)
			}
		})
	}
}
