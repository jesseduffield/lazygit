// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/internal/ucd"
)

var verbose = flag.Bool("verbose", false, "set to true to print the internal tables of matchers")

func TestCompliance(t *testing.T) {
	filepath.Walk("testdata", func(file string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		r, err := os.Open(file)
		if err != nil {
			t.Fatal(err)
		}
		ucd.Parse(r, func(p *ucd.Parser) {
			name := strings.Replace(path.Join(p.String(0), p.String(1)), " ", "", -1)
			if skip[name] {
				return
			}
			t.Run(info.Name()+"/"+name, func(t *testing.T) {
				supported := makeTagList(p.String(0))
				desired := makeTagList(p.String(1))
				gotCombined, index, conf := NewMatcher(supported).Match(desired...)

				gotMatch := supported[index]
				wantMatch := mk(p.String(2))
				if gotMatch != wantMatch {
					t.Fatalf("match: got %q; want %q (%v)", gotMatch, wantMatch, conf)
				}
				wantCombined, err := Raw.Parse(p.String(3))
				if err == nil && gotCombined != wantCombined {
					t.Errorf("combined: got %q; want %q (%v)", gotCombined, wantCombined, conf)
				}
			})
		})
		return nil
	})
}

var skip = map[string]bool{
	// TODO: bugs
	// Honor the wildcard match. This may only be useful to select non-exact
	// stuff.
	"mul,af/nl": true, // match: got "af"; want "mul"

	// TODO: include other extensions.
	// combined: got "en-GB-u-ca-buddhist-nu-arab"; want "en-GB-fonipa-t-m0-iso-i0-pinyin-u-ca-buddhist-nu-arab"
	"und,en-GB-u-sd-gbsct/en-fonipa-u-nu-Arab-ca-buddhist-t-m0-iso-i0-pinyin": true,

	// Inconsistencies with Mark Davis' implementation where it is not clear
	// which is better.

	// Inconsistencies in combined. I think the Go approach is more appropriate.
	// We could use -u-rg- and -u-va- as alternative.
	"und,fr/fr-BE-fonipa":              true, // combined: got "fr"; want "fr-BE-fonipa"
	"und,fr-CA/fr-BE-fonipa":           true, // combined: got "fr-CA"; want "fr-BE-fonipa"
	"und,fr-fonupa/fr-BE-fonipa":       true, // combined: got "fr-fonupa"; want "fr-BE-fonipa"
	"und,no/nn-BE-fonipa":              true, // combined: got "no"; want "no-BE-fonipa"
	"50,und,fr-CA-fonupa/fr-BE-fonipa": true, // combined: got "fr-CA-fonupa"; want "fr-BE-fonipa"

	// The initial number is a threshold. As we don't use scoring, we will not
	// implement this.
	"50,und,fr-Cyrl-CA-fonupa/fr-BE-fonipa": true,
	// match: got "und"; want "fr-Cyrl-CA-fonupa"
	// combined: got "und"; want "fr-Cyrl-BE-fonipa"

	// Other interesting cases to test:
	// - Should same language or same script have the preference if there is
	//   usually no understanding of the other script?
	// - More specific region in desired may replace enclosing supported.
}

func makeTagList(s string) (tags []Tag) {
	for _, s := range strings.Split(s, ",") {
		tags = append(tags, mk(strings.TrimSpace(s)))
	}
	return tags
}

func TestMatchStrings(t *testing.T) {
	testCases := []struct {
		supported string
		desired   string // strings separted by |
		tag       string
		index     int
	}{{
		supported: "en",
		desired:   "",
		tag:       "en",
		index:     0,
	}, {
		supported: "en",
		desired:   "nl",
		tag:       "en",
		index:     0,
	}, {
		supported: "en,nl",
		desired:   "nl",
		tag:       "nl",
		index:     1,
	}, {
		supported: "en,nl",
		desired:   "nl|en",
		tag:       "nl",
		index:     1,
	}, {
		supported: "en-GB,nl",
		desired:   "en ; q=0.1,nl",
		tag:       "nl",
		index:     1,
	}, {
		supported: "en-GB,nl",
		desired:   "en;q=0.005 | dk; q=0.1,nl ",
		tag:       "en-GB",
		index:     0,
	}, {
		// do not match faulty tags with und
		supported: "en,und",
		desired:   "|en",
		tag:       "en",
		index:     0,
	}}
	for _, tc := range testCases {
		t.Run(path.Join(tc.supported, tc.desired), func(t *testing.T) {
			m := NewMatcher(makeTagList(tc.supported))
			tag, index := MatchStrings(m, strings.Split(tc.desired, "|")...)
			if tag.String() != tc.tag || index != tc.index {
				t.Errorf("got %v, %d; want %v, %d", tag, index, tc.tag, tc.index)
			}
		})
	}
}

