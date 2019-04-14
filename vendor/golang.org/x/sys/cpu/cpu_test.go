// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpu_test

import (
	"runtime"
	"testing"

	"golang.org/x/sys/cpu"
)

func TestAMD64minimalFeatures(t *testing.T) {
	if runtime.GOARCH == "amd64" {
		if !cpu.Initialized {
			t.Fatal("Initialized expected true, got false")
		}
		if !cpu.X86.HasSSE2 {
			t.Fatal("HasSSE2 expected true, got false")
		}
	}
}

func TestAVX2hasAVX(t *testing.T) {
	if runtime.GOARCH == "amd64" {
		if cpu.X86.HasAVX2 && !cpu.X86.HasAVX {
			t.Fatal("HasAVX expected true, got false")
		}
	}
}

func TestARM64minimalFeatures(t *testing.T) {
	if runtime.GOARCH != "arm64" || runtime.GOOS != "linux" {
		return
	}
	if !cpu.ARM64.HasASIMD {
		t.Fatal("HasASIMD expected true, got false")
	}
	if !cpu.ARM64.HasFP {
		t.Fatal("HasFP expected true, got false")
	}
}

// On ppc64x, the ISA bit for POWER8 should always be set on POWER8 and beyond.
func TestPPC64minimalFeatures(t *testing.T) {
	// Do not run this with gccgo on ppc64, as it doesn't have POWER8 as a minimum
	// requirement.
	if runtime.Compiler == "gccgo" && runtime.GOARCH == "ppc64" {
		t.Skip("gccgo does not require POWER8 on ppc64; skipping")
	}
	if runtime.GOARCH == "ppc64" || runtime.GOARCH == "ppc64le" {
		if !cpu.PPC64.IsPOWER8 {
			t.Fatal("IsPOWER8 expected true, got false")
		}
	}
}

func TestS390X(t *testing.T) {
	if runtime.GOARCH != "s390x" {
		return
	}
	if testing.Verbose() {
		t.Logf("%+v\n", cpu.S390X)
	}
	// z/Architecture is mandatory
	if !cpu.S390X.HasZARCH {
		t.Error("HasZARCH expected true, got false")
	}
	// vector-enhancements require vector facility to be enabled
	if cpu.S390X.HasVXE && !cpu.S390X.HasVX {
		t.Error("HasVX expected true, got false (VXE is true)")
	}
}
