// +build ignore

package termbox

/*
#include <termios.h>
#include <sys/ioctl.h>
*/
import "C"

type syscall_Termios C.struct_termios

const (
	syscall_IGNBRK = C.IGNBRK
	syscall_BRKINT = C.BRKINT
	syscall_PARMRK = C.PARMRK
	syscall_ISTRIP = C.ISTRIP
	syscall_INLCR  = C.INLCR
	syscall_IGNCR  = C.IGNCR
	syscall_ICRNL  = C.ICRNL
	syscall_IXON   = C.IXON
	syscall_OPOST  = C.OPOST
	syscall_ECHO   = C.ECHO
	syscall_ECHONL = C.ECHONL
	syscall_ICANON = C.ICANON
	syscall_ISIG   = C.ISIG
	syscall_IEXTEN = C.IEXTEN
	syscall_CSIZE  = C.CSIZE
	syscall_PARENB = C.PARENB
	syscall_CS8    = C.CS8
	syscall_VMIN   = C.VMIN
	syscall_VTIME  = C.VTIME

	// on darwin change these to (on *bsd too?):
	// C.TIOCGETA
	// C.TIOCSETA
	syscall_TCGETS = C.TCGETS
	syscall_TCSETS = C.TCSETS
)
