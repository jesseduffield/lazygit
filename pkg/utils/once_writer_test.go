package utils

import (
	"bytes"
	"testing"
)

func TestOnceWriter(t *testing.T) {
	innerWriter := bytes.NewBuffer(nil)
	counter := 0
	onceWriter := NewOnceWriter(innerWriter, func() {
		counter += 1
	})
	_, _ = onceWriter.Write([]byte("hello"))
	_, _ = onceWriter.Write([]byte("hello"))
	if counter != 1 {
		t.Errorf("expected counter to be 1, got %d", counter)
	}
}
