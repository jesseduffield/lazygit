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
	beginRender, getBeginRenderCallCount := getCounter()
	swapInRender, getSwapInRenderCallCount := getCounter()
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
		beginRender,
		swapInRender,
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

	fn := manager.NewCmdTask(start, "prefix\n", LinesToRead{Total: 20, InitialRefreshAfter: -1}, onDone)

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
		{0, getBeginRenderCallCount(), "beginRender"},
		{0, getSwapInRenderCallCount(), "swapInRender"},
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
	beginRender, getBeginRenderCallCount := getCounter()
	swapInRender, getSwapInRenderCallCount := getCounter()
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
		beginRender,
		swapInRender,
		newTask,
	)

	stop := make(chan struct{})
	reader := bytes.NewBufferString("test")
	start := func() (*exec.Cmd, io.Reader) {
		// not actually starting this because it's not necessary
		cmd := exec.Command("blah")

		return cmd, reader
	}

	fn := manager.NewCmdTask(start, "prefix\n", LinesToRead{Total: 20, InitialRefreshAfter: -1}, onDone)
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
		{0, getBeforeStartCallCount(), "beforeStart"},
		{1, getRefreshViewCallCount(), "refreshView"},
		{1, getOnEndOfInputCallCount(), "onEndOfInput"},
		{0, getOnNewKeyCallCount(), "onNewKey"},
		{1, getBeginRenderCallCount(), "beginRender"},
		{1, getSwapInRenderCallCount(), "swapInRender"},
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
// When a RenderRestore is set, the first paint is driven by its FirstPaintReady
// predicate rather than the InitialRefreshAfter line count, and Apply runs exactly
// once. Apply controls when the off-screen render is swapped in (it calls swapIn
// after locating the target), so its scan runs while the previous content is still
// displayed rather than revealing the new content at the old scroll. This is the
// read-loop half of the escape restore. See RenderRestore.
func TestNewCmdTaskRestore(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	linesWritten := func() int { return strings.Count(writer.String(), "\n") }

	swapped := false
	applyCount := 0
	applyAtLines := -1
	swappedBeforeApply := false
	swappedByApply := false

	task := gocui.NewFakeTask()
	manager := NewViewBufferManager(
		utils.NewDummyLog(),
		writer,
		func() {},                 // beforeStart
		func() {},                 // refreshView
		func() {},                 // onEndOfInput
		func() {},                 // onNewKey
		func() {},                 // beginRender
		func() { swapped = true }, // swapInRender
		func() gocui.Task { return task },
	)

	restore := &RenderRestore{
		// Ready once five lines have loaded — well before InitialRefreshAfter (30).
		FirstPaintReady: func() bool { return linesWritten() >= 5 },
		Apply: func(swapIn func()) {
			applyCount++
			applyAtLines = linesWritten()
			// The render must not be swapped in until Apply asks for it.
			if swapped {
				swappedBeforeApply = true
			}
			swapIn()
			swappedByApply = swapped
		},
	}

	stop := make(chan struct{})
	reader := BlankLineReader{totalLinesToYield: 50}
	start := func() (*exec.Cmd, io.Reader) {
		return exec.Command("blah"), &reader
	}
	fn := manager.NewCmdTask(start, "", LinesToRead{Total: 50, InitialRefreshAfter: 30, Restore: restore}, func() {})

	wg := sync.WaitGroup{}
	wg.Go(func() {
		time.Sleep(100 * time.Millisecond)
		close(stop)
	})
	_ = fn(TaskOpts{Stop: stop, InitialContentLoaded: func() { task.Done() }})
	wg.Wait()

	assert.Equal(t, 1, applyCount, "Apply should run exactly once")
	assert.False(t, swappedBeforeApply, "the off-screen render should not be swapped in before Apply runs")
	assert.True(t, swappedByApply, "Apply should swap the off-screen render in via swapIn")
	// The first paint was driven by FirstPaintReady (>=5 lines), not by
	// InitialRefreshAfter (30).
	assert.GreaterOrEqual(t, applyAtLines, 5)
	assert.Less(t, applyAtLines, 30)
}

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
			LinesToRead{Total: 100, InitialRefreshAfter: 120},
			[]int{100},
		},
		{
			"total == initialRefreshAfter",
			150,
			LinesToRead{Total: 100, InitialRefreshAfter: 100},
			[]int{100},
		},
		{
			"total > initialRefreshAfter",
			150,
			LinesToRead{Total: 100, InitialRefreshAfter: 50},
			[]int{50, 100},
		},
		{
			"initialRefreshAfter == -1",
			150,
			LinesToRead{Total: 100, InitialRefreshAfter: -1},
			[]int{100},
		},
		{
			"totalTaskLines < initialRefreshAfter",
			25,
			LinesToRead{Total: 100, InitialRefreshAfter: 50},
			[]int{25},
		},
		{
			"totalTaskLines between total and initialRefreshAfter",
			75,
			LinesToRead{Total: 100, InitialRefreshAfter: 50},
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
