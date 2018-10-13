package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/scbizu/cmd"
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

// PrintCmdOutput collect the command output
// then print it in status panel
func (cs *CommandStatus) PrintCmdOutput(gui *Gui) {

	go func() {
		for {
			select {
			case cmdStdOutStr := <-cs.command.Stdout:
				gui.Log.Infof("collect message:%v", cmdStdOutStr)
				if err := cs.refreshCommandStatus(cs.gui, cmdStdOutStr); err != nil {
					gui.Log.Errorf("get commandStatus Panel view failed:%v", err.Error())
					return
				}
			case cmdErrOutStr := <-cs.command.Stderr:
				gui.Log.Infof("collect error message:%v", cmdErrOutStr)
				if err := cs.refreshCommandStatus(cs.gui, cmdErrOutStr); err != nil {
					gui.Log.Errorf("get commandStatus Panel view failed:%v", err.Error())
					return
				}
			}
		}
	}()

	cs.command.Start()
}

func (cs *CommandStatus) refreshCommandStatus(g *gocui.Gui, cmdStr string) error {
	v, err := g.View("commandStatus")
	if err != nil {
		return err
	}

	g.Update(func(*gocui.Gui) error {

		// DO NOT invoke v.Clear() HERE:
		// The command status history will be buffered
		// And will be finally destroyed when user quits lazygit

		fmt.Fprintf(v, "%s\n", cmdStr)
		return nil

	})
	return nil
}
