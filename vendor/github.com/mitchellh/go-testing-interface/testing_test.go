package testing

import (
	"testing"
)

func TestT(t *testing.T) {
	testTFunc(t) // Just verify this doesn't give a compiler error
}

func TestRuntimeT(t *testing.T) {
	var _ T = new(RuntimeT) // Another compiler check
}

func testTFunc(t T) {}
