package gui

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMenuItemHandling(t *testing.T) {

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

	t.Run("formatListItems", func(t *testing.T) {

		type scenario struct {
			testName  string
			menuItems []*menuItem
			assert    func(data string, err error)
		}

		scenarios := []scenario{
			{
				testName:  "should return empty formatted result on empty input",
				menuItems: []*menuItem{},
				assert: func(data string, err error) {
					assert.Nil(t, err)
					assert.Equal(t, "", data)
				},
			},
			{
				testName: "should return formatted version of single item",
				menuItems: []*menuItem{
					{
						displayString: "FAKE_DISPLAY_STRING",
					}},
				assert: func(data string, err error) {
					assert.Nil(t, err)
					assert.Equal(t, "FAKE_DISPLAY_STRING", data)
				},
			}, {
				testName:  "should return formatted version of multiple items",
				menuItems: fakeMenuItems,
				assert: func(data string, err error) {
					assert.Nil(t, err)
					assert.Equal(t, "FAKE_DISPLAY_STRING_ONE\nFAKE_DISPLAY_STRING_TWO\nFAKE_DISPLAY_STRING_THREE", data)
				},
			},
		}

		for _, scenario := range scenarios {
			t.Run(scenario.testName, func(t *testing.T) {
				result, err := formatListItems(scenario.menuItems)
				scenario.assert(result, err)
			})
		}

	})
}
