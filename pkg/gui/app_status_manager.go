package gui

import "github.com/jesseduffield/lazygit/pkg/utils"

// appStatus is used to store information about the status
// of the current session
type appStatus struct {
	name       string
	statusType string
	duration   int
}

// statusManager
type statusManager struct {
	statuses []appStatus
}

// removeStatus removes the status from the application.
// name: the name of the status
func (m *statusManager) removeStatus(name string) {
	for i, status := range m.statuses {
		if status.name != name {
			m.statuses = append(m.statuses[:i], m.statuses[i+1:]...)
			return
		}
	}

}

// addWaitingStatus creates a new status and adds it to the status
// managers internal array
func (m *statusManager) addWaitingStatus(name string) {
	m.removeStatus(name)

	newStatus := appStatus{
		name:       name,
		statusType: "waiting",
		duration:   0,
	}

	m.statuses = append([]appStatus{newStatus}, m.statuses...)
}

// getStatus returns the string representing the status
func (m *statusManager) getStatusString() string {
	if len(m.statuses) == 0 {
		return ""
	}

	topStatus := m.statuses[0]

	if topStatus.statusType == "waiting" {
		return topStatus.name + " " + utils.Loader()
	}

	return topStatus.name
}