func TestAddLikelySubtags(t *testing.T) {
	tests := []struct{ in, out string }{
		{"aa", "aa-Latn-ET"},
		{"aa-Latn", "aa-Latn-ET"},
		{"aa-Arab", "aa-Arab-ET"},
		{"aa-Arab-ER", "aa-Arab-ER"},
		{"kk", "kk-Cyrl-KZ"},
		{"kk-CN", "kk-Arab-CN"},
		{"cmn", "cmn"},
		{"zh-AU", "zh-Hant-AU"},
		{"zh-VN", "zh-Hant-VN"},
		{"zh-SG", "zh-Hans-SG"},
		{"zh-Hant", "zh-Hant-TW"},
		{"zh-Hani", "zh-Hani-CN"},
		{"und-Hani", "zh-Hani-CN"},
		{"und", "en-Latn-US"},
		{"und-GB", "en-Latn-GB"},
		{"und-CW", "pap-Latn-CW"},
		{"und-YT", "fr-Latn-YT"},
		{"und-Arab", "ar-Arab-EG"},
		{"und-AM", "hy-Armn-AM"},
		{"und-TW", "zh-Hant-TW"},
		{"und-002", "en-Latn-NG"},
		{"und-Latn-002", "en-Latn-NG"},
		{"en-Latn-002", "en-Latn-NG"},
		{"en-002", "en-Latn-NG"},
		{"en-001", "en-Latn-US"},
		{"und-003", "en-Latn-US"},
		{"und-GB", "en-Latn-GB"},
		{"Latn-001", "en-Latn-US"},
		{"en-001", "en-Latn-US"},
		{"es-419", "es-Latn-419"},
		{"he-145", "he-Hebr-IL"},
		{"ky-145", "ky-Latn-TR"},
		{"kk", "kk-Cyrl-KZ"},
		// Don't specialize duplicate and ambiguous matches.
		{"kk-034", "kk-Arab-034"}, // Matches IR and AF. Both are Arab.
		{"ku-145", "ku-Latn-TR"},  // Matches IQ, TR, and LB, but kk -> TR.
		{"und-Arab-CC", "ms-Arab-CC"},
		{"und-Arab-GB", "ks-Arab-GB"},
		{"und-Hans-CC", "zh-Hans-CC"},
		{"und-CC", "en-Latn-CC"},
		{"sr", "sr-Cyrl-RS"},
		{"sr-151", "sr-Latn-151"}, // Matches RO and RU.
		// We would like addLikelySubtags to generate the same results if the input
		// only changes by adding tags that would otherwise have been added
		// by the expansion.
		// In other words:
		//     und-AA -> xx-Scrp-AA   implies und-Scrp-AA -> xx-Scrp-AA
		//     und-AA -> xx-Scrp-AA   implies xx-AA -> xx-Scrp-AA
		//     und-Scrp -> xx-Scrp-AA implies und-Scrp-AA -> xx-Scrp-AA
		//     und-Scrp -> xx-Scrp-AA implies xx-Scrp -> xx-Scrp-AA
		//     xx -> xx-Scrp-AA       implies xx-Scrp -> xx-Scrp-AA
		//     xx -> xx-Scrp-AA       implies xx-AA -> xx-Scrp-AA
		//
		// The algorithm specified in
		//   http://unicode.org/reports/tr35/tr35-9.html#Supplemental_Data,
		// Section C.10, does not handle the first case. For example,
		// the CLDR data contains an entry und-BJ -> fr-Latn-BJ, but not
		// there is no rule for und-Latn-BJ.  According to spec, und-Latn-BJ
		// would expand to en-Latn-BJ, violating the aforementioned principle.
		// We deviate from the spec by letting und-Scrp-AA expand to xx-Scrp-AA
		// if a rule of the form und-AA -> xx-Scrp-AA is defined.
		// Note that as of version 23, CLDR has some explicitly specified
		// entries that do not conform to these rules. The implementation
		// will not correct these explicit inconsistencies. A later versions of CLDR
		// is supposed to fix this.
		{"und-Latn-BJ", "fr-Latn-BJ"},
		{"und-Bugi-ID", "bug-Bugi-ID"},
		// regions, scripts and languages without definitions
		{"und-Arab-AA", "ar-Arab-AA"},
		{"und-Afak-RE", "fr-Afak-RE"},
		{"und-Arab-GB", "ks-Arab-GB"},
		{"abp-Arab-GB", "abp-Arab-GB"},
		// script has preference over region
		{"und-Arab-NL", "ar-Arab-NL"},
		{"zza", "zza-Latn-TR"},
		// preserve variants and extensions
		{"de-1901", "de-Latn-DE-1901"},
		{"de-x-abc", "de-Latn-DE-x-abc"},
		{"de-1901-x-abc", "de-Latn-DE-1901-x-abc"},
		{"x-abc", "x-abc"}, // TODO: is this the desired behavior?
	}
	for i, tt := range tests {
		in, _ := Parse(tt.in)
		out, _ := Parse(tt.out)
		in, _ = in.addLikelySubtags()
		if in.String() != out.String() {
			t.Errorf("%d: add(%s) was %s; want %s", i, tt.in, in, tt.out)
		}
	}
}
func TestMinimize(t *testing.T) {
	tests := []struct{ in, out string }{
		{"aa", "aa"},
		{"aa-Latn", "aa"},
		{"aa-Latn-ET", "aa"},
		{"aa-ET", "aa"},
		{"aa-Arab", "aa-Arab"},
		{"aa-Arab-ER", "aa-Arab-ER"},
		{"aa-Arab-ET", "aa-Arab"},
		{"und", "und"},
		{"und-Latn", "und"},
		{"und-Latn-US", "und"},
		{"en-Latn-US", "en"},
		{"cmn", "cmn"},
		{"cmn-Hans", "cmn-Hans"},
		{"cmn-Hant", "cmn-Hant"},
		{"zh-AU", "zh-AU"},
		{"zh-VN", "zh-VN"},
		{"zh-SG", "zh-SG"},
		{"zh-Hant", "zh-Hant"},
		{"zh-Hant-TW", "zh-TW"},
		{"zh-Hans", "zh"},
		{"zh-Hani", "zh-Hani"},
		{"und-Hans", "und-Hans"},
		{"und-Hani", "und-Hani"},

		{"und-CW", "und-CW"},
		{"und-YT", "und-YT"},
		{"und-Arab", "und-Arab"},
		{"und-AM", "und-AM"},
		{"und-Arab-CC", "und-Arab-CC"},
		{"und-CC", "und-CC"},
		{"und-Latn-BJ", "und-BJ"},
		{"und-Bugi-ID", "und-Bugi"},
		{"bug-Bugi-ID", "bug-Bugi"},
		// regions, scripts and languages without definitions
		{"und-Arab-AA", "und-Arab-AA"},
		// preserve variants and extensions
		{"de-Latn-1901", "de-1901"},
		{"de-Latn-x-abc", "de-x-abc"},
		{"de-DE-1901-x-abc", "de-1901-x-abc"},
		{"x-abc", "x-abc"}, // TODO: is this the desired behavior?
	}
	for i, tt := range tests {
		in, _ := Parse(tt.in)
		out, _ := Parse(tt.out)
		min, _ := in.minimize()
		if min.String() != out.String() {
			t.Errorf("%d: min(%s) was %s; want %s", i, tt.in, min, tt.out)
		}
		max, _ := min.addLikelySubtags()
		if x, _ := in.addLikelySubtags(); x.String() != max.String() {
			t.Errorf("%d: max(min(%s)) = %s; want %s", i, tt.in, max, x)
		}
	}
}

