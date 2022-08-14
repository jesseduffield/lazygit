package tasks

import (
	"bytes"
	"io"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func getCounter() (func(), func() int) {
	counter := 0
	return func() { counter++ }, func() int { return counter }
}

func TestNewCmdTaskInstantStop(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	beforeStart, getBeforeStartCallCount := getCounter()
	refreshView, getRefreshViewCallCount := getCounter()
	onEndOfInput, getOnEndOfInputCallCount := getCounter()
	onNewKey, getOnNewKeyCallCount := getCounter()
	onDone, getOnDoneCallCount := getCounter()

	manager := NewViewBufferManager(
		utils.NewDummyLog(),
		writer,
		beforeStart,
		refreshView,
		onEndOfInput,
		onNewKey,
	)

	stop := make(chan struct{})
	reader := bytes.NewBufferString("test")
	start := func() (*exec.Cmd, io.Reader) {
		// not actually starting this because it's not necessary
		cmd := secureexec.Command("blah blah")

		close(stop)

		return cmd, reader
	}

	fn := manager.NewCmdTask(start, "prefix\n", 20, onDone)

	_ = fn(stop)

	callCountExpectations := []struct {
		expected int
		actual   int
		name     string
	}{
		{0, getBeforeStartCallCount(), "beforeStart"},
		{1, getRefreshViewCallCount(), "refreshView"},
		{0, getOnEndOfInputCallCount(), "onEndOfInput"},
		{0, getOnNewKeyCallCount(), "onNewKey"},
		{1, getOnDoneCallCount(), "onDone"},
	}
	for _, expectation := range callCountExpectations {
		if expectation.actual != expectation.expected {
			t.Errorf("expected %s to be called %d times, got %d", expectation.name, expectation.expected, expectation.actual)
		}
	}

	expectedContent := ""
	actualContent := writer.String()
	if actualContent != expectedContent {
		t.Errorf("expected writer to receive the following content: \n%s\n. But instead it received: %s", expectedContent, actualContent)
	}
}

func TestNewCmdTask(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	beforeStart, getBeforeStartCallCount := getCounter()
	refreshView, getRefreshViewCallCount := getCounter()
	onEndOfInput, getOnEndOfInputCallCount := getCounter()
	onNewKey, getOnNewKeyCallCount := getCounter()
	onDone, getOnDoneCallCount := getCounter()

	manager := NewViewBufferManager(
		utils.NewDummyLog(),
		writer,
		beforeStart,
		refreshView,
		onEndOfInput,
		onNewKey,
	)

	stop := make(chan struct{})
	reader := bytes.NewBufferString("test")
	start := func() (*exec.Cmd, io.Reader) {
		// not actually starting this because it's not necessary
		cmd := secureexec.Command("blah blah")

		return cmd, reader
	}

	fn := manager.NewCmdTask(start, "prefix\n", 20, onDone)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(stop)
		wg.Done()
	}()
	_ = fn(stop)

	wg.Wait()

	callCountExpectations := []struct {
		expected int
		actual   int
		name     string
	}{
		{1, getBeforeStartCallCount(), "beforeStart"},
		{1, getRefreshViewCallCount(), "refreshView"},
		{1, getOnEndOfInputCallCount(), "onEndOfInput"},
		{0, getOnNewKeyCallCount(), "onNewKey"},
		{1, getOnDoneCallCount(), "onDone"},
	}
	for _, expectation := range callCountExpectations {
		if expectation.actual != expectation.expected {
			t.Errorf("expected %s to be called %d times, got %d", expectation.name, expectation.expected, expectation.actual)
		}
	}

	expectedContent := "prefix\ntest\n"
	actualContent := writer.String()
	if actualContent != expectedContent {
		t.Errorf("expected writer to receive the following content: \n%s\n. But instead it received: %s", expectedContent, actualContent)
	}
}
