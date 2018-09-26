package gui

import (
	"fmt"

	"github.com/go-cmd/cmd"
	"github.com/jesseduffield/gocui"
)

// CommandStatus defines the go-cmd's cmd
type CommandStatus struct {
	command *cmd.Cmd
	gui     *gocui.Gui
}

// NewCommandStatus new command status
func NewCommandStatus(c *cmd.Cmd, gui *gocui.Gui) *CommandStatus {
	return &CommandStatus{
		command: c,
		gui:     gui,
	}
}

// Update updates the command status panel
func (cs *CommandStatus) Update(ui *Gui) {

	go func() {
		for {
			select {
			case status := <-cs.command.Stdout:
				ui.Log.Infof("collect message:%v", status)
				if err := cs.refreshCommandStatus(cs.gui, status); err != nil {
					ui.Log.Errorf("get commandStatus Panel view failed:%v", err.Error())
					return
				}
			case status := <-cs.command.Stderr:
				ui.Log.Infof("collect error message:%v", status)
				if err := cs.refreshCommandStatus(cs.gui, status); err != nil {
					ui.Log.Errorf("get commandStatus Panel view failed:%v", err.Error())
					return
				}
			}
		}
	}()

	<-cs.command.Start()
}

func (cs *CommandStatus) refreshCommandStatus(g *gocui.Gui, str string) error {
	v, err := g.View("commandStatus")
	if err != nil {
		return err
	}

	g.Update(func(*gocui.Gui) error {

		// DO NOT invoke v.Clear() HERE:
		// The command status history will be buffered
		// And will be finally destroyed when user quits lazygit

		fmt.Fprintf(v, "%s\n", str)
		return nil

	})
	return nil
}
