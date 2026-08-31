package gui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinimumScreenHeight(t *testing.T) {
	assert.Equal(t, 9, minimumScreenHeight(5, false))
	assert.Equal(t, 12, minimumScreenHeight(8, false))
	assert.Equal(t, 11, minimumScreenHeight(5, true))
	assert.Equal(t, 12, minimumScreenHeight(8, true))
}
