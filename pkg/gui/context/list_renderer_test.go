package context

import (
	"fmt"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// wrapping string in my own type to give it an ID method which is required for list items
type mystring string

func (self mystring) ID() string {
	return string(self)
}

func TestListRenderer_renderLines(t *testing.T) {
	scenarios := []struct {
		name            string
		modelStrings    []mystring
		nonModelIndices []int
		startIdx        int
		endIdx          int
		expectedOutput  string
	}{
		{
			name:         "Render whole list",
			modelStrings: []mystring{"a", "b", "c"},
			startIdx:     0,
			endIdx:       3,
			expectedOutput: `
				a
				b
				c`,
		},
		{
			name:         "Partial list, beginning",
			modelStrings: []mystring{"a", "b", "c"},
			startIdx:     0,
			endIdx:       2,
			expectedOutput: `
				a
				b`,
		},
		{
			name:         "Partial list, end",
			modelStrings: []mystring{"a", "b", "c"},
			startIdx:     1,
			endIdx:       3,
			expectedOutput: `
				b
				c`,
		},
		{
			name:         "Pass an endIdx greater than the model length",
			modelStrings: []mystring{"a", "b", "c"},
			startIdx:     2,
			endIdx:       5,
			expectedOutput: `
				c`,
		},
		{
			name:            "Whole list with section headers",
			modelStrings:    []mystring{"a", "b", "c"},
			nonModelIndices: []int{1, 3},
			startIdx:        0,
			endIdx:          5,
			expectedOutput: `
				a
				--- 1 (0) ---
				b
				c
				--- 3 (1) ---`,
		},
		{
			name:            "Multiple consecutive headers",
			modelStrings:    []mystring{"a", "b", "c"},
			nonModelIndices: []int{0, 0, 2, 2, 2},
			startIdx:        0,
			endIdx:          8,
			expectedOutput: `
				--- 0 (0) ---
				--- 0 (1) ---
				a
				b
				--- 2 (2) ---
				--- 2 (3) ---
				--- 2 (4) ---
				c`,
		},
		{
			name:            "Partial list with headers, beginning",
			modelStrings:    []mystring{"a", "b", "c"},
			nonModelIndices: []int{1, 3},
			startIdx:        0,
			endIdx:          3,
			expectedOutput: `
				a
				--- 1 (0) ---
				b`,
		},
		{
			name:            "Partial list with headers, end (beyond end index)",
			modelStrings:    []mystring{"a", "b", "c"},
			nonModelIndices: []int{1, 3},
			startIdx:        2,
			endIdx:          7,
			expectedOutput: `
				b
				c
				--- 3 (1) ---`,
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			viewModel := NewListViewModel[mystring](func() []mystring { return s.modelStrings })
			var getNonModelItems func() []*NonModelItem
			if s.nonModelIndices != nil {
				getNonModelItems = func() []*NonModelItem {
					return lo.Map(s.nonModelIndices, func(modelIndex int, nonModelIndex int) *NonModelItem {
						return &NonModelItem{
							Index:   modelIndex,
							Content: fmt.Sprintf("--- %d (%d) ---", modelIndex, nonModelIndex),
						}
					})
				}
			}
			self := &ListRenderer{
				list: viewModel,
				getDisplayStrings: func(startIdx int, endIdx int) [][]string {
					return lo.Map(s.modelStrings[startIdx:endIdx],
						func(s mystring, _ int) []string { return []string{string(s)} })
				},
				getNonModelItems: getNonModelItems,
			}

			expectedOutput := strings.Join(lo.Map(
				strings.Split(strings.TrimPrefix(s.expectedOutput, "\n"), "\n"),
				func(line string, _ int) string { return strings.TrimSpace(line) }), "\n")

			assert.Equal(t, expectedOutput, self.renderLines(s.startIdx, s.endIdx))
		})
	}
}

type myint int

func (self myint) ID() string {
	return fmt.Sprint(int(self))
}

