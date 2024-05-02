//go:build tcell_minimal || nacl || zos || plan9 || windows || android || js
// +build tcell_minimal nacl zos plan9 windows android js

// Copyright 2019 The TCell Authors
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

import (
	"errors"

	"github.com/gdamore/tcell/v2/terminfo"
)

func loadDynamicTerminfo(_ string) (*terminfo.Terminfo, error) {
	return nil, errors.New("terminal type unsupported")
}
