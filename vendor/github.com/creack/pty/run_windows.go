//go:build windows
// +build windows

package pty

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"unsafe"
)

type startupInfoEx struct {
	startupInfo syscall.StartupInfo
	lpAttrList  syscall.Handle
}

const (
	_EXTENDED_STARTUPINFO_PRESENT = 0x00080000

	_PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE = 0x00020016
)

// StartWithSize assigns a pseudo-terminal Tty to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding Pty.
//
// This will resize the Pty to the specified size before starting the command.
// Starts the process in a new session and sets the controlling terminal.
func StartWithSize(c *exec.Cmd, sz *Winsize) (Pty, error) {
	return StartWithAttrs(c, sz, c.SysProcAttr)
}

// StartWithAttrs assigns a pseudo-terminal Tty to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding Pty.
//
// This will resize the Pty to the specified size before starting the command if a size is provided.
// The `attrs` parameter overrides the one set in c.SysProcAttr.
//
// This should generally not be needed. Used in some edge cases where it is needed to create a pty
// without a controlling terminal.
func StartWithAttrs(c *exec.Cmd, sz *Winsize, attrs *syscall.SysProcAttr) (_ Pty, err error) {
	pty, tty, err := open()
	if err != nil {
		return nil, err
	}

	defer func() {
		// unlike unix command exec, do not close tty unless error happened
		if err != nil {
			_ = tty.Close()
			_ = pty.Close()
		}
	}()

	if sz != nil {
		if err = Setsize(pty, sz); err != nil {
			return nil, err
		}
	}

	// unlike unix command exec, do not set stdin/stdout/stderr

	c.SysProcAttr = attrs

	// do not use os/exec.Start since we need to append console handler to startup info

	err = start((*cmd)(unsafe.Pointer(c)), syscall.Handle(tty.Fd()))
	if err != nil {
		return nil, err
	}

	return pty, err
}

func createExtendedStartupInfo(consoleHandle syscall.Handle) (_ *startupInfoEx, err error) {
	// append console handler to new process
	var (
		attrBufSize uint64
		si          startupInfoEx
	)

	si.startupInfo.Cb = uint32(unsafe.Sizeof(si))

	// get size of attr list
	err = initializeProcThreadAttributeList.Find()
	if err != nil {
		return nil, err
	}

	r1, _, err := initializeProcThreadAttributeList.Call(
		0, // list ptr
		1, // list item count
		0, // dwFlags: reserved, MUST be 0
		uintptr(unsafe.Pointer(&attrBufSize)),
	)
	if r1 == 0 {
		// according to
		// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-initializeprocthreadattributelist
		// which says: This initial call will return an error by design. This is expected behavior.
		//
		// so here we check the returned value of the attr buf size, if it's zero, we cannot update attribute list
		if attrBufSize == 0 {
			return nil, os.NewSyscallError("InitializeProcThreadAttributeList (size)", err)
		}
	}

	attrListBuf := make([]byte, attrBufSize)
	si.lpAttrList = syscall.Handle(unsafe.Pointer(&attrListBuf[0]))
	// create attr list with console handler
	r1, _, err = initializeProcThreadAttributeList.Call(
		uintptr(si.lpAttrList),                // attr list buf
		1,                                     // list item count
		0,                                     // dwFlags: reserved, MUST be 0
		uintptr(unsafe.Pointer(&attrBufSize)), // size of the list
	)
	if r1 == 0 {
		// false
		return nil, os.NewSyscallError("InitializeProcThreadAttributeList (create)", err)
	}

	err = updateProcThreadAttribute.Find()
	if err != nil {
		return nil, err
	}

	r1, _, err = updateProcThreadAttribute.Call(
		uintptr(si.lpAttrList), // buf list
		0,                      // dwFlags: reserved, MUST be 0
		_PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE,
		uintptr(consoleHandle),
		unsafe.Sizeof(consoleHandle),
		0,
		0,
	)
	if r1 == 0 {
		// false
		if deleteProcThreadAttributeList.Find() == nil {
			_, _, _ = deleteProcThreadAttributeList.Call(uintptr(si.lpAttrList))
		}
		return nil, os.NewSyscallError("UpdateProcThreadAttribute", err)
	}

	return &si, nil
}

