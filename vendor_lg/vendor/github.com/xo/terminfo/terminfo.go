// Package terminfo implements reading terminfo files in pure go.
package terminfo

import (
	"io"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
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

	// ErrDatabaseDirectoryNotFound is the database directory not found error.
	ErrDatabaseDirectoryNotFound Error = "database directory not found"

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

	// ExtBoolsNames is the map of extended bool capabilities to their index.
	ExtBoolNames map[int][]byte

	// ExtNums are the extended num capabilities.
	ExtNums map[int]int

	// ExtNumsNames is the map of extended num capabilities to their index.
	ExtNumNames map[int][]byte

	// ExtStrings are the extended string capabilities.
	ExtStrings map[int][]byte

	// ExtStringsNames is the map of extended string capabilities to their index.
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
		len: len(buf),
	}

	// read header
	h, err := d.readInts(6, 16)
	if err != nil {
		return nil, err
	}

	var numWidth int

	// check magic
	if h[fieldMagic] == magic {
		numWidth = 16
	} else if h[fieldMagic] == magicExtended {
		numWidth = 32
	} else {
		return nil, ErrInvalidMagic
	}

	// check header
	if hasInvalidCaps(h) {
		return nil, ErrInvalidHeader
	}

	// check remaining length
	if d.len-d.pos < capLength(h) {
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
	if d.pos >= d.len {
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
	if d.len-d.pos != extCapLength(eh, numWidth) {
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
	extIndexes, err := d.readInts(eh[fieldExtOffsetCount], 16)
	if err != nil {
		return nil, err
	}

	// read string data table
	extData, err := d.readBytes(eh[fieldExtTableSize])
	if err != nil {
		return nil, err
	}

	// precautionary check that exactly at end of file
	if d.pos != d.len {
		return nil, ErrUnexpectedFileEnd
	}

	var last int
	// read extended string caps
	ti.ExtStrings, last, err = readStrings(extIndexes, extData, eh[fieldExtStringCount])
	if err != nil {
		return nil, err
	}
	extIndexes, extData = extIndexes[eh[fieldExtStringCount]:], extData[last:]

	// read extended bool names
	ti.ExtBoolNames, _, err = readStrings(extIndexes, extData, eh[fieldExtBoolCount])
	if err != nil {
		return nil, err
	}
	extIndexes = extIndexes[eh[fieldExtBoolCount]:]

	// read extended num names
	ti.ExtNumNames, _, err = readStrings(extIndexes, extData, eh[fieldExtNumCount])
	if err != nil {
		return nil, err
	}
	extIndexes = extIndexes[eh[fieldExtNumCount]:]

	// read extended string names
	ti.ExtStringNames, _, err = readStrings(extIndexes, extData, eh[fieldExtStringCount])
	if err != nil {
		return nil, err
	}
	//extIndexes = extIndexes[eh[fieldExtStringCount]:]

	return ti, nil
}

// Open reads the terminfo file name from the specified directory dir.
func Open(dir, name string) (*Terminfo, error) {
	var err error
	var buf []byte
	var filename string
	for _, f := range []string{
		path.Join(dir, name[0:1], name),
		path.Join(dir, strconv.FormatUint(uint64(name[0]), 16), name),
	} {
		buf, err = ioutil.ReadFile(f)
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

// stringCaps returns all string and extended capabilities using f to format the
// index key.
func (ti *Terminfo) stringCaps(f func(int) string, extended bool) map[string][]byte {
	m := make(map[string][]byte, len(ti.Strings)+len(ti.ExtStrings))
	if !extended {
		for k, v := range ti.Strings {
			m[f(k)] = v
		}
	} else {
		for k, v := range ti.ExtStrings {
			m[string(ti.ExtStringNames[k])] = v
		}
	}
	return m
}

// StringCaps returns all string capabilities.
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
func (ti *Terminfo) Printf(i int, v ...interface{}) string {
	return Printf(ti.Strings[i], v...)
}

// Fprintf prints the string cap i to writer w, interpolating parameters v.
func (ti *Terminfo) Fprintf(w io.Writer, i int, v ...interface{}) {
	Fprintf(w, ti.Strings[i], v...)
}

// Color takes a foreground and background color and returns string that sets
// them for this terminal.
func (ti *Terminfo) Colorf(fg, bg int, str string) string {
	maxColors := int(ti.Nums[MaxColors])

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

// Puts emits the string to the writer, but expands inline padding indications
// (of the form $<[delay]> where [delay] is msec) to a suitable number of
// padding characters (usually null bytes) based upon the supplied baud. At
// high baud rates, more padding characters will be inserted.
/*func (ti *Terminfo) Puts(w io.Writer, s string, lines, baud int) (int, error) {
	var err error
	for {
		start := strings.Index(s, "$<")
		if start == -1 {
			// most strings don't need padding, which is good news!
			return io.WriteString(w, s)
		}

		end := strings.Index(s, ">")
		if end == -1 {
			// unterminated... just emit bytes unadulterated.
			return io.WriteString(w, "$<"+s)
		}

		var c int
		c, err = io.WriteString(w, s[:start])
		if err != nil {
			return n + c, err
		}
		n += c

		s = s[start+2:]
		val := s[:end]
		s = s[end+1:]
		var ms int
		var dot, mandatory, asterisk bool
		unit := 1000
		for _, ch := range val {
			switch {
			case ch >= '0' && ch <= '9':
				ms = (ms * 10) + int(ch-'0')
				if dot {
					unit *= 10
				}
			case ch == '.' && !dot:
				dot = true
			case ch == '*' && !asterisk:
				ms *= lines
				asterisk = true
			case ch == '/':
				mandatory = true
			default:
				break
			}
		}

		z, pad := ((baud/8)/unit)*ms, ti.Strings[PadChar]
		b := make([]byte, len(pad)*z)
		for bp := copy(b, pad); bp < len(b); bp *= 2 {
			copy(b[bp:], b[:bp])
		}

		if (!ti.Bools[XonXoff] && baud > int(ti.Nums[PaddingBaudRate])) || mandatory {
			c, err = w.Write(b)
			if err != nil {
				return n + c, err
			}
			n += c
		}
	}

	return n, nil
}*/
