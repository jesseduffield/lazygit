package tasks

import (
	"bytes"
	"io"
	"os/exec"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
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
	resetOrigin, getResetOriginCallCount := getCounter()
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
		resetOrigin,
		beginRender,
		swapInRender,
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
		{0, getResetOriginCallCount(), "resetOrigin"},
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
	resetOrigin, getResetOriginCallCount := getCounter()
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
		resetOrigin,
		beginRender,
		swapInRender,
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
		{0, getBeforeStartCallCount(), "beforeStart"},
		{1, getRefreshViewCallCount(), "refreshView"},
		{1, getOnEndOfInputCallCount(), "onEndOfInput"},
		{0, getResetOriginCallCount(), "resetOrigin"},
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

	assert.True(t, thenCalled)
}

// A task rendering content the view wasn't already showing resets the scroll
// position to the top, at its first paint. If it is stopped and replaced before
// it ever paints — a background refresh landing just after the user clicked a
// different item, say — the replacement renders the same content and so decides
// on no reset of its own; it has to perform the one the stopped task was owed,
// or the view keeps the scroll position of the content it showed before.
func TestResetOriginSurvivesTaskReplacement(t *testing.T) {
	resetOrigin, getResetOriginCallCount := getCounter()

	manager := NewViewBufferManager(
		utils.NewDummyLog(),
		bytes.NewBuffer(nil),
		func() {},
		func() {},
		func() {},
		resetOrigin,
		func() {},
		func() {},
		func() gocui.Task { return gocui.NewFakeTask() },
		// no UI thread in the test; run the view mutations inline
		func(f func()) error { f(); return nil },
	)

	startTask := func(key string, reader io.Reader, onDone func()) {
		start := func() (Cmd, io.Reader) {
			// not actually starting this because it's not necessary
			return ExecCmd{Cmd: exec.Command("blah")}, reader
		}
		// The first-paint point is far beyond what any of these readers yield, so
		// only reaching EOF paints.
		_ = manager.NewTask(manager.NewCmdTask(start, "", LinesToRead{100, 50, nil}, onDone), key)
	}
	runTaskToCompletion := func(key string) {
		done := make(chan struct{})
		startTask(key, &BlankLineReader{totalLinesToYield: 3}, func() { close(done) })
		<-done
	}

	// A render of content the view wasn't showing resets the scroll position.
	runTaskToCompletion("cmd1")
	assert.Equal(t, 1, getResetOriginCallCount())

	// Different content again, but this task stalls before it can paint.
	stalled := BlockingLineReader{
		linesToYield: 3,
		blocked:      make(chan struct{}),
		unblock:      make(chan struct{}),
	}
	defer close(stalled.unblock)
	startTask("cmd2", &stalled, nil)
	<-stalled.blocked

	// The replacement shows the same content as the stalled task, so it has no
	// reset of its own to do — but it must still do that task's.
	runTaskToCompletion("cmd2")
	assert.Equal(t, 2, getResetOriginCallCount())
}

// A render that takes long enough to start takes the view over to say
// "loading...", which means blanking whatever it was showing. That is only worth
// doing when the content coming is different from what's on screen: re-rendering
// the same content (a background refresh, say) would otherwise blank the view and
// paint the same thing back, a visible flicker for nothing.
func TestLoadingIndicatorOnlyTakesOverForNewContent(t *testing.T) {
	var beforeStartCount atomic.Int32

	manager := NewViewBufferManager(
		utils.NewDummyLog(),
		io.Discard,
		func() { beforeStartCount.Add(1) },
		func() {},
		func() {},
		func() {},
		func() {},
		func() {},
		func() gocui.Task { return gocui.NewFakeTask() },
		// no UI thread in the test; run the view mutations inline
		func(f func()) error { f(); return nil },
	)

	startTask := func(key string, reader io.Reader, onDone func()) {
		start := func() (Cmd, io.Reader) {
			// not actually starting this because it's not necessary
			return ExecCmd{Cmd: exec.Command("blah")}, reader
		}
		_ = manager.NewTask(manager.NewCmdTask(start, "", LinesToRead{100, 50, nil}, onDone), key)
	}
	// Starts a task whose command produces nothing at all, so that it is still
	// waiting for its first line when the loading indicator falls due. Returns
	// the reader so the caller can let it finish.
	startStalledTask := func(key string) *BlockingLineReader {
		reader := &BlockingLineReader{
			blocked: make(chan struct{}),
			unblock: make(chan struct{}),
		}
		startTask(key, reader, nil)
		<-reader.blocked
		return reader
	}

	// Get some content on screen first: the indicator is only due when a render
	// is slow, and this one isn't.
	done := make(chan struct{})
	startTask("cmd1", &BlankLineReader{totalLinesToYield: 3}, func() { close(done) })
	<-done
	assert.EqualValues(t, 0, beforeStartCount.Load())

	// A slow re-render of that same content must leave the view alone however
	// long it takes. The indicator is due 200ms in, so give it well past that.
	sameContent := startStalledTask("cmd1")
	defer close(sameContent.unblock)
	time.Sleep(500 * time.Millisecond)
	assert.EqualValues(t, 0, beforeStartCount.Load())

	// Different content, though, is worth taking the view over for.
	newContent := startStalledTask("cmd2")
	defer close(newContent.unblock)
	assert.Eventually(t,
		func() bool { return beforeStartCount.Load() == 1 },
		2*time.Second, 10*time.Millisecond)
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
