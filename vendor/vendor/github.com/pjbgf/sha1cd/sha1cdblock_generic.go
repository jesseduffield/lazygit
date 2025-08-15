// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Originally from: https://github.com/go/blob/master/src/crypto/sha1/sha1block.go
// It has been modified to support collision detection.

package sha1cd

import (
	"fmt"
	"math/bits"

	shared "github.com/pjbgf/sha1cd/internal"
	"github.com/pjbgf/sha1cd/ubc"
)

// blockGeneric is a portable, pure Go version of the SHA-1 block step.
// It's used by sha1block_generic.go and tests.
func blockGeneric(dig *digest, p []byte) {
	var w [16]uint32

	// cs stores the pre-step compression state for only the steps required for the
	// collision detection, which are 0, 58 and 65.
	// Refer to ubc/const.go for more details.
	cs := [shared.PreStepState][shared.WordBuffers]uint32{}

	h0, h1, h2, h3, h4 := dig.h[0], dig.h[1], dig.h[2], dig.h[3], dig.h[4]
	for len(p) >= shared.Chunk {
		m1 := [shared.Rounds]uint32{}
		hi := 1

		// Collision attacks are thwarted by hashing a detected near-collision block 3 times.
		// Think of it as extending SHA-1 from 80-steps to 240-steps for such blocks:
		// 		The best collision attacks against SHA-1 have complexity about 2^60,
		// 		thus for 240-steps an immediate lower-bound for the best cryptanalytic attacks would be 2^180.
		// 		An attacker would be better off using a generic birthday search of complexity 2^80.
	rehash:
		a, b, c, d, e := h0, h1, h2, h3, h4

		// Each of the four 20-iteration rounds
		// differs only in the computation of f and
		// the choice of K (K0, K1, etc).
		i := 0

		// Store pre-step compression state for the collision detection.
		cs[0] = [shared.WordBuffers]uint32{a, b, c, d, e}

		for ; i < 16; i++ {
			// load step
			j := i * 4
			w[i] = uint32(p[j])<<24 | uint32(p[j+1])<<16 | uint32(p[j+2])<<8 | uint32(p[j+3])

			f := b&c | (^b)&d
			t := bits.RotateLeft32(a, 5) + f + e + w[i&0xf] + shared.K0
			a, b, c, d, e = t, a, bits.RotateLeft32(b, 30), c, d

			// Store compression state for the collision detection.
			m1[i] = w[i&0xf]
		}
		for ; i < 20; i++ {
			tmp := w[(i-3)&0xf] ^ w[(i-8)&0xf] ^ w[(i-14)&0xf] ^ w[(i)&0xf]
			w[i&0xf] = tmp<<1 | tmp>>(32-1)

			f := b&c | (^b)&d
			t := bits.RotateLeft32(a, 5) + f + e + w[i&0xf] + shared.K0
			a, b, c, d, e = t, a, bits.RotateLeft32(b, 30), c, d

			// Store compression state for the collision detection.
			m1[i] = w[i&0xf]
		}
		for ; i < 40; i++ {
			tmp := w[(i-3)&0xf] ^ w[(i-8)&0xf] ^ w[(i-14)&0xf] ^ w[(i)&0xf]
			w[i&0xf] = tmp<<1 | tmp>>(32-1)

			f := b ^ c ^ d
			t := bits.RotateLeft32(a, 5) + f + e + w[i&0xf] + shared.K1
			a, b, c, d, e = t, a, bits.RotateLeft32(b, 30), c, d

			// Store compression state for the collision detection.
			m1[i] = w[i&0xf]
		}
		for ; i < 60; i++ {
			if i == 58 {
				// Store pre-step compression state for the collision detection.
				cs[1] = [shared.WordBuffers]uint32{a, b, c, d, e}
			}

			tmp := w[(i-3)&0xf] ^ w[(i-8)&0xf] ^ w[(i-14)&0xf] ^ w[(i)&0xf]
			w[i&0xf] = tmp<<1 | tmp>>(32-1)

			f := ((b | c) & d) | (b & c)
			t := bits.RotateLeft32(a, 5) + f + e + w[i&0xf] + shared.K2
			a, b, c, d, e = t, a, bits.RotateLeft32(b, 30), c, d

			// Store compression state for the collision detection.
			m1[i] = w[i&0xf]
		}
		for ; i < 80; i++ {
			if i == 65 {
				// Store pre-step compression state for the collision detection.
				cs[2] = [shared.WordBuffers]uint32{a, b, c, d, e}
			}

			tmp := w[(i-3)&0xf] ^ w[(i-8)&0xf] ^ w[(i-14)&0xf] ^ w[(i)&0xf]
			w[i&0xf] = tmp<<1 | tmp>>(32-1)

			f := b ^ c ^ d
			t := bits.RotateLeft32(a, 5) + f + e + w[i&0xf] + shared.K3
			a, b, c, d, e = t, a, bits.RotateLeft32(b, 30), c, d

			// Store compression state for the collision detection.
			m1[i] = w[i&0xf]
		}

		h0 += a
		h1 += b
		h2 += c
		h3 += d
		h4 += e

		if hi == 2 {
			hi++
			goto rehash
		}

		if hi == 1 {
			col := checkCollision(m1, cs, [shared.WordBuffers]uint32{h0, h1, h2, h3, h4})
			if col {
				dig.col = true
				hi++
				goto rehash
			}
		}

		p = p[shared.Chunk:]
	}

	dig.h[0], dig.h[1], dig.h[2], dig.h[3], dig.h[4] = h0, h1, h2, h3, h4
}

