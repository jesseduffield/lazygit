package boxlayout

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrangeWindows(t *testing.T) {
	type scenario struct {
		testName string
		root     *Box
		x0       int
		y0       int
		width    int
		height   int
		test     func(result map[string]Dimensions)
	}

	scenarios := []scenario{
		{
			testName: "Empty box",
			root:     &Box{},
			x0:       0,
			y0:       0,
			width:    10,
			height:   10,
			test: func(result map[string]Dimensions) {
				assert.EqualValues(t, result, map[string]Dimensions{})
			},
		},
		{
			testName: "Box with static and dynamic panel",
			root:     &Box{Children: []*Box{{Size: 1, Window: "static"}, {Weight: 1, Window: "dynamic"}}},
			x0:       0,
			y0:       0,
			width:    10,
			height:   10,
			test: func(result map[string]Dimensions) {
				assert.EqualValues(
					t,
					result,
					map[string]Dimensions{
						"dynamic": {X0: 0, X1: 9, Y0: 1, Y1: 9},
						"static":  {X0: 0, X1: 9, Y0: 0, Y1: 0},
					},
				)
			},
		},
		{
			testName: "Box with static and two dynamic panels",
			root:     &Box{Children: []*Box{{Size: 1, Window: "static"}, {Weight: 1, Window: "dynamic1"}, {Weight: 2, Window: "dynamic2"}}},
			x0:       0,
			y0:       0,
			width:    10,
			height:   10,
			test: func(result map[string]Dimensions) {
				assert.EqualValues(
					t,
					result,
					map[string]Dimensions{
						"static":   {X0: 0, X1: 9, Y0: 0, Y1: 0},
						"dynamic1": {X0: 0, X1: 9, Y0: 1, Y1: 3},
						"dynamic2": {X0: 0, X1: 9, Y0: 4, Y1: 9},
					},
				)
			},
		},
		{
			testName: "Box with COLUMN direction",
			root:     &Box{Direction: COLUMN, Children: []*Box{{Size: 1, Window: "static"}, {Weight: 1, Window: "dynamic1"}, {Weight: 2, Window: "dynamic2"}}},
			x0:       0,
			y0:       0,
			width:    10,
			height:   10,
			test: func(result map[string]Dimensions) {
				assert.EqualValues(
					t,
					result,
					map[string]Dimensions{
						"static":   {X0: 0, X1: 0, Y0: 0, Y1: 9},
						"dynamic1": {X0: 1, X1: 3, Y0: 0, Y1: 9},
						"dynamic2": {X0: 4, X1: 9, Y0: 0, Y1: 9},
					},
				)
			},
		},
		{
			testName: "Box with COLUMN direction only on wide boxes with narrow box",
			root: &Box{ConditionalDirection: func(width int, height int) Direction {
				if width > 4 {
					return COLUMN
				} else {
					return ROW
				}
			}, Children: []*Box{{Weight: 1, Window: "dynamic1"}, {Weight: 1, Window: "dynamic2"}}},
			x0:     0,
			y0:     0,
			width:  4,
			height: 4,
			test: func(result map[string]Dimensions) {
				assert.EqualValues(
					t,
					result,
					map[string]Dimensions{
						"dynamic1": {X0: 0, X1: 3, Y0: 0, Y1: 1},
						"dynamic2": {X0: 0, X1: 3, Y0: 2, Y1: 3},
					},
				)
			},
		},
		{
			testName: "Box with COLUMN direction only on wide boxes with wide box",
			root: &Box{ConditionalDirection: func(width int, height int) Direction {
				if width > 4 {
					return COLUMN
				} else {
					return ROW
				}
			}, Children: []*Box{{Weight: 1, Window: "dynamic1"}, {Weight: 1, Window: "dynamic2"}}},
			// 5 / 2 = 2 remainder 1. That remainder goes to the first box.
			x0:     0,
			y0:     0,
			width:  5,
			height: 5,
			test: func(result map[string]Dimensions) {
				assert.EqualValues(
					t,
					result,
					map[string]Dimensions{
						"dynamic1": {X0: 0, X1: 2, Y0: 0, Y1: 4},
						"dynamic2": {X0: 3, X1: 4, Y0: 0, Y1: 4},
					},
				)
			},
		},
		{
			testName: "Box with conditional children where box is wide",
			root: &Box{ConditionalChildren: func(width int, height int) []*Box {
				if width > 4 {
					return []*Box{{Window: "wide", Weight: 1}}
				} else {
					return []*Box{{Window: "narrow", Weight: 1}}
				}
			}},
			x0:     0,
			y0:     0,
			width:  5,
			height: 5,
			test: func(result map[string]Dimensions) {
				assert.EqualValues(
					t,
					result,
					map[string]Dimensions{
						"wide": {X0: 0, X1: 4, Y0: 0, Y1: 4},
					},
				)
			},
		},
		{
			testName: "Box with conditional children where box is narrow",
			root: &Box{ConditionalChildren: func(width int, height int) []*Box {
				if width > 4 {
					return []*Box{{Window: "wide", Weight: 1}}
				} else {
					return []*Box{{Window: "narrow", Weight: 1}}
				}
			}},
			x0:     0,
			y0:     0,
			width:  4,
			height: 4,
			test: func(result map[string]Dimensions) {
				assert.EqualValues(
					t,
					result,
					map[string]Dimensions{
						"narrow": {X0: 0, X1: 3, Y0: 0, Y1: 3},
					},
				)
			},
		},
		{
			testName: "Box with static child with size too large",
			root:     &Box{Direction: COLUMN, Children: []*Box{{Size: 11, Window: "static"}, {Weight: 1, Window: "dynamic1"}, {Weight: 2, Window: "dynamic2"}}},
			x0:       0,
			y0:       0,
			width:    10,
			height:   10,
			test: func(result map[string]Dimensions) {
				assert.EqualValues(
					t,
					result,
					map[string]Dimensions{
						"static": {X0: 0, X1: 9, Y0: 0, Y1: 9},
						// not sure if X0: 10, X1: 9 makes any sense, but testing this in the
						// actual GUI it seems harmless
						"dynamic1": {X0: 10, X1: 9, Y0: 0, Y1: 9},
						"dynamic2": {X0: 10, X1: 9, Y0: 0, Y1: 9},
					},
				)
			},
		},
		{
			// 10 total space minus 2 from the status box leaves us with 8.
			// Total weight is 3, 8 / 3 = 2 with 2 remainder.
			// We want to end up with 2, 3, 5 (one unit from remainder to each dynamic box)
			testName: "Distributing remainder across weighted boxes",
			root:     &Box{Direction: COLUMN, Children: []*Box{{Size: 2, Window: "static"}, {Weight: 1, Window: "dynamic1"}, {Weight: 2, Window: "dynamic2"}}},
			x0:       0,
			y0:       0,
			width:    10,
			height:   10,
			test: func(result map[string]Dimensions) {
				assert.EqualValues(
					t,
					result,
					map[string]Dimensions{
						"static":   {X0: 0, X1: 1, Y0: 0, Y1: 9}, // 2
						"dynamic1": {X0: 2, X1: 4, Y0: 0, Y1: 9}, // 3
						"dynamic2": {X0: 5, X1: 9, Y0: 0, Y1: 9}, // 5
					},
				)
			},
		},
		{
			// 9 total space.
			// total weight is 5, 9 / 5 = 1 with 4 remainder
			// we want to give 2 of that remainder to the first, 1 to the second, and 1 to the last.
			// Reason being that we just give units to each box evenly and consider weight in subsequent passes.
			testName: "Distributing remainder across weighted boxes 2",
			root:     &Box{Direction: COLUMN, Children: []*Box{{Weight: 2, Window: "dynamic1"}, {Weight: 2, Window: "dynamic2"}, {Weight: 1, Window: "dynamic3"}}},
			x0:       0,
			y0:       0,
			width:    9,
			height:   10,
			test: func(result map[string]Dimensions) {
				assert.EqualValues(
					t,
					result,
					map[string]Dimensions{
						"dynamic1": {X0: 0, X1: 3, Y0: 0, Y1: 9}, // 4
						"dynamic2": {X0: 4, X1: 6, Y0: 0, Y1: 9}, // 3
						"dynamic3": {X0: 7, X1: 8, Y0: 0, Y1: 9}, // 2
					},
				)
			},
		},
		{
			// 9 total space.
			// total weight is 5, 9 / 5 = 1 with 4 remainder
			// we want to give 2 of that remainder to the first, 1 to the second, and 1 to the last.
			// Reason being that we just give units to each box evenly and consider weight in subsequent passes.
			testName: "Distributing remainder across weighted boxes with unnormalized weights",
			root:     &Box{Direction: COLUMN, Children: []*Box{{Weight: 4, Window: "dynamic1"}, {Weight: 4, Window: "dynamic2"}, {Weight: 2, Window: "dynamic3"}}},
			x0:       0,
			y0:       0,
			width:    9,
			height:   10,
			test: func(result map[string]Dimensions) {
				assert.EqualValues(
					t,
					result,
					map[string]Dimensions{
						"dynamic1": {X0: 0, X1: 3, Y0: 0, Y1: 9}, // 4
						"dynamic2": {X0: 4, X1: 6, Y0: 0, Y1: 9}, // 3
						"dynamic3": {X0: 7, X1: 8, Y0: 0, Y1: 9}, // 2
					},
				)
			},
		},
		{
			testName: "Another distribution test",
			root: &Box{Direction: COLUMN, Children: []*Box{
				{Weight: 3, Window: "dynamic1"},
				{Weight: 1, Window: "dynamic2"},
				{Weight: 1, Window: "dynamic3"},
			}},
			x0:     0,
			y0:     0,
			width:  9,
			height: 10,
			test: func(result map[string]Dimensions) {
				assert.EqualValues(
					t,
					result,
					map[string]Dimensions{
						"dynamic1": {X0: 0, X1: 4, Y0: 0, Y1: 9}, // 5
						"dynamic2": {X0: 5, X1: 6, Y0: 0, Y1: 9}, // 2
						"dynamic3": {X0: 7, X1: 8, Y0: 0, Y1: 9}, // 2
					},
				)
			},
		},
		{
			testName: "Box with zero weight",
			root: &Box{Direction: COLUMN, Children: []*Box{
				{Weight: 1, Window: "dynamic1"},
				{Weight: 0, Window: "dynamic2"},
			}},
			x0:     0,
			y0:     0,
			width:  10,
			height: 10,
			test: func(result map[string]Dimensions) {
				assert.EqualValues(
					t,
					result,
					map[string]Dimensions{
						"dynamic1": {X0: 0, X1: 9, Y0: 0, Y1: 9},
						"dynamic2": {X0: 10, X1: 9, Y0: 0, Y1: 9}, // when X0 > X1, we will hide the window
					},
				)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			s.test(ArrangeWindows(s.root, s.x0, s.y0, s.width, s.height))
		})
	}
}

func TestNormalizeWeights(t *testing.T) {
	scenarios := []struct {
		testName string
		input    []int
		expected []int
	}{
		{
			testName: "empty",
			input:    []int{},
			expected: []int{},
		},
		{
			testName: "one item of value 1",
			input:    []int{1},
			expected: []int{1},
		},
		{
			testName: "one item of value greater than 1",
			input:    []int{2},
			expected: []int{1},
		},
		{
			testName: "slice contains 1",
			input:    []int{2, 1},
			expected: []int{2, 1},
		},
		{
			testName: "slice contains 2 and 2",
			input:    []int{2, 2},
			expected: []int{1, 1},
		},
		{
			testName: "no common multiple",
			input:    []int{2, 3},
			expected: []int{2, 3},
		},
		{
			testName: "complex case",
			input:    []int{10, 10, 20},
			expected: []int{1, 1, 2},
		},
		{
			testName: "when a zero weight is included it is ignored",
			input:    []int{10, 10, 20, 0},
			expected: []int{1, 1, 2, 0},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			assert.EqualValues(t, s.expected, normalizeWeights(s.input))
		})
	}
}
