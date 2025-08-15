package lo

import (
	"os"
	"testing"
	"time"
)

// https://github.com/stretchr/testify/issues/1101
func testWithTimeout(t *testing.T, timeout time.Duration) {
	t.Helper()

	testFinished := make(chan struct{})
	t.Cleanup(func() { close(testFinished) })

	go func() {
		select {
		case <-testFinished:
		case <-time.After(timeout):
			t.Errorf("test timed out after %s", timeout)
			os.Exit(1)
		}
	}()
}

type foo struct {
	bar string
}

func (f foo) Clone() foo {
	return foo{f.bar}
}