func TestRegionGroups(t *testing.T) {
	testCases := []struct {
		a, b     string
		distance uint8
	}{
		{"zh-TW", "zh-HK", 5},
		{"zh-MO", "zh-HK", 4},
		{"es-ES", "es-AR", 5},
		{"es-ES", "es", 4},
		{"es-419", "es-MX", 4},
		{"es-AR", "es-MX", 4},
		{"es-ES", "es-MX", 5},
		{"es-PT", "es-MX", 5},
	}
	for _, tc := range testCases {
		a := MustParse(tc.a)
		aScript, _ := a.Script()
		b := MustParse(tc.b)
		bScript, _ := b.Script()

		if aScript != bScript {
			t.Errorf("scripts differ: %q vs %q", aScript, bScript)
			continue
		}
		d, _ := regionGroupDist(a.region, b.region, aScript.scriptID, a.lang)
		if d != tc.distance {
			t.Errorf("got %q; want %q", d, tc.distance)
		}
	}
}

func TestIsParadigmLocale(t *testing.T) {
	testCases := map[string]bool{
		"en-US":  true,
		"en-GB":  true,
		"en-VI":  false,
		"es-GB":  false,
		"es-ES":  true,
		"es-419": true,
	}
	for str, want := range testCases {
		tag := Make(str)
		got := isParadigmLocale(tag.lang, tag.region)
		if got != want {
			t.Errorf("isPL(%q) = %v; want %v", str, got, want)
		}
	}
}

