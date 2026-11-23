// Package terminfo implements reading terminfo files in pure go.
package terminfo

//go:generate go run gen.go

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Error is a terminfo error.
type Error string

// Error satisfies the error interface.
func (err Error) Error() string {
	return string(err)
}

const (
	// ErrInvalidFileSize is the invalid file size error.
	ErrInvalidFileSize Error = "invalid file size"
	// ErrUnexpectedFileEnd is the unexpected file end error.
	ErrUnexpectedFileEnd Error = "unexpected file end"
	// ErrInvalidStringTable is the invalid string table error.
	ErrInvalidStringTable Error = "invalid string table"
	// ErrInvalidMagic is the invalid magic error.
	ErrInvalidMagic Error = "invalid magic"
	// ErrInvalidHeader is the invalid header error.
	ErrInvalidHeader Error = "invalid header"
	// ErrInvalidNames is the invalid names error.
	ErrInvalidNames Error = "invalid names"
	// ErrInvalidExtendedHeader is the invalid extended header error.
	ErrInvalidExtendedHeader Error = "invalid extended header"
	// ErrEmptyTermName is the empty term name error.
	ErrEmptyTermName Error = "empty term name"
	// ErrFileNotFound is the file not found error.
	ErrFileNotFound Error = "file not found"
	// ErrInvalidTermProgramVersion is the invalid TERM_PROGRAM_VERSION error.
	ErrInvalidTermProgramVersion Error = "invalid TERM_PROGRAM_VERSION"
)

// Terminfo describes a terminal's capabilities.
type Terminfo struct {
	// File is the original source file.
	File string
	// Names are the provided cap names.
	Names []string
	// Bools are the bool capabilities.
	Bools map[int]bool
	// BoolsM are the missing bool capabilities.
	BoolsM map[int]bool
	// Nums are the num capabilities.
	Nums map[int]int
	// NumsM are the missing num capabilities.
	NumsM map[int]bool
	// Strings are the string capabilities.
	Strings map[int][]byte
	// StringsM are the missing string capabilities.
	StringsM map[int]bool
	// ExtBools are the extended bool capabilities.
	ExtBools map[int]bool
	// ExtBoolNames is the map of extended bool capabilities to their index.
	ExtBoolNames map[int][]byte
	// ExtNums are the extended num capabilities.
	ExtNums map[int]int
	// ExtNumNames is the map of extended num capabilities to their index.
	ExtNumNames map[int][]byte
	// ExtStrings are the extended string capabilities.
	ExtStrings map[int][]byte
	// ExtStringsM are the missing extended string capabilities.
	ExtStringsM map[int]bool
	// ExtStringNames is the map of extended string capabilities to their index.
	ExtStringNames map[int][]byte
}

