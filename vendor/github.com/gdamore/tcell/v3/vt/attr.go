// Copyright 2025 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package vt

// Attr is a synthetic combination of display attributes for cells, apart
// from color which is handled separately.
type Attr uint16

const (
	Plain         = Attr(0)      // basic, plain style
	Bold          = Attr(1 << 0) // maybe double strike, or just brighter
	Blink         = Attr(1 << 1) // NB: many terminals do not support it
	Reverse       = Attr(1 << 2) // foreground and background reversed
	Dim           = Attr(1 << 3) // fainter, may also be lower alpha
	Italic        = Attr(1 << 4) // italicized
	StrikeThrough = Attr(1 << 5) // crossed-out
	Underline     = Attr(1 << 6) // any underline style
	Overline      = Attr(1 << 7) // rarely supported

	// Underline styles, always mixed with underline, only one can be selected
	PlainUnderline  = Underline
	DoubleUnderline = Underline | 1<<13
	CurlyUnderline  = Underline | 2<<13
	DottedUnderline = Underline | 3<<13
	DashedUnderline = Underline | 4<<13
	UnderlineMask   = Underline | 7<<13 // bits 13, 14, 15
)
