package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_dropDiffPrefix(t *testing.T) {
	scenarios := []struct {
		name           string
		diff           string
		expectedResult string
	}{
		{
			name:           "empty string",
			diff:           "",
			expectedResult: "",
		},
		{
			name: "only added lines",
			diff: `+line1
+line2
`,
			expectedResult: `line1
line2
`,
		},
		{
			name: "added lines with context",
			diff: ` line1
+line2
`,
			expectedResult: `line1
line2
`,
		},
		{
			name: "only deleted lines",
			diff: `-line1
-line2
`,
			expectedResult: `line1
line2
`,
		},
		{
			name: "deleted lines with context",
			diff: `-line1
 line2
`,
			expectedResult: `line1
line2
`,
		},
		{
			name: "only context",
			diff: ` line1
 line2
`,
			expectedResult: `line1
line2
`,
		},
		{
			name: "added and deleted lines",
			diff: `+line1
-line2
`,
			expectedResult: `+line1
-line2
`,
		},
		{
			name: "hunk header lines",
			diff: `@@ -1,8 +1,11 @@
 line1
`,
			expectedResult: `@@ -1,8 +1,11 @@
 line1
`,
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expectedResult, dropDiffPrefix(s.diff))
		})
	}
}
