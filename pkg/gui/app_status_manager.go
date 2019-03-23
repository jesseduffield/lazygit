package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type appStatus struct {
	name       string
	statusType string
	duration   int
}

type statusManager struct {
	statuses []appStatus
}

func (m *statusManager) removeStatus(name string) {
	newStatuses := []appStatus{}
	for _, status := range m.statuses {
		if status.name != name {
			newStatuses = append(newStatuses, status)
		}
	}
	m.statuses = newStatuses
}

func (m *statusManager) addWaitingStatus(name string) {
	m.removeStatus(name)
	newStatus := appStatus{
		name:       name,
		statusType: "waiting",
		duration:   0,
	}
	m.statuses = append([]appStatus{newStatus}, m.statuses...)
}

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

// WithWaitingStatus wraps a function and shows a waiting status while the function is still executing
func (gui *Gui) WithWaitingStatus(name string, f func() error) error {
	go func() {
		gui.g.Update(func(g *gocui.Gui) error {
			gui.statusManager.addWaitingStatus(name)
			return nil
		})

		defer gui.g.Update(func(g *gocui.Gui) error {
			gui.statusManager.removeStatus(name)
			return nil
		})

		if err := f(); err != nil {
			gui.g.Update(func(g *gocui.Gui) error {
				return gui.createErrorPanel(gui.g, err.Error())
			})
		}
	}()

	return nil
}
