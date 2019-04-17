// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package idna

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/internal/ucd"
)

func TestAllocToUnicode(t *testing.T) {
	avg := testtext.AllocsPerRun(1000, func() {
		ToUnicode("www.golang.org")
	})
	if avg > 0 {
		t.Errorf("got %f; want 0", avg)
	}
}

func TestAllocToASCII(t *testing.T) {
	avg := testtext.AllocsPerRun(1000, func() {
		ToASCII("www.golang.org")
	})
	if avg > 0 {
		t.Errorf("got %f; want 0", avg)
	}
}

func TestProfiles(t *testing.T) {
	testCases := []struct {
		name      string
		want, got *Profile
	}{
		{"Punycode", punycode, New()},
		{"Registration", registration, New(ValidateForRegistration())},
		{"Registration", registration, New(
			ValidateForRegistration(),
			VerifyDNSLength(true),
			BidiRule(),
		)},
		{"Lookup", lookup, New(MapForLookup(), BidiRule(), Transitional(true))},
		{"Display", display, New(MapForLookup(), BidiRule())},
	}
	for _, tc := range testCases {
		// Functions are not comparable, but the printed version will include
		// their pointers.
		got := fmt.Sprintf("%#v", tc.got)
		want := fmt.Sprintf("%#v", tc.want)
		if got != want {
			t.Errorf("%s: \ngot  %#v,\nwant %#v", tc.name, got, want)
		}
	}
}

// doTest performs a single test f(input) and verifies that the output matches
// out and that the returned error is expected. The errors string contains
// all allowed error codes as categorized in
// http://www.unicode.org/Public/idna/9.0.0/IdnaTest.txt:
// P: Processing
// V: Validity
// A: to ASCII
// B: Bidi
// C: Context J
func doTest(t *testing.T, f func(string) (string, error), name, input, want, errors string) {
	errors = strings.Trim(errors, "[]")
	test := "ok"
	if errors != "" {
		test = "err:" + errors
	}
	// Replace some of the escape sequences to make it easier to single out
	// tests on the command name.
	in := strings.Trim(strconv.QuoteToASCII(input), `"`)
	in = strings.Replace(in, `\u`, "#", -1)
	in = strings.Replace(in, `\U`, "#", -1)
	name = fmt.Sprintf("%s/%s/%s", name, in, test)

	testtext.Run(t, name, func(t *testing.T) {
		got, err := f(input)

		if err != nil {
			code := err.(interface {
				code() string
			}).code()
			if strings.Index(errors, code) == -1 {
				t.Errorf("error %q not in set of expected errors {%v}", code, errors)
			}
		} else if errors != "" {
			t.Errorf("no errors; want error in {%v}", errors)
		}

		if want != "" && got != want {
			t.Errorf(`string: got %+q; want %+q`, got, want)
		}
	})
}

func TestConformance(t *testing.T) {
	testtext.SkipIfNotLong(t)

	r := gen.OpenUnicodeFile("idna", "", "IdnaTest.txt")
	defer r.Close()

	section := "main"
	started := false
	p := ucd.New(r, ucd.CommentHandler(func(s string) {
		if started {
			section = strings.ToLower(strings.Split(s, " ")[0])
		}
	}))
	transitional := New(Transitional(true), VerifyDNSLength(true), BidiRule(), MapForLookup())
	nonTransitional := New(VerifyDNSLength(true), BidiRule(), MapForLookup())
	for p.Next() {
		started = true

		// What to test
		profiles := []*Profile{}
		switch p.String(0) {
		case "T":
			profiles = append(profiles, transitional)
		case "N":
			profiles = append(profiles, nonTransitional)
		case "B":
			profiles = append(profiles, transitional)
			profiles = append(profiles, nonTransitional)
		}

		src := unescape(p.String(1))

		wantToUnicode := unescape(p.String(2))
		if wantToUnicode == "" {
			wantToUnicode = src
		}
		wantToASCII := unescape(p.String(3))
		if wantToASCII == "" {
			wantToASCII = wantToUnicode
		}
		wantErrToUnicode := ""
		if strings.HasPrefix(wantToUnicode, "[") {
			wantErrToUnicode = wantToUnicode
			wantToUnicode = ""
		}
		wantErrToASCII := ""
		if strings.HasPrefix(wantToASCII, "[") {
			wantErrToASCII = wantToASCII
			wantToASCII = ""
		}

		// TODO: also do IDNA tests.
		// invalidInIDNA2008 := p.String(4) == "NV8"

		for _, p := range profiles {
			name := fmt.Sprintf("%s:%s", section, p)
			doTest(t, p.ToUnicode, name+":ToUnicode", src, wantToUnicode, wantErrToUnicode)
			doTest(t, p.ToASCII, name+":ToASCII", src, wantToASCII, wantErrToASCII)
		}
	}
}

func unescape(s string) string {
	s, err := strconv.Unquote(`"` + s + `"`)
	if err != nil {
		panic(err)
	}
	return s
}

func BenchmarkProfile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Lookup.ToASCII("www.yahoogle.com")
	}
}
