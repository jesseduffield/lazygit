package gui

import (
	"sync"
	"time"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

type appStatus struct {
	message    string
	statusType string
	id         int
}

type statusManager struct {
	statuses []appStatus
	nextId   int
	mutex    sync.Mutex
}

func (m *statusManager) removeStatus(id int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	newStatuses := []appStatus{}
	for _, status := range m.statuses {
		if status.id != id {
			newStatuses = append(newStatuses, status)
		}
	}
	m.statuses = newStatuses
}

func (m *statusManager) addWaitingStatus(message string) int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.nextId += 1
	id := m.nextId

	newStatus := appStatus{
		message:    message,
		statusType: "waiting",
		id:         id,
	}
	m.statuses = append([]appStatus{newStatus}, m.statuses...)

	return id
}

func (m *statusManager) addToastStatus(message string) int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.nextId++
	id := m.nextId

	newStatus := appStatus{
		message:    message,
		statusType: "toast",
		id:         id,
	}
	m.statuses = append([]appStatus{newStatus}, m.statuses...)

	go func() {
		time.Sleep(time.Second * 2)

		m.removeStatus(id)
	}()

	return id
}

func (m *statusManager) getStatusString() string {
	if len(m.statuses) == 0 {
		return ""
	}
	topStatus := m.statuses[0]
	if topStatus.statusType == "waiting" {
		return topStatus.message + " " + utils.Loader()
	}
	return topStatus.message
}

func (gui *Gui) raiseToast(message string) {
	gui.statusManager.addToastStatus(message)

	gui.renderAppStatus()
}

func (gui *Gui) renderAppStatus() {
	go utils.Safe(func() {
		ticker := time.NewTicker(time.Millisecond * 50)
		defer ticker.Stop()
		for range ticker.C {
			appStatus := gui.statusManager.getStatusString()
			gui.OnUIThread(func() error {
				return gui.renderString(gui.Views.AppStatus, appStatus)
			})

			if appStatus == "" {
				return
			}
		}
	})
}

// withWaitingStatus wraps a function and shows a waiting status while the function is still executing
func (gui *Gui) withWaitingStatus(message string, f func() error) error {
	go utils.Safe(func() {
		id := gui.statusManager.addWaitingStatus(message)

		defer func() {
			gui.statusManager.removeStatus(id)
		}()

		gui.renderAppStatus()

		if err := f(); err != nil {
			gui.OnUIThread(func() error {
				return gui.PopupHandler.Error(err)
			})
		}
	})

	return nil
}