// Implementation of String methods for various types for debugging purposes.

func (m *matcher) String() string {
	w := &bytes.Buffer{}
	fmt.Fprintln(w, "Default:", m.default_)
	for tag, h := range m.index {
		fmt.Fprintf(w, "  %s: %v\n", tag, h)
	}
	return w.String()
}

func (h *matchHeader) String() string {
	w := &bytes.Buffer{}
	fmt.Fprint(w, "haveTag: ")
	for _, h := range h.haveTags {
		fmt.Fprintf(w, "%v, ", h)
	}
	return w.String()
}

func (t haveTag) String() string {
	return fmt.Sprintf("%v:%d:%v:%v-%v|%v", t.tag, t.index, t.conf, t.maxRegion, t.maxScript, t.altScript)
}

func TestBestMatchAlloc(t *testing.T) {
	m := NewMatcher(makeTagList("en sr nl"))
	// Go allocates when creating a list of tags from a single tag!
	list := []Tag{English}
	avg := testtext.AllocsPerRun(1, func() {
		m.Match(list...)
	})
	if avg > 0 {
		t.Errorf("got %f; want 0", avg)
	}
}

var benchHave = []Tag{
	mk("en"),
	mk("en-GB"),
	mk("za"),
	mk("zh-Hant"),
	mk("zh-Hans-CN"),
	mk("zh"),
	mk("zh-HK"),
	mk("ar-MK"),
	mk("en-CA"),
	mk("fr-CA"),
	mk("fr-US"),
	mk("fr-CH"),
	mk("fr"),
	mk("lt"),
	mk("lv"),
	mk("iw"),
	mk("iw-NL"),
	mk("he"),
	mk("he-IT"),
	mk("tlh"),
	mk("ja"),
	mk("ja-Jpan"),
	mk("ja-Jpan-JP"),
	mk("de"),
	mk("de-CH"),
	mk("de-AT"),
	mk("de-DE"),
	mk("sr"),
	mk("sr-Latn"),
	mk("sr-Cyrl"),
	mk("sr-ME"),
}

var benchWant = [][]Tag{
	[]Tag{
		mk("en"),
	},
	[]Tag{
		mk("en-AU"),
		mk("de-HK"),
		mk("nl"),
		mk("fy"),
		mk("lv"),
	},
	[]Tag{
		mk("en-AU"),
		mk("de-HK"),
		mk("nl"),
		mk("fy"),
	},
	[]Tag{
		mk("ja-Hant"),
		mk("da-HK"),
		mk("nl"),
		mk("zh-TW"),
	},
	[]Tag{
		mk("ja-Hant"),
		mk("da-HK"),
		mk("nl"),
		mk("hr"),
	},
}

func BenchmarkMatch(b *testing.B) {
	m := newMatcher(benchHave, nil)
	for i := 0; i < b.N; i++ {
		for _, want := range benchWant {
			m.getBest(want...)
		}
	}
}

func BenchmarkMatchExact(b *testing.B) {
	want := mk("en")
	m := newMatcher(benchHave, nil)
	for i := 0; i < b.N; i++ {
		m.getBest(want)
	}
}

func BenchmarkMatchAltLanguagePresent(b *testing.B) {
	want := mk("hr")
	m := newMatcher(benchHave, nil)
	for i := 0; i < b.N; i++ {
		m.getBest(want)
	}
}

func BenchmarkMatchAltLanguageNotPresent(b *testing.B) {
	want := mk("nn")
	m := newMatcher(benchHave, nil)
	for i := 0; i < b.N; i++ {
		m.getBest(want)
	}
}

func BenchmarkMatchAltScriptPresent(b *testing.B) {
	want := mk("zh-Hant-CN")
	m := newMatcher(benchHave, nil)
	for i := 0; i < b.N; i++ {
		m.getBest(want)
	}
}

func BenchmarkMatchAltScriptNotPresent(b *testing.B) {
	want := mk("fr-Cyrl")
	m := newMatcher(benchHave, nil)
	for i := 0; i < b.N; i++ {
		m.getBest(want)
	}
}

func BenchmarkMatchLimitedExact(b *testing.B) {
	want := []Tag{mk("he-NL"), mk("iw-NL")}
	m := newMatcher(benchHave, nil)
	for i := 0; i < b.N; i++ {
		m.getBest(want...)
	}
}
