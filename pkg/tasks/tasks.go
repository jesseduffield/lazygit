package tasks

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sasha-s/go-deadlock"
	"github.com/sirupsen/logrus"
)

// Cmd abstracts over a started external process. *exec.Cmd satisfies the bulk
// of it via ExecCmd, but pty implementations can supply their own types — on
// Windows, ConPTY has to spawn via CreateProcess directly and can't use
// *exec.Cmd (see golang/go#62708).
type Cmd interface {
	Wait() error
	String() string
	// Terminate makes the process stop early, as gracefully as the platform
	// allows. It doesn't wait for the process to exit.
	Terminate() error
}

// ExecCmd adapts *exec.Cmd to Cmd.
type ExecCmd struct {
	*exec.Cmd
}

// Terminate sends SIGTERM on Unix. On Windows it does nothing, so a stopped
// command keeps running until it next writes to its (by then closed) output
// pipe.
func (c ExecCmd) Terminate() error {
	return oscommands.TerminateProcessGracefully(c.Process)
}

// This file revolves around running commands that will be output to the main panel
// in the gui. If we're flicking through the commits panel, we want to invoke a
// `git show` command for each commit, but we don't want to read the entire output
// at once (because that would slow things down); we just want to fill the panel
// and then read more as the user scrolls down. We also want to ensure that we're only
// ever running one `git show` command at time, and that we only have one command
// writing its output to the main panel at a time.

const THROTTLE_TIME = time.Millisecond * 30

// we use this to check if the system is under stress right now. Hopefully this makes sense on other machines
const COMMAND_START_THRESHOLD = time.Millisecond * 10

type ViewBufferManager struct {
	// this blocks until the task has been properly stopped
	stopCurrentTask func()

	// this is what we write the output of the task to. It's typically a view
	writer io.Writer

	waitingMutex deadlock.Mutex
	// Guards newTaskID and taskKey, which identify the most recently requested
	// task. Both are written on the goroutine NewTask spawns, and taskKey is
	// read from the UI thread (GetTaskKey), so neither may be touched without
	// holding this.
	taskIDMutex deadlock.Mutex
	Log         *logrus.Entry
	newTaskID   int
	// The channel by which the currently-running task is told to read more
	// lines (e.g. as the user scrolls). Held in an atomic because it's swapped
	// out as tasks come and go while ReadLines/ReadToEnd read it from the UI
	// thread; nil when no task is running.
	readLines atomic.Pointer[chan LinesToRead]
	taskKey   string

	// Resets the view's scroll position to the top. A render whose content is
	// different from what the view last showed (a different command key) calls
	// this — but at its *first paint*, not when the task starts: the off-screen
	// render leaves the previous content displayed until the swap, so resetting
	// the origin up front would scroll that still-displayed content to the top
	// before the new content replaces it. See newContentPending.
	resetOrigin func()

	// Whether the content the running task is rendering differs from what the
	// view is currently showing (i.e. the command key changed). Two things key
	// off it: the loading indicator only takes the view over when it is set,
	// since there is no point clearing content we are about to render
	// identically; and the first paint that reveals the content resets the
	// scroll to the top and clears it.
	//
	// It deliberately outlives the task that set it: a task can be stopped and
	// replaced before it ever paints — a background refresh landing just after
	// the user clicked a different item, say — and the replacement, which
	// renders the same content and so sets nothing of its own, still has to do
	// what that task was owed.
	newContentPending atomic.Bool

	// When set, the next command task puts the view back where it was once it has
	// re-rendered the content, instead of showing the new render from the top (see
	// RenderRestore). It is installed just before the re-render is triggered.
	//
	// Like newContentPending it outlives the task it was installed for, and for the
	// same reason: that task can be stopped and replaced before it ever paints, and
	// the replacement, rendering the same content, is then the one that owes the
	// user their position. It is cleared by whichever task applies it. Guarded by
	// taskIDMutex, like the task key.
	restoreForNextTask *RenderRestore

	// When set, the next command task leaves the view's scroll position alone even
	// though it renders a different command's output, that output being the same
	// content laid out differently (see SetKeepScrollPositionForNextTask). The task
	// that starts consumes it, in place of noting that new content is on its way.
	// Guarded by taskIDMutex, like the task key.
	keepScrollForNextTask bool

	// Whether a command task is currently reading content into the view. While
	// this is true the content is still growing, so callers (e.g. the layout)
	// must not clamp the view's scroll position to the amount loaded so far.
	loading atomic.Bool

	// beforeStart is the function that is called before starting a new task
	beforeStart  func()
	refreshView  func()
	onEndOfInput func()

	// beginRender starts an off-screen render: the new content is built without
	// disturbing what's displayed. swapInRender then promotes it to the display
	// in one step. Together they keep the view showing the previous render until
	// the new one has read enough to paint, instead of revealing it line by line.
	beginRender  func()
	swapInRender func()

	// see docs/dev/Busy.md
	// A gocui task is not the same thing as the tasks defined in this file.
	// A gocui task simply represents the fact that lazygit is busy doing something,
	// whereas the tasks in this file are about rendering content to a view.
	newGocuiTask func() gocui.Task

	// Runs f on the UI thread and blocks until it has completed. All mutations
	// of the view happen through this, so that the view is only ever touched on
	// the UI thread (where it is also laid out and drawn), never on the task's
	// own goroutine.
	onUIThread func(f func()) error

	// if the user flicks through a heap of items, with each one
	// spawning a process to render something to the main view,
	// it can slow things down quite a bit. In these situations we
	// want to throttle the spawning of processes. Atomic because it's set
	// from one task's stop goroutine and read when the next task starts.
	throttle atomic.Bool
}

