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

	"github.com/jesseduffield/gocui"
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
	)

	stop := make(chan struct{})
	reader := bytes.NewBufferString("test")
	start := func() (*exec.Cmd, io.Reader) {
		// not actually starting this because it's not necessary
		cmd := exec.Command("blah")

		close(stop)

		return cmd, reader
	}

	fn := manager.NewCmdTask(start, "prefix\n", LinesToRead{20, -1}, onDone)

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
	)

	stop := make(chan struct{})
	reader := bytes.NewBufferString("test")
	start := func() (*exec.Cmd, io.Reader) {
		// not actually starting this because it's not necessary
		cmd := exec.Command("blah")

		return cmd, reader
	}

	fn := manager.NewCmdTask(start, "prefix\n", LinesToRead{20, -1}, onDone)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(stop)
		wg.Done()
	}()
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

	d.linesYielded += 1
	p[0] = '\n'
	return 1, nil
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
			LinesToRead{100, 120},
			[]int{100},
		},
		{
			"total == initialRefreshAfter",
			150,
			LinesToRead{100, 100},
			[]int{100},
		},
		{
			"total > initialRefreshAfter",
			150,
			LinesToRead{100, 50},
			[]int{50, 100},
		},
		{
			"initialRefreshAfter == -1",
			150,
			LinesToRead{100, -1},
			[]int{100},
		},
		{
			"totalTaskLines < initialRefreshAfter",
			25,
			LinesToRead{100, 50},
			[]int{25},
		},
		{
			"totalTaskLines between total and initialRefreshAfter",
			75,
			LinesToRead{100, 50},
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
		)

		stop := make(chan struct{})
		reader := BlankLineReader{totalLinesToYield: s.totalTaskLines}
		start := func() (*exec.Cmd, io.Reader) {
			// not actually starting this because it's not necessary
			cmd := exec.Command("blah")

			return cmd, &reader
		}

		fn := manager.NewCmdTask(start, "", s.linesToRead, func() {})
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			time.Sleep(100 * time.Millisecond)
			close(stop)
			wg.Done()
		}()
		_ = fn(TaskOpts{Stop: stop, InitialContentLoaded: func() { task.Done() }})

		wg.Wait()

		if !reflect.DeepEqual(lineCountsOnRefresh, s.expectedLineCountsOnRefresh) {
			t.Errorf("%s: expected line counts on refresh: %v, got %v",
				s.name, s.expectedLineCountsOnRefresh, lineCountsOnRefresh)
		}
	}
}
