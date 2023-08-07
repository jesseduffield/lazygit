// Copyright 2020 The TCell Authors
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

// This is just a "minimalist" set of the base terminal descriptions.
// It should be sufficient for most applications.

// Package base contains the base terminal descriptions that are likely
// to be needed by any stock application.  It is imported by default in the
// terminfo package, so terminal types listed here will be available to any
// tcell application.
package base

import (
	// The following imports just register themselves --
	// thse are the terminal types we aggregate in this package.
	_ "github.com/stefanhaller/tcell/v2/terminfo/a/ansi"
	_ "github.com/stefanhaller/tcell/v2/terminfo/v/vt100"
	_ "github.com/stefanhaller/tcell/v2/terminfo/v/vt102"
	_ "github.com/stefanhaller/tcell/v2/terminfo/v/vt220"
	_ "github.com/stefanhaller/tcell/v2/terminfo/x/xterm"
)
