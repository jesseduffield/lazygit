package gocui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyIsPrintable(t *testing.T) {
	assert.True(t, NewKeyRune('x').IsPrintable())
	assert.True(t, NewKeyRune('界').IsPrintable())
	assert.True(t, NewKeyRune(' ').IsPrintable())
	assert.False(t, NewKeyStrMod("x", ModCtrl).IsPrintable())
	assert.False(t, NewKeyName(KeyEnter).IsPrintable())
	assert.False(t, Key{}.IsPrintable())
}
