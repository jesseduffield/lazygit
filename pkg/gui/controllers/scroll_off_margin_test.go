package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_calculateLinesToScrollUp(t *testing.T) {
	scenarios := []struct {
		name                  string
		viewPortStart         int
		viewPortHeight        int
		scrollOffMargin       int
		lineIdxBefore         int
		lineIdxAfter          int
		expectedLinesToScroll int
	}{
		{
			name:                  "before position is above viewport - don't scroll",
			viewPortStart:         10,
			viewPortHeight:        10,
			scrollOffMargin:       3,
			lineIdxBefore:         9,
			lineIdxAfter:          8,
			expectedLinesToScroll: 0,
		},
		{
			name:                  "before position is below viewport - don't scroll",
			viewPortStart:         10,
			viewPortHeight:        10,
			scrollOffMargin:       3,
			lineIdxBefore:         20,
			lineIdxAfter:          19,
			expectedLinesToScroll: 0,
		},
		{
			name:                  "before and after positions are outside scroll-off margin - don't scroll",
			viewPortStart:         10,
			viewPortHeight:        10,
			scrollOffMargin:       3,
			lineIdxBefore:         14,
			lineIdxAfter:          13,
			expectedLinesToScroll: 0,
		},
		{
			name:                  "before outside, after inside scroll-off margin - scroll by 1",
			viewPortStart:         10,
			viewPortHeight:        10,
			scrollOffMargin:       3,
			lineIdxBefore:         13,
			lineIdxAfter:          12,
			expectedLinesToScroll: 1,
		},
		{
			name:                  "scroll-off margin is zero - scroll by 1 at end of view",
			viewPortStart:         10,
			viewPortHeight:        10,
			scrollOffMargin:       0,
			lineIdxBefore:         10,
			lineIdxAfter:          9,
			expectedLinesToScroll: 1,
		},
		{
			name:                  "before inside scroll-off margin - scroll by more than 1",
			viewPortStart:         10,
			viewPortHeight:        10,
			scrollOffMargin:       3,
			lineIdxBefore:         11,
			lineIdxAfter:          10,
			expectedLinesToScroll: 3,
		},
		{
			name:                  "very large scroll-off margin - keep view centered (even viewport height)",
			viewPortStart:         10,
			viewPortHeight:        10,
			scrollOffMargin:       999,
			lineIdxBefore:         15,
			lineIdxAfter:          14,
			expectedLinesToScroll: 1,
		},
		{
			name:                  "very large scroll-off margin - keep view centered (odd viewport height)",
			viewPortStart:         10,
			viewPortHeight:        9,
			scrollOffMargin:       999,
			lineIdxBefore:         14,
			lineIdxAfter:          13,
			expectedLinesToScroll: 1,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			linesToScroll := calculateLinesToScrollUp(scenario.viewPortStart, scenario.viewPortHeight, scenario.scrollOffMargin, scenario.lineIdxBefore, scenario.lineIdxAfter)
			assert.Equal(t, scenario.expectedLinesToScroll, linesToScroll)
		})
	}
}

func Test_calculateLinesToScrollDown(t *testing.T) {
	scenarios := []struct {
		name                  string
		viewPortStart         int
		viewPortHeight        int
		scrollOffMargin       int
		lineIdxBefore         int
		lineIdxAfter          int
		expectedLinesToScroll int
	}{
		{
			name:                  "before position is above viewport - don't scroll",
			viewPortStart:         10,
			viewPortHeight:        10,
			scrollOffMargin:       3,
			lineIdxBefore:         9,
			lineIdxAfter:          10,
			expectedLinesToScroll: 0,
		},
		{
			name:                  "before position is below viewport - don't scroll",
			viewPortStart:         10,
			viewPortHeight:        10,
			scrollOffMargin:       3,
			lineIdxBefore:         20,
			lineIdxAfter:          21,
			expectedLinesToScroll: 0,
		},
		{
			name:                  "before and after positions are outside scroll-off margin - don't scroll",
			viewPortStart:         10,
			viewPortHeight:        10,
			scrollOffMargin:       3,
			lineIdxBefore:         15,
			lineIdxAfter:          16,
			expectedLinesToScroll: 0,
		},
		{
			name:                  "before outside, after inside scroll-off margin - scroll by 1",
			viewPortStart:         10,
			viewPortHeight:        10,
			scrollOffMargin:       3,
			lineIdxBefore:         16,
			lineIdxAfter:          17,
			expectedLinesToScroll: 1,
		},
		{
			name:                  "scroll-off margin is zero - scroll by 1 at end of view",
			viewPortStart:         10,
			viewPortHeight:        10,
			scrollOffMargin:       0,
			lineIdxBefore:         19,
			lineIdxAfter:          20,
			expectedLinesToScroll: 1,
		},
		{
			name:                  "before inside scroll-off margin - scroll by more than 1",
			viewPortStart:         10,
			viewPortHeight:        10,
			scrollOffMargin:       3,
			lineIdxBefore:         18,
			lineIdxAfter:          19,
			expectedLinesToScroll: 3,
		},
		{
			name:                  "very large scroll-off margin - keep view centered (even viewport height)",
			viewPortStart:         10,
			viewPortHeight:        10,
			scrollOffMargin:       999,
			lineIdxBefore:         15,
			lineIdxAfter:          16,
			expectedLinesToScroll: 1,
		},
		{
			name:                  "very large scroll-off margin - keep view centered (odd viewport height)",
			viewPortStart:         10,
			viewPortHeight:        9,
			scrollOffMargin:       999,
			lineIdxBefore:         14,
			lineIdxAfter:          15,
			expectedLinesToScroll: 1,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			linesToScroll := calculateLinesToScrollDown(scenario.viewPortStart, scenario.viewPortHeight, scenario.scrollOffMargin, scenario.lineIdxBefore, scenario.lineIdxAfter)
			assert.Equal(t, scenario.expectedLinesToScroll, linesToScroll)
		})
	}
}
