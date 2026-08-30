// Copyright 2022 The TCell Authors
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

package tcell

import (
	"strings"
	"sync"

	"golang.org/x/text/encoding"

	gencoding "github.com/gdamore/encoding"
)

var encodings map[string]encoding.Encoding
var encodingLk sync.Mutex
var encodingFallback EncodingFallback = EncodingFallbackASCII

// RegisterEncoding may be called by the application to register an encoding.
// The presence of additional encodings will facilitate application usage with
// terminal environments where the I/O subsystem does not support Unicode.
//
// Modern systems and terminal emulators usually use UTF-8, and for those
// systems, this API is also unnecessary.  For example, Windows, macOS, and
// modern Linux systems generally will work out of the box without any of this.
//
// Use of UTF-8 is recommended when possible, as it saves quite a lot processing overhead.
//
// Aliases can be registered as well, for example "8859-15" could be an alias
// for "ISO8859-15".
//
// For POSIX systems, this package will check the environment variables
// LC_ALL, LC_CTYPE,  and LANG (in that order) to determine the character set.
// These are expected to have the following pattern:
//
//	$language[.$codeset[@$variant]
//
// We extract only the $codeset part, which will usually be something like
// UTF-8 or ISO8859-15 or KOI8-R.  Note that if the locale is either "POSIX"
// or "C", then we assume US-ASCII (the POSIX 'portable character set'
// and assume all other characters are somehow invalid.)
//
// Please see the Go documentation for golang.org/x/text/encoding -- most of
// the common ones exist already as stock variables.  For example, ISO8859-15
// can be registered using the following code:
//
// Note that some encodings are quite large (for example GB18030 which is a
// superset of Unicode) and so the application size can be expected to
// increase quite a bit as each encoding is added.
//
// The East Asian encodings have been seen to add 100-200K per encoding to the
// size of the resulting binary.
func RegisterEncoding(charset string, enc encoding.Encoding) {
	encodingLk.Lock()
	charset = strings.ToLower(charset)
	encodings[charset] = enc
	encodingLk.Unlock()
}

// EncodingFallback describes how the system behaves when the locale
// requires a character set that we do not support.  The system always
// supports UTF-8 and US-ASCII. On Windows consoles, UTF-16LE is also
// supported automatically.  Other character sets must be added using the
// RegisterEncoding API.  (A large group of nearly all of them can be
// added using the RegisterAll function in the encoding sub package.)
// The default action will be to fallback to UTF-8.
type EncodingFallback int

const (
	// EncodingFallbackUTF8 behavior causes GetEncoding to assume
	// UTF8 can pass unmodified upon failure.
	EncodingFallbackUTF8 = iota

	// EncodingFallbackFail behavior causes GetEncoding to fail
	// when it cannot find an encoding.
	EncodingFallbackFail

	// EncodingFallbackASCII behavior causes GetEncoding to fall back
	// to a 7-bit ASCII encoding, if no other encoding can be found.
	EncodingFallbackASCII
)

// SetEncodingFallback changes the behavior of GetEncoding when a suitable
// encoding is not found.  The default is EncodingFallbackFail, which
// causes GetEncoding to simply return nil.
func SetEncodingFallback(fb EncodingFallback) {
	encodingLk.Lock()
	encodingFallback = fb
	encodingLk.Unlock()
}

// GetEncoding is used by Screen implementors who want to locate an encoding
// for the given character set name.  Note that this will return nil for
// either the Unicode (UTF-8) or ASCII encodings, since we don't use
// encodings for them but instead have our own native methods.
func GetEncoding(charset string) encoding.Encoding {
	charset = strings.ToLower(charset)
	encodingLk.Lock()
	defer encodingLk.Unlock()
	if enc, ok := encodings[charset]; ok {
		return enc
	}
	switch encodingFallback {
	case EncodingFallbackASCII:
		return gencoding.ASCII
	case EncodingFallbackUTF8:
		return encoding.Nop
	}
	return nil
}

func init() {
	// We always support UTF-8 and ASCII.
	encodings = make(map[string]encoding.Encoding)
	encodings["utf-8"] = gencoding.UTF8
	encodings["utf8"] = gencoding.UTF8
	encodings["us-ascii"] = gencoding.ASCII
	encodings["ascii"] = gencoding.ASCII
	encodings["iso646"] = gencoding.ASCII
}
