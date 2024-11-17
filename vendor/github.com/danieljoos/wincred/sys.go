//go:build windows
// +build windows

package wincred

import (
	"reflect"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modadvapi32            = windows.NewLazySystemDLL("advapi32.dll")
	procCredRead           = modadvapi32.NewProc("CredReadW")
	procCredWrite     proc = modadvapi32.NewProc("CredWriteW")
	procCredDelete    proc = modadvapi32.NewProc("CredDeleteW")
	procCredFree      proc = modadvapi32.NewProc("CredFree")
	procCredEnumerate      = modadvapi32.NewProc("CredEnumerateW")
)

// Interface for syscall.Proc: helps testing
type proc interface {
	Call(a ...uintptr) (r1, r2 uintptr, lastErr error)
}

// https://docs.microsoft.com/en-us/windows/desktop/api/wincred/ns-wincred-_credentialw
type sysCREDENTIAL struct {
	Flags              uint32
	Type               uint32
	TargetName         *uint16
	Comment            *uint16
	LastWritten        windows.Filetime
	CredentialBlobSize uint32
	CredentialBlob     uintptr
	Persist            uint32
	AttributeCount     uint32
	Attributes         uintptr
	TargetAlias        *uint16
	UserName           *uint16
}

// https://docs.microsoft.com/en-us/windows/desktop/api/wincred/ns-wincred-_credential_attributew
type sysCREDENTIAL_ATTRIBUTE struct {
	Keyword   *uint16
	Flags     uint32
	ValueSize uint32
	Value     uintptr
}

// https://docs.microsoft.com/en-us/windows/desktop/api/wincred/ns-wincred-_credentialw
type sysCRED_TYPE uint32

const (
	sysCRED_TYPE_GENERIC                 sysCRED_TYPE = 0x1
	sysCRED_TYPE_DOMAIN_PASSWORD         sysCRED_TYPE = 0x2
	sysCRED_TYPE_DOMAIN_CERTIFICATE      sysCRED_TYPE = 0x3
	sysCRED_TYPE_DOMAIN_VISIBLE_PASSWORD sysCRED_TYPE = 0x4
	sysCRED_TYPE_GENERIC_CERTIFICATE     sysCRED_TYPE = 0x5
	sysCRED_TYPE_DOMAIN_EXTENDED         sysCRED_TYPE = 0x6

	// https://docs.microsoft.com/en-us/windows/desktop/Debug/system-error-codes
	sysERROR_NOT_FOUND         = windows.Errno(1168)
	sysERROR_INVALID_PARAMETER = windows.Errno(87)
	sysERROR_BAD_USERNAME      = windows.Errno(2202)
)

// https://docs.microsoft.com/en-us/windows/desktop/api/wincred/nf-wincred-credreadw
func sysCredRead(targetName string, typ sysCRED_TYPE) (*Credential, error) {
	var pcred *sysCREDENTIAL
	targetNamePtr, _ := windows.UTF16PtrFromString(targetName)
	ret, _, err := syscall.SyscallN(
		procCredRead.Addr(),
		uintptr(unsafe.Pointer(targetNamePtr)),
		uintptr(typ),
		0,
		uintptr(unsafe.Pointer(&pcred)),
	)
	if ret == 0 {
		return nil, err
	}
	defer procCredFree.Call(uintptr(unsafe.Pointer(pcred)))

	return sysToCredential(pcred), nil
}

// https://docs.microsoft.com/en-us/windows/desktop/api/wincred/nf-wincred-credwritew
func sysCredWrite(cred *Credential, typ sysCRED_TYPE) error {
	ncred := sysFromCredential(cred)
	ncred.Type = uint32(typ)
	ret, _, err := procCredWrite.Call(
		uintptr(unsafe.Pointer(ncred)),
		0,
	)
	if ret == 0 {
		return err
	}

	return nil
}

// https://docs.microsoft.com/en-us/windows/desktop/api/wincred/nf-wincred-creddeletew
func sysCredDelete(cred *Credential, typ sysCRED_TYPE) error {
	targetNamePtr, _ := windows.UTF16PtrFromString(cred.TargetName)
	ret, _, err := procCredDelete.Call(
		uintptr(unsafe.Pointer(targetNamePtr)),
		uintptr(typ),
		0,
	)
	if ret == 0 {
		return err
	}

	return nil
}

// https://docs.microsoft.com/en-us/windows/desktop/api/wincred/nf-wincred-credenumeratew
func sysCredEnumerate(filter string, all bool) ([]*Credential, error) {
	var count int
	var pcreds uintptr
	var filterPtr *uint16
	if !all {
		filterPtr, _ = windows.UTF16PtrFromString(filter)
	}
	ret, _, err := syscall.SyscallN(
		procCredEnumerate.Addr(),
		uintptr(unsafe.Pointer(filterPtr)),
		0,
		uintptr(unsafe.Pointer(&count)),
		uintptr(unsafe.Pointer(&pcreds)),
	)
	if ret == 0 {
		return nil, err
	}
	defer procCredFree.Call(pcreds)
	credsSlice := *(*[]*sysCREDENTIAL)(unsafe.Pointer(&reflect.SliceHeader{
		Data: pcreds,
		Len:  count,
		Cap:  count,
	}))
	creds := make([]*Credential, count, count)
	for i, cred := range credsSlice {
		creds[i] = sysToCredential(cred)
	}

	return creds, nil
}