func TestListRenderer_ModelIndexToViewIndex_and_back(t *testing.T) {
	scenarios := []struct {
		name            string
		numModelItems   int
		nonModelIndices []int

		modelIndices        []int
		expectedViewIndices []int

		viewIndices          []int
		expectedModelIndices []int
	}{
		{
			name:            "no headers (no getNonModelItems provided)",
			numModelItems:   3,
			nonModelIndices: nil, // no get

			modelIndices:        []int{-1, 0, 1, 2, 3, 4},
			expectedViewIndices: []int{0, 0, 1, 2, 3, 3},

			viewIndices:          []int{-1, 0, 1, 2, 3, 4},
			expectedModelIndices: []int{0, 0, 1, 2, 3, 3},
		},
		{
			name:            "no headers (getNonModelItems returns zero items)",
			numModelItems:   3,
			nonModelIndices: []int{},

			modelIndices:        []int{-1, 0, 1, 2, 3, 4},
			expectedViewIndices: []int{0, 0, 1, 2, 3, 3},

			viewIndices:          []int{-1, 0, 1, 2, 3, 4},
			expectedModelIndices: []int{0, 0, 1, 2, 3, 3},
		},
		{
			name:            "basic",
			numModelItems:   3,
			nonModelIndices: []int{1, 2},

			/*
				0: model 0
				1: --- header 0 ---
				2: model 1
				3: --- header 1 ---
				4: model 2
			*/

			modelIndices:        []int{-1, 0, 1, 2, 3, 4},
			expectedViewIndices: []int{0, 0, 2, 4, 5, 5},

			viewIndices:          []int{-1, 0, 1, 2, 3, 4, 5, 6},
			expectedModelIndices: []int{0, 0, 1, 1, 2, 2, 3, 3},
		},
		{
			name:            "consecutive section headers",
			numModelItems:   3,
			nonModelIndices: []int{0, 0, 2, 2, 2, 3, 3},

			/*
				0: --- header 0 ---
				1: --- header 1 ---
				2: model 0
				3: model 1
				4: --- header 2 ---
				5: --- header 3 ---
				6: --- header 4 ---
				7: model 2
				8: --- header 5 ---
				9: --- header 6 ---
			*/
			modelIndices:        []int{-1, 0, 1, 2, 3, 4},
			expectedViewIndices: []int{2, 2, 3, 7, 10, 10},

			viewIndices:          []int{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			expectedModelIndices: []int{0, 0, 0, 0, 1, 2, 2, 2, 2, 3, 3, 3, 3},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			// Expect lists of equal length for each test:
			assert.Equal(t, len(s.modelIndices), len(s.expectedViewIndices))
			assert.Equal(t, len(s.viewIndices), len(s.expectedModelIndices))

			modelInts := lo.Map(lo.Range(s.numModelItems), func(i int, _ int) myint { return myint(i) })
			viewModel := NewListViewModel[myint](func() []myint { return modelInts })
			var getNonModelItems func() []*NonModelItem
			if s.nonModelIndices != nil {
				getNonModelItems = func() []*NonModelItem {
					return lo.Map(s.nonModelIndices, func(modelIndex int, _ int) *NonModelItem {
						return &NonModelItem{Index: modelIndex, Content: ""}
					})
				}
			}
			self := &ListRenderer{
				list: viewModel,
				getDisplayStrings: func(startIdx int, endIdx int) [][]string {
					return lo.Map(modelInts[startIdx:endIdx],
						func(i myint, _ int) []string { return []string{fmt.Sprint(i)} })
				},
				getNonModelItems: getNonModelItems,
			}

			// Need to render first so that it knows the non-model items
			self.renderLines(-1, -1)

			for i := 0; i < len(s.modelIndices); i++ {
				assert.Equal(t, s.expectedViewIndices[i], self.ModelIndexToViewIndex(s.modelIndices[i]))
			}

			for i := 0; i < len(s.viewIndices); i++ {
				assert.Equal(t, s.expectedModelIndices[i], self.ViewIndexToModelIndex(s.viewIndices[i]))
			}
		})
	}
}
