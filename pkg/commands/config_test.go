package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGitCommandUsingGpg is a function.
func TestGitCommandUsingGpg(t *testing.T) {
	type scenario struct {
		testName          string
		getGitConfigValue func(string) (string, error)
		test              func(bool)
	}

	scenarios := []scenario{
		{
			"Option global and local config commit.gpgsign is not set",
			func(string) (string, error) { return "", nil },
			func(gpgEnabled bool) {
				assert.False(t, gpgEnabled)
			},
		},
		{
			"Option commit.gpgsign is true",
			func(string) (string, error) {
				return "True", nil
			},
			func(gpgEnabled bool) {
				assert.True(t, gpgEnabled)
			},
		},
		{
			"Option commit.gpgsign is on",
			func(string) (string, error) {
				return "ON", nil
			},
			func(gpgEnabled bool) {
				assert.True(t, gpgEnabled)
			},
		},
		{
			"Option commit.gpgsign is yes",
			func(string) (string, error) {
				return "YeS", nil
			},
			func(gpgEnabled bool) {
				assert.True(t, gpgEnabled)
			},
		},
		{
			"Option commit.gpgsign is 1",
			func(string) (string, error) {
				return "1", nil
			},
			func(gpgEnabled bool) {
				assert.True(t, gpgEnabled)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGit()
			gitCmd.getGitConfigValue = s.getGitConfigValue
			s.test(gitCmd.UsingGpg())
		})
	}
}
