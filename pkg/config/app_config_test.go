package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestCommitPrefixMigrations(t *testing.T) {
	scenarios := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"Empty String",
			"",
			"",
		}, {
			"Single CommitPrefix Rename",
			`
git:
  commitPrefix:
     pattern: "^\\w+-\\w+.*"
     replace: '[JIRA $0] '`,
			`
git:
  commitPrefix:
    - pattern: "^\\w+-\\w+.*"
      replace: '[JIRA $0] '`,
		}, {
			"Complicated CommitPrefixes Rename",
			`
git:
  commitPrefixes:
    foo:
      pattern: "^\\w+-\\w+.*"
      replace: '[OTHER $0] '
    CrazyName!@#$^*&)_-)[[}{f{[]:
      pattern: "^foo.bar*"
      replace: '[FUN $0] '`,
			`
git:
  commitPrefixes:
     foo:
       - pattern: "^\\w+-\\w+.*"
         replace: '[OTHER $0] '
     CrazyName!@#$^*&)_-)[[}{f{[]:
       - pattern: "^foo.bar*"
         replace: '[FUN $0] '`,
		}, {
			"Incomplete Configuration",
			"git:",
			"git:",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			expectedConfig := GetDefaultConfig()
			err := yaml.Unmarshal([]byte(s.expected), expectedConfig)
			if err != nil {
				t.Error(err)
			}
			actual, err := computeMigratedConfig("path doesn't matter", []byte(s.input))
			if err != nil {
				t.Error(err)
			}
			actualConfig := GetDefaultConfig()
			err = yaml.Unmarshal(actual, actualConfig)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, expectedConfig, actualConfig)
		})
	}
}