type LinesToRead struct {
	// The total number of lines the task should have read once this request is
	// satisfied. This is an absolute count from the start of the task, not a
	// delta: the task keeps track of how many lines it has already read and only
	// reads the shortfall, so a request for a total at or below what has already
	// been read reads nothing. -1 means read all the way to the end.
	Total int

	// Number of lines after which we have read enough to fill the view, and can
	// do an initial refresh. Only set for the initial read request; -1 for
	// subsequent requests.
	InitialRefreshAfter int

	// Function to call after reading the lines is done
	Then func()
}

// RenderRestore puts a view back where it was when it re-renders content the user
// is already looking at, laid out differently — a different context size, whitespace
// ignored, another diff renderer — instead of showing the new render from the top.
//
// The task reads the new content into an off-screen buffer; the restore says when
// enough of it has arrived to show the remembered position (FirstPaintReady), and
// then finds that position and reveals it (Apply). It is a pair of callbacks rather
// than a scroll position because a different layout of the same content puts the
// remembered line somewhere else, which only looking at the new content can answer.
type RenderRestore struct {
	// FirstPaintReady reports whether enough of the new content has been read for
	// the restore to show what it is looking for. It is consulted after each line
	// is read, on the task's own goroutine.
	FirstPaintReady func() bool

	// Apply runs once, on the UI thread, at the first paint. It finds its target in
	// the off-screen content, calls swapIn to promote that content to the display,
	// and places the view on the target — in that order, so that the search runs
	// while the previous content is still displayed, and the new content is never
	// drawn at the previous render's scroll position.
	//
	// It must call swapIn even when it finds nothing to place the view on, in which
	// case the view keeps the position the paint gave it: the offset it had, or the
	// top for content the view hasn't seen.
	Apply func(swapIn func())

	// Done is called once the restore has had its render — after Apply, or when it
	// is given up because the view is being shown something other than a re-render
	// of what it was remembered from. It is how a caller that has to wait for the
	// view to be back where it belongs knows that it either is, or never will be.
	// Optional, and called on the UI thread, as Apply is.
	Done func()
}

// resolved reports that this restore's render has happened, or that there will not be
// one. Called on the UI thread, from wherever the restore ends: once, whichever way it
// ended.
func (self *RenderRestore) resolved() {
	if self.Done != nil {
		done := self.Done
		self.Done = nil
		done()
	}
}

// SetRestoreForNextTask arranges for the next command task to put the view back
// where it is now once it has re-rendered. Call it right before triggering a
// re-render of the content the view is showing; see RenderRestore.
func (self *ViewBufferManager) SetRestoreForNextTask(restore *RenderRestore) {
	self.taskIDMutex.Lock()
	defer self.taskIDMutex.Unlock()

	self.restoreForNextTask = restore
}

