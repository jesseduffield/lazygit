package controllers

import (
	"testing"
)

func TestReplaceForkUsername_SSH_OK(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		forkUser string
		expected string
	}{
		{
			name:     "github ssh basic",
			in:       "git@github.com:old/repo.git",
			forkUser: "new",
			expected: "git@github.com:new/repo.git",
		},
		{
			name:     "ssh no .git",
			in:       "git@github.com:old/repo",
			forkUser: "new",
			expected: "git@github.com:new/repo",
		},
		{
			name:     "gitlab subgroup ssh",
			in:       "git@gitlab.com:group/sub/repo.git",
			forkUser: "alice",
			expected: "git@gitlab.com:alice/repo.git",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := replaceForkUsername(c.in, c.forkUser)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.expected {
				t.Fatalf("expected %q, got %q", c.expected, got)
			}
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
			got, err := replaceForkUsername(c.in, c.forkUser)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.expected {
				t.Fatalf("expected %q, got %q", c.expected, got)
			}
		})
	}
}

func TestReplaceForkUsername_Errors(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		forkUser string
	}{
		{
			name:     "empty fork user",
			in:       "git@github.com:old/repo.git",
			forkUser: "",
		},
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
			in:       "ssh://git@github.com/old/repo.git", // explicit ssh:// not supported here
			forkUser: "x",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := replaceForkUsername(c.in, c.forkUser)
			if err == nil {
				t.Fatalf("expected error but got nil")
			}
		})
	}
}
