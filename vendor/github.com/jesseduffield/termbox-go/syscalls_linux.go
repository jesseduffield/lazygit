// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs syscalls.go

package termbox

import "syscall"

type syscall_Termios syscall.Termios

const (
	syscall_IGNBRK = syscall.IGNBRK
	syscall_BRKINT = syscall.BRKINT
	syscall_PARMRK = syscall.PARMRK
	syscall_ISTRIP = syscall.ISTRIP
	syscall_INLCR  = syscall.INLCR
	syscall_IGNCR  = syscall.IGNCR
	syscall_ICRNL  = syscall.ICRNL
	syscall_IXON   = syscall.IXON
	syscall_OPOST  = syscall.OPOST
	syscall_ECHO   = syscall.ECHO
	syscall_ECHONL = syscall.ECHONL
	syscall_ICANON = syscall.ICANON
	syscall_ISIG   = syscall.ISIG
	syscall_IEXTEN = syscall.IEXTEN
	syscall_CSIZE  = syscall.CSIZE
	syscall_PARENB = syscall.PARENB
	syscall_CS8    = syscall.CS8
	syscall_VMIN   = syscall.VMIN
	syscall_VTIME  = syscall.VTIME

	syscall_TCGETS = syscall.TCGETS
	syscall_TCSETS = syscall.TCSETS
)
