// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build openbsd
// +build 386 amd64 arm

package unix

import (
	"errors"
	"fmt"
	"strconv"
	"syscall"
	"unsafe"
)

const (
	_SYS_PLEDGE = 108
)

// Pledge implements the pledge syscall.
//
// The pledge syscall does not accept execpromises on OpenBSD releases
// before 6.3.
//
// execpromises must be empty when Pledge is called on OpenBSD
// releases predating 6.3, otherwise an error will be returned.
//
// For more information see pledge(2).
func Pledge(promises, execpromises string) error {
	maj, min, err := majmin()
	if err != nil {
		return err
	}

	// If OpenBSD <= 5.9, pledge is not available.
	if (maj == 5 && min != 9) || maj < 5 {
		return fmt.Errorf("pledge syscall is not available on OpenBSD %d.%d", maj, min)
	}

	// If OpenBSD <= 6.2 and execpromises is not empty
	// return an error - execpromises is not available before 6.3
	if (maj < 6 || (maj == 6 && min <= 2)) && execpromises != "" {
		return fmt.Errorf("cannot use execpromises on OpenBSD %d.%d", maj, min)
	}

	pptr, err := syscall.BytePtrFromString(promises)
	if err != nil {
		return err
	}

	// This variable will hold either a nil unsafe.Pointer or
	// an unsafe.Pointer to a string (execpromises).
	var expr unsafe.Pointer

	// If we're running on OpenBSD > 6.2, pass execpromises to the syscall.
	if maj > 6 || (maj == 6 && min > 2) {
		exptr, err := syscall.BytePtrFromString(execpromises)
		if err != nil {
			return err
		}
		expr = unsafe.Pointer(exptr)
	}

	_, _, e := syscall.Syscall(_SYS_PLEDGE, uintptr(unsafe.Pointer(pptr)), uintptr(expr), 0)
	if e != 0 {
		return e
	}

	return nil
}

// majmin returns major and minor version number for an OpenBSD system.
func majmin() (major int, minor int, err error) {
	var v Utsname
	err = Uname(&v)
	if err != nil {
		return
	}

	major, err = strconv.Atoi(string(v.Release[0]))
	if err != nil {
		err = errors.New("cannot parse major version number returned by uname")
		return
	}

	minor, err = strconv.Atoi(string(v.Release[2]))
	if err != nil {
		err = errors.New("cannot parse minor version number returned by uname")
		return
	}

	return
}
