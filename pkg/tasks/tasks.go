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
	taskIDMutex  deadlock.Mutex
	Log          *logrus.Entry
	newTaskID    int
	readLines    chan LinesToRead
	taskKey      string
	onNewKey     func()

	// When non-nil, the next cmd/pty task re-establishes the view's scroll
	// position and selection once it has re-rendered the content (returning to a
	// focused main view on escape). The task does not reset the view's origin to
	// the top at start: instead it keeps the placeholder showing at its current
	// scroll until the restore can show the saved position as part of the first
	// paint, rather than flicking to the top. See RenderRestore.
	//
	// It is set just before triggering the re-render. Unlike a per-task field, it
	// is *not* cleared when a task starts: it must survive a task being stopped
	// and replaced (a periodic refresh can stop the escape's re-render before it
	// first-paints), so it is kept until a task successfully applies it, or until
	// the content key changes (a different item was selected, so the saved
	// position no longer applies). Guarded by taskIDMutex, like the task key.
	restoreForNextTask *RenderRestore

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

	// if the user flicks through a heap of items, with each one
	// spawning a process to render something to the main view,
	// it can slow things down quite a bit. In these situations we
	// want to throttle the spawning of processes.
	throttle bool
}

type LinesToRead struct {
	// Total number of lines to read
	Total int

	// Number of lines after which we have read enough to fill the view, and can
	// do an initial refresh. Only set for the initial read request; -1 for
	// subsequent requests.
	InitialRefreshAfter int

	// When set, re-establishes the view's scroll position and selection as the
	// re-render first paints (returning to a focused main view on escape). When
	// set, it — rather than InitialRefreshAfter — also decides when the first
	// paint happens, since the saved position may be reachable only after more
	// than a screenful has loaded. Only set for the initial read request. See
	// RenderRestore.
	Restore *RenderRestore

	// Function to call after reading the lines is done
	Then func()
}

// RenderRestore re-establishes a view's scroll position and selection after it
// re-renders content the user was already looking at (returning to a focused main
// view on escape). The render task reads the new content into an off-screen
// buffer; RenderRestore decides, as that buffer fills, when the task has read far
// enough to show the saved position (FirstPaintReady), and then — once the
// off-screen buffer has been swapped in — scrolls there and restores the
// selection (Apply). It is a predicate rather than a fixed scroll position so the
// target can be a row matching a patch identity, located by scanning the loading
// content, rather than a line number that the changed content may have moved.
type RenderRestore struct {
	// FirstPaintReady reports whether the task has now read enough of the new
	// (off-screen) content to first-paint at the saved position. Evaluated after
	// each line is read.
	FirstPaintReady func() bool

	// Apply runs once, just after the off-screen render is swapped in at the first
	// paint, to scroll to the saved position and restore the selection.
	Apply func()
}

func (self *ViewBufferManager) GetTaskKey() string {
	return self.taskKey
}

func NewViewBufferManager(
	log *logrus.Entry,
	writer io.Writer,
	beforeStart func(),
	refreshView func(),
	onEndOfInput func(),
	onNewKey func(),
	beginRender func(),
	swapInRender func(),
	newGocuiTask func() gocui.Task,
) *ViewBufferManager {
	return &ViewBufferManager{
		Log:          log,
		writer:       writer,
		beforeStart:  beforeStart,
		refreshView:  refreshView,
		onEndOfInput: onEndOfInput,
		readLines:    nil,
		onNewKey:     onNewKey,
		beginRender:  beginRender,
		swapInRender: swapInRender,
		newGocuiTask: newGocuiTask,
	}
}

func (self *ViewBufferManager) ReadLines(n int) {
	if self.readLines != nil {
		go utils.Safe(func() {
			self.readLines <- LinesToRead{Total: n, InitialRefreshAfter: -1}
		})
	}
}

// SetRestoreForNextTask makes the next cmd/pty task re-establish the view's
// scroll position and selection once it has re-rendered the content. Call this
// right before triggering a re-render of content the view was already showing
// (returning to a focused main view on escape). See the field doc and RenderRestore.
func (self *ViewBufferManager) SetRestoreForNextTask(restore *RenderRestore) {
	self.taskIDMutex.Lock()
	defer self.taskIDMutex.Unlock()

	self.restoreForNextTask = restore
}

// GetRestoreForNextTask returns the pending restore, or nil. It is not gated on
// the command key: returning to a focused main view after staging changes the
// command (e.g. the unstaged diff becomes the staged one once the last unstaged
// hunk is gone), yet the line to land on is still in the new content. The restore
// validates itself instead — its scan simply doesn't find the target line when
// the content no longer contains it (a different item was selected), in which case
// applying it is a no-op. So it is safe to hand to whatever renders next.
func (self *ViewBufferManager) GetRestoreForNextTask() *RenderRestore {
	self.taskIDMutex.Lock()
	defer self.taskIDMutex.Unlock()

	return self.restoreForNextTask
}

// ClearRestoreForNextTask drops the pending restore. A task clears it once it has
// first-painted and applied it (whether or not it found its target line), so that
// it lives for exactly one re-render — surviving a task being stopped and replaced
// before it could paint, but not re-applying on every later render.
func (self *ViewBufferManager) ClearRestoreForNextTask() {
	self.taskIDMutex.Lock()
	defer self.taskIDMutex.Unlock()

	self.restoreForNextTask = nil
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
	if self.readLines != nil {
		go utils.Safe(func() {
			self.readLines <- LinesToRead{Total: -1, InitialRefreshAfter: -1, Then: then}
		})
	} else if then != nil {
		then()
	}
}

