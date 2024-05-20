//go:build amd64 && dragonfly
//+build amd64,dragonfly

// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs types_dragonfly.go

package pty

const (
	_C_SPECNAMELEN = 0x3f
)

type fiodgnameArg struct {
	Name      *byte
	Len       uint32
	Pad_cgo_0 [4]byte
}
