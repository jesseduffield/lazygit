//go:build solaris
// +build solaris

package pty

/* based on:
http://src.illumos.org/source/xref/illumos-gate/usr/src/lib/libc/port/gen/pt.c
*/

import (
	"errors"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

func open() (pty, tty *os.File, err error) {
	ptmxfd, err := syscall.Open("/dev/ptmx", syscall.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}
	p := os.NewFile(uintptr(ptmxfd), "/dev/ptmx")
	// In case of error after this point, make sure we close the ptmx fd.
	defer func() {
		if err != nil {
			_ = p.Close() // Best effort.
		}
	}()

	sname, err := ptsname(p)
	if err != nil {
		return nil, nil, err
	}

	if err := grantpt(p); err != nil {
		return nil, nil, err
	}

	if err := unlockpt(p); err != nil {
		return nil, nil, err
	}

	ptsfd, err := syscall.Open(sname, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}
	t := os.NewFile(uintptr(ptsfd), sname)

	// In case of error after this point, make sure we close the pts fd.
	defer func() {
		if err != nil {
			_ = t.Close() // Best effort.
		}
	}()

	// pushing terminal driver STREAMS modules as per pts(7)
	for _, mod := range []string{"ptem", "ldterm", "ttcompat"} {
		if err := streamsPush(t, mod); err != nil {
			return nil, nil, err
		}
	}

	return p, t, nil
}

func ptsname(f *os.File) (string, error) {
	dev, err := ptsdev(f)
	if err != nil {
		return "", err
	}
	fn := "/dev/pts/" + strconv.FormatInt(int64(dev), 10)

	if err := syscall.Access(fn, 0); err != nil {
		return "", err
	}
	return fn, nil
}

func unlockpt(f *os.File) error {
	istr := strioctl{
		icCmd:     UNLKPT,
		icTimeout: 0,
		icLen:     0,
		icDP:      nil,
	}
	return ioctl(f, I_STR, uintptr(unsafe.Pointer(&istr)))
}

func minor(x uint64) uint64 { return x & 0377 }

func ptsdev(f *os.File) (uint64, error) {
	istr := strioctl{
		icCmd:     ISPTM,
		icTimeout: 0,
		icLen:     0,
		icDP:      nil,
	}

	if err := ioctl(f, I_STR, uintptr(unsafe.Pointer(&istr))); err != nil {
		return 0, err
	}
	var errors = make(chan error, 1)
	var results = make(chan uint64, 1)
	defer close(errors)
	defer close(results)

	var err error
	var sc syscall.RawConn
	sc, err = f.SyscallConn()
	if err != nil {
		return 0, err
	}
	err = sc.Control(func(fd uintptr) {
		var status syscall.Stat_t
		if err := syscall.Fstat(int(fd), &status); err != nil {
			results <- 0
			errors <- err
		}
		results <- uint64(minor(status.Rdev))
		errors <- nil
	})
	if err != nil {
		return 0, err
	}
	return <-results, <-errors
}

type ptOwn struct {
	rUID int32
	rGID int32
}

func grantpt(f *os.File) error {
	if _, err := ptsdev(f); err != nil {
		return err
	}
	pto := ptOwn{
		rUID: int32(os.Getuid()),
		// XXX should first attempt to get gid of DEFAULT_TTY_GROUP="tty"
		rGID: int32(os.Getgid()),
	}
	istr := strioctl{
		icCmd:     OWNERPT,
		icTimeout: 0,
		icLen:     int32(unsafe.Sizeof(strioctl{})),
		icDP:      unsafe.Pointer(&pto),
	}
	if err := ioctl(f, I_STR, uintptr(unsafe.Pointer(&istr))); err != nil {
		return errors.New("access denied")
	}
	return nil
}

// streamsPush pushes STREAMS modules if not already done so.
func streamsPush(f *os.File, mod string) error {
	buf := []byte(mod)

	// XXX I_FIND is not returning an error when the module
	// is already pushed even though truss reports a return
	// value of 1. A bug in the Go Solaris syscall interface?
	// XXX without this we are at risk of the issue
	// https://www.illumos.org/issues/9042
	// but since we are not using libc or XPG4.2, we should not be
	// double-pushing modules

	if err := ioctl(f, I_FIND, uintptr(unsafe.Pointer(&buf[0]))); err != nil {
		return nil
	}
	return ioctl(f, I_PUSH, uintptr(unsafe.Pointer(&buf[0])))
}