// copied from os/exec.(*Cmd).Start
// start starts the specified command but does not wait for it to complete.
//
// If Start returns successfully, the c.Process field will be set.
//
// The Wait method will return the exit code and release associated resources
// once the command exits.
func start(c *cmd, consoleHandle syscall.Handle) error {
	if c.lookPathErr != nil {
		_cmd_closeDescriptors(c, c.closeAfterStart)
		_cmd_closeDescriptors(c, c.closeAfterWait)
		return c.lookPathErr
	}
	if runtime.GOOS == "windows" {
		lp, err := lookExtensions(c.Path, c.Dir)
		if err != nil {
			_cmd_closeDescriptors(c, c.closeAfterStart)
			_cmd_closeDescriptors(c, c.closeAfterWait)
			return err
		}
		c.Path = lp
	}
	if c.Process != nil {
		return errors.New("exec: already started")
	}
	if c.ctx != nil {
		select {
		case <-c.ctx.Done():
			_cmd_closeDescriptors(c, c.closeAfterStart)
			_cmd_closeDescriptors(c, c.closeAfterWait)
			return c.ctx.Err()
		default:
		}
	}

	//c.childFiles = make([]*os.File, 0, 3+len(c.ExtraFiles))
	//type F func() (*os.File, error)
	//for _, setupFd := range []F{c.stdin, c.stdout, c.stderr} {
	//	fd, err := setupFd()
	//	if err != nil {
	//		closeDescriptors(c, c.closeAfterStart)
	//		closeDescriptors(c, c.closeAfterWait)
	//		return err
	//	}
	//	c.childFiles = append(c.childFiles, fd)
	//}
	//c.childFiles = append(c.childFiles, c.ExtraFiles...)

	envv, err := _cmd_envv(c)
	if err != nil {
		return err
	}

	c.Process, err = startProcess(c.Path, _cmd_argv(c), &os.ProcAttr{
		Dir:   c.Dir,
		Files: c.childFiles,
		Env:   addCriticalEnv(dedupEnv(envv)),
		Sys:   c.SysProcAttr,
	}, consoleHandle)
	if err != nil {
		_cmd_closeDescriptors(c, c.closeAfterStart)
		_cmd_closeDescriptors(c, c.closeAfterWait)
		return err
	}

	_cmd_closeDescriptors(c, c.closeAfterStart)

	// Don't allocate the channel unless there are goroutines to fire.
	if len(c.goroutine) > 0 {
		c.errch = make(chan error, len(c.goroutine))
		for _, fn := range c.goroutine {
			go func(fn func() error) {
				c.errch <- fn()
			}(fn)
		}
	}

	if c.ctx != nil {
		c.waitDone = make(chan struct{})
		go func() {
			select {
			case <-c.ctx.Done():
				_ = c.Process.Kill()
			case <-c.waitDone:
			}
		}()
	}

	return nil
}

// copied from os.startProcess, add consoleHandle arg
func startProcess(name string, argv []string, attr *os.ProcAttr, consoleHandle syscall.Handle) (p *os.Process, err error) {
	// If there is no SysProcAttr (ie. no Chroot or changed
	// UID/GID), double-check existence of the directory we want
	// to chdir into. We can make the error clearer this way.
	if attr != nil && attr.Sys == nil && attr.Dir != "" {
		if _, err := os.Stat(attr.Dir); err != nil {
			pe := err.(*os.PathError)
			pe.Op = "chdir"
			return nil, pe
		}
	}

	sysattr := &syscall.ProcAttr{
		Dir: attr.Dir,
		Env: attr.Env,
		Sys: attr.Sys,
	}
	if sysattr.Env == nil {
		sysattr.Env, err = execEnvDefault(sysattr.Sys)
		if err != nil {
			return nil, err
		}
	}
	sysattr.Files = make([]uintptr, 0, len(attr.Files))
	for _, f := range attr.Files {
		sysattr.Files = append(sysattr.Files, f.Fd())
	}

	pid, h, e := syscallStartProcess(name, argv, sysattr, consoleHandle)

	// Make sure we don't run the finalizers of attr.Files.
	runtime.KeepAlive(attr)

	if e != nil {
		return nil, &os.PathError{Op: "fork/exec", Path: name, Err: e}
	}

	return newProcess(pid, h), nil
}

//go:linkname zeroProcAttr syscall.zeroProcAttr
var zeroProcAttr syscall.ProcAttr

//go:linkname zeroSysProcAttr syscall.zeroSysProcAttr
var zeroSysProcAttr syscall.SysProcAttr

