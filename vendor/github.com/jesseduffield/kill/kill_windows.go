// adapted from https://blog.csdn.net/fyxichen/article/details/51857864

package kill

import (
	"golang.org/x/sys/windows"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

const PROCESS_ALL_ACCESS = windows.STANDARD_RIGHTS_REQUIRED | windows.SYNCHRONIZE | 0xffff

func GetWindowsHandle(pid int) (handle windows.Handle, err error) {
	handle, err = windows.OpenProcess(PROCESS_ALL_ACCESS, false, uint32(pid))
	return
}

func GetCreationTime(pid int) (time int64, err error) {
	handle, err := GetWindowsHandle(pid)
	if err != nil {
		return
	}
	defer closeHandle(HANDLE(handle))

	var u syscall.Rusage
	err = syscall.GetProcessTimes(syscall.Handle(handle), &u.CreationTime, &u.ExitTime, &u.KernelTime, &u.UserTime)
	if err != nil {
		return
	}

	time = u.CreationTime.Nanoseconds()
	return
}

// Kill kills a process, along with any child processes it may have spawned.
func Kill(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		// You can't kill a person with no body
		return nil
	}

	ptime, err := GetCreationTime(cmd.Process.Pid)
	if err != nil {
		return err
	}

	pids := Getppids(uint32(cmd.Process.Pid), ptime)
	for _, pid := range pids {
		pro, err := os.FindProcess(int(pid))
		if err != nil {
			continue
		}

		pro.Kill()
	}

	return nil
}

// PrepareForChildren ensures that child processes of this parent process will share the same group id
// as the parent, meaning when the call Kill on the parent process, we'll kill
// the whole group, parent and children both. Gruesome when you think about it.
func PrepareForChildren(cmd *exec.Cmd) {
	// do nothing because on windows our Kill function handles children by default.
}

const (
	MAX_PATH           = 260
	TH32CS_SNAPPROCESS = 0x00000002
)

type ProcessInfo struct {
	Name string
	Pid  uint32
	PPid uint32
}

type PROCESSENTRY32 struct {
	DwSize              uint32
	CntUsage            uint32
	Th32ProcessID       uint32
	Th32DefaultHeapID   uintptr
	Th32ModuleID        uint32
	CntThreads          uint32
	Th32ParentProcessID uint32
	PcPriClassBase      int32
	DwFlags             uint32
	SzExeFile           [MAX_PATH]uint16
}

type HANDLE uintptr

var (
	modkernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procCreateToolhelp32Snapshot = modkernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modkernel32.NewProc("Process32FirstW")
	procProcess32Next            = modkernel32.NewProc("Process32NextW")
	procCloseHandle              = modkernel32.NewProc("CloseHandle")
)

func Getppids(pid uint32, ptime int64) []uint32 {
	infos, err := GetProcs()
	if err != nil {
		return []uint32{pid}
	}
	var pids []uint32 = make([]uint32, 0, len(infos))
	var index int = 0
	pids = append(pids, pid)

	var length int = len(pids)
	for index < length {
		for _, info := range infos {
			if info.PPid == pids[index] {
				ctime, err := GetCreationTime(int(info.Pid))
				if err != nil {
					continue
				}

				if ctime >= ptime {
					// Only appending if child is newer than parent, otherwise PPid was reused
					pids = append(pids, info.Pid)
				}
			}
		}
		index += 1
		length = len(pids)
	}
	return pids
}

func GetProcs() (procs []ProcessInfo, err error) {
	snap := createToolhelp32Snapshot(TH32CS_SNAPPROCESS, uint32(0))
	if snap == 0 {
		err = syscall.GetLastError()
		return
	}
	defer closeHandle(snap)
	var pe32 PROCESSENTRY32
	pe32.DwSize = uint32(unsafe.Sizeof(pe32))
	if process32First(snap, &pe32) == false {
		err = syscall.GetLastError()
		return
	}
	procs = append(procs, ProcessInfo{syscall.UTF16ToString(pe32.SzExeFile[:260]), pe32.Th32ProcessID, pe32.Th32ParentProcessID})
	for process32Next(snap, &pe32) {
		procs = append(procs, ProcessInfo{syscall.UTF16ToString(pe32.SzExeFile[:260]), pe32.Th32ProcessID, pe32.Th32ParentProcessID})
	}
	return
}

func createToolhelp32Snapshot(flags, processId uint32) HANDLE {
	ret, _, _ := procCreateToolhelp32Snapshot.Call(uintptr(flags), uintptr(processId))
	if ret <= 0 {
		return HANDLE(0)
	}
	return HANDLE(ret)
}

func process32First(snapshot HANDLE, pe *PROCESSENTRY32) bool {
	ret, _, _ := procProcess32First.Call(uintptr(snapshot), uintptr(unsafe.Pointer(pe)))
	return ret != 0
}

func process32Next(snapshot HANDLE, pe *PROCESSENTRY32) bool {
	ret, _, _ := procProcess32Next.Call(uintptr(snapshot), uintptr(unsafe.Pointer(pe)))
	return ret != 0
}

func closeHandle(object HANDLE) bool {
	ret, _, _ := procCloseHandle.Call(uintptr(object))
	return ret != 0
}
