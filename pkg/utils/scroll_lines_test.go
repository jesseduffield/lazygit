package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScrollHeight(t *testing.T) {
	type scenario struct {
		windowHeight           int
		scrollHeightFromConfig float64
		expected               int
	}

	scenarios := []scenario{
		{
			10,
			5,
			5,
		},
		{
			1000,
			-5,
			-5,
		},
		{
			10,
			-5.2,
			-5,
		},
		{
			10,
			0.5,
			5,
		},
		{
			12,
			-0.25,
			-3,
		},
		{
			9,
			0.5,
			4,
		},
		{
			9,
			-0.5,
			-4,
		},
		{
			1,
			-0.5,
			-1,
		},
		{
			1,
			0.5,
			1,
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, ScrollHeight(s.windowHeight, s.scrollHeightFromConfig))
	}
}
