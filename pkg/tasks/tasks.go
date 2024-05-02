package tasks

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
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

	// beforeStart is the function that is called before starting a new task
	beforeStart  func()
	refreshView  func()
	onEndOfInput func()

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
}

func (m *ViewBufferManager) GetTaskKey() string {
	return m.taskKey
}

func NewViewBufferManager(
	log *logrus.Entry,
	writer io.Writer,
	beforeStart func(),
	refreshView func(),
	onEndOfInput func(),
	onNewKey func(),
	newGocuiTask func() gocui.Task,
) *ViewBufferManager {
	return &ViewBufferManager{
		Log:          log,
		writer:       writer,
		beforeStart:  beforeStart,
		refreshView:  refreshView,
		onEndOfInput: onEndOfInput,
		readLines:    make(chan LinesToRead, 1024),
		onNewKey:     onNewKey,
		newGocuiTask: newGocuiTask,
	}
}

func (self *ViewBufferManager) ReadLines(n int) {
	go utils.Safe(func() {
		self.readLines <- LinesToRead{Total: n, InitialRefreshAfter: -1}
	})
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

		go utils.Safe(func() {
			<-opts.Stop
			// we use the time it took to start the program as a way of checking if things
			// are running slow at the moment. This is admittedly a crude estimate, but
			// the point is that we only want to throttle when things are running slow
			// and the user is flicking through a bunch of items.
			self.throttle = time.Since(startTime) < THROTTLE_TIME && timeToStart > COMMAND_START_THRESHOLD
			if err := oscommands.Kill(cmd); err != nil {
				if !strings.Contains(err.Error(), "process already finished") {
					self.Log.Errorf("error when running cmd task: %v", err)
				}
			}

			// for pty's we need to call onDone here so that cmd.Wait() doesn't block forever
			onDone()
		})

		loadingMutex := deadlock.Mutex{}

		// not sure if it's the right move to redefine this or not
		self.readLines = make(chan LinesToRead, 1024)

		done := make(chan struct{})

		scanner := bufio.NewScanner(r)
		scanner.Split(bufio.ScanLines)

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

		outer:
			for {
				select {
				case <-opts.Stop:
					break outer
				case linesToRead := <-self.readLines:
					for i := 0; i < linesToRead.Total; i++ {
						var ok bool
						var line []byte
						select {
						case <-opts.Stop:
							break outer
						case line, ok = <-lineChan:
							break
						}

						loadingMutex.Lock()
						if !loaded {
							self.beforeStart()
							if prefix != "" {
								writeToView([]byte(prefix))
							}
							loaded = true
						}
						loadingMutex.Unlock()

						if !ok {
							// if we're here then there's nothing left to scan from the source
							// so we're at the EOF and can flush the stale content
							self.onEndOfInput()
							break outer
						}
						writeToView(append(line, '\n'))
						lineWrittenChan <- struct{}{}

						if i+1 == linesToRead.InitialRefreshAfter {
							// We have read enough lines to fill the view, so do a first refresh
							// here to show what we have. Continue reading and refresh again at
							// the end to make sure the scrollbar has the right size.
							refreshViewIfStale()
						}
					}
					refreshViewIfStale()
					onFirstPageShown()
				}
			}

			refreshViewIfStale()

			if err := cmd.Wait(); err != nil {
				// it's fine if we've killed this program ourselves
				if !strings.Contains(err.Error(), "signal: killed") {
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

		if self.GetTaskKey() != key && self.onNewKey != nil {
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
