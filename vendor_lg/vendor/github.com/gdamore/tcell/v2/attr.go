// Copyright 2024 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tcell

// AttrMask represents a mask of text attributes, apart from color.
// Note that support for attributes may vary widely across terminals.
type AttrMask uint

// Attributes are not colors, but affect the display of text.  They can
// be combined, in some cases, but not others. (E.g. you can have Dim Italic,
// but only CurlyUnderline cannot be mixed with DottedUnderline.)
const (
	AttrBold AttrMask = 1 << iota
	AttrBlink
	AttrReverse
	AttrUnderline // Deprecated: Use UnderlineStyle
	AttrDim
	AttrItalic
	AttrStrikeThrough
	AttrInvalid AttrMask = 1 << 31 // Mark the style or attributes invalid
	AttrNone    AttrMask = 0       // Just normal text.
)
