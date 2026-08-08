package tasks

import (
	"bytes"
	"io"
	"os/exec"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
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
	task := gocui.NewFakeTask()
	newTask := func() gocui.Task {
		return task
	}

	manager := NewViewBufferManager(
		utils.NewDummyLog(),
		writer,
		beforeStart,
		refreshView,
		onEndOfInput,
		onNewKey,
		newTask,
		// no UI thread in the test; run the view mutations inline
		func(f func()) error { f(); return nil },
	)

	stop := make(chan struct{})
	reader := bytes.NewBufferString("test")
	start := func() (Cmd, io.Reader) {
		// not actually starting this because it's not necessary
		cmd := exec.Command("blah")

		close(stop)

		return ExecCmd{Cmd: cmd}, reader
	}

	fn := manager.NewCmdTask(start, "prefix\n", LinesToRead{20, -1, nil}, onDone)

	_ = fn(TaskOpts{Stop: stop, InitialContentLoaded: func() { task.Done() }})

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

	if task.Status() != gocui.TaskStatusDone {
		t.Errorf("expected task status to be 'done', got '%s'", task.FormatStatus())
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
	task := gocui.NewFakeTask()
	newTask := func() gocui.Task {
		return task
	}

	manager := NewViewBufferManager(
		utils.NewDummyLog(),
		writer,
		beforeStart,
		refreshView,
		onEndOfInput,
		onNewKey,
		newTask,
		// no UI thread in the test; run the view mutations inline
		func(f func()) error { f(); return nil },
	)

	stop := make(chan struct{})
	reader := bytes.NewBufferString("test")
	start := func() (Cmd, io.Reader) {
		// not actually starting this because it's not necessary
		cmd := exec.Command("blah")

		return ExecCmd{Cmd: cmd}, reader
	}

	fn := manager.NewCmdTask(start, "prefix\n", LinesToRead{20, -1, nil}, onDone)
	wg := sync.WaitGroup{}
	wg.Go(func() {
		time.Sleep(100 * time.Millisecond)
		close(stop)
	})
	_ = fn(TaskOpts{Stop: stop, InitialContentLoaded: func() { task.Done() }})

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

	if task.Status() != gocui.TaskStatusDone {
		t.Errorf("expected task status to be 'done', got '%s'", task.FormatStatus())
	}

	expectedContent := "prefix\ntest\n"
	actualContent := writer.String()
	if actualContent != expectedContent {
		t.Errorf("expected writer to receive the following content: \n%s\n. But instead it received: %s", expectedContent, actualContent)
	}
}

// A dummy reader that simply yields as many blank lines as requested. The only
// thing we want to do with the output is count the number of lines.
type BlankLineReader struct {
	totalLinesToYield int
	linesYielded      int
}

func (d *BlankLineReader) Read(p []byte) (n int, err error) {
	if d.totalLinesToYield == d.linesYielded {
		return 0, io.EOF
	}

	d.linesYielded++
	p[0] = '\n'
	return 1, nil
}

// A dummy reader that yields the given number of blank lines and then blocks
// until unblock is closed, at which point it reports EOF. This lets a test hold
// a task in its "still loading" state for as long as it needs to.
type BlockingLineReader struct {
	linesToYield int
	linesYielded int
	reachedEnd   bool
	blocked      chan struct{}
	unblock      chan struct{}
}

func (d *BlockingLineReader) Read(p []byte) (n int, err error) {
	if d.linesYielded == d.linesToYield {
		if !d.reachedEnd {
			d.reachedEnd = true
			close(d.blocked)
		}
		<-d.unblock
		return 0, io.EOF
	}

	d.linesYielded++
	p[0] = '\n'
	return 1, nil
}

func TestNewCmdTaskQueuedReadAtEndOfInput(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	task := gocui.NewFakeTask()

	manager := NewViewBufferManager(
		utils.NewDummyLog(),
		writer,
		func() {},
		func() {},
		func() {},
		func() {},
		func() gocui.Task { return task },
		// no UI thread in the test; run the view mutations inline
		func(f func()) error { f(); return nil },
	)

	reader := BlockingLineReader{
		linesToYield: 5,
		blocked:      make(chan struct{}),
		unblock:      make(chan struct{}),
	}
	start := func() (Cmd, io.Reader) {
		// not actually starting this because it's not necessary
		return ExecCmd{Cmd: exec.Command("blah")}, &reader
	}

	// The initial request asks for far more lines than the reader has, so the
	// task reaches EOF while that request is still the one being served.
	fn := manager.NewCmdTask(start, "", LinesToRead{100, -1, nil}, func() {})

	thenCalled := false
	wg := sync.WaitGroup{}
	wg.Go(func() {
		_ = fn(TaskOpts{Stop: make(chan struct{}), InitialContentLoaded: func() { task.Done() }})
	})

	<-reader.blocked
	manager.ReadToEnd(func() { thenCalled = true })
	// ReadToEnd queues its request from a goroutine; wait for it to land so that
	// it is definitely outstanding by the time we let the task reach EOF.
	for len(*manager.readLines.Load()) == 0 {
		time.Sleep(time.Millisecond)
	}
	close(reader.unblock)

	wg.Wait()

	/* EXPECTED:
	assert.True(t, thenCalled)
	ACTUAL: */
	assert.False(t, thenCalled)
}

func TestNewCmdTaskRefresh(t *testing.T) {
	type scenario struct {
		name                        string
		totalTaskLines              int
		linesToRead                 LinesToRead
		expectedLineCountsOnRefresh []int
	}

	scenarios := []scenario{
		{
			"total < initialRefreshAfter",
			150,
			LinesToRead{100, 120, nil},
			[]int{100},
		},
		{
			"total == initialRefreshAfter",
			150,
			LinesToRead{100, 100, nil},
			[]int{100},
		},
		{
			"total > initialRefreshAfter",
			150,
			LinesToRead{100, 50, nil},
			[]int{50, 100},
		},
		{
			"initialRefreshAfter == -1",
			150,
			LinesToRead{100, -1, nil},
			[]int{100},
		},
		{
			"totalTaskLines < initialRefreshAfter",
			25,
			LinesToRead{100, 50, nil},
			[]int{25},
		},
		{
			"totalTaskLines between total and initialRefreshAfter",
			75,
			LinesToRead{100, 50, nil},
			[]int{50, 75},
		},
	}

	for _, s := range scenarios {
		writer := bytes.NewBuffer(nil)
		lineCountsOnRefresh := []int{}
		refreshView := func() {
			lineCountsOnRefresh = append(lineCountsOnRefresh, strings.Count(writer.String(), "\n"))
		}

		task := gocui.NewFakeTask()
		newTask := func() gocui.Task {
			return task
		}

		manager := NewViewBufferManager(
			utils.NewDummyLog(),
			writer,
			func() {},
			refreshView,
			func() {},
			func() {},
			newTask,
			// no UI thread in the test; run the view mutations inline
			func(f func()) error { f(); return nil },
		)

		stop := make(chan struct{})
		reader := BlankLineReader{totalLinesToYield: s.totalTaskLines}
		start := func() (Cmd, io.Reader) {
			// not actually starting this because it's not necessary
			cmd := exec.Command("blah")

			return ExecCmd{Cmd: cmd}, &reader
		}

		fn := manager.NewCmdTask(start, "", s.linesToRead, func() {})
		wg := sync.WaitGroup{}
		wg.Go(func() {
			time.Sleep(100 * time.Millisecond)
			close(stop)
		})
		_ = fn(TaskOpts{Stop: stop, InitialContentLoaded: func() { task.Done() }})

		wg.Wait()

		if !reflect.DeepEqual(lineCountsOnRefresh, s.expectedLineCountsOnRefresh) {
			t.Errorf("%s: expected line counts on refresh: %v, got %v",
				s.name, s.expectedLineCountsOnRefresh, lineCountsOnRefresh)
		}
	}
}
