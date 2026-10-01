package terminfo

import (
	"slices"
)

const (
	// maxFileLength is the max file length.
	maxFileLength = 32768
	// magic is the file magic for terminfo files.
	magic = 0o432
	// magicExtended is the file magic for terminfo files with the extended
	// number format.
	magicExtended = 0o1036
)

// header fields.
const (
	fieldMagic = iota
	fieldNameSize
	fieldBoolCount
	fieldNumCount
	fieldStringCount
	fieldTableSize
)

// header extended fields.
const (
	fieldExtBoolCount = iota
	fieldExtNumCount
	fieldExtStringCount
	fieldExtOffsetUsed
	fieldExtTableSize
)

// hasInvalidCaps determines if the capabilities in h are invalid.
func hasInvalidCaps(h []int) bool {
	return h[fieldBoolCount] > CapCountBool ||
		h[fieldNumCount] > CapCountNum ||
		h[fieldStringCount] > CapCountString
}

// capLength returns the total length of the capabilities in bytes.
func capLength(h []int) int {
	return h[fieldNameSize] +
		h[fieldBoolCount] +
		(h[fieldNameSize]+h[fieldBoolCount])%2 + // account for word align
		h[fieldNumCount]*2 +
		h[fieldStringCount]*2 +
		h[fieldTableSize]
}

// extOffsetCount returns the number of offsets stored in the extended string
// data table index.
//
// An offset is written for every extended string cap, followed by one for the
// name of every extended bool, num and string cap.
func extOffsetCount(h []int) int {
	return h[fieldExtBoolCount] +
		h[fieldExtNumCount] +
		h[fieldExtStringCount]*2
}

// hasInvalidExtOffset determines if the extended offset field is valid.
//
// The field holds the number of offsets that are actually used (ie, that point
// into the string data table), which excludes any absent or cancelled extended
// string cap, and so is only ever less than or equal to the number of offsets
// present in the file.
func hasInvalidExtOffset(h []int) bool {
	return h[fieldExtBoolCount] < 0 ||
		h[fieldExtNumCount] < 0 ||
		h[fieldExtStringCount] < 0 ||
		h[fieldExtOffsetUsed] < 0 ||
		h[fieldExtOffsetUsed] > extOffsetCount(h)
}

// extCapLength returns the total length of extended capabilities in bytes.
func extCapLength(h []int, numWidth int) int {
	return h[fieldExtBoolCount] +
		h[fieldExtBoolCount]%2 + // account for word align
		h[fieldExtNumCount]*(numWidth/8) +
		extOffsetCount(h)*2 +
		h[fieldExtTableSize]
}

// findNull finds the position of null in buf.
func findNull(buf []byte, i int) int {
	for ; i < len(buf); i++ {
		if buf[i] == 0 {
			return i
		}
	}
	return -1
}

// readStrings decodes n strings from string data table buf using the indexes in
// idx, returning the decoded strings, the cancelled strings, and the offset of
// the end of the last decoded string.
func readStrings(idx []int, buf []byte, n int) (map[int][]byte, map[int]bool, int, error) {
	var last int
	m, mm := make(map[int][]byte), make(map[int]bool)
	for i := range n {
		start := idx[i]
		if start == -2 {
			mm[i] = true
			continue
		}
		if start < 0 {
			continue
		}
		if end := findNull(buf, start); end != -1 {
			m[i], last = buf[start:end], end+1
		} else {
			return nil, nil, 0, ErrInvalidStringTable
		}
	}
	return m, mm, last, nil
}

// decoder holds state info while decoding a terminfo file.
type decoder struct {
	buf []byte
	pos int
	n   int
}

// readBytes reads the next n bytes of buf, incrementing pos by n.
func (d *decoder) readBytes(n int) ([]byte, error) {
	if d.n < d.pos+n {
		return nil, ErrUnexpectedFileEnd
	}
	n, d.pos = d.pos, d.pos+n
	return d.buf[n:d.pos], nil
}

