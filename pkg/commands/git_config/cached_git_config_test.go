package git_config

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetBool(t *testing.T) {
	type scenario struct {
		testName      string
		mockResponses map[string]string
		expected      bool
	}

	scenarios := []scenario{
		{
			"Option global and local config commit.gpgsign is not set",
			map[string]string{},
			false,
		},
		{
			"Some other random key is set",
			map[string]string{"blah": "blah"},
			false,
		},
		{
			"Option commit.gpgsign is true",
			map[string]string{"commit.gpgsign": "True"},
			true,
		},
		{
			"Option commit.gpgsign is on",
			map[string]string{"commit.gpgsign": "ON"},
			true,
		},
		{
			"Option commit.gpgsign is yes",
			map[string]string{"commit.gpgsign": "YeS"},
			true,
		},
		{
			"Option commit.gpgsign is 1",
			map[string]string{"commit.gpgsign": "1"},
			true,
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			fake := NewFakeGitConfig(s.mockResponses)
			real := NewCachedGitConfig(
				func(cmd *exec.Cmd) (string, error) {
					assert.Equal(t, "config --get --null commit.gpgsign", strings.Join(cmd.Args[1:], " "))
					return fake.Get("commit.gpgsign"), nil
				},
				utils.NewDummyLog(),
			)
			result := real.GetBool("commit.gpgsign")
			assert.Equal(t, s.expected, result)
		})
	}
}

func TestGet(t *testing.T) {
	type scenario struct {
		testName      string
		mockResponses map[string]string
		expected      string
	}

	scenarios := []scenario{
		{
			"not set",
			map[string]string{},
			"",
		},
		{
			"is set",
			map[string]string{"commit.gpgsign": "blah"},
			"blah",
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			fake := NewFakeGitConfig(s.mockResponses)
			real := NewCachedGitConfig(
				func(cmd *exec.Cmd) (string, error) {
					assert.Equal(t, "config --get --null commit.gpgsign", strings.Join(cmd.Args[1:], " "))
					return fake.Get("commit.gpgsign"), nil
				},
				utils.NewDummyLog(),
			)
			result := real.Get("commit.gpgsign")
			assert.Equal(t, s.expected, result)
		})
	}

	// verifying that the cache is used
	count := 0
	real := NewCachedGitConfig(
		func(cmd *exec.Cmd) (string, error) {
			count++
			assert.Equal(t, "config --get --null commit.gpgsign", strings.Join(cmd.Args[1:], " "))
			return "blah", nil
		},
		utils.NewDummyLog(),
	)
	result := real.Get("commit.gpgsign")
	assert.Equal(t, "blah", result)
	result = real.Get("commit.gpgsign")
	assert.Equal(t, "blah", result)
	assert.Equal(t, 1, count)
}