// Decode decodes the terminfo data contained in buf.
func Decode(buf []byte) (*Terminfo, error) {
	var err error
	// check max file length
	if len(buf) >= maxFileLength {
		return nil, ErrInvalidFileSize
	}
	d := &decoder{
		buf: buf,
		n:   len(buf),
	}
	// read header
	h, err := d.readInts(6, 16)
	if err != nil {
		return nil, err
	}
	var numWidth int
	// check magic
	switch h[fieldMagic] {
	case magic:
		numWidth = 16
	case magicExtended:
		numWidth = 32
	default:
		return nil, ErrInvalidMagic
	}
	// check header
	if hasInvalidCaps(h) {
		return nil, ErrInvalidHeader
	}
	// check remaining length
	if d.n-d.pos < capLength(h) {
		return nil, ErrUnexpectedFileEnd
	}
	// read names
	names, err := d.readBytes(h[fieldNameSize])
	if err != nil {
		return nil, err
	}
	// check name is terminated properly
	i := findNull(names, 0)
	if i == -1 {
		return nil, ErrInvalidNames
	}
	names = names[:i]
	// read bool caps
	bools, boolsM, err := d.readBools(h[fieldBoolCount])
	if err != nil {
		return nil, err
	}
	// read num caps
	nums, numsM, err := d.readNums(h[fieldNumCount], numWidth)
	if err != nil {
		return nil, err
	}
	// read string caps
	strs, strsM, err := d.readStrings(h[fieldStringCount], h[fieldTableSize])
	if err != nil {
		return nil, err
	}
	ti := &Terminfo{
		Names:    strings.Split(string(names), "|"),
		Bools:    bools,
		BoolsM:   boolsM,
		Nums:     nums,
		NumsM:    numsM,
		Strings:  strs,
		StringsM: strsM,
	}
	// at the end of file, so no extended caps
	if d.pos >= d.n {
		return ti, nil
	}
	// decode extended header
	eh, err := d.readInts(5, 16)
	if err != nil {
		return nil, err
	}
	// check extended offset field
	if hasInvalidExtOffset(eh) {
		return nil, ErrInvalidExtendedHeader
	}
	// check extended cap lengths
	if d.n-d.pos != extCapLength(eh, numWidth) {
		return nil, ErrInvalidExtendedHeader
	}
	// read extended bool caps
	ti.ExtBools, _, err = d.readBools(eh[fieldExtBoolCount])
	if err != nil {
		return nil, err
	}
	// read extended num caps
	ti.ExtNums, _, err = d.readNums(eh[fieldExtNumCount], numWidth)
	if err != nil {
		return nil, err
	}
	// read extended string data table indexes
	extIndexes, err := d.readInts(extOffsetCount(eh), 16)
	if err != nil {
		return nil, err
	}
	// read string data table
	extData, err := d.readBytes(eh[fieldExtTableSize])
	if err != nil {
		return nil, err
	}
	// precautionary check that exactly at end of file
	if d.pos != d.n {
		return nil, ErrUnexpectedFileEnd
	}
	var last int
	// read extended string caps
	ti.ExtStrings, ti.ExtStringsM, last, err = readStrings(extIndexes, extData, eh[fieldExtStringCount])
	if err != nil {
		return nil, err
	}
	extIndexes, extData = extIndexes[eh[fieldExtStringCount]:], extData[last:]
	// read extended bool names
	ti.ExtBoolNames, _, _, err = readStrings(extIndexes, extData, eh[fieldExtBoolCount])
	if err != nil {
		return nil, err
	}
	extIndexes = extIndexes[eh[fieldExtBoolCount]:]
	// read extended num names
	ti.ExtNumNames, _, _, err = readStrings(extIndexes, extData, eh[fieldExtNumCount])
	if err != nil {
		return nil, err
	}
	extIndexes = extIndexes[eh[fieldExtNumCount]:]
	// read extended string names
	ti.ExtStringNames, _, _, err = readStrings(extIndexes, extData, eh[fieldExtStringCount])
	if err != nil {
		return nil, err
	}
	return ti, nil
}

// Open reads the terminfo file name from the specified directory dir. Open
// adds the entry to the cache of Load, so a caller must not change the
// returned value, as for Load.
func Open(dir, name string) (*Terminfo, error) {
	var err error
	var buf []byte
	var filename string
	for _, f := range []string{
		filepath.Join(dir, name[0:1], name),
		filepath.Join(dir, strconv.FormatUint(uint64(name[0]), 16), name),
	} {
		buf, err = os.ReadFile(f)
		if err == nil {
			filename = f
			break
		}
	}
	if buf == nil {
		return nil, ErrFileNotFound
	}
	// decode
	ti, err := Decode(buf)
	if err != nil {
		return nil, err
	}
	// save original file name
	ti.File = filename
	// add to cache
	termCache.Lock()
	for _, n := range ti.Names {
		termCache.db[n] = ti
	}
	termCache.Unlock()
	return ti, nil
}

// BoolCaps returns all bool capabilities.
func (ti *Terminfo) BoolCaps() map[string]bool {
	return ti.boolCaps(BoolCapName, false)
}

// BoolCapsShort returns all bool capabilities, using the short name as the
// index.
func (ti *Terminfo) BoolCapsShort() map[string]bool {
	return ti.boolCaps(BoolCapNameShort, false)
}

// ExtBoolCaps returns all extended bool capabilities.
func (ti *Terminfo) ExtBoolCaps() map[string]bool {
	return ti.boolCaps(BoolCapName, true)
}

// ExtBoolCapsShort returns all extended bool capabilities, using the short
// name as the index.
func (ti *Terminfo) ExtBoolCapsShort() map[string]bool {
	return ti.boolCaps(BoolCapNameShort, true)
}

