package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceForkUsername_SSH_OK(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		forkUser string
		expected string
	}{
		{
			name:     "github ssh scp-like basic",
			in:       "git@github.com:old/repo.git",
			forkUser: "new",
			expected: "git@github.com:new/repo.git",
		},
		{
			name:     "ssh scp-like no .git",
			in:       "git@github.com:old/repo",
			forkUser: "new",
			expected: "git@github.com:new/repo",
		},
		{
			name:     "gitlab subgroup ssh scp-like",
			in:       "git@gitlab.com:group/sub/repo.git",
			forkUser: "alice",
			expected: "git@gitlab.com:alice/repo.git",
		},
		{
			name:     "ssh url style basic",
			in:       "ssh://git@github.com/old/repo.git",
			forkUser: "new",
			expected: "ssh://git@github.com/new/repo.git",
		},
		{
			name:     "ssh url style with port",
			in:       "ssh://git@github.com:2222/old/repo.git",
			forkUser: "bob",
			expected: "ssh://git@github.com:2222/bob/repo.git",
		},
		{
			name:     "ssh url style multi subgroup",
			in:       "ssh://git@gitlab.com/group/sub/repo.git",
			forkUser: "alice",
			expected: "ssh://git@gitlab.com/alice/repo.git",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := replaceForkUsername(c.in, c.forkUser, false)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, got)
		})
	}
}

func TestReplaceForkUsername_HTTPS_OK(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		forkUser string
		expected string
	}{
		{
			name:     "github https basic",
			in:       "https://github.com/old/repo.git",
			forkUser: "new",
			expected: "https://github.com/new/repo.git",
		},
		{
			name:     "https no .git",
			in:       "https://github.com/old/repo",
			forkUser: "new",
			expected: "https://github.com/new/repo",
		},
		{
			name:     "https with port",
			in:       "https://git.example.com:8443/group/repo",
			forkUser: "me",
			expected: "https://git.example.com:8443/me/repo",
		},
		{
			name:     "gitlab multi subgroup https",
			in:       "https://gitlab.com/group/sub/sub2/repo",
			forkUser: "bob",
			expected: "https://gitlab.com/bob/repo",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := replaceForkUsername(c.in, c.forkUser, false)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, got)
		})
	}
}

func TestReplaceForkUsername_IntegrationTest_OK(t *testing.T) {
	got, err := replaceForkUsername("../origin", "bob", true)
	assert.NoError(t, err)
	assert.Equal(t, "../bob", got)
}

func TestReplaceForkUsername_Errors(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		forkUser string
	}{
		{
			name:     "https host only",
			in:       "https://github.com",
			forkUser: "x",
		},
		{
			name:     "https host slash only",
			in:       "https://github.com/",
			forkUser: "x",
		},
		{
			name:     "https only repo (no owner)",
			in:       "https://github.com/repo.git",
			forkUser: "x",
		},
		{
			name:     "ssh missing path",
			in:       "git@github.com",
			forkUser: "x",
		},
		{
			name:     "ssh one segment only",
			in:       "git@github.com:repo.git",
			forkUser: "x",
		},
		{
			name:     "unsupported scheme",
			in:       "ftp://github.com/old/repo.git",
			forkUser: "x",
		},
		{
			name:     "empty url",
			in:       "",
			forkUser: "x",
		},
		{
			name:     "integration test URL outside of integration test",
			in:       "../origin",
			forkUser: "x",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := replaceForkUsername(c.in, c.forkUser, false)
			assert.EqualError(t, err, "unsupported or invalid remote URL: "+c.in)
		})
	}
}
