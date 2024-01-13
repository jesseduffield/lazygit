// Package util provides utility functions for the goldmark.
package util

import (
	"bytes"
	"io"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// A CopyOnWriteBuffer is a byte buffer that copies buffer when
// it need to be changed.
type CopyOnWriteBuffer struct {
	buffer []byte
	copied bool
}

// NewCopyOnWriteBuffer returns a new CopyOnWriteBuffer.
func NewCopyOnWriteBuffer(buffer []byte) CopyOnWriteBuffer {
	return CopyOnWriteBuffer{
		buffer: buffer,
		copied: false,
	}
}

// Write writes given bytes to the buffer.
// Write allocate new buffer and clears it at the first time.
func (b *CopyOnWriteBuffer) Write(value []byte) {
	if !b.copied {
		b.buffer = make([]byte, 0, len(b.buffer)+20)
		b.copied = true
	}
	b.buffer = append(b.buffer, value...)
}

// WriteString writes given string to the buffer.
// WriteString allocate new buffer and clears it at the first time.
func (b *CopyOnWriteBuffer) WriteString(value string) {
	b.Write(StringToReadOnlyBytes(value))
}

// Append appends given bytes to the buffer.
// Append copy buffer at the first time.
func (b *CopyOnWriteBuffer) Append(value []byte) {
	if !b.copied {
		tmp := make([]byte, len(b.buffer), len(b.buffer)+20)
		copy(tmp, b.buffer)
		b.buffer = tmp
		b.copied = true
	}
	b.buffer = append(b.buffer, value...)
}

// AppendString appends given string to the buffer.
// AppendString copy buffer at the first time.
func (b *CopyOnWriteBuffer) AppendString(value string) {
	b.Append(StringToReadOnlyBytes(value))
}

// WriteByte writes the given byte to the buffer.
// WriteByte allocate new buffer and clears it at the first time.
func (b *CopyOnWriteBuffer) WriteByte(c byte) {
	if !b.copied {
		b.buffer = make([]byte, 0, len(b.buffer)+20)
		b.copied = true
	}
	b.buffer = append(b.buffer, c)
}

// AppendByte appends given bytes to the buffer.
// AppendByte copy buffer at the first time.
func (b *CopyOnWriteBuffer) AppendByte(c byte) {
	if !b.copied {
		tmp := make([]byte, len(b.buffer), len(b.buffer)+20)
		copy(tmp, b.buffer)
		b.buffer = tmp
		b.copied = true
	}
	b.buffer = append(b.buffer, c)
}

// Bytes returns bytes of this buffer.
func (b *CopyOnWriteBuffer) Bytes() []byte {
	return b.buffer
}

// IsCopied returns true if buffer has been copied, otherwise false.
func (b *CopyOnWriteBuffer) IsCopied() bool {
	return b.copied
}

// IsEscapedPunctuation returns true if character at a given index i
// is an escaped punctuation, otherwise false.
func IsEscapedPunctuation(source []byte, i int) bool {
	return source[i] == '\\' && i < len(source)-1 && IsPunct(source[i+1])
}

// ReadWhile read the given source while pred is true.
func ReadWhile(source []byte, index [2]int, pred func(byte) bool) (int, bool) {
	j := index[0]
	ok := false
	for ; j < index[1]; j++ {
		c1 := source[j]
		if pred(c1) {
			ok = true
			continue
		}
		break
	}
	return j, ok
}

// IsBlank returns true if the given string is all space characters.
func IsBlank(bs []byte) bool {
	for _, b := range bs {
		if !IsSpace(b) {
			return false
		}
	}
	return true
}

// VisualizeSpaces visualize invisible space characters.
func VisualizeSpaces(bs []byte) []byte {
	bs = bytes.Replace(bs, []byte(" "), []byte("[SPACE]"), -1)
	bs = bytes.Replace(bs, []byte("\t"), []byte("[TAB]"), -1)
	bs = bytes.Replace(bs, []byte("\n"), []byte("[NEWLINE]\n"), -1)
	bs = bytes.Replace(bs, []byte("\r"), []byte("[CR]"), -1)
	bs = bytes.Replace(bs, []byte("\v"), []byte("[VTAB]"), -1)
	bs = bytes.Replace(bs, []byte("\x00"), []byte("[NUL]"), -1)
	bs = bytes.Replace(bs, []byte("\ufffd"), []byte("[U+FFFD]"), -1)
	return bs
}

// TabWidth calculates actual width of a tab at the given position.
func TabWidth(currentPos int) int {
	return 4 - currentPos%4
}

// IndentPosition searches an indent position with the given width for the given line.
// If the line contains tab characters, paddings may be not zero.
// currentPos==0 and width==2:
//
//     position: 0    1
//               [TAB]aaaa
//     width:    1234 5678
//
// width=2 is in the tab character. In this case, IndentPosition returns
// (pos=1, padding=2)
func IndentPosition(bs []byte, currentPos, width int) (pos, padding int) {
	return IndentPositionPadding(bs, currentPos, 0, width)
}

// IndentPositionPadding searches an indent position with the given width for the given line.
// This function is mostly same as IndentPosition except this function
// takes account into additional paddings.
func IndentPositionPadding(bs []byte, currentPos, paddingv, width int) (pos, padding int) {
	if width == 0 {
		return 0, paddingv
	}
	w := 0
	i := 0
	l := len(bs)
	for ; i < l; i++ {
		if bs[i] == '\t' && w < width {
			w += TabWidth(currentPos + w)
		} else if bs[i] == ' ' && w < width {
			w++
		} else {
			break
		}
	}
	if w >= width {
		return i - paddingv, w - width
	}
	return -1, -1
}

// DedentPosition dedents lines by the given width.
//
// Deprecated: This function has bugs. Use util.IndentPositionPadding and util.FirstNonSpacePosition.
func DedentPosition(bs []byte, currentPos, width int) (pos, padding int) {
	if width == 0 {
		return 0, 0
	}
	w := 0
	l := len(bs)
	i := 0
	for ; i < l; i++ {
		if bs[i] == '\t' {
			w += TabWidth(currentPos + w)
		} else if bs[i] == ' ' {
			w++
		} else {
			break
		}
	}
	if w >= width {
		return i, w - width
	}
	return i, 0
}

// DedentPositionPadding dedents lines by the given width.
// This function is mostly same as DedentPosition except this function
// takes account into additional paddings.
//
// Deprecated: This function has bugs. Use util.IndentPositionPadding and util.FirstNonSpacePosition.
func DedentPositionPadding(bs []byte, currentPos, paddingv, width int) (pos, padding int) {
	if width == 0 {
		return 0, paddingv
	}

	w := 0
	i := 0
	l := len(bs)
	for ; i < l; i++ {
		if bs[i] == '\t' {
			w += TabWidth(currentPos + w)
		} else if bs[i] == ' ' {
			w++
		} else {
			break
		}
	}
	if w >= width {
		return i - paddingv, w - width
	}
	return i - paddingv, 0
}

// IndentWidth calculate an indent width for the given line.
func IndentWidth(bs []byte, currentPos int) (width, pos int) {
	l := len(bs)
	for i := 0; i < l; i++ {
		b := bs[i]
		if b == ' ' {
			width++
			pos++
		} else if b == '\t' {
			width += TabWidth(currentPos + width)
			pos++
		} else {
			break
		}
	}
	return
}

// FirstNonSpacePosition returns a position line that is a first nonspace
// character.
func FirstNonSpacePosition(bs []byte) int {
	i := 0
	for ; i < len(bs); i++ {
		c := bs[i]
		if c == ' ' || c == '\t' {
			continue
		}
		if c == '\n' {
			return -1
		}
		return i
	}
	return -1
}

// FindClosure returns a position that closes the given opener.
// If codeSpan is set true, it ignores characters in code spans.
// If allowNesting is set true, closures correspond to nested opener will be
// ignored.
//
// Deprecated: This function can not handle newlines. Many elements
// can be existed over multiple lines(e.g. link labels).
// Use text.Reader.FindClosure.
func FindClosure(bs []byte, opener, closure byte, codeSpan, allowNesting bool) int {
	i := 0
	opened := 1
	codeSpanOpener := 0
	for i < len(bs) {
		c := bs[i]
		if codeSpan && codeSpanOpener != 0 && c == '`' {
			codeSpanCloser := 0
			for ; i < len(bs); i++ {
				if bs[i] == '`' {
					codeSpanCloser++
				} else {
					i--
					break
				}
			}
			if codeSpanCloser == codeSpanOpener {
				codeSpanOpener = 0
			}
		} else if codeSpanOpener == 0 && c == '\\' && i < len(bs)-1 && IsPunct(bs[i+1]) {
			i += 2
			continue
		} else if codeSpan && codeSpanOpener == 0 && c == '`' {
			for ; i < len(bs); i++ {
				if bs[i] == '`' {
					codeSpanOpener++
				} else {
					i--
					break
				}
			}
		} else if (codeSpan && codeSpanOpener == 0) || !codeSpan {
			if c == closure {
				opened--
				if opened == 0 {
					return i
				}
			} else if c == opener {
				if !allowNesting {
					return -1
				}
				opened++
			}
		}
		i++
	}
	return -1
}

// TrimLeft trims characters in the given s from head of the source.
// bytes.TrimLeft offers same functionalities, but bytes.TrimLeft
// allocates new buffer for the result.
func TrimLeft(source, b []byte) []byte {
	i := 0
	for ; i < len(source); i++ {
		c := source[i]
		found := false
		for j := 0; j < len(b); j++ {
			if c == b[j] {
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	return source[i:]
}

// TrimRight trims characters in the given s from tail of the source.
func TrimRight(source, b []byte) []byte {
	i := len(source) - 1
	for ; i >= 0; i-- {
		c := source[i]
		found := false
		for j := 0; j < len(b); j++ {
			if c == b[j] {
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	return source[:i+1]
}

// TrimLeftLength returns a length of leading specified characters.
func TrimLeftLength(source, s []byte) int {
	return len(source) - len(TrimLeft(source, s))
}

// TrimRightLength returns a length of trailing specified characters.
func TrimRightLength(source, s []byte) int {
	return len(source) - len(TrimRight(source, s))
}

// TrimLeftSpaceLength returns a length of leading space characters.
func TrimLeftSpaceLength(source []byte) int {
	i := 0
	for ; i < len(source); i++ {
		if !IsSpace(source[i]) {
			break
		}
	}
	return i
}

// TrimRightSpaceLength returns a length of trailing space characters.
func TrimRightSpaceLength(source []byte) int {
	l := len(source)
	i := l - 1
	for ; i >= 0; i-- {
		if !IsSpace(source[i]) {
			break
		}
	}
	if i < 0 {
		return l
	}
	return l - 1 - i
}

// TrimLeftSpace returns a subslice of the given string by slicing off all leading
// space characters.
func TrimLeftSpace(source []byte) []byte {
	return TrimLeft(source, spaces)
}

// TrimRightSpace returns a subslice of the given string by slicing off all trailing
// space characters.
func TrimRightSpace(source []byte) []byte {
	return TrimRight(source, spaces)
}

// DoFullUnicodeCaseFolding performs full unicode case folding to given bytes.
func DoFullUnicodeCaseFolding(v []byte) []byte {
	var rbuf []byte
	cob := NewCopyOnWriteBuffer(v)
	n := 0
	for i := 0; i < len(v); i++ {
		c := v[i]
		if c < 0xb5 {
			if c >= 0x41 && c <= 0x5a {
				// A-Z to a-z
				cob.Write(v[n:i])
				cob.WriteByte(c + 32)
				n = i + 1
			}
			continue
		}

		if !utf8.RuneStart(c) {
			continue
		}
		r, length := utf8.DecodeRune(v[i:])
		if r == utf8.RuneError {
			continue
		}
		folded, ok := unicodeCaseFoldings[r]
		if !ok {
			continue
		}

		cob.Write(v[n:i])
		if rbuf == nil {
			rbuf = make([]byte, 4)
		}
		for _, f := range folded {
			l := utf8.EncodeRune(rbuf, f)
			cob.Write(rbuf[:l])
		}
		i += length - 1
		n = i + 1
	}
	if cob.IsCopied() {
		cob.Write(v[n:])
	}
	return cob.Bytes()
}

// ReplaceSpaces replaces sequence of spaces with the given repl.
func ReplaceSpaces(source []byte, repl byte) []byte {
	var ret []byte
	start := -1
	for i, c := range source {
		iss := IsSpace(c)
		if start < 0 && iss {
			start = i
			continue
		} else if start >= 0 && iss {
			continue
		} else if start >= 0 {
			if ret == nil {
				ret = make([]byte, 0, len(source))
				ret = append(ret, source[:start]...)
			}
			ret = append(ret, repl)
			start = -1
		}
		if ret != nil {
			ret = append(ret, c)
		}
	}
	if start >= 0 && ret != nil {
		ret = append(ret, repl)
	}
	if ret == nil {
		return source
	}
	return ret
}

// ToRune decode given bytes start at pos and returns a rune.
func ToRune(source []byte, pos int) rune {
	i := pos
	for ; i >= 0; i-- {
		if utf8.RuneStart(source[i]) {
			break
		}
	}
	r, _ := utf8.DecodeRune(source[i:])
	return r
}

// ToValidRune returns 0xFFFD if the given rune is invalid, otherwise v.
func ToValidRune(v rune) rune {
	if v == 0 || !utf8.ValidRune(v) {
		return rune(0xFFFD)
	}
	return v
}

// ToLinkReference converts given bytes into a valid link reference string.
// ToLinkReference performs unicode case folding, trims leading and trailing spaces,  converts into lower
// case and replace spaces with a single space character.
func ToLinkReference(v []byte) string {
	v = TrimLeftSpace(v)
	v = TrimRightSpace(v)
	v = DoFullUnicodeCaseFolding(v)
	return string(ReplaceSpaces(v, ' '))
}

var htmlEscapeTable = [256][]byte{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, []byte("&quot;"), nil, nil, nil, []byte("&amp;"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, []byte("&lt;"), nil, []byte("&gt;"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}

// EscapeHTMLByte returns HTML escaped bytes if the given byte should be escaped,
// otherwise nil.
func EscapeHTMLByte(b byte) []byte {
	return htmlEscapeTable[b]
}

// EscapeHTML escapes characters that should be escaped in HTML text.
func EscapeHTML(v []byte) []byte {
	cob := NewCopyOnWriteBuffer(v)
	n := 0
	for i := 0; i < len(v); i++ {
		c := v[i]
		escaped := htmlEscapeTable[c]
		if escaped != nil {
			cob.Write(v[n:i])
			cob.Write(escaped)
			n = i + 1
		}
	}
	if cob.IsCopied() {
		cob.Write(v[n:])
	}
	return cob.Bytes()
}

// UnescapePunctuations unescapes blackslash escaped punctuations.
func UnescapePunctuations(source []byte) []byte {
	cob := NewCopyOnWriteBuffer(source)
	limit := len(source)
	n := 0
	for i := 0; i < limit; {
		c := source[i]
		if i < limit-1 && c == '\\' && IsPunct(source[i+1]) {
			cob.Write(source[n:i])
			cob.WriteByte(source[i+1])
			i += 2
			n = i
			continue
		}
		i++
	}
	if cob.IsCopied() {
		cob.Write(source[n:])
	}
	return cob.Bytes()
}

// ResolveNumericReferences resolve numeric references like '&#1234;" .
func ResolveNumericReferences(source []byte) []byte {
	cob := NewCopyOnWriteBuffer(source)
	buf := make([]byte, 6, 6)
	limit := len(source)
	ok := false
	n := 0
	for i := 0; i < limit; i++ {
		if source[i] == '&' {
			pos := i
			next := i + 1
			if next < limit && source[next] == '#' {
				nnext := next + 1
				if nnext < limit {
					nc := source[nnext]
					// code point like #x22;
					if nnext < limit && nc == 'x' || nc == 'X' {
						start := nnext + 1
						i, ok = ReadWhile(source, [2]int{start, limit}, IsHexDecimal)
						if ok && i < limit && source[i] == ';' {
							v, _ := strconv.ParseUint(BytesToReadOnlyString(source[start:i]), 16, 32)
							cob.Write(source[n:pos])
							n = i + 1
							runeSize := utf8.EncodeRune(buf, ToValidRune(rune(v)))
							cob.Write(buf[:runeSize])
							continue
						}
						// code point like #1234;
					} else if nc >= '0' && nc <= '9' {
						start := nnext
						i, ok = ReadWhile(source, [2]int{start, limit}, IsNumeric)
						if ok && i < limit && i-start < 8 && source[i] == ';' {
							v, _ := strconv.ParseUint(BytesToReadOnlyString(source[start:i]), 0, 32)
							cob.Write(source[n:pos])
							n = i + 1
							runeSize := utf8.EncodeRune(buf, ToValidRune(rune(v)))
							cob.Write(buf[:runeSize])
							continue
						}
					}
				}
			}
			i = next - 1
		}
	}
	if cob.IsCopied() {
		cob.Write(source[n:])
	}
	return cob.Bytes()
}

// ResolveEntityNames resolve entity references like '&ouml;" .
func ResolveEntityNames(source []byte) []byte {
	cob := NewCopyOnWriteBuffer(source)
	limit := len(source)
	ok := false
	n := 0
	for i := 0; i < limit; i++ {
		if source[i] == '&' {
			pos := i
			next := i + 1
			if !(next < limit && source[next] == '#') {
				start := next
				i, ok = ReadWhile(source, [2]int{start, limit}, IsAlphaNumeric)
				if ok && i < limit && source[i] == ';' {
					name := BytesToReadOnlyString(source[start:i])
					entity, ok := LookUpHTML5EntityByName(name)
					if ok {
						cob.Write(source[n:pos])
						n = i + 1
						cob.Write(entity.Characters)
						continue
					}
				}
			}
			i = next - 1
		}
	}
	if cob.IsCopied() {
		cob.Write(source[n:])
	}
	return cob.Bytes()
}

var htmlSpace = []byte("%20")

// URLEscape escape the given URL.
// If resolveReference is set true:
//   1. unescape punctuations
//   2. resolve numeric references
//   3. resolve entity references
//
// URL encoded values (%xx) are kept as is.
func URLEscape(v []byte, resolveReference bool) []byte {
	if resolveReference {
		v = UnescapePunctuations(v)
		v = ResolveNumericReferences(v)
		v = ResolveEntityNames(v)
	}
	cob := NewCopyOnWriteBuffer(v)
	limit := len(v)
	n := 0

	for i := 0; i < limit; {
		c := v[i]
		if urlEscapeTable[c] == 1 {
			i++
			continue
		}
		if c == '%' && i+2 < limit && IsHexDecimal(v[i+1]) && IsHexDecimal(v[i+1]) {
			i += 3
			continue
		}
		u8len := utf8lenTable[c]
		if u8len == 99 { // invalid utf8 leading byte, skip it
			i++
			continue
		}
		if c == ' ' {
			cob.Write(v[n:i])
			cob.Write(htmlSpace)
			i++
			n = i
			continue
		}
		if int(u8len) > len(v) {
			u8len = int8(len(v) - 1)
		}
		if u8len == 0 {
			i++
			n = i
			continue
		}
		cob.Write(v[n:i])
		stop := i + int(u8len)
		if stop > len(v) {
			i++
			n = i
			continue
		}
		cob.Write(StringToReadOnlyBytes(url.QueryEscape(string(v[i:stop]))))
		i += int(u8len)
		n = i
	}
	if cob.IsCopied() && n < limit {
		cob.Write(v[n:])
	}
	return cob.Bytes()
}

// FindURLIndex returns a stop index value if the given bytes seem an URL.
// This function is equivalent to [A-Za-z][A-Za-z0-9.+-]{1,31}:[^<>\x00-\x20]* .
func FindURLIndex(b []byte) int {
	i := 0
	if !(len(b) > 0 && urlTable[b[i]]&7 == 7) {
		return -1
	}
	i++
	for ; i < len(b); i++ {
		c := b[i]
		if urlTable[c]&4 != 4 {
			break
		}
	}
	if i == 1 || i > 33 || i >= len(b) {
		return -1
	}
	if b[i] != ':' {
		return -1
	}
	i++
	for ; i < len(b); i++ {
		c := b[i]
		if urlTable[c]&1 != 1 {
			break
		}
	}
	return i
}

var emailDomainRegexp = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*`)

// FindEmailIndex returns a stop index value if the given bytes seem an email address.
func FindEmailIndex(b []byte) int {
	// TODO: eliminate regexps
	i := 0
	for ; i < len(b); i++ {
		c := b[i]
		if emailTable[c]&1 != 1 {
			break
		}
	}
	if i == 0 {
		return -1
	}
	if i >= len(b) || b[i] != '@' {
		return -1
	}
	i++
	if i >= len(b) {
		return -1
	}
	match := emailDomainRegexp.FindSubmatchIndex(b[i:])
	if match == nil {
		return -1
	}
	return i + match[1]
}

var spaces = []byte(" \t\n\x0b\x0c\x0d")

var spaceTable = [256]int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

var punctTable = [256]int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// a-zA-Z0-9, ;/?:@&=+$,-_.!~*'()#
var urlEscapeTable = [256]int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

var utf8lenTable = [256]int8{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4, 99, 99, 99, 99, 99, 99, 99, 99}

var urlTable = [256]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 5, 1, 5, 5, 1, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 1, 1, 0, 1, 0, 1, 1, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 1, 1, 1, 1, 1, 1, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

var emailTable = [256]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// UTF8Len returns a byte length of the utf-8 character.
func UTF8Len(b byte) int8 {
	return utf8lenTable[b]
}

// IsPunct returns true if the given character is a punctuation, otherwise false.
func IsPunct(c byte) bool {
	return punctTable[c] == 1
}

// IsPunctRune returns true if the given rune is a punctuation, otherwise false.
func IsPunctRune(r rune) bool {
	return int32(r) <= 256 && IsPunct(byte(r)) || unicode.IsPunct(r)
}

// IsSpace returns true if the given character is a space, otherwise false.
func IsSpace(c byte) bool {
	return spaceTable[c] == 1
}

// IsSpaceRune returns true if the given rune is a space, otherwise false.
func IsSpaceRune(r rune) bool {
	return int32(r) <= 256 && IsSpace(byte(r)) || unicode.IsSpace(r)
}

// IsNumeric returns true if the given character is a numeric, otherwise false.
func IsNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}

// IsHexDecimal returns true if the given character is a hexdecimal, otherwise false.
func IsHexDecimal(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

// IsAlphaNumeric returns true if the given character is a alphabet or a numeric, otherwise false.
func IsAlphaNumeric(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9'
}

// IsEastAsianWideRune returns trhe if the given rune is an east asian wide character, otherwise false.
func IsEastAsianWideRune(r rune) bool {
	return unicode.Is(unicode.Hiragana, r) ||
		unicode.Is(unicode.Katakana, r) ||
		unicode.Is(unicode.Han, r) ||
		unicode.Is(unicode.Lm, r) ||
		unicode.Is(unicode.Hangul, r)
}

// A BufWriter is a subset of the bufio.Writer .
type BufWriter interface {
	io.Writer
	Available() int
	Buffered() int
	Flush() error
	WriteByte(c byte) error
	WriteRune(r rune) (size int, err error)
	WriteString(s string) (int, error)
}

// A PrioritizedValue struct holds pair of an arbitrary value and a priority.
type PrioritizedValue struct {
	// Value is an arbitrary value that you want to prioritize.
	Value interface{}
	// Priority is a priority of the value.
	Priority int
}

// PrioritizedSlice is a slice of the PrioritizedValues
type PrioritizedSlice []PrioritizedValue

// Sort sorts the PrioritizedSlice in ascending order.
func (s PrioritizedSlice) Sort() {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Priority < s[j].Priority
	})
}

// Remove removes the given value from this slice.
func (s PrioritizedSlice) Remove(v interface{}) PrioritizedSlice {
	i := 0
	found := false
	for ; i < len(s); i++ {
		if s[i].Value == v {
			found = true
			break
		}
	}
	if !found {
		return s
	}
	return append(s[:i], s[i+1:]...)
}

// Prioritized returns a new PrioritizedValue.
func Prioritized(v interface{}, priority int) PrioritizedValue {
	return PrioritizedValue{v, priority}
}

func bytesHash(b []byte) uint64 {
	var hash uint64 = 5381
	for _, c := range b {
		hash = ((hash << 5) + hash) + uint64(c)
	}
	return hash
}

// BytesFilter is a efficient data structure for checking whether bytes exist or not.
// BytesFilter is thread-safe.
type BytesFilter interface {
	// Add adds given bytes to this set.
	Add([]byte)

	// Contains return true if this set contains given bytes, otherwise false.
	Contains([]byte) bool

	// Extend copies this filter and adds given bytes to new filter.
	Extend(...[]byte) BytesFilter
}

type bytesFilter struct {
	chars     [256]uint8
	threshold int
	slots     [][][]byte
}

// NewBytesFilter returns a new BytesFilter.
func NewBytesFilter(elements ...[]byte) BytesFilter {
	s := &bytesFilter{
		threshold: 3,
		slots:     make([][][]byte, 64),
	}
	for _, element := range elements {
		s.Add(element)
	}
	return s
}

func (s *bytesFilter) Add(b []byte) {
	l := len(b)
	m := s.threshold
	if l < s.threshold {
		m = l
	}
	for i := 0; i < m; i++ {
		s.chars[b[i]] |= 1 << uint8(i)
	}
	h := bytesHash(b) % uint64(len(s.slots))
	slot := s.slots[h]
	if slot == nil {
		slot = [][]byte{}
	}
	s.slots[h] = append(slot, b)
}

func (s *bytesFilter) Extend(bs ...[]byte) BytesFilter {
	newFilter := NewBytesFilter().(*bytesFilter)
	newFilter.chars = s.chars
	newFilter.threshold = s.threshold
	for k, v := range s.slots {
		newSlot := make([][]byte, len(v))
		copy(newSlot, v)
		newFilter.slots[k] = v
	}
	for _, b := range bs {
		newFilter.Add(b)
	}
	return newFilter
}

func (s *bytesFilter) Contains(b []byte) bool {
	l := len(b)
	m := s.threshold
	if l < s.threshold {
		m = l
	}
	for i := 0; i < m; i++ {
		if (s.chars[b[i]] & (1 << uint8(i))) == 0 {
			return false
		}
	}
	h := bytesHash(b) % uint64(len(s.slots))
	slot := s.slots[h]
	if slot == nil || len(slot) == 0 {
		return false
	}
	for _, element := range slot {
		if bytes.Equal(element, b) {
			return true
		}
	}
	return false
}
