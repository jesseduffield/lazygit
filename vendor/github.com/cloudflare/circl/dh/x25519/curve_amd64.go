//go:build amd64 && !purego
// +build amd64,!purego

package x25519

import (
	fp "github.com/cloudflare/circl/math/fp25519"
	"golang.org/x/sys/cpu"
)

var hasBmi2Adx = cpu.X86.HasBMI2 && cpu.X86.HasADX

var _ = hasBmi2Adx

func double(x, z *fp.Elt)             { doubleAmd64(x, z) }
func diffAdd(w *[5]fp.Elt, b uint)    { diffAddAmd64(w, b) }
func ladderStep(w *[5]fp.Elt, b uint) { ladderStepAmd64(w, b) }
func mulA24(z, x *fp.Elt)             { mulA24Amd64(z, x) }

//go:noescape
func ladderStepAmd64(w *[5]fp.Elt, b uint)

//go:noescape
func diffAddAmd64(w *[5]fp.Elt, b uint)

//go:noescape
func doubleAmd64(x, z *fp.Elt)

//go:noescape
func mulA24Amd64(z, x *fp.Elt)
