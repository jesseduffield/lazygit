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

// Package tcell provides a lower-level, portable API for building
// programs that interact with terminals or consoles.  It works with
// both common (and many uncommon!) terminals or terminal emulators,
// and Windows console implementations.
//
// It supports rich color, and modern terminal capabilities such as
// rich key reporting, mouse tracking, bracketed paste, desktop notifications
// and 24-bit color, when the underlying terminal supports it.
package tcell
