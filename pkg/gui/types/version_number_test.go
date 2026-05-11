package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVersionNumber(t *testing.T) {
	tests := []struct {
		versionStr string
		expected   *VersionNumber
		err        error
	}{
		{
			versionStr: "1.2.3",
			expected: &VersionNumber{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			err: nil,
		},
		{
			versionStr: "v1.2.3",
			expected: &VersionNumber{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			err: nil,
		},
		{
			versionStr: "12.34.56",
			expected: &VersionNumber{
				Major: 12,
				Minor: 34,
				Patch: 56,
			},
			err: nil,
		},
		{
			versionStr: "1.2",
			expected: &VersionNumber{
				Major: 1,
				Minor: 2,
				Patch: 0,
			},
			err: nil,
		},
		{
			versionStr: "1",
			expected:   nil,
			err:        errors.New("unexpected version format: 1"),
		},
		{
			versionStr: "invalid",
			expected:   nil,
			err:        errors.New("unexpected version format: invalid"),
		},
		{
			versionStr: "junk_before 1.2.3",
			expected:   nil,
			err:        errors.New("unexpected version format: junk_before 1.2.3"),
		},
		{
			versionStr: "1.2.3 junk_after",
			expected:   nil,
			err:        errors.New("unexpected version format: 1.2.3 junk_after"),
		},
	}

	for _, test := range tests {
		t.Run(test.versionStr, func(t *testing.T) {
			actual, err := ParseVersionNumber(test.versionStr)
			assert.Equal(t, test.expected, actual)
			assert.Equal(t, test.err, err)
		})
	}
}
