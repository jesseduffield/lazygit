package tasks

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

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
) *ViewBufferManager {
	return &ViewBufferManager{
		Log:          log,
		writer:       writer,
		beforeStart:  beforeStart,
		refreshView:  refreshView,
		onEndOfInput: onEndOfInput,
		readLines:    make(chan LinesToRead, 1024),
		onNewKey:     onNewKey,
	}
}

func (self *ViewBufferManager) ReadLines(n int) {
	go utils.Safe(func() {
		self.readLines <- LinesToRead{Total: n, InitialRefreshAfter: -1}
	})
}

// note: onDone may be called twice
func (self *ViewBufferManager) NewCmdTask(start func() (*exec.Cmd, io.Reader), prefix string, linesToRead LinesToRead, onDone func()) func(chan struct{}) error {
	return func(stop chan struct{}) error {
		var once sync.Once
		var onDoneWrapper func()
		if onDone != nil {
			onDoneWrapper = func() { once.Do(onDone) }
		}

		if self.throttle {
			self.Log.Info("throttling task")
			time.Sleep(THROTTLE_TIME)
		}

		select {
		case <-stop:
			return nil
		default:
		}

		startTime := time.Now()
		cmd, r := start()
		timeToStart := time.Since(startTime)

		go utils.Safe(func() {
			<-stop
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
			if onDoneWrapper != nil {
				onDoneWrapper()
			}
		})

		loadingMutex := deadlock.Mutex{}

		// not sure if it's the right move to redefine this or not
		self.readLines = make(chan LinesToRead, 1024)

		done := make(chan struct{})

		scanner := bufio.NewScanner(r)
		scanner.Split(bufio.ScanLines)

		loaded := false

		go utils.Safe(func() {
			ticker := time.NewTicker(time.Millisecond * 200)
			defer ticker.Stop()
			select {
			case <-stop:
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
				_, _ = self.writer.Write(content)
				isViewStale = true
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
				case <-stop:
					break outer
				case linesToRead := <-self.readLines:
					for i := 0; i < linesToRead.Total; i++ {
						select {
						case <-stop:
							break outer
						default:
						}

						ok := scanner.Scan()
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
						writeToView(append(scanner.Bytes(), '\n'))

						if i+1 == linesToRead.InitialRefreshAfter {
							// We have read enough lines to fill the view, so do a first refresh
							// here to show what we have. Continue reading and refresh again at
							// the end to make sure the scrollbar has the right size.
							refreshViewIfStale()
						}
					}
					refreshViewIfStale()
				}
			}

			refreshViewIfStale()

			if err := cmd.Wait(); err != nil {
				// it's fine if we've killed this program ourselves
				if !strings.Contains(err.Error(), "signal: killed") {
					self.Log.Errorf("Unexpected error when running cmd task: %v", err)
				}
			}

			// calling onDoneWrapper here again in case the program ended on its own accord
			if onDoneWrapper != nil {
				onDoneWrapper()
			}

			close(done)
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

func (self *ViewBufferManager) NewTask(f func(stop chan struct{}) error, key string) error {
	go utils.Safe(func() {
		self.taskIDMutex.Lock()
		self.newTaskID++
		taskID := self.newTaskID

		if self.GetTaskKey() != key && self.onNewKey != nil {
			self.onNewKey()
		}
		self.taskKey = key

		self.taskIDMutex.Unlock()

		self.waitingMutex.Lock()
		defer self.waitingMutex.Unlock()

		if taskID < self.newTaskID {
			return
		}

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

		go utils.Safe(func() {
			if err := f(stop); err != nil {
				self.Log.Error(err) // might need an onError callback
			}

			close(notifyStopped)
		})
	})

	return nil
}
