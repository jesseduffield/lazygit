// Copyright 2026 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vt

// graphemeJoinerBlocks has one bit per 64 runes, set if any rune in that
// block can extend a grapheme cluster.  A clear bit means the segmenter
// cannot join anything in the block, so putRune may skip it.  Only 2619 of
// the 1113984 runes above U+007F can join, so nearly every bit is clear.
//
// The table is derived from the segmenter, not from Go's unicode tables,
// and TestGraphemeJoinerBlocks re-derives it to catch drift.
var graphemeJoinerBlocks = [272]uint64{
	0:   0xfffffffffbc43000,
	1:   0x0089ff15f0002007,
	2:   0x00a8000000000009,
	3:   0x0000000000000005,
	10:  0x00008ffd0e000000,
	15:  0x4100100000000000,
	16:  0x6c30090000002880,
	17:  0x387527b117cffbff,
	19:  0x0000000000020000,
	22:  0xe000180000000010,
	27:  0x0004000000000000,
	28:  0x3000000000000000,
	29:  0x0000070000000260,
	30:  0x0000002808880c15,
	31:  0x0000000000008000,
	224: 0x00000000000000f3,
}