// HasRestoreForNextTask reports whether the next command task already has a position
// waiting to be put back, for a caller that would otherwise install one of its own
// over it.
func (self *ViewBufferManager) HasRestoreForNextTask() bool {
	self.taskIDMutex.Lock()
	defer self.taskIDMutex.Unlock()

	return self.restoreForNextTask != nil
}

func (self *ViewBufferManager) getRestoreForNextTask() *RenderRestore {
	self.taskIDMutex.Lock()
	defer self.taskIDMutex.Unlock()

	return self.restoreForNextTask
}

// SetKeepScrollPositionForNextTask arranges for the next command task to leave the
// view's scroll position alone, rather than showing its content from the top the way a
// render of different content does. Call it right before triggering a re-render of the
// content the view is showing, when the command producing it is not the one that
// produced what is on screen — a different context size, another diff renderer.
//
// It is the coarser sibling of SetRestoreForNextTask, for the same moment: the restore
// puts the view back on the line it remembers, which is only possible where the lines
// of the new rendering can be told apart; this says merely "the content is a
// rearrangement of what is there, so the offset into it is nearer to where the user was
// than the top is". Both can be set at once, and then the restore has the first say.
func (self *ViewBufferManager) SetKeepScrollPositionForNextTask() {
	self.taskIDMutex.Lock()
	defer self.taskIDMutex.Unlock()

	self.keepScrollForNextTask = true
}

// clearRestore drops a restore once a task has applied it, so that it rides exactly
// one re-render. One installed since — the user pressing the key again while this
// task was still reading — is left alone: it belongs to the render on its way.
func (self *ViewBufferManager) clearRestore(restore *RenderRestore) {
	self.taskIDMutex.Lock()
	defer self.taskIDMutex.Unlock()

	if self.restoreForNextTask == restore {
		self.restoreForNextTask = nil
	}
}

// DropRestoreForNextTask gives up a restore that has no render to ride, because the
// view is being given something other than a re-render of the content it was
// remembered from — a message where a diff was. Without this the restore would sit
// there and claim some later render of that view, putting the user somewhere they
// haven't been for a while.
func (self *ViewBufferManager) DropRestoreForNextTask() {
	self.taskIDMutex.Lock()
	restore := self.restoreForNextTask
	self.restoreForNextTask = nil
	self.taskIDMutex.Unlock()

	if restore != nil {
		restore.resolved()
	}
}

func (self *ViewBufferManager) GetTaskKey() string {
	self.taskIDMutex.Lock()
	defer self.taskIDMutex.Unlock()

	return self.taskKey
}

// ForgetRenderedContent records that the view no longer shows the render whose key it
// is holding, because it has been emptied. The key says what the view is showing, and
// is what the next task is compared against to know whether it is rendering something
// new; a view with nothing in it is showing nothing, so whatever comes next is.
func (self *ViewBufferManager) ForgetRenderedContent() {
	self.taskIDMutex.Lock()
	defer self.taskIDMutex.Unlock()

	self.taskKey = ""
}

func NewViewBufferManager(
	log *logrus.Entry,
	writer io.Writer,
	beforeStart func(),
	refreshView func(),
	onEndOfInput func(),
	resetOrigin func(),
	beginRender func(),
	swapInRender func(),
	newGocuiTask func() gocui.Task,
	onUIThread func(f func()) error,
) *ViewBufferManager {
	return &ViewBufferManager{
		Log:          log,
		writer:       writer,
		beforeStart:  beforeStart,
		refreshView:  refreshView,
		onEndOfInput: onEndOfInput,
		resetOrigin:  resetOrigin,
		beginRender:  beginRender,
		swapInRender: swapInRender,
		newGocuiTask: newGocuiTask,
		onUIThread:   onUIThread,
	}
}

// ReadLines asks the task to ensure it has read at least totalLines lines in
// total. Because the count is absolute rather than a delta, repeated requests
// (e.g. as the user scrolls down, back up, and down again) don't re-read lines
// that have already been read: the task only ever reads the shortfall.
func (self *ViewBufferManager) ReadLines(totalLines int) {
	if ch := self.readLines.Load(); ch != nil {
		readLines := *ch
		go utils.Safe(func() {
			readLines <- LinesToRead{Total: totalLines, InitialRefreshAfter: -1}
		})
	}
}

