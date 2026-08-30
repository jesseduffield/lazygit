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

// Separate Row and Col types are used to reduce the chance of mixing up the coordinate axes.
// We use zero-based coordinates (although VT hardware underneath uses one-based in escape sequences).
// The upper left corner of the screen is at coordinate (0, 0).

// Row indicates a row number (y position). We use zero based indices, although the VT
// standard mostly communicates using 1-based offsets.
type Row int

// Col indicates a column number (x position).  We use zero based indices.
type Col int

// Coord indicates a coordinate.  This can also be used for window sizes.
type Coord struct {
	X Col // Column number, or X position.
	Y Row // Row number, or Y position.
}
