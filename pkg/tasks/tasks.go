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
	"github.com/sirupsen/logrus"
)

const THROTTLE_TIME = time.Millisecond * 30

// we use this to check if the system is under stress right now. Hopefully this makes sense on other machines
const COMMAND_START_THRESHOLD = time.Millisecond * 10

type Task struct {
	stop          chan struct{}
	stopped       bool
	stopMutex     sync.Mutex
	notifyStopped chan struct{}
	Log           *logrus.Entry
	f             func(chan struct{}) error
}

type ViewBufferManager struct {
	writer       io.Writer
	currentTask  *Task
	waitingMutex sync.Mutex
	taskIDMutex  sync.Mutex
	Log          *logrus.Entry
	newTaskId    int
	readLines    chan int
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
		readLines:    make(chan int, 1024),
		onNewKey:     onNewKey,
	}
}

func (m *ViewBufferManager) ReadLines(n int) {
	go utils.Safe(func() {
		m.readLines <- n
	})
}

func (m *ViewBufferManager) NewCmdTask(start func() (*exec.Cmd, io.Reader), prefix string, linesToRead int, onDone func()) func(chan struct{}) error {
	return func(stop chan struct{}) error {
		if m.throttle {
			m.Log.Info("throttling task")
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
			m.throttle = time.Since(startTime) < THROTTLE_TIME && timeToStart > COMMAND_START_THRESHOLD
			if err := oscommands.Kill(cmd); err != nil {
				if !strings.Contains(err.Error(), "process already finished") {
					m.Log.Errorf("error when running cmd task: %v", err)
				}
			}
		})

		loadingMutex := sync.Mutex{}

		// not sure if it's the right move to redefine this or not
		m.readLines = make(chan int, 1024)

		done := make(chan struct{})

		go utils.Safe(func() {
			scanner := bufio.NewScanner(r)
			scanner.Split(bufio.ScanLines)

			loaded := false

			go utils.Safe(func() {
				ticker := time.NewTicker(time.Millisecond * 200)
				defer ticker.Stop()
				select {
				case <-ticker.C:
					loadingMutex.Lock()
					if !loaded {
						m.beforeStart()
						_, _ = m.writer.Write([]byte("loading..."))
						m.refreshView()
					}
					loadingMutex.Unlock()
				case <-stop:
					return
				}
			})

		outer:
			for {
				select {
				case linesToRead := <-m.readLines:
					for i := 0; i < linesToRead; i++ {
						ok := scanner.Scan()
						loadingMutex.Lock()
						if !loaded {
							m.beforeStart()
							if prefix != "" {
								_, _ = m.writer.Write([]byte(prefix))
							}
							loaded = true
						}
						loadingMutex.Unlock()

						select {
						case <-stop:
							break outer
						default:
						}
						if !ok {
							// if we're here then there's nothing left to scan from the source
							// so we're at the EOF and can flush the stale content
							m.onEndOfInput()
							break outer
						}
						_, _ = m.writer.Write(append(scanner.Bytes(), '\n'))
					}
					m.refreshView()
				case <-stop:
					break outer
				}
			}

			m.refreshView()

			if err := cmd.Wait(); err != nil {
				// it's fine if we've killed this program ourselves
				if !strings.Contains(err.Error(), "signal: killed") {
					m.Log.Error(err)
				}
			}

			if onDone != nil {
				onDone()
			}

			close(done)
		})

		m.readLines <- linesToRead

		<-done

		return nil
	}
}

// Close closes the task manager, killing whatever task may currently be running
func (t *ViewBufferManager) Close() {
	if t.currentTask == nil {
		return
	}

	c := make(chan struct{})

	go utils.Safe(func() {
		t.currentTask.Stop()
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

func (m *ViewBufferManager) NewTask(f func(stop chan struct{}) error, key string) error {
	go utils.Safe(func() {
		m.taskIDMutex.Lock()
		m.newTaskId++
		taskID := m.newTaskId

		if m.GetTaskKey() != key && m.onNewKey != nil {
			m.onNewKey()
		}
		m.taskKey = key

		m.taskIDMutex.Unlock()

		m.waitingMutex.Lock()
		defer m.waitingMutex.Unlock()

		if taskID < m.newTaskId {
			return
		}

		stop := make(chan struct{})
		notifyStopped := make(chan struct{})

		if m.currentTask != nil {
			m.currentTask.Stop()
		}

		m.currentTask = &Task{
			stop:          stop,
			notifyStopped: notifyStopped,
			Log:           m.Log,
			f:             f,
		}

		go utils.Safe(func() {
			if err := f(stop); err != nil {
				m.Log.Error(err) // might need an onError callback
			}

			close(notifyStopped)
		})
	})

	return nil
}

func (t *Task) Stop() {
	t.stopMutex.Lock()
	defer t.stopMutex.Unlock()
	if t.stopped {
		return
	}
	close(t.stop)
	<-t.notifyStopped
	t.stopped = true
}