// IsLoading reports whether a command task is currently reading content into the
// view, meaning the content is still growing.
func (self *ViewBufferManager) IsLoading() bool {
	return self.loading.Load()
}

// StartLoading marks the view as loading content. It must be called
// synchronously when a command/pty task is started, before the task's goroutine
// runs, so that a layout pass happening in between doesn't clamp the scroll
// position to the not-yet-loaded content. It is cleared when the task reaches
// the end of its input.
func (self *ViewBufferManager) StartLoading() {
	self.loading.Store(true)
}

func (self *ViewBufferManager) ReadToEnd(then func()) {
	if ch := self.readLines.Load(); ch != nil {
		readLines := *ch
		go utils.Safe(func() {
			readLines <- LinesToRead{Total: -1, InitialRefreshAfter: -1, Then: then}
		})
	} else if then != nil {
		then()
	}
}

func (self *ViewBufferManager) NewCmdTask(start func() (Cmd, io.Reader), prefix string, linesToRead LinesToRead, onDoneFn func()) func(TaskOpts) error {
	return func(opts TaskOpts) error {
		var onDoneOnce sync.Once
		var onFirstPageShownOnce sync.Once

		onFirstPageShown := func() {
			onFirstPageShownOnce.Do(func() {
				opts.InitialContentLoaded()
			})
		}

		onDone := func() {
			if onDoneFn != nil {
				onDoneOnce.Do(onDoneFn)
			}
			onFirstPageShown()
		}

		// Whatever position is owed to the user belongs to this render: it was
		// remembered just before the re-render that led here was triggered.
		restore := self.getRestoreForNextTask()

		if self.throttle.Load() {
			self.Log.Info("throttling task")
			time.Sleep(THROTTLE_TIME)
		}

		select {
		case <-opts.Stop:
			onDone()
			return nil
		default:
		}

		startTime := time.Now()
		cmd, r := start()
		timeToStart := time.Since(startTime)

		done := make(chan struct{})

		go utils.Safe(func() {
			select {
			case <-done:
				// The command finished and did not have to be preemptively stopped before the next command.
				// No need to throttle.
				self.throttle.Store(false)
			case <-opts.Stop:
				// we use the time it took to start the program as a way of checking if things
				// are running slow at the moment. This is admittedly a crude estimate, but
				// the point is that we only want to throttle when things are running slow
				// and the user is flicking through a bunch of items.
				self.throttle.Store(time.Since(startTime) < THROTTLE_TIME && timeToStart > COMMAND_START_THRESHOLD)

				// Kill the still-running command. The only reason to do this is to save CPU usage
				// when flicking through several very long diffs when diff.algorithm = histogram is
				// being used, in which case multiple git processes continue to calculate expensive
				// diffs in the background even though they have been stopped already.
				if err := cmd.Terminate(); err != nil {
					self.Log.Errorf("error when trying to terminate cmd task: %v; Command: %v", err, cmd.String())
				}

				// close the task's stdout pipe (or the pty if we're using one) to make the command terminate
				onDone()
			}
		})

		loadingMutex := deadlock.Mutex{}

		readLines := make(chan LinesToRead, 1024)
		self.readLines.Store(&readLines)

		scanner := bufio.NewScanner(r)
		scanner.Split(utils.ScanLinesAndTruncateWhenLongerThanBuffer(bufio.MaxScanTokenSize))

		lineChan := make(chan []byte)
		lineWrittenChan := make(chan struct{})

		// We're reading from the scanner in a separate goroutine because on windows
		// if running git through a shim, we sometimes kill the parent process without
		// killing its children, meaning the scanner blocks forever. This solution
		// leaves us with a dead goroutine, but it's better than blocking all
		// rendering to main views.
		go utils.Safe(func() {
			defer close(lineChan)
			for scanner.Scan() {
				select {
				case <-opts.Stop:
					return
				case lineChan <- scanner.Bytes():
					// We need to confirm the data has been fed into the view before we
					// pull more from the scanner because the scanner uses the same backing
					// array and we don't want to be mutating that while it's being written
					<-lineWrittenChan
				}
			}

			if err := scanner.Err(); err != nil {
				self.Log.Error(err)
			}
		})

		loaded := false

		go utils.Safe(func() {
			ticker := time.NewTicker(time.Millisecond * 200)
			defer ticker.Stop()
			select {
			case <-opts.Stop:
				return
			case <-ticker.C:
				loadingMutex.Lock()
				// Only take the view over to say "loading..." when the content coming
				// is different from what's on screen. A re-render of the same content
				// leaves the view showing exactly what it should already, so clearing
				// it for the message and then rendering the same thing back is a
				// visible flicker for nothing — and a slow re-render of unchanged
				// content is common (a background refresh over a repo with submodules
				// that have uncommitted changes, say). The pending flag isn't consumed
				// here; the first paint still owes the scroll reset.
				//
				// A restore keeps the view too: it is there to make a re-render of what
				// the user is looking at seamless, and blanking the view for a message
				// before putting them back where they were is the flicker it exists to
				// avoid.
				if !loaded && restore == nil && self.newContentPending.Load() {
					self.beforeStart()
					// beforeStart cleared the previous content to show "loading...", so
					// put the view back at the top for it (beforeStart doesn't touch the
					// origin). The origin is view state the UI thread reads while laying
					// out, so write it there.
					_ = self.onUIThread(self.resetOrigin)
					_, _ = self.writer.Write([]byte("loading..."))
					self.refreshView()
				}
				loadingMutex.Unlock()
			}
		})

		go utils.Safe(func() {
			isViewStale := true
			writeToView := func(content []byte) {
				isViewStale = true
				_, _ = self.writer.Write(content)
			}
			refreshViewIfStale := func() {
				if isViewStale {
					self.refreshView()
					isViewStale = false
				}
			}

			// Go's select picks randomly among ready cases, so once opts.Stop is
			// closed the selects below could still service a ready data channel
			// instead of bailing. Check stop explicitly first to give it priority:
			// a task that's been stopped (it's being replaced by a newer one) must
			// not touch the view here — it would start an off-screen render and
			// write the prefix into it, clobbering what the incoming task is about
			// to render.
			stopped := func() bool {
				select {
				case <-opts.Stop:
					return true
				default:
					return false
				}
			}

			// The total number of lines we have read so far. Requests specify an
			// absolute target total (see LinesToRead.Total), so we compare against
			// this to work out how many more lines, if any, we still need to read.
			linesRead := 0

			// The first paint swaps the off-screen render in to reveal the new
			// content, and settles the scroll position in the same step — so the new
			// content first appears already where it belongs, and no draw can land
			// between the two and show it at the previous render's scroll. It happens
			// once, either when we've read far enough (below) or at end of input for
			// content shorter than that. Callers run it on the UI thread: it writes
			// the view's origin.
			painted := false
			firstPaint := func() {
				if painted {
					return
				}
				painted = true
				// Content the view hasn't seen is shown from the top, and this is where
				// the view goes there — before the restore below, which decides where to
				// put the view from where it is. The position the paint settles on is the
				// restore's to move from, so it has to be the one the new content is
				// about to be revealed at.
				if self.newContentPending.Swap(false) {
					self.resetOrigin()
				}
				if restore != nil {
					// The restore does the swap itself, so that it can find where the
					// user was in the new content before it is revealed.
					restore.Apply(self.swapInRender)
					self.clearRestore(restore)
					restore.resolved()
					return
				}
				self.swapInRender()
			}

			// Set LAZYGIT_SLOW_RENDER=<milliseconds> to sleep that long after each
			// line is written to the view, stretching async loads out so the frames
			// of a re-render become visible. Useful for debugging scroll/flicker
			// behaviour; has no effect when the variable is unset.
			var slowRenderPerLine time.Duration
			if v := os.Getenv("LAZYGIT_SLOW_RENDER"); v != "" {
				if ms, err := strconv.Atoi(v); err == nil {
					slowRenderPerLine = time.Duration(ms) * time.Millisecond
				}
			}

		outer:
			for {
				if stopped() {
					break outer
				}
				select {
				case <-opts.Stop:
					break outer
				case linesToRead := <-readLines:
					callThen := func() {
						if linesToRead.Then != nil {
							linesToRead.Then()
						}
					}
					// A restore that hasn't painted yet keeps us reading past the lines
					// asked for, all the way to the end of the input if need be. What it
					// is looking for may be anywhere in the new content, and a rendering
					// that has to be parsed as a diff to be searched at all can only be
					// parsed whole — so stopping early would leave it nothing to find,
					// and the view somewhere the user didn't put it.
					for linesToRead.Total == -1 || linesRead < linesToRead.Total || (restore != nil && !painted) {
						if stopped() {
							callThen()
							break outer
						}
						var ok bool
						var line []byte
						select {
						case <-opts.Stop:
							callThen()
							break outer
						case line, ok = <-lineChan:
							// process line below
						}

						loadingMutex.Lock()
						if !loaded {
							// Build the new content off-screen, leaving the previous render
							// displayed until we swap in below; this is what keeps an async
							// re-render from showing a half-loaded buffer.
							self.beginRender()
							if prefix != "" {
								writeToView([]byte(prefix))
							}
							loaded = true
						}
						loadingMutex.Unlock()

						if !ok {
							// lineChan is closed. At a genuine end of input we swap in what we
							// read and finalize. But lineChan is also closed when this task has
							// been stopped to make way for a newer one: stopping closes
							// opts.Stop, and the scanner goroutine then closes lineChan, so the
							// select above can land here instead of on the opts.Stop case. A
							// stopped task is being replaced and must leave the view to the
							// incoming task — swapping in its half-read buffer, clamping the
							// origin, or clearing `loading` would all corrupt what that task is
							// about to render. So bail out here, the same as the explicit stop
							// case above.
							select {
							case <-opts.Stop:
								callThen()
								break outer
							default:
							}
							// Genuine end of input: do the first paint now if it hasn't happened
							// yet (the content was shorter than a screenful, so we never reached
							// the point below), and flush the stale content. onEndOfInput reads
							// the view's dimensions (to decide whether to scroll) and sets the
							// origin, both of which are UI-thread-only, so run it there — as is
							// firstPaint, which also writes the origin.
							_ = self.onUIThread(func() {
								firstPaint()
								self.onEndOfInput()
							})
							// The content is fully loaded now, so it's safe again for the
							// layout to clamp the scroll position to it. We deliberately
							// don't clear this when stopped (rather than EOF'd), because that
							// means a newer task is taking over and is still loading.
							self.loading.Store(false)
							callThen()
							// Any read requests that were queued while we were reading are
							// now trivially satisfied, since we've read everything. Fire
							// their callbacks instead of dropping them when we break out of
							// the loop below (and nil out readLines).
						drain:
							for {
								select {
								case queued := <-readLines:
									if queued.Then != nil {
										queued.Then()
									}
								default:
									break drain
								}
							}
							break outer
						}
						writeToView(append(line, '\n'))
						lineWrittenChan <- struct{}{}
						linesRead++

						if slowRenderPerLine > 0 {
							time.Sleep(slowRenderPerLine)
						}

						if !painted {
							// Do the first paint once we have read enough lines to fill the
							// view — or, when a position is waiting to be restored, once the
							// restore says it can show it, since where the view should be is
							// its call. Continue reading afterwards and refresh again at the
							// end to make sure the scrollbar has the right size.
							var ready bool
							if restore != nil {
								ready = restore.FirstPaintReady()
							} else {
								ready = linesRead == linesToRead.InitialRefreshAfter
							}
							if ready {
								_ = self.onUIThread(firstPaint)
								refreshViewIfStale()
							}
						}
					}
					refreshViewIfStale()
					onFirstPageShown()
					callThen()
				}
			}

			self.readLines.Store(nil)

			refreshViewIfStale()

			select {
			case <-opts.Stop:
				// If we stopped the task, don't block waiting for it; this could cause a delay if
				// the process takes a while until it actually terminates. We still want to call
				// Wait to reclaim any resources, but do it on a background goroutine, and ignore
				// any errors.
				go func() { _ = cmd.Wait() }()
			default:
				if err := cmd.Wait(); err != nil {
					self.Log.Errorf("Unexpected error when running cmd task: %v; Failed command: %v", err, cmd.String())
				}
			}

			// calling this here again in case the program ended on its own accord
			onDone()

			close(done)
			close(lineWrittenChan)
		})

		readLines <- linesToRead

		<-done

		return nil
	}
}