// NumCaps returns all num capabilities.
func (ti *Terminfo) NumCaps() map[string]int {
	return ti.numCaps(NumCapName, false)
}

// NumCapsShort returns all num capabilities, using the short name as the
// index.
func (ti *Terminfo) NumCapsShort() map[string]int {
	return ti.numCaps(NumCapNameShort, false)
}

// ExtNumCaps returns all extended num capabilities.
func (ti *Terminfo) ExtNumCaps() map[string]int {
	return ti.numCaps(NumCapName, true)
}

// ExtNumCapsShort returns all extended num capabilities, using the short
// name as the index.
func (ti *Terminfo) ExtNumCapsShort() map[string]int {
	return ti.numCaps(NumCapNameShort, true)
}

// StringCaps returns all string capabilities. It leaves out each absent
// capability.
func (ti *Terminfo) StringCaps() map[string][]byte {
	return ti.stringCaps(StringCapName, false)
}

// StringCapsShort returns all string capabilities, using the short name as the
// index.
func (ti *Terminfo) StringCapsShort() map[string][]byte {
	return ti.stringCaps(StringCapNameShort, false)
}

// ExtStringCaps returns all extended string capabilities.
func (ti *Terminfo) ExtStringCaps() map[string][]byte {
	return ti.stringCaps(StringCapName, true)
}

// ExtStringCapsShort returns all extended string capabilities, using the short
// name as the index.
func (ti *Terminfo) ExtStringCapsShort() map[string][]byte {
	return ti.stringCaps(StringCapNameShort, true)
}

// Has determines if the bool cap i is present.
func (ti *Terminfo) Has(i int) bool {
	return ti.Bools[i]
}

// Num returns the num cap i, or -1 if not present.
func (ti *Terminfo) Num(i int) int {
	n, ok := ti.Nums[i]
	if !ok {
		return -1
	}
	return n
}

// Printf formats the string cap i, interpolating parameters v.
func (ti *Terminfo) Printf(i int, v ...any) string {
	return Printf(ti.Strings[i], v...)
}

// Fprintf prints the string cap i to writer w, interpolating parameters v.
func (ti *Terminfo) Fprintf(w io.Writer, i int, v ...any) {
	Fprintf(w, ti.Strings[i], v...)
}

// Colorf returns str with the escape sequences that set the foreground color
// fg and the background color bg before it, and reset them after it.
func (ti *Terminfo) Colorf(fg, bg int, str string) string {
	maxColors := ti.Nums[MaxColors]
	// map bright colors to lower versions if the color table only holds 8.
	if maxColors == 8 {
		if fg > 7 && fg < 16 {
			fg -= 8
		}
		if bg > 7 && bg < 16 {
			bg -= 8
		}
	}
	var s string
	if maxColors > fg && fg >= 0 {
		s += ti.Printf(SetAForeground, fg)
	}
	if maxColors > bg && bg >= 0 {
		s += ti.Printf(SetABackground, bg)
	}
	return s + str + ti.Printf(ExitAttributeMode)
}

// Goto returns a string suitable for addressing the cursor at the given
// row and column. The origin 0, 0 is in the upper left corner of the screen.
func (ti *Terminfo) Goto(row, col int) string {
	return Printf(ti.Strings[CursorAddress], row, col)
}

