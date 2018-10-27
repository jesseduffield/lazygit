package procs

import (
	"fmt"
	"sync"
)

// Manager manages a set of Processes.
type Manager struct {
	Processes map[string]*Process

	lock sync.Mutex
}

// NewManager creates a new *Manager.
func NewManager() *Manager {
	return &Manager{
		Processes: make(map[string]*Process),
	}

}

// StdoutHandler returns an OutHandler that will ensure the underlying
// process has an empty stdout buffer and logs to stdout a prefixed value
// of "$name | $line".
func (m *Manager) StdoutHandler(name string) OutHandler {
	return func(line string) string {
		fmt.Printf("%s | %s\n", name, line)
		return ""
	}
}

// StderrHandler returns an OutHandler that will ensure the underlying
// process has an empty stderr buffer and logs to stdout a prefixed value
// of "$name | $line".
func (m *Manager) StderrHandler(name string) OutHandler {
	return func(line string) string {
		fmt.Printf("%s | %s\n", name, line)
		return ""
	}
}

// Start and managed a new process using the default handlers from a
// string.
func (m *Manager) Start(name, cmd string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	p := NewProcess(cmd)
	p.OutputHandler = m.StdoutHandler(name)
	p.ErrHandler = m.StderrHandler(name)
	err := p.Start()
	if err != nil {
		return err
	}

	m.Processes[name] = p
	return nil
}

// StartProcess starts and manages a predifined process.
func (m *Manager) StartProcess(name string, p *Process) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	err := p.Start()
	if err != nil {
		return err
	}

	m.Processes[name] = p
	return nil
}

// Stop will try to stop a managed process. If the process does not
// exist, no error is returned.
func (m *Manager) Stop(name string) error {
	p, ok := m.Processes[name]
	// We don't mind stopping a process that doesn't exist.
	if !ok {
		return nil
	}

	return p.Stop()
}

// Remove will try to stop and remove a managed process.
func (m *Manager) Remove(name string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	err := m.Stop(name)
	if err != nil {
		return err
	}

	// Note that if the stop fails we don't remove it from the map of
	// processes to avoid losing the reference.
	delete(m.Processes, name)

	return nil
}

// Wait will block until all managed processes have finished.
func (m *Manager) Wait() error {
	wg := &sync.WaitGroup{}
	wg.Add(len(m.Processes))

	for _, p := range m.Processes {
		go func(proc *Process) {
			defer wg.Done()
			proc.Wait()
		}(p)
	}

	wg.Wait()

	return nil
}