func (self *ViewBufferManager) NewCmdTask(start func() (*exec.Cmd, io.Reader), prefix string, linesToRead LinesToRead, onDoneFn func()) func(TaskOpts) error {
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

		if self.throttle {
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
				self.throttle = false
			case <-opts.Stop:
				// we use the time it took to start the program as a way of checking if things
				// are running slow at the moment. This is admittedly a crude estimate, but
				// the point is that we only want to throttle when things are running slow
				// and the user is flicking through a bunch of items.
				self.throttle = time.Since(startTime) < THROTTLE_TIME && timeToStart > COMMAND_START_THRESHOLD

				// Kill the still-running command. The only reason to do this is to save CPU usage
				// when flicking through several very long diffs when diff.algorithm = histogram is
				// being used, in which case multiple git processes continue to calculate expensive
				// diffs in the background even though they have been stopped already.
				//
				// Unfortunately this will do nothing on Windows, so Windows users will have to live
				// with the higher CPU usage.
				if err := oscommands.TerminateProcessGracefully(cmd); err != nil {
					self.Log.Errorf("error when trying to terminate cmd task: %v; Command: %v %v", err, cmd.Path, cmd.Args)
				}

				// close the task's stdout pipe (or the pty if we're using one) to make the command terminate
				onDone()
			}
		})

		loadingMutex := deadlock.Mutex{}

		self.readLines = make(chan LinesToRead, 1024)

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
				if !loaded {
					self.beforeStart()
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
			// not touch the view here — beforeStart clears it and the prefix gets
			// written, clobbering what the incoming task is about to render.
			stopped := func() bool {
				select {
				case <-opts.Stop:
					return true
				default:
					return false
				}
			}

			// The first paint swaps the off-screen render in to reveal the new
			// content. When a restore is pending (returning to a focused main view on
			// escape, see RenderRestore), it scrolls to the saved position and
			// restores the selection in the same step, so the real content first
			// appears already at the right place rather than at the top. firstPaint
			// happens once, either when we've read far enough (below) or at end of
			// input for content shorter than that.
			restore := linesToRead.Restore
			painted := false
			firstPaint := func() {
				if painted {
					return
				}
				painted = true
				self.swapInRender()
				if restore != nil {
					restore.Apply()
				}
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
				case linesToRead := <-self.readLines:
					callThen := func() {
						if linesToRead.Then != nil {
							linesToRead.Then()
						}
					}
					for i := 0; linesToRead.Total == -1 || i < linesToRead.Total; i++ {
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
							// incoming task — swapping in its half-read buffer, applying the
							// saved scroll, clamping the origin, or clearing `loading` would all
							// corrupt what that task is about to render. So bail out here, the
							// same as the explicit stop case above.
							select {
							case <-opts.Stop:
								break outer
							default:
							}
							// Genuine end of input: do the first paint now if it hasn't happened
							// yet — the content was shorter than the first-paint point, or a
							// restore's target line was never found. firstPaint swaps in whatever
							// we read and, for a restore, scrolls to the saved position before
							// onEndOfInput clamps the origin back into range for short content.
							firstPaint()
							self.onEndOfInput()
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
								case queued := <-self.readLines:
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

						if slowRenderPerLine > 0 {
							time.Sleep(slowRenderPerLine)
						}

						// Do the first paint as soon as we've read far enough: for a restore,
						// when it can show the saved position (RenderRestore.FirstPaintReady,
						// e.g. its target line plus a screenful below it have loaded); otherwise
						// when we've read enough lines to fill the view (InitialRefreshAfter).
						// This swaps the off-screen content in and refreshes; we keep reading
						// afterwards and refresh again at the end so the scrollbar ends up the
						// right size.
						if !painted {
							var ready bool
							if restore != nil {
								ready = restore.FirstPaintReady()
							} else {
								ready = linesToRead.InitialRefreshAfter > 0 && i+1 >= linesToRead.InitialRefreshAfter
							}
							if ready {
								firstPaint()
								refreshViewIfStale()
							}
						}
					}
					refreshViewIfStale()
					onFirstPageShown()
					callThen()
				}
			}

			self.readLines = nil

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
					self.Log.Errorf("Unexpected error when running cmd task: %v; Failed command: %v %v", err, cmd.Path, cmd.Args)
				}
			}

			// calling this here again in case the program ended on its own accord
			onDone()

			close(done)
			close(lineWrittenChan)
		})

		self.readLines <- linesToRead

		<-done

		return nil
	}
}

// Close closes the task manager, killing whatever task may currently be running
func (self *ViewBufferManager) Close() {
	if self.stopCurrentTask == nil {
		return
	}

	c := make(chan struct{})

	go utils.Safe(func() {
		self.stopCurrentTask()
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

	go utils.Safe(func() {
		defer completeGocuiTask()

		self.taskIDMutex.Lock()
		self.newTaskID++
		taskID := self.newTaskID

		// Reset the origin to the top when the command changed, unless a restore is
		// pending: a restore re-establishes the scroll position itself (and we keep
		// the placeholder showing at its current scroll until it does), so resetting
		// to the top would just flicker. The restore isn't cleared here: it must
		// outlive a task being stopped and replaced before it could paint, and it
		// validates itself against the content it lands in (see restoreForNextTask),
		// so a stale one can't apply to the wrong place — it's cleared once a task
		// has applied it.
		if self.GetTaskKey() != key && self.onNewKey != nil && self.restoreForNextTask == nil {
			self.onNewKey()
		}
		self.taskKey = key

		self.taskIDMutex.Unlock()

		self.waitingMutex.Lock()

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

		self.readLines = nil

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
