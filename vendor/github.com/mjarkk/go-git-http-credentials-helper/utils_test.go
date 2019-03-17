package gitcredentialhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFreePort(t *testing.T) {
	out := getFreePort()
	assert.Equal(t, len(out), 4)
}