// Close closes the task manager, killing whatever task may currently be running
func (self *ViewBufferManager) Close() {
	// stopCurrentTask is written by NewTask's goroutine under waitingMutex (and
	// so is the sync.Once it closes over), so read it under the lock and call
	// the captured value; a task starting on shutdown must not race us here.
	self.waitingMutex.Lock()
	stopCurrentTask := self.stopCurrentTask
	self.waitingMutex.Unlock()

	if stopCurrentTask == nil {
		return
	}

	c := make(chan struct{})

	go utils.Safe(func() {
		stopCurrentTask()
		c <- struct{}{}
	})

	select {
	case <-c:
		return
	case <-time.After(3 * time.Second):
		fmt.Println("cannot kill child process")
	}
}

// different kinds of tasks:
// 1) command based, where the manager can be asked to read more lines,  but the command can be killed
// 2) string based, where the manager can also be asked to read more lines

type TaskOpts struct {
	// Channel that tells the task to stop, because another task wants to run.
	Stop chan struct{}

	// Only for tasks which are long-running, where we read more lines sporadically.
	// We use this to keep track of when a user's action is complete (i.e. all views
	// have been refreshed to display the results of their action)
	InitialContentLoaded func()
}

func (self *ViewBufferManager) NewTask(f func(TaskOpts) error, key string) error {
	gocuiTask := self.newGocuiTask()

	var completeTaskOnce sync.Once

	completeGocuiTask := func() {
		completeTaskOnce.Do(func() {
			gocuiTask.Done()
		})
	}

	// Assign the taskID synchronously so it reflects NewTask call order
	// rather than the order in which the spawned goroutines happen to be
	// scheduled. Otherwise two NewTask calls in quick succession can have
	// their goroutines race, with the later-called task ending up with the
	// lower taskID and losing the staleness check below.
	self.taskIDMutex.Lock()
	self.newTaskID++
	taskID := self.newTaskID
	self.taskIDMutex.Unlock()

	go utils.Safe(func() {
		defer completeGocuiTask()

		self.taskIDMutex.Lock()

		// Bail out before touching shared view state if a newer task has
		// already been queued: if we reset the view here we'd do it for a task
		// that's about to exit, potentially wiping output the winning task has
		// already written.
		if taskID < self.newTaskID {
			self.taskIDMutex.Unlock()
			return
		}

		// Note we don't reset the origin here even when the command key changed:
		// that's deferred to the first paint that reveals the new content (see
		// newContentPending), so the previous content — left displayed until the
		// swap — doesn't visibly jump to the top before the new content appears.
		// Read taskKey directly: we already hold the mutex that guards it, and
		// GetTaskKey would take it again. A pending restore isn't dropped here
		// either, even for a different command: the re-renders it rides are all
		// different commands (a different context size, another diff renderer), and
		// it validates itself against the content it lands in anyway.
		// A task told to keep the scroll position renders the content the view is
		// already showing, laid out differently, so the reset it would otherwise owe
		// would take the user away from what they are reading — and the loading
		// message, which the same flag governs, would blank content that is about to
		// come back looking much the same.
		if self.taskKey != key && self.resetOrigin != nil && !self.keepScrollForNextTask {
			self.newContentPending.Store(true)
		}
		self.keepScrollForNextTask = false
		self.taskKey = key

		self.taskIDMutex.Unlock()

		self.waitingMutex.Lock()

		// Re-check staleness after acquiring waitingMutex: a newer task
		// may have arrived while we were blocked here.
		self.taskIDMutex.Lock()
		if taskID < self.newTaskID {
			self.waitingMutex.Unlock()
			self.taskIDMutex.Unlock()
			return
		}
		self.taskIDMutex.Unlock()

		if self.stopCurrentTask != nil {
			self.stopCurrentTask()
		}

		self.readLines.Store(nil)

		stop := make(chan struct{})
		notifyStopped := make(chan struct{})

		var once sync.Once
		onStop := func() {
			close(stop)
			<-notifyStopped
		}

		self.stopCurrentTask = func() { once.Do(onStop) }

		self.waitingMutex.Unlock()

		if err := f(TaskOpts{Stop: stop, InitialContentLoaded: completeGocuiTask}); err != nil {
			self.Log.Error(err) // might need an onError callback
		}

		close(notifyStopped)
	})

	return nil
}
