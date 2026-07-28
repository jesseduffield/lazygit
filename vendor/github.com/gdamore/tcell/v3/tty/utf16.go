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

package tty

import (
	"unicode/utf16"
	"unicode/utf8"
)

// decodeUTF16Rune decodes one UTF-16 code unit at a time while preserving
// malformed input as replacement characters instead of silently discarding it.
func decodeUTF16Rune(surrogate *rune, wc rune) []rune {
	switch {
	case wc >= 0xD800 && wc <= 0xDBFF:
		if *surrogate != 0 {
			*surrogate = wc
			return []rune{utf8.RuneError}
		}
		*surrogate = wc
		return nil
	case wc >= 0xDC00 && wc <= 0xDFFF:
		if *surrogate == 0 {
			return []rune{utf8.RuneError}
		}
		decoded := utf16.DecodeRune(*surrogate, wc)
		*surrogate = 0
		return []rune{decoded}
	default:
		if *surrogate != 0 {
			*surrogate = 0
			return []rune{utf8.RuneError, wc}
		}
		return []rune{wc}
	}
}
