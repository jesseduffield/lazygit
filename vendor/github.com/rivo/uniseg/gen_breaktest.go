//go:build generate

// This program generates a Go containing a slice of test cases based on the
// Unicode Character Database auxiliary data files. The command line arguments
// are as follows:
//
//   1. The name of the Unicode data file (just the filename, without extension).
//   2. The name of the locally generated Go file.
//   3. The name of the slice containing the test cases.
//   4. The name of the generator, for logging purposes.
//
//go:generate go run gen_breaktest.go GraphemeBreakTest graphemebreak_test.go graphemeBreakTestCases graphemes
//go:generate go run gen_breaktest.go WordBreakTest wordbreak_test.go wordBreakTestCases words
//go:generate go run gen_breaktest.go SentenceBreakTest sentencebreak_test.go sentenceBreakTestCases sentences
//go:generate go run gen_breaktest.go LineBreakTest linebreak_test.go lineBreakTestCases lines

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// We want to test against a specific version rather than the latest. When the
// package is upgraded to a new version, change these to generate new tests.
const (
	testCaseURL = `https://www.unicode.org/Public/14.0.0/ucd/auxiliary/%s.txt`
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Not enough arguments, see code for details")
		os.Exit(1)
	}

	log.SetPrefix("gen_breaktest (" + os.Args[4] + "): ")
	log.SetFlags(0)

	// Read text of testcases and parse into Go source code.
	src, err := parse(fmt.Sprintf(testCaseURL, os.Args[1]))
	if err != nil {
		log.Fatal(err)
	}

	// Format the Go code.
	formatted, err := format.Source(src)
	if err != nil {
		log.Fatalln("gofmt:", err)
	}

	// Write it out.
	log.Print("Writing to ", os.Args[2])
	if err := ioutil.WriteFile(os.Args[2], formatted, 0644); err != nil {
		log.Fatal(err)
	}
}

// parse reads a break text file, either from a local file or from a URL. It
// parses the file data into Go source code representing the test cases.
func parse(url string) ([]byte, error) {
	log.Printf("Parsing %s", url)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body := res.Body
	defer body.Close()

	buf := new(bytes.Buffer)
	buf.Grow(120 << 10)
	buf.WriteString(`package uniseg

// Code generated via go generate from gen_breaktest.go. DO NOT EDIT.

// ` + os.Args[3] + ` are Grapheme testcases taken from
// ` + url + `
// on ` + time.Now().Format("January 2, 2006") + `. See
// https://www.unicode.org/license.html for the Unicode license agreement.
var ` + os.Args[3] + ` = []testCase {
`)

	sc := bufio.NewScanner(body)
	num := 1
	var line []byte
	original := make([]byte, 0, 64)
	expected := make([]byte, 0, 64)
	for sc.Scan() {
		num++
		line = sc.Bytes()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		var comment []byte
		if i := bytes.IndexByte(line, '#'); i >= 0 {
			comment = bytes.TrimSpace(line[i+1:])
			line = bytes.TrimSpace(line[:i])
		}
		original, expected, err := parseRuneSequence(line, original[:0], expected[:0])
		if err != nil {
			return nil, fmt.Errorf(`line %d: %v: %q`, num, err, line)
		}
		fmt.Fprintf(buf, "\t{original: \"%s\", expected: %s}, // %s\n", original, expected, comment)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	// Check for final "# EOF", useful check if we're streaming via HTTP
	if !bytes.Equal(line, []byte("# EOF")) {
		return nil, fmt.Errorf(`line %d: exected "# EOF" as final line, got %q`, num, line)
	}
	buf.WriteString("}\n")
	return buf.Bytes(), nil
}

// Used by parseRuneSequence to match input via bytes.HasPrefix.
var (
	prefixBreak     = []byte("÷ ")
	prefixDontBreak = []byte("× ")
	breakOk         = []byte("÷")
	breakNo         = []byte("×")
)

// parseRuneSequence parses a rune + breaking opportunity sequence from b
// and appends the Go code for testcase.original to orig
// and appends the Go code for testcase.expected to exp.
// It retuns the new orig and exp slices.
//
// E.g. for the input b="÷ 0020 × 0308 ÷ 1F1E6 ÷"
// it will append
//     "\u0020\u0308\U0001F1E6"
// and "[][]rune{{0x0020,0x0308},{0x1F1E6},}"
// to orig and exp respectively.
//
// The formatting of exp is expected to be cleaned up by gofmt or format.Source.
// Note we explicitly require the sequence to start with ÷ and we implicitly
// require it to end with ÷.
func parseRuneSequence(b, orig, exp []byte) ([]byte, []byte, error) {
	// Check for and remove first ÷ or ×.
	if !bytes.HasPrefix(b, prefixBreak) && !bytes.HasPrefix(b, prefixDontBreak) {
		return nil, nil, errors.New("expected ÷ or × as first character")
	}
	if bytes.HasPrefix(b, prefixBreak) {
		b = b[len(prefixBreak):]
	} else {
		b = b[len(prefixDontBreak):]
	}

	boundary := true
	exp = append(exp, "[][]rune{"...)
	for len(b) > 0 {
		if boundary {
			exp = append(exp, '{')
		}
		exp = append(exp, "0x"...)
		// Find end of hex digits.
		var i int
		for i = 0; i < len(b) && b[i] != ' '; i++ {
			if d := b[i]; ('0' <= d || d <= '9') ||
				('A' <= d || d <= 'F') ||
				('a' <= d || d <= 'f') {
				continue
			}
			return nil, nil, errors.New("bad hex digit")
		}
		switch i {
		case 4:
			orig = append(orig, "\\u"...)
		case 5:
			orig = append(orig, "\\U000"...)
		default:
			return nil, nil, errors.New("unsupport code point hex length")
		}
		orig = append(orig, b[:i]...)
		exp = append(exp, b[:i]...)
		b = b[i:]

		// Check for space between hex and ÷ or ×.
		if len(b) < 1 || b[0] != ' ' {
			return nil, nil, errors.New("bad input")
		}
		b = b[1:]

		// Check for next boundary.
		switch {
		case bytes.HasPrefix(b, breakOk):
			boundary = true
			b = b[len(breakOk):]
		case bytes.HasPrefix(b, breakNo):
			boundary = false
			b = b[len(breakNo):]
		default:
			return nil, nil, errors.New("missing ÷ or ×")
		}
		if boundary {
			exp = append(exp, '}')
		}
		exp = append(exp, ',')
		if len(b) > 0 && b[0] == ' ' {
			b = b[1:]
		}
	}
	exp = append(exp, '}')
	return orig, exp, nil
}