// readInts reads n number of ints with width w.
func (d *decoder) readInts(n, w int) ([]int, error) {
	w /= 8
	l := n * w
	buf, err := d.readBytes(l)
	if err != nil {
		return nil, err
	}
	// align
	d.pos += d.pos % 2
	z := make([]int, n)
	for i, j := 0, 0; i < l; i, j = i+w, j+1 {
		switch w {
		case 1:
			z[j] = int(buf[i])
		case 2:
			z[j] = int(int16(uint16(buf[i+1])<<8 | uint16(buf[i])))
		case 4:
			z[j] = int(int32(uint32(buf[i+3])<<24 | uint32(buf[i+2])<<16 | uint32(buf[i+1])<<8 | uint32(buf[i])))
		}
	}
	return z, nil
}

// readBools reads the next n bools.
func (d *decoder) readBools(n int) (map[int]bool, map[int]bool, error) {
	buf, err := d.readInts(n, 8)
	if err != nil {
		return nil, nil, err
	}
	// process
	bools, boolsM := make(map[int]bool), make(map[int]bool)
	for i, b := range buf {
		bools[i] = b == 1
		if int8(b) == -2 {
			boolsM[i] = true
		}
	}
	return bools, boolsM, nil
}

// readNums reads the next n nums.
func (d *decoder) readNums(n, w int) (map[int]int, map[int]bool, error) {
	buf, err := d.readInts(n, w)
	if err != nil {
		return nil, nil, err
	}
	// process
	nums, numsM := make(map[int]int), make(map[int]bool)
	for i := range n {
		nums[i] = buf[i]
		if buf[i] == -2 {
			numsM[i] = true
		}
	}
	return nums, numsM, nil
}

// readStringTable reads the string data for n strings and the accompanying data
// table of length sz.
func (d *decoder) readStringTable(n, sz int) ([][]byte, []int, error) {
	buf, err := d.readInts(n, 16)
	if err != nil {
		return nil, nil, err
	}
	// read string data table
	data, err := d.readBytes(sz)
	if err != nil {
		return nil, nil, err
	}
	// align
	d.pos += d.pos % 2
	// process
	s := make([][]byte, n)
	var m []int
	for i := range n {
		start := buf[i]
		if start == -2 {
			m = append(m, i)
		} else if start >= 0 {
			if end := findNull(data, start); end != -1 {
				s[i] = data[start:end]
			} else {
				return nil, nil, ErrInvalidStringTable
			}
		}
	}
	return s, m, nil
}

// readStrings reads the next n strings and processes the string data table of
// length sz.
func (d *decoder) readStrings(n, sz int) (map[int][]byte, map[int]bool, error) {
	s, m, err := d.readStringTable(n, sz)
	if err != nil {
		return nil, nil, err
	}
	strs := make(map[int][]byte)
	for k, v := range s {
		if k == AcsChars {
			v = canonicalizeAscChars(v)
		}
		strs[k] = v
	}
	strsM := make(map[int]bool, len(m))
	for _, k := range m {
		strsM[k] = true
	}
	return strs, strsM, nil
}

// canonicalizeAscChars reorders chars to be unique, in order.
//
// see repair_ascc in ncurses-6.3/progs/dump_entry.c
func canonicalizeAscChars(z []byte) []byte {
	var c []byte
	enc := make(map[byte]byte, len(z)/2)
	for i := 0; i < len(z); i += 2 {
		if _, ok := enc[z[i]]; !ok {
			a, b := z[i], z[i+1]
			// log.Printf(">>> a: %d %c, b: %d %c", a, a, b, b)
			c, enc[a] = append(c, a), b
		}
	}
	slices.Sort(c)
	r := make([]byte, 2*len(c))
	for i := 0; i < len(c); i++ {
		r[i*2], r[i*2+1] = c[i], enc[c[i]]
	}
	return r
}
