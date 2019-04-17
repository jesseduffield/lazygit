// +build !js,!appengine

package runewidth

import (
	"crypto/sha256"
	"fmt"
	"os"
	"sort"
	"testing"
	"unicode/utf8"
)

var _ sort.Interface = (*table)(nil)

func init() {
	os.Setenv("RUNEWIDTH_EASTASIAN", "")
	handleEnv()
}

func (t table) Len() int {
	return len(t)
}

func (t table) Less(i, j int) bool {
	return t[i].first < t[j].first
}

func (t *table) Swap(i, j int) {
	(*t)[i], (*t)[j] = (*t)[j], (*t)[i]
}

var tables = []table{
	private,
	nonprint,
	combining,
	doublewidth,
	ambiguous,
	emoji,
	notassigned,
	neutral,
}

func TestTableChecksums(t *testing.T) {
	check := func(tbl table, wantN int, wantSHA string) {
		gotN := 0
		buf := make([]byte, utf8.MaxRune+1)
		for r := rune(0); r <= utf8.MaxRune; r++ {
			if inTable(r, tbl) {
				gotN++
				buf[r] = 1
			}
		}
		gotSHA := fmt.Sprintf("%x", sha256.Sum256(buf))
		if gotN != wantN || gotSHA != wantSHA {
			t.Errorf("n = %d want %d, sha256 = %s want %s", gotN, wantN, gotSHA, wantSHA)
		}
	}

	check(private, 137468, "a4a641206dc8c5de80bd9f03515a54a706a5a4904c7684dc6a33d65c967a51b2")
	check(nonprint, 2143, "288904683eb225e7c4c0bd3ee481b53e8dace404ec31d443afdbc4d13729fe95")
	check(combining, 2097, "b1dabe5f35b7ccf868999bf6df6134f346ae14a4eb16f22e1dc8a98240ba1b53")
	check(doublewidth, 180993, "06f5d5d5ebb8b9ee74fdf6003ecfbb313f9c042eb3cb4fce2a9e06089eb68dda")
	check(ambiguous, 138739, "d05e339a10f296de6547ff3d6c5aee32f627f6555477afebd4a3b7e3cf74c9e3")
	check(emoji, 1236, "9b2d75cf8ca48c5075c525a92ce5cf2608fa451c589f33d7d153e9df93f4e2f7")
	check(notassigned, 846357, "b06b7acc03725de394d92b09306aa7a9c0c0b53f36884db4c835cbb04971e421")
	check(neutral, 25561, "87fffca79a3a6d413d23adf1c591bdcc1ea5d906d0d466b12a76357bbbb74607")
}

func isCompact(t *testing.T, tbl table) bool {
	for i := range tbl {
		if tbl[i].last < tbl[i].first { // sanity check
			t.Errorf("table invalid: %v", tbl[i])
			return false
		}
		if i+1 < len(tbl) && tbl[i].last+1 >= tbl[i+1].first { // can be combined into one entry
			t.Errorf("table not compact: %v %v", tbl[i-1], tbl[i])
			return false
		}
	}
	return true
}

// This is a utility function in case that a table has changed.
func printCompactTable(tbl table) {
	counter := 0
	printEntry := func(first, last rune) {
		if counter%3 == 0 {
			fmt.Printf("\t")
		}
		fmt.Printf("{0x%04X, 0x%04X},", first, last)
		if (counter+1)%3 == 0 {
			fmt.Printf("\n")
		} else {
			fmt.Printf(" ")
		}
		counter++
	}

	sort.Sort(&tbl) // just in case
	first := rune(-1)
	for i := range tbl {
		if first < 0 {
			first = tbl[i].first
		}
		if i+1 < len(tbl) && tbl[i].last+1 >= tbl[i+1].first { // can be combined into one entry
			continue
		}
		printEntry(first, tbl[i].last)
		first = -1
	}
	fmt.Printf("\n\n")
}

func TestSorted(t *testing.T) {
	for _, tbl := range tables {
		if !sort.IsSorted(&tbl) {
			t.Errorf("table not sorted")
		}
		if !isCompact(t, tbl) {
			t.Errorf("table not compact")
			// printCompactTable(tbl)
		}
	}
}

