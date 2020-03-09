package tasks

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/sirupsen/logrus"
)

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

	// beforeStart is the function that is called before starting a new task
	beforeStart func()
	refreshView func()
}

func NewViewBufferManager(log *logrus.Entry, writer io.Writer, beforeStart func(), refreshView func()) *ViewBufferManager {
	return &ViewBufferManager{Log: log, writer: writer, beforeStart: beforeStart, refreshView: refreshView, readLines: make(chan int, 1024)}
}

func (m *ViewBufferManager) ReadLines(n int) {
	go func() {
		m.readLines <- n
	}()
}

func (m *ViewBufferManager) NewCmdTask(r io.Reader, cmd *exec.Cmd, linesToRead int, onDone func()) func(chan struct{}) error {
	return func(stop chan struct{}) error {
		go func() {
			<-stop
			if err := commands.Kill(cmd); err != nil {
				m.Log.Warn(err)
			}
			if onDone != nil {
				onDone()
			}
		}()

		loadingMutex := sync.Mutex{}

		// not sure if it's the right move to redefine this or not
		m.readLines = make(chan int, 1024)

		done := make(chan struct{})

		go func() {
			scanner := bufio.NewScanner(r)
			scanner.Split(bufio.ScanLines)

			loaded := false

			go func() {
				ticker := time.NewTicker(time.Millisecond * 100)
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
			}()

		outer:
			for {
				select {
				case linesToRead := <-m.readLines:
					for i := 0; i < linesToRead; i++ {
						ok := scanner.Scan()
						loadingMutex.Lock()
						if !loaded {
							m.beforeStart()
							loaded = true
						}
						loadingMutex.Unlock()

						select {
						case <-stop:
							m.refreshView()
							break outer
						default:
						}
						if !ok {
							m.refreshView()
							break outer
						}
						_, _ = m.writer.Write(append(scanner.Bytes(), []byte("\n")...))
					}
					m.refreshView()
				case <-stop:
					m.refreshView()
					break outer
				}
			}

			if err := cmd.Wait(); err != nil {
				m.Log.Warn(err)
			}

			m.refreshView()

			if onDone != nil {
				onDone()
			}

			close(done)
		}()

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

	go func() {
		t.currentTask.Stop()
		c <- struct{}{}
	}()

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

func (m *ViewBufferManager) NewTask(f func(stop chan struct{}) error) error {
	go func() {
		m.taskIDMutex.Lock()
		m.newTaskId++
		taskID := m.newTaskId
		m.Log.Infof("starting task %d", taskID)
		m.taskIDMutex.Unlock()

		m.waitingMutex.Lock()
		defer m.waitingMutex.Unlock()

		m.Log.Infof("done waiting")
		if taskID < m.newTaskId {
			m.Log.Infof("returning cos the task is obsolete")
			return
		}

		stop := make(chan struct{})
		notifyStopped := make(chan struct{})

		if m.currentTask != nil {
			m.Log.Info("asking task to stop")
			m.currentTask.Stop()
			m.Log.Info("task stopped")
		}

		m.currentTask = &Task{
			stop:          stop,
			notifyStopped: notifyStopped,
			Log:           m.Log,
			f:             f,
		}

		go func() {
			if err := f(stop); err != nil {
				m.Log.Error(err) // might need an onError callback
			}

			m.Log.Infof("returning from task %d", taskID)
			close(notifyStopped)
		}()
	}()

	return nil
}

func (t *Task) Stop() {
	t.stopMutex.Lock()
	defer t.stopMutex.Unlock()
	if t.stopped {
		return
	}
	close(t.stop)
	t.Log.Info("closed stop channel, waiting for notifyStopped message")
	<-t.notifyStopped
	t.Log.Info("received notifystopped message")
	t.stopped = true
}