// Puts writes s to w, and applies each padding marker in s, such as the
// $<100/> in the flash capability of xterm. It follows tputs in ncurses.
//
// A marker holds a delay in milliseconds, with one optional decimal place. A
// '*' after the number multiplies the delay by lines, the count of lines that
// the output affects. A '/' after the number makes the delay mandatory. Puts
// applies a delay when it is mandatory, when s is the bell or the flash
// capability, or when baud is at least the padding baud rate and the terminal
// has no xon/xoff flow control. Puts writes no marker to w.
//
// Puts writes a delay as pad characters, and the baud rate sets how many. If
// the terminal has no pad character, Puts sleeps for the delay instead. Puts
// limits each delay to 30 seconds.
func (ti *Terminfo) Puts(w io.Writer, s string, lines, baud int) (int, error) {
	always := ti.isString(Bell, s) || ti.isString(FlashScreen, s)
	normal := !ti.Bools[XonXoff] && ti.Nums[PaddingBaudRate] != 0 && baud >= ti.Nums[PaddingBaudRate]
	var buf []byte
	var n int
	for i := 0; i < len(s); i++ {
		if s[i] != '$' || i+1 == len(s) || s[i+1] != '<' {
			buf = append(buf, s[i])
			continue
		}
		// a marker starts with a digit or a '.', and ends with a '>'
		j := i + 2
		if j == len(s) || (!isDigit(s[j]) && s[j] != '.') || !strings.Contains(s[j:], ">") {
			buf = append(buf, '$', '<')
			i++
			continue
		}
		d, end, mandatory := readDelay(s, j, lines)
		// skip the '>', as ncurses skips one byte
		i = end
		if d <= 0 || (!always && !normal && !mandatory) {
			continue
		}
		ms := min(d/10, maxDelay)
		if ti.Bools[NoPadChar] {
			c, err := w.Write(buf)
			n, buf = n+c, buf[:0]
			if err != nil {
				return n, err
			}
			time.Sleep(time.Duration(ms) * time.Millisecond)
			continue
		}
		pad := byte(0)
		if p := ti.Strings[PadChar]; len(p) != 0 {
			pad = p[0]
		}
		buf = append(buf, bytes.Repeat([]byte{pad}, max(ms*baud/(baudByte*1000), 0))...)
	}
	c, err := w.Write(buf)
	return n + c, err
}

// boolCaps returns all bool and extended capabilities using f to format the
// index key.
func (ti *Terminfo) boolCaps(f func(int) string, extended bool) map[string]bool {
	m := make(map[string]bool, len(ti.Bools)+len(ti.ExtBools))
	if !extended {
		for k, v := range ti.Bools {
			m[f(k)] = v
		}
	} else {
		for k, v := range ti.ExtBools {
			m[string(ti.ExtBoolNames[k])] = v
		}
	}
	return m
}

// numCaps returns all num and extended capabilities using f to format the
// index key.
func (ti *Terminfo) numCaps(f func(int) string, extended bool) map[string]int {
	m := make(map[string]int, len(ti.Nums)+len(ti.ExtNums))
	if !extended {
		for k, v := range ti.Nums {
			m[f(k)] = v
		}
	} else {
		for k, v := range ti.ExtNums {
			m[string(ti.ExtNumNames[k])] = v
		}
	}
	return m
}

// stringCaps returns all string and extended capabilities using f to format the
// index key. It leaves out each absent capability, which has a nil value.
func (ti *Terminfo) stringCaps(f func(int) string, extended bool) map[string][]byte {
	m := make(map[string][]byte, len(ti.Strings)+len(ti.ExtStrings))
	if !extended {
		for k, v := range ti.Strings {
			if v != nil {
				m[f(k)] = v
			}
		}
	} else {
		for k, v := range ti.ExtStrings {
			if v != nil {
				m[string(ti.ExtStringNames[k])] = v
			}
		}
	}
	return m
}

// readDelay reads the padding marker in s that starts at j, after the "$<".
// It returns the delay in tenths of a millisecond, the index of the byte after
// the delay, and whether the delay is mandatory.
func readDelay(s string, j, lines int) (int, int, bool) {
	var d int
	for ; j < len(s) && isDigit(s[j]); j++ {
		d = d*10 + int(s[j]-'0')
	}
	d *= 10
	if j < len(s) && s[j] == '.' {
		if j++; j < len(s) && isDigit(s[j]) {
			d += int(s[j] - '0')
			j++
		}
		for j < len(s) && isDigit(s[j]) {
			j++
		}
	}
	var mandatory bool
	for ; j < len(s) && (s[j] == '*' || s[j] == '/'); j++ {
		if s[j] == '*' {
			d *= lines
		} else {
			mandatory = true
		}
	}
	return d, j, mandatory
}

// maxDelay is the longest delay that Puts applies, in milliseconds. It is
// MAX_DELAY_MSECS in ncurses.
const maxDelay = 30000

// baudByte is the count of bits that one byte takes on the line: 7 bits, 1
// parity bit and 1 stop bit. It is BAUDBYTE in ncurses.
const baudByte = 9

// isString reports whether the string capability i is present and equal to s.
func (ti *Terminfo) isString(i int, s string) bool {
	v, ok := ti.Strings[i]
	return ok && v != nil && string(v) == s
}

// isDigit reports whether c is an ASCII digit.
func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
