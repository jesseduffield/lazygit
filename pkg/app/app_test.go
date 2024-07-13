package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidGhVersion(t *testing.T) {
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
			`gh version 1.0.0 (2020-08-23)
			https://github.com/cli/cli/releases/tag/v1.0.0`,
			false,
		},
		{
			`gh version 2.0.0 (2021-08-23)
			https://github.com/cli/cli/releases/tag/v2.0.0`,
			true,
		},
		{
			`gh version 1.1.0 (2021-10-14)
			https://github.com/cli/cli/releases/tag/v1.1.0

			A new release of gh is available: 1.1.0 → v2.2.0
			To upgrade, run: brew update && brew upgrade gh
			https://github.com/cli/cli/releases/tag/v2.2.0`,
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.versionStr, func(t *testing.T) {
			result := isGhVersionValid(s.versionStr)
			assert.Equal(t, result, s.expectedResult)
		})
	}
}