var runewidthtests = []struct {
	in    rune
	out   int
	eaout int
}{
	{'世', 2, 2},
	{'界', 2, 2},
	{'ｾ', 1, 1},
	{'ｶ', 1, 1},
	{'ｲ', 1, 1},
	{'☆', 1, 2}, // double width in ambiguous
	{'\x00', 0, 0},
	{'\x01', 0, 0},
	{'\u0300', 0, 0},
	{'\u2028', 0, 0},
	{'\u2029', 0, 0},
}

func TestRuneWidth(t *testing.T) {
	c := NewCondition()
	c.EastAsianWidth = false
	for _, tt := range runewidthtests {
		if out := c.RuneWidth(tt.in); out != tt.out {
			t.Errorf("RuneWidth(%q) = %d, want %d", tt.in, out, tt.out)
		}
	}
	c.EastAsianWidth = true
	for _, tt := range runewidthtests {
		if out := c.RuneWidth(tt.in); out != tt.eaout {
			t.Errorf("RuneWidth(%q) = %d, want %d", tt.in, out, tt.eaout)
		}
	}
}

var isambiguouswidthtests = []struct {
	in  rune
	out bool
}{
	{'世', false},
	{'■', true},
	{'界', false},
	{'○', true},
	{'㈱', false},
	{'①', true},
	{'②', true},
	{'③', true},
	{'④', true},
	{'⑤', true},
	{'⑥', true},
	{'⑦', true},
	{'⑧', true},
	{'⑨', true},
	{'⑩', true},
	{'⑪', true},
	{'⑫', true},
	{'⑬', true},
	{'⑭', true},
	{'⑮', true},
	{'⑯', true},
	{'⑰', true},
	{'⑱', true},
	{'⑲', true},
	{'⑳', true},
	{'☆', true},
}

func TestIsAmbiguousWidth(t *testing.T) {
	for _, tt := range isambiguouswidthtests {
		if out := IsAmbiguousWidth(tt.in); out != tt.out {
			t.Errorf("IsAmbiguousWidth(%q) = %v, want %v", tt.in, out, tt.out)
		}
	}
}

var stringwidthtests = []struct {
	in    string
	out   int
	eaout int
}{
	{"■㈱の世界①", 10, 12},
	{"スター☆", 7, 8},
	{"つのだ☆HIRO", 11, 12},
}

func TestStringWidth(t *testing.T) {
	c := NewCondition()
	c.EastAsianWidth = false
	for _, tt := range stringwidthtests {
		if out := c.StringWidth(tt.in); out != tt.out {
			t.Errorf("StringWidth(%q) = %d, want %d", tt.in, out, tt.out)
		}
	}
	c.EastAsianWidth = true
	for _, tt := range stringwidthtests {
		if out := c.StringWidth(tt.in); out != tt.eaout {
			t.Errorf("StringWidth(%q) = %d, want %d", tt.in, out, tt.eaout)
		}
	}
}

func TestStringWidthInvalid(t *testing.T) {
	s := "こんにちわ\x00世界"
	if out := StringWidth(s); out != 14 {
		t.Errorf("StringWidth(%q) = %d, want %d", s, out, 14)
	}
}

func TestTruncateSmaller(t *testing.T) {
	s := "あいうえお"
	expected := "あいうえお"

	if out := Truncate(s, 10, "..."); out != expected {
		t.Errorf("Truncate(%q) = %q, want %q", s, out, expected)
	}
}

func TestTruncate(t *testing.T) {
	s := "あいうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおおおおお"
	expected := "あいうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおお..."
	out := Truncate(s, 80, "...")
	if out != expected {
		t.Errorf("Truncate(%q) = %q, want %q", s, out, expected)
	}
	width := StringWidth(out)
	if width != 79 {
		t.Errorf("width of Truncate(%q) should be %d, but %d", s, 79, width)
	}
}

