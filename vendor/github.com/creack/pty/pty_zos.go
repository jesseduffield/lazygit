//go:build zos
// +build zos

package pty

import (
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

const (
	SYS_UNLOCKPT     = 0x37B
	SYS_GRANTPT      = 0x37A
	SYS_POSIX_OPENPT = 0xC66
	SYS_FCNTL        = 0x18C
	SYS___PTSNAME_A  = 0x718

	SETCVTON = 1

	O_NONBLOCK = 0x04

	F_SETFL       = 4
	F_CONTROL_CVT = 13
)

type f_cnvrt struct {
	Cvtcmd int32
	Pccsid int16
	Fccsid int16
}

func open() (pty, tty *os.File, err error) {
	ptmxfd, err := openpt(os.O_RDWR | syscall.O_NOCTTY)
	if err != nil {
		return nil, nil, err
	}

	// Needed for z/OS so that the characters are not garbled if ptyp* is untagged
	cvtreq := f_cnvrt{Cvtcmd: SETCVTON, Pccsid: 0, Fccsid: 1047}
	if _, err = fcntl(uintptr(ptmxfd), F_CONTROL_CVT, uintptr(unsafe.Pointer(&cvtreq))); err != nil {
		return nil, nil, err
	}

	p := os.NewFile(uintptr(ptmxfd), "/dev/ptmx")
	if p == nil {
		return nil, nil, err
	}

	// In case of error after this point, make sure we close the ptmx fd.
	defer func() {
		if err != nil {
			_ = p.Close() // Best effort.
		}
	}()

	sname, err := ptsname(ptmxfd)
	if err != nil {
		return nil, nil, err
	}

	_, err = grantpt(ptmxfd)
	if err != nil {
		return nil, nil, err
	}

	if _, err = unlockpt(ptmxfd); err != nil {
		return nil, nil, err
	}

	ptsfd, err := syscall.Open(sname, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}

	if _, err = fcntl(uintptr(ptsfd), F_CONTROL_CVT, uintptr(unsafe.Pointer(&cvtreq))); err != nil {
		return nil, nil, err
	}

	t := os.NewFile(uintptr(ptsfd), sname)
	if err != nil {
		return nil, nil, err
	}

	return p, t, nil
}

func openpt(oflag int) (fd int, err error) {
	r0, _, e1 := runtime.CallLeFuncWithErr(runtime.GetZosLibVec()+SYS_POSIX_OPENPT<<4, uintptr(oflag))
	fd = int(r0)
	if e1 != 0 {
		err = syscall.Errno(e1)
	}
	return
}

func fcntl(fd uintptr, cmd int, arg uintptr) (val int, err error) {
	r0, _, e1 := runtime.CallLeFuncWithErr(runtime.GetZosLibVec()+SYS_FCNTL<<4, uintptr(fd), uintptr(cmd), arg)
	val = int(r0)
	if e1 != 0 {
		err = syscall.Errno(e1)
	}
	return
}

func ptsname(fd int) (name string, err error) {
	r0, _, e1 := runtime.CallLeFuncWithPtrReturn(runtime.GetZosLibVec()+SYS___PTSNAME_A<<4, uintptr(fd))
	name = u2s(unsafe.Pointer(r0))
	if e1 != 0 {
		err = syscall.Errno(e1)
	}
	return
}

func grantpt(fildes int) (rc int, err error) {
	r0, _, e1 := runtime.CallLeFuncWithErr(runtime.GetZosLibVec()+SYS_GRANTPT<<4, uintptr(fildes))
	rc = int(r0)
	if e1 != 0 {
		err = syscall.Errno(e1)
	}
	return
}

func unlockpt(fildes int) (rc int, err error) {
	r0, _, e1 := runtime.CallLeFuncWithErr(runtime.GetZosLibVec()+SYS_UNLOCKPT<<4, uintptr(fildes))
	rc = int(r0)
	if e1 != 0 {
		err = syscall.Errno(e1)
	}
	return
}

func u2s(cstr unsafe.Pointer) string {
	str := (*[1024]uint8)(cstr)
	i := 0
	for str[i] != 0 {
		i++
	}
	return string(str[:i])
}
