package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHistoryBuffer(t *testing.T) {
	hb := NewHistoryBuffer[int](5)
	assert.NotNil(t, hb)
	assert.Equal(t, 5, hb.maxSize)
	assert.Equal(t, 0, len(hb.items))
}

func TestPush(t *testing.T) {
	hb := NewHistoryBuffer[int](3)
	hb.Push(1)
	hb.Push(2)
	hb.Push(3)
	hb.Push(4)

	assert.Equal(t, 3, len(hb.items))
	assert.Equal(t, []int{4, 3, 2}, hb.items)
}

func TestPeekAt(t *testing.T) {
	hb := NewHistoryBuffer[int](3)
	hb.Push(1)
	hb.Push(2)
	hb.Push(3)

	item, err := hb.PeekAt(0)
	assert.Nil(t, err)
	assert.Equal(t, 3, item)

	item, err = hb.PeekAt(1)
	assert.Nil(t, err)
	assert.Equal(t, 2, item)

	item, err = hb.PeekAt(2)
	assert.Nil(t, err)
	assert.Equal(t, 1, item)

	item, err = hb.PeekAt(-1)
	assert.Nil(t, err)
	assert.Equal(t, 0, item)

	_, err = hb.PeekAt(3)
	assert.NotNil(t, err)
	assert.Equal(t, "Index out of range", err.Error())

	_, err = hb.PeekAt(-2)
	assert.NotNil(t, err)
	assert.Equal(t, "Index out of range", err.Error())
}

func TestPeekAtEmptyBuffer(t *testing.T) {
	hb := NewHistoryBuffer[int](3)

	_, err := hb.PeekAt(0)
	assert.NotNil(t, err)
	assert.Equal(t, "Buffer is empty", err.Error())
}
