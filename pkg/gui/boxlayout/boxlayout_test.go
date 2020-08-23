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
			"Empty box",
			&Box{},
			0,
			0,
			10,
			10,
			func(result map[string]Dimensions) {
				assert.EqualValues(t, result, map[string]Dimensions{})
			},
		},
		{
			"Box with static and dynamic panel",
			&Box{Children: []*Box{{Size: 1, Window: "static"}, {Weight: 1, Window: "dynamic"}}},
			0,
			0,
			10,
			10,
			func(result map[string]Dimensions) {
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
			"Box with static and two dynamic panels",
			&Box{Children: []*Box{{Size: 1, Window: "static"}, {Weight: 1, Window: "dynamic1"}, {Weight: 2, Window: "dynamic2"}}},
			0,
			0,
			10,
			10,
			func(result map[string]Dimensions) {
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
			"Box with COLUMN direction",
			&Box{Direction: COLUMN, Children: []*Box{{Size: 1, Window: "static"}, {Weight: 1, Window: "dynamic1"}, {Weight: 2, Window: "dynamic2"}}},
			0,
			0,
			10,
			10,
			func(result map[string]Dimensions) {
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
			"Box with COLUMN direction only on wide boxes with narrow box",
			&Box{ConditionalDirection: func(width int, height int) int {
				if width > 4 {
					return COLUMN
				} else {
					return ROW
				}
			}, Children: []*Box{{Weight: 1, Window: "dynamic1"}, {Weight: 1, Window: "dynamic2"}}},
			0,
			0,
			4,
			4,
			func(result map[string]Dimensions) {
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
			"Box with COLUMN direction only on wide boxes with wide box",
			&Box{ConditionalDirection: func(width int, height int) int {
				if width > 4 {
					return COLUMN
				} else {
					return ROW
				}
			}, Children: []*Box{{Weight: 1, Window: "dynamic1"}, {Weight: 1, Window: "dynamic2"}}},
			0,
			0,
			5,
			5,
			func(result map[string]Dimensions) {
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
			"Box with conditional children where box is wide",
			&Box{ConditionalChildren: func(width int, height int) []*Box {
				if width > 4 {
					return []*Box{{Window: "wide", Weight: 1}}
				} else {
					return []*Box{{Window: "narrow", Weight: 1}}
				}
			}},
			0,
			0,
			5,
			5,
			func(result map[string]Dimensions) {
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
			"Box with conditional children where box is narrow",
			&Box{ConditionalChildren: func(width int, height int) []*Box {
				if width > 4 {
					return []*Box{{Window: "wide", Weight: 1}}
				} else {
					return []*Box{{Window: "narrow", Weight: 1}}
				}
			}},
			0,
			0,
			4,
			4,
			func(result map[string]Dimensions) {
				assert.EqualValues(
					t,
					result,
					map[string]Dimensions{
						"narrow": {X0: 0, X1: 3, Y0: 0, Y1: 3},
					},
				)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			s.test(ArrangeWindows(s.root, s.x0, s.y0, s.width, s.height))
		})
	}
}