// copied from syscall.StartProcess, add consoleHandle arg
func syscallStartProcess(argv0 string, argv []string, attr *syscall.ProcAttr, consoleHandle syscall.Handle) (pid int, handle uintptr, err error) {
	if len(argv0) == 0 {
		return 0, 0, syscall.EWINDOWS
	}
	if attr == nil {
		attr = &zeroProcAttr
	}
	sys := attr.Sys
	if sys == nil {
		sys = &zeroSysProcAttr
	}

	//if len(attr.Files) > 3 {
	//	return 0, 0, syscall.EWINDOWS
	//}
	//if len(attr.Files) < 3 {
	//	return 0, 0, syscall.EINVAL
	//}

	if len(attr.Dir) != 0 {
		// StartProcess assumes that argv0 is relative to attr.Dir,
		// because it implies Chdir(attr.Dir) before executing argv0.
		// Windows CreateProcess assumes the opposite: it looks for
		// argv0 relative to the current directory, and, only once the new
		// process is started, it does Chdir(attr.Dir). We are adjusting
		// for that difference here by making argv0 absolute.
		var err error
		argv0, err = joinExeDirAndFName(attr.Dir, argv0)
		if err != nil {
			return 0, 0, err
		}
	}
	argv0p, err := syscall.UTF16PtrFromString(argv0)
	if err != nil {
		return 0, 0, err
	}

	var cmdline string
	// Windows CreateProcess takes the command line as a single string:
	// use attr.CmdLine if set, else build the command line by escaping
	// and joining each argument with spaces
	if sys.CmdLine != "" {
		cmdline = sys.CmdLine
	} else {
		cmdline = makeCmdLine(argv)
	}

	var argvp *uint16
	if len(cmdline) != 0 {
		argvp, err = syscall.UTF16PtrFromString(cmdline)
		if err != nil {
			return 0, 0, err
		}
	}

	var dirp *uint16
	if len(attr.Dir) != 0 {
		dirp, err = syscall.UTF16PtrFromString(attr.Dir)
		if err != nil {
			return 0, 0, err
		}
	}

	// Acquire the fork lock so that no other threads
	// create new fds that are not yet close-on-exec
	// before we fork.
	syscall.ForkLock.Lock()
	defer syscall.ForkLock.Unlock()

	//p, _ := syscall.GetCurrentProcess()
	//fd := make([]syscall.Handle, len(attr.Files))
	//for i := range attr.Files {
	//	if attr.Files[i] > 0 {
	//		err := syscall.DuplicateHandle(p, syscall.Handle(attr.Files[i]), p, &fd[i], 0, true, syscall.DUPLICATE_SAME_ACCESS)
	//		if err != nil {
	//			return 0, 0, err
	//		}
	//		defer syscall.CloseHandle(syscall.Handle(fd[i]))
	//	}
	//}

	// replaced default syscall.StartupInfo with custom startupInfEx for console handle
	//si := new(syscall.StartupInfo)
	//si.Cb = uint32(unsafe.Sizeof(*si))
	si, err := createExtendedStartupInfo(consoleHandle)
	if err != nil {
		return 0, 0, err
	}
	// add finalizer for attribute list cleanup, best effort
	runtime.SetFinalizer(si, func(si *startupInfoEx) {
		if deleteProcThreadAttributeList.Find() == nil {
			_, _, _ = deleteProcThreadAttributeList.Call(uintptr(si.lpAttrList))
		}
	})

	si.startupInfo.Flags = syscall.STARTF_USESTDHANDLES
	if sys.HideWindow {
		si.startupInfo.Flags |= syscall.STARTF_USESHOWWINDOW
		si.startupInfo.ShowWindow = syscall.SW_HIDE
	}
	//si.StdInput = fd[0]
	//si.StdOutput = fd[1]
	//si.StdErr = fd[2]

	pi := new(syscall.ProcessInformation)

	flags := sys.CreationFlags | syscall.CREATE_UNICODE_ENVIRONMENT

	// add startupInfoEx flag
	flags = flags | _EXTENDED_STARTUPINFO_PRESENT

	// ignore security attrs since both Process and Thread handles are not inheritable for conPty
	if sys.Token != 0 {
		err = syscall.CreateProcessAsUser(sys.Token, argv0p, argvp, nil, nil, false, flags, createEnvBlock(attr.Env), dirp, &si.startupInfo, pi)
	} else {
		err = syscall.CreateProcess(argv0p, argvp, nil, nil, false, flags, createEnvBlock(attr.Env), dirp, &si.startupInfo, pi)
	}
	if err != nil {
		return 0, 0, err
	}
	defer syscall.CloseHandle(syscall.Handle(pi.Thread))

	return int(pi.ProcessId), uintptr(pi.Process), nil
}
