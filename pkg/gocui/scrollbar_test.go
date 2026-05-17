package gocui

import "testing"

func TestCalcScrollbar(t *testing.T) {
	tests := []struct {
		testName       string
		listSize       int
		pageSize       int
		position       int
		scrollAreaSize int

		expectedStart  int
		expectedHeight int
	}{
		{
			testName:       "page size greater than list size",
			listSize:       5,
			pageSize:       10,
			position:       0,
			scrollAreaSize: 20,

			expectedStart:  0,
			expectedHeight: 20,
		},
		{
			testName:       "page size matches list size",
			listSize:       10,
			pageSize:       10,
			position:       0,
			scrollAreaSize: 20,

			expectedStart:  0,
			expectedHeight: 20,
		},
		{
			testName:       "page size half of list size",
			listSize:       10,
			pageSize:       5,
			position:       0,
			scrollAreaSize: 20,

			expectedStart:  0,
			expectedHeight: 10,
		},
		{
			testName:       "page size half of list size at scroll end",
			listSize:       10,
			pageSize:       5,
			position:       5,
			scrollAreaSize: 20,

			expectedStart:  10,
			expectedHeight: 10,
		},
		{
			testName: "page size third of list size having scrolled half the way",
			listSize: 15,
			// Recall that my max position is listSize - pageSize i.e 15 - 5 i.e. 10.
			// So if I've scrolled to position 5 that means I've done one page and I've got
			// one page to go which means by scrollbar should take up a third of the available
			// space and appear in the centre of the scrollbar area
			pageSize:       5,
			position:       5,
			scrollAreaSize: 21,

			expectedStart:  7,
			expectedHeight: 7,
		},
		{
			testName:       "page size third of list size having scrolled the full way",
			listSize:       15,
			pageSize:       5,
			position:       10,
			scrollAreaSize: 21,

			expectedStart:  14,
			expectedHeight: 7,
		},
		{
			testName:       "page size third of list size having scrolled by one",
			listSize:       15,
			pageSize:       5,
			position:       1,
			scrollAreaSize: 21,

			expectedStart:  2,
			expectedHeight: 7,
		},
		{
			testName:       "page size third of list size having scrolled up from the bottom by one",
			listSize:       15,
			pageSize:       5,
			position:       9,
			scrollAreaSize: 21,

			expectedStart:  12,
			expectedHeight: 7,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			start, height := calcScrollbar(test.listSize, test.pageSize, test.position, test.scrollAreaSize)
			if start != test.expectedStart {
				t.Errorf("expected start to be %d, got %d", test.expectedStart, start)
			}

			if height != test.expectedHeight {
				t.Errorf("expected height to be %d, got %d", test.expectedHeight, height)
			}
		})
	}
}
