package context

import (
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestListRenderer_renderLines(t *testing.T) {
	scenarios := []struct {
		name           string
		modelStrings   []string
		startIdx       int
		endIdx         int
		expectedOutput string
	}{
		{
			name:         "Render whole list",
			modelStrings: []string{"a", "b", "c"},
			startIdx:     0,
			endIdx:       3,
			expectedOutput: `
				a
				b
				c`,
		},
		{
			name:         "Partial list, beginning",
			modelStrings: []string{"a", "b", "c"},
			startIdx:     0,
			endIdx:       2,
			expectedOutput: `
				a
				b`,
		},
		{
			name:         "Partial list, end",
			modelStrings: []string{"a", "b", "c"},
			startIdx:     1,
			endIdx:       3,
			expectedOutput: `
				b
				c`,
		},
		{
			name:         "Pass an endIdx greater than the model length",
			modelStrings: []string{"a", "b", "c"},
			startIdx:     2,
			endIdx:       5,
			expectedOutput: `
				c`,
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			viewModel := NewListViewModel[string](func() []string { return s.modelStrings })
			self := &ListRenderer{
				list: viewModel,
				getDisplayStrings: func(startIdx int, endIdx int) [][]string {
					return lo.Map(s.modelStrings[startIdx:endIdx],
						func(s string, _ int) []string { return []string{s} })
				},
			}

			expectedOutput := strings.Join(lo.Map(
				strings.Split(strings.TrimPrefix(s.expectedOutput, "\n"), "\n"),
				func(line string, _ int) string { return strings.TrimSpace(line) }), "\n")

			assert.Equal(t, expectedOutput, self.renderLines(s.startIdx, s.endIdx))
		})
	}
}
