package gui

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMenuPanelFilter(t *testing.T) {

	fakeMenuItems :=
		[]*menuItem{
			{
				displayStrings: []string{"a", "FAKE_DISPLAY_STRING_ONE"},
			}, {
				displayStrings: []string{"a", "FAKE_DISPLAY_STRING_TWO"},
			}, {
				displayStrings: []string{"a", "FAKE_DISPLAY_STRING_THREE"},
			},
		}

	t.Run("filterListItems", func(t *testing.T) {
		type scenario struct {
			testName  string
			menuItems []*menuItem
			filter    string
			assert    func(filteredItems []*menuItem, err error)
		}

		scenarios := []scenario{
			{
				testName:  "should return input data when no filter provided",
				menuItems: fakeMenuItems,
				filter:    "",
				assert: func(filteredItems []*menuItem, err error) {
					assert.Equal(t, fakeMenuItems, filteredItems)
				},
			}, {
				testName:  "should return no data when filter matches no elements",
				menuItems: fakeMenuItems,
				filter:    "NOT_IN_THAT_HAYSTACK",
				assert: func(filteredItems []*menuItem, err error) {
					assert.Empty(t, filteredItems)
				},
			}, {
				testName:  "should return item when single match",
				menuItems: fakeMenuItems,
				filter:    "ONE",
				assert: func(filteredItems []*menuItem, err error) {
					assert.Len(t, filteredItems, 1)
					assert.Equal(t, filteredItems[0], fakeMenuItems[0])
				},
			},
		}

		for _, scenario := range scenarios {
			t.Run(scenario.testName, func(t *testing.T) {
				result, err := filterListItems(scenario.menuItems, scenario.filter)
				scenario.assert(result, err)
			})
		}

	})
}
