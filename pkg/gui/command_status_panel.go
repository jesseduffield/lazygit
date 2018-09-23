package gui

import (
	"fmt"

	"github.com/go-cmd/cmd"
	"github.com/golang-collections/collections/stack"
	"github.com/jesseduffield/gocui"
)

const (
	// cmdStatusLen is a magic number
	// that defines the max columns of CommandStatus panel
	cmdStatusLen = 12
)

// CommandStatus defines the go-cmd's cmd with a ticker
type CommandStatus struct {
	command *cmd.Cmd
	gui     *gocui.Gui
	stack   *stack.Stack
}

// NewCommandStatus new command status
func NewCommandStatus(c *cmd.Cmd, gui *gocui.Gui) *CommandStatus {
	return &CommandStatus{
		command: c,
		gui:     gui,
		stack:   stack.New(),
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
		// v.Clear()

		fmt.Fprintf(v, "%s\n", str)
		return nil

	})
	return nil
}