func TestTruncateFit(t *testing.T) {
	s := "aあいうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおおおおお"
	expected := "aあいうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおお..."

	out := Truncate(s, 80, "...")
	if out != expected {
		t.Errorf("Truncate(%q) = %q, want %q", s, out, expected)
	}
	width := StringWidth(out)
	if width != 80 {
		t.Errorf("width of Truncate(%q) should be %d, but %d", s, 80, width)
	}
}

func TestTruncateJustFit(t *testing.T) {
	s := "あいうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおおおお"
	expected := "あいうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおおおお"

	out := Truncate(s, 80, "...")
	if out != expected {
		t.Errorf("Truncate(%q) = %q, want %q", s, out, expected)
	}
	width := StringWidth(out)
	if width != 80 {
		t.Errorf("width of Truncate(%q) should be %d, but %d", s, 80, width)
	}
}

func TestWrap(t *testing.T) {
	s := `東京特許許可局局長はよく柿喰う客だ/東京特許許可局局長はよく柿喰う客だ
123456789012345678901234567890

END`
	expected := `東京特許許可局局長はよく柿喰う
客だ/東京特許許可局局長はよく
柿喰う客だ
123456789012345678901234567890

END`

	if out := Wrap(s, 30); out != expected {
		t.Errorf("Wrap(%q) = %q, want %q", s, out, expected)
	}
}

func TestTruncateNoNeeded(t *testing.T) {
	s := "あいうえおあい"
	expected := "あいうえおあい"

	if out := Truncate(s, 80, "..."); out != expected {
		t.Errorf("Truncate(%q) = %q, want %q", s, out, expected)
	}
}

var isneutralwidthtests = []struct {
	in  rune
	out bool
}{
	{'→', false},
	{'┊', false},
	{'┈', false},
	{'～', false},
	{'└', false},
	{'⣀', true},
	{'⣀', true},
}

func TestIsNeutralWidth(t *testing.T) {
	for _, tt := range isneutralwidthtests {
		if out := IsNeutralWidth(tt.in); out != tt.out {
			t.Errorf("IsNeutralWidth(%q) = %v, want %v", tt.in, out, tt.out)
		}
	}
}

func TestFillLeft(t *testing.T) {
	s := "あxいうえお"
	expected := "    あxいうえお"

	if out := FillLeft(s, 15); out != expected {
		t.Errorf("FillLeft(%q) = %q, want %q", s, out, expected)
	}
}

func TestFillLeftFit(t *testing.T) {
	s := "あいうえお"
	expected := "あいうえお"

	if out := FillLeft(s, 10); out != expected {
		t.Errorf("FillLeft(%q) = %q, want %q", s, out, expected)
	}
}

func TestFillRight(t *testing.T) {
	s := "あxいうえお"
	expected := "あxいうえお    "

	if out := FillRight(s, 15); out != expected {
		t.Errorf("FillRight(%q) = %q, want %q", s, out, expected)
	}
}

func TestFillRightFit(t *testing.T) {
	s := "あいうえお"
	expected := "あいうえお"

	if out := FillRight(s, 10); out != expected {
		t.Errorf("FillRight(%q) = %q, want %q", s, out, expected)
	}
}

func TestEnv(t *testing.T) {
	old := os.Getenv("RUNEWIDTH_EASTASIAN")
	defer os.Setenv("RUNEWIDTH_EASTASIAN", old)

	os.Setenv("RUNEWIDTH_EASTASIAN", "0")
	handleEnv()

	if w := RuneWidth('│'); w != 1 {
		t.Errorf("RuneWidth('│') = %d, want %d", w, 1)
	}
}

func TestZeroWidthJointer(t *testing.T) {
	c := NewCondition()
	c.ZeroWidthJoiner = true

	var tests = []struct {
		in   string
		want int
	}{
		{"👩", 2},
		{"👩‍", 2},
		{"👩‍🍳", 2},
		{"‍🍳", 2},
		{"👨‍👨", 2},
		{"👨‍👨‍👧", 2},
		{"🏳️‍🌈", 2},
		{"あ👩‍🍳い", 6},
		{"あ‍🍳い", 6},
		{"あ‍い", 4},
	}

	for _, tt := range tests {
		if got := c.StringWidth(tt.in); got != tt.want {
			t.Errorf("StringWidth(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}
