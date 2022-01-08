package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGitVersionValid(t *testing.T) {
	type scenario struct {
		versionStr     string
		expectedResult bool
	}

	scenarios := []scenario{
		{
			"",
			false,
		},
		{
			"git version 1.9.0",
			false,
		},
		{
			"git version 1.9.0 (Apple Git-128)",
			false,
		},
		{
			"git version 2.4.0",
			true,
		},
		{
			"git version 2.24.3 (Apple Git-128)",
			true,
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.versionStr, func(t *testing.T) {
			result := isGitVersionValid(s.versionStr)
			assert.Equal(t, result, s.expectedResult)
		})
	}
}
