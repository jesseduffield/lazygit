// +build !windows

package gocui

import (
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"github.com/go-errors/errors"
)

type windowSize struct {
	rows    uint16
	cols    uint16
	xpixels uint16
	ypixels uint16
}

// getTermWindowSize is get terminal window size on linux or unix.
// When gocui run inside the docker contaienr need to check and get the window size.
func (g *Gui) getTermWindowSize() (int, int, error) {
	out, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return 0, 0, err
	}
	defer out.Close()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGWINCH, syscall.SIGINT)
	defer signal.Stop(signalCh)

	var sz windowSize

	for {
		_, _, err = syscall.Syscall(
			syscall.SYS_IOCTL,
			out.Fd(),
			uintptr(syscall.TIOCGWINSZ),
			uintptr(unsafe.Pointer(&sz)),
		)

		// check terminal window size
		if sz.cols > 0 && sz.rows > 0 {
			return int(sz.cols), int(sz.rows), nil
		}

		select {
		case signal := <-signalCh:
			switch signal {
			// when the terminal window size is changed
			case syscall.SIGWINCH:
				continue
			// ctrl + c to cancel
			case syscall.SIGINT:
				return 0, 0, errors.New("There was not enough window space to start the application")
			}
		}
	}
}
