//go:build !windows
// +build !windows

package wincred

import (
	"errors"
	"syscall"
)

const (
	sysCRED_TYPE_GENERIC                 = 0
	sysCRED_TYPE_DOMAIN_PASSWORD         = 0
	sysCRED_TYPE_DOMAIN_CERTIFICATE      = 0
	sysCRED_TYPE_DOMAIN_VISIBLE_PASSWORD = 0
	sysCRED_TYPE_GENERIC_CERTIFICATE     = 0
	sysCRED_TYPE_DOMAIN_EXTENDED         = 0

	sysERROR_NOT_FOUND         = syscall.Errno(1)
	sysERROR_INVALID_PARAMETER = syscall.Errno(1)
	sysERROR_BAD_USERNAME      = syscall.Errno(1)
)

func sysCredRead(...interface{}) (*Credential, error) {
	return nil, errors.New("Operation not supported")
}

func sysCredWrite(...interface{}) error {
	return errors.New("Operation not supported")
}

func sysCredDelete(...interface{}) error {
	return errors.New("Operation not supported")
}

func sysCredEnumerate(...interface{}) ([]*Credential, error) {
	return nil, errors.New("Operation not supported")
}
