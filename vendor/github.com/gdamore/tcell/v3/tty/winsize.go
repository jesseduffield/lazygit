// Copyright 2025 The TCell Authors
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

package tty

// WindowSize represents the dimensions of the window of a terminal.
type WindowSize struct {
	Width       int // Width in characters
	Height      int // Height in characters
	PixelWidth  int // Width in pixels (zero if not available or known)
	PixelHeight int // Height in pixels (zero if not available or known)
}

// CellDimensions returns the dimensions of a single cell, in pixels
func (ws WindowSize) CellDimensions() (int, int) {
	if ws.PixelWidth == 0 || ws.PixelHeight == 0 {
		return 0, 0
	}
	return (ws.PixelWidth / ws.Width), (ws.PixelHeight / ws.Height)
}
