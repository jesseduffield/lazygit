// +build darwin

package ps

import (
	"testing"
)

func TestDarwinProcess_impl(t *testing.T) {
	var _ Process = new(DarwinProcess)
}
