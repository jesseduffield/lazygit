// Copyright 2026 The TCell Authors
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

// CursorStyle represents the style of cursor, and covers the shape, whether it
// blinks, and whether it is visible.  Cursor color is handled separately, if at all.
type CursorStyle byte

const (
	SteadyBlock = CursorStyle(iota) // The default
	SteadyBar
	SteadyUnderline
	BlinkingBlock     = SteadyBlock | blinkingCursor
	BlinkingBar       = SteadyBar | blinkingCursor
	BlinkingUnderline = SteadyUnderline | blinkingCursor

	hiddenCursor   = CursorStyle(1 << 7) // If set, cursor should be hidden
	blinkingCursor = CursorStyle(1 << 6) // If set, cursor should blink
)

func (cs CursorStyle) IsVisible() bool {
	return cs&hiddenCursor == 0
}

func (cs CursorStyle) IsBlinking() bool {
	return cs&blinkingCursor != 0
}

func (cs CursorStyle) Hide() CursorStyle {
	return cs | hiddenCursor
}

func (cs CursorStyle) Show() CursorStyle {
	return cs &^ hiddenCursor
}

func (cs CursorStyle) Blink() CursorStyle {
	return cs | blinkingCursor
}

func (cs CursorStyle) Steady() CursorStyle {
	return cs &^ blinkingCursor
}
