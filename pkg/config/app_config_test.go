package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitPrefixMigrations(t *testing.T) {
	scenarios := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty String",
			input:    "",
			expected: "",
		}, {
			name: "Single CommitPrefix Rename",
			input: `git:
  commitPrefix:
     pattern: "^\\w+-\\w+.*"
     replace: '[JIRA $0] '
`,
			expected: `git:
  commitPrefix:
    - pattern: "^\\w+-\\w+.*"
      replace: '[JIRA $0] '
`,
		}, {
			name: "Complicated CommitPrefixes Rename",
			input: `git:
  commitPrefixes:
    foo:
      pattern: "^\\w+-\\w+.*"
      replace: '[OTHER $0] '
    CrazyName!@#$^*&)_-)[[}{f{[]:
      pattern: "^foo.bar*"
      replace: '[FUN $0] '
`,
			expected: `git:
  commitPrefixes:
    foo:
      - pattern: "^\\w+-\\w+.*"
        replace: '[OTHER $0] '
    CrazyName!@#$^*&)_-)[[}{f{[]:
      - pattern: "^foo.bar*"
        replace: '[FUN $0] '
`,
		}, {
			name:     "Incomplete Configuration",
			input:    "git:",
			expected: "git:",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			actual, err := computeMigratedConfig("path doesn't matter", []byte(s.input))
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, s.expected, string(actual))
		})
	}
}