func checkCollision(
	m1 [shared.Rounds]uint32,
	cs [shared.PreStepState][shared.WordBuffers]uint32,
	state [shared.WordBuffers]uint32) bool {

	if mask := ubc.CalculateDvMask(m1); mask != 0 {
		dvs := ubc.SHA1_dvs()

		for i := 0; dvs[i].DvType != 0; i++ {
			if (mask & ((uint32)(1) << uint32(dvs[i].MaskB))) != 0 {
				var csState [shared.WordBuffers]uint32
				switch dvs[i].TestT {
				case 58:
					csState = cs[1]
				case 65:
					csState = cs[2]
				case 0:
					csState = cs[0]
				default:
					panic(fmt.Sprintf("dvs data is trying to use a testT that isn't available: %d", dvs[i].TestT))
				}

				col := hasCollided(
					dvs[i].TestT, // testT is the step number
					// m2 is a secondary message created XORing with
					// ubc's DM prior to the SHA recompression step.
					m1, dvs[i].Dm,
					csState,
					state)

				if col {
					return true
				}
			}
		}
	}
	return false
}

func hasCollided(step uint32, m1, dm [shared.Rounds]uint32,
	state [shared.WordBuffers]uint32, h [shared.WordBuffers]uint32) bool {
	// Intermediary Hash Value.
	ihv := [shared.WordBuffers]uint32{}

	a, b, c, d, e := state[0], state[1], state[2], state[3], state[4]

	// Walk backwards from current step to undo previous compression.
	// The existing collision detection does not have dvs higher than 65,
	// start value of i accordingly.
	for i := uint32(64); i >= 60; i-- {
		a, b, c, d, e = b, c, d, e, a
		if step > i {
			b = bits.RotateLeft32(b, -30)
			f := b ^ c ^ d
			e -= bits.RotateLeft32(a, 5) + f + shared.K3 + (m1[i] ^ dm[i]) // m2 = m1 ^ dm.
		}
	}
	for i := uint32(59); i >= 40; i-- {
		a, b, c, d, e = b, c, d, e, a
		if step > i {
			b = bits.RotateLeft32(b, -30)
			f := ((b | c) & d) | (b & c)
			e -= bits.RotateLeft32(a, 5) + f + shared.K2 + (m1[i] ^ dm[i])
		}
	}
	for i := uint32(39); i >= 20; i-- {
		a, b, c, d, e = b, c, d, e, a
		if step > i {
			b = bits.RotateLeft32(b, -30)
			f := b ^ c ^ d
			e -= bits.RotateLeft32(a, 5) + f + shared.K1 + (m1[i] ^ dm[i])
		}
	}
	for i := uint32(20); i > 0; i-- {
		j := i - 1
		a, b, c, d, e = b, c, d, e, a
		if step > j {
			b = bits.RotateLeft32(b, -30) // undo the rotate left
			f := b&c | (^b)&d
			// subtract from e
			e -= bits.RotateLeft32(a, 5) + f + shared.K0 + (m1[j] ^ dm[j])
		}
	}

	ihv[0] = a
	ihv[1] = b
	ihv[2] = c
	ihv[3] = d
	ihv[4] = e
	a = state[0]
	b = state[1]
	c = state[2]
	d = state[3]
	e = state[4]

	// Recompress blocks based on the current step.
	// The existing collision detection does not have dvs below 58, so they have been removed
	// from the source code. If new dvs are added which target rounds below 40, that logic
	// will need to be readded here.
	for i := uint32(40); i < 60; i++ {
		if step <= i {
			f := ((b | c) & d) | (b & c)
			t := bits.RotateLeft32(a, 5) + f + e + shared.K2 + (m1[i] ^ dm[i])
			a, b, c, d, e = t, a, bits.RotateLeft32(b, 30), c, d
		}
	}
	for i := uint32(60); i < 80; i++ {
		if step <= i {
			f := b ^ c ^ d
			t := bits.RotateLeft32(a, 5) + f + e + shared.K3 + (m1[i] ^ dm[i])
			a, b, c, d, e = t, a, bits.RotateLeft32(b, 30), c, d
		}
	}

	ihv[0] += a
	ihv[1] += b
	ihv[2] += c
	ihv[3] += d
	ihv[4] += e

	if ((ihv[0] ^ h[0]) | (ihv[1] ^ h[1]) |
		(ihv[2] ^ h[2]) | (ihv[3] ^ h[3]) | (ihv[4] ^ h[4])) == 0 {
		return true
	}

	return false
}
