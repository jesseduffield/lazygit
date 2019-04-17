package language

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

const onePlusEpsilon = "1.00000000000000000000000000000001"

func TestGetPluralSpec(t *testing.T) {
	tests := []struct {
		src  string
		spec *PluralSpec
	}{
		{"pl", pluralSpecs["pl"]},
		{"en", pluralSpecs["en"]},
		{"en-US", pluralSpecs["en"]},
		{"en_US", pluralSpecs["en"]},
		{"en-GB", pluralSpecs["en"]},
		{"zh-CN", pluralSpecs["zh"]},
		{"zh-TW", pluralSpecs["zh"]},
		{"pt-BR", pluralSpecs["pt"]},
		{"pt_BR", pluralSpecs["pt"]},
		{"pt-PT", pluralSpecs["pt"]},
		{"pt_PT", pluralSpecs["pt"]},
		{"zh-Hans-CN", pluralSpecs["zh"]},
		{"zh-Hant-TW", pluralSpecs["zh"]},
		{"zh-CN", pluralSpecs["zh"]},
		{"zh-TW", pluralSpecs["zh"]},
		{"zh-Hans", pluralSpecs["zh"]},
		{"zh-Hant", pluralSpecs["zh"]},
		{"ko-KR", pluralSpecs["ko"]},
		{"ko_KR", pluralSpecs["ko"]},
		{"ko-KP", pluralSpecs["ko"]},
		{"ko_KP", pluralSpecs["ko"]},
		{"en-US-en-US", pluralSpecs["en"]},
		{"th", pluralSpecs["th"]},
		{"th-TH", pluralSpecs["th"]},
		{"hr", pluralSpecs["hr"]},
		{"bs", pluralSpecs["bs"]},
		{"sr", pluralSpecs["sr"]},
		{"ti", pluralSpecs["ti"]},
		{"vi", pluralSpecs["vi"]},
		{"vi-VN", pluralSpecs["vi"]},
		{"mk", pluralSpecs["mk"]},
		{"mk-MK", pluralSpecs["mk"]},
		{"lv", pluralSpecs["lv"]},
		{"lv-LV", pluralSpecs["lv"]},
		{".en-US..en-US.", nil},
		{"zh, en-gb;q=0.8, en;q=0.7", nil},
		{"zh,en-gb;q=0.8,en;q=0.7", nil},
		{"xx, en-gb;q=0.8, en;q=0.7", nil},
		{"xx,en-gb;q=0.8,en;q=0.7", nil},
		{"xx-YY,xx;q=0.8,en-US,en;q=0.8,de;q=0.6,nl;q=0.4", nil},
		{"/foo/es/en.json", nil},
		{"xx-Yyen-US", nil},
		{"en US", nil},
		{"", nil},
		{"-", nil},
		{"_", nil},
		{".", nil},
		{"-en", nil},
		{"_en", nil},
		{"-en-", nil},
		{"_en_", nil},
		{"xx", nil},
	}
	for _, test := range tests {
		spec := GetPluralSpec(test.src)
		if spec != test.spec {
			t.Errorf("getPluralSpec(%q) = %v expected %v", test.src, spec, test.spec)
		}
	}
}

type pluralTest struct {
	num    interface{}
	plural Plural
}

func appendIntegerTests(tests []pluralTest, plural Plural, examples []string) []pluralTest {
	for _, ex := range expandExamples(examples) {
		i, err := strconv.ParseInt(ex, 10, 64)
		if err != nil {
			panic(err)
		}
		tests = append(tests, pluralTest{ex, plural}, pluralTest{i, plural})
	}
	return tests
}

func appendDecimalTests(tests []pluralTest, plural Plural, examples []string) []pluralTest {
	for _, ex := range expandExamples(examples) {
		tests = append(tests, pluralTest{ex, plural})
	}
	return tests
}

func expandExamples(examples []string) []string {
	var expanded []string
	for _, ex := range examples {
		if parts := strings.Split(ex, "~"); len(parts) == 2 {
			for ex := parts[0]; ; ex = increment(ex) {
				expanded = append(expanded, ex)
				if ex == parts[1] {
					break
				}
			}
		} else {
			expanded = append(expanded, ex)
		}
	}
	return expanded
}

func increment(dec string) string {
	runes := []rune(dec)
	carry := true
	for i := len(runes) - 1; carry && i >= 0; i-- {
		switch runes[i] {
		case '.':
			continue
		case '9':
			runes[i] = '0'
		default:
			runes[i]++
			carry = false
		}
	}
	if carry {
		runes = append([]rune{'1'}, runes...)
	}
	return string(runes)
}

//
// Below here are tests that were manually written before tests were automatically generated.
// These are kept around as sanity checks for our code generation.
//

func TestArabic(t *testing.T) {
	tests := []pluralTest{
		{0, Zero},
		{"0", Zero},
		{"0.0", Zero},
		{"0.00", Zero},
		{1, One},
		{"1", One},
		{"1.0", One},
		{"1.00", One},
		{onePlusEpsilon, Other},
		{2, Two},
		{"2", Two},
		{"2.0", Two},
		{"2.00", Two},
		{3, Few},
		{"3", Few},
		{"3.0", Few},
		{"3.00", Few},
		{10, Few},
		{"10", Few},
		{"10.0", Few},
		{"10.00", Few},
		{103, Few},
		{"103", Few},
		{"103.0", Few},
		{"103.00", Few},
		{110, Few},
		{"110", Few},
		{"110.0", Few},
		{"110.00", Few},
		{11, Many},
		{"11", Many},
		{"11.0", Many},
		{"11.00", Many},
		{99, Many},
		{"99", Many},
		{"99.0", Many},
		{"99.00", Many},
		{111, Many},
		{"111", Many},
		{"111.0", Many},
		{"111.00", Many},
		{199, Many},
		{"199", Many},
		{"199.0", Many},
		{"199.00", Many},
		{100, Other},
		{"100", Other},
		{"100.0", Other},
		{"100.00", Other},
		{102, Other},
		{"102", Other},
		{"102.0", Other},
		{"102.00", Other},
		{200, Other},
		{"200", Other},
		{"200.0", Other},
		{"200.00", Other},
		{202, Other},
		{"202", Other},
		{"202.0", Other},
		{"202.00", Other},
	}
	tests = appendFloatTests(tests, 0.1, 0.9, Other)
	tests = appendFloatTests(tests, 1.1, 1.9, Other)
	tests = appendFloatTests(tests, 2.1, 2.9, Other)
	tests = appendFloatTests(tests, 3.1, 3.9, Other)
	tests = appendFloatTests(tests, 4.1, 4.9, Other)
	runTests(t, "ar", tests)
}

func TestBelarusian(t *testing.T) {
	tests := []pluralTest{
		{0, Many},
		{1, One},
		{2, Few},
		{3, Few},
		{4, Few},
		{5, Many},
		{19, Many},
		{20, Many},
		{21, One},
		{11, Many},
		{52, Few},
		{101, One},
		{"0.1", Other},
		{"0.7", Other},
		{"1.5", Other},
		{"1.0", One},
		{onePlusEpsilon, Other},
		{"2.0", Few},
		{"10.0", Many},
	}
	runTests(t, "be", tests)
}

func TestBurmese(t *testing.T) {
	tests := appendIntTests(nil, 0, 10, Other)
	tests = appendFloatTests(tests, 0, 10, Other)
	runTests(t, "my", tests)
}

func TestCatalan(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{"0", Other},
		{1, One},
		{"1", One},
		{"1.0", Other},
		{onePlusEpsilon, Other},
		{2, Other},
		{"2", Other},
	}
	tests = appendIntTests(tests, 2, 10, Other)
	tests = appendFloatTests(tests, 0, 10, Other)
	runTests(t, "ca", tests)
}

func TestChinese(t *testing.T) {
	tests := appendIntTests(nil, 0, 10, Other)
	tests = appendFloatTests(tests, 0, 10, Other)
	runTests(t, "zh", tests)
}

func TestCzech(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{"0", Other},
		{1, One},
		{"1", One},
		{onePlusEpsilon, Many},
		{2, Few},
		{"2", Few},
		{3, Few},
		{"3", Few},
		{4, Few},
		{"4", Few},
		{5, Other},
		{"5", Other},
	}
	tests = appendFloatTests(tests, 0, 10, Many)
	runTests(t, "cs", tests)
}

func TestDanish(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{onePlusEpsilon, One},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.1, 1.9, One)
	tests = appendFloatTests(tests, 2.0, 10.0, Other)
	runTests(t, "da", tests)
}

func TestDutch(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, Other)
	runTests(t, "nl", tests)
}

func TestEnglish(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, Other)
	runTests(t, "en", tests)
}

func TestFrench(t *testing.T) {
	tests := []pluralTest{
		{0, One},
		{1, One},
		{onePlusEpsilon, One},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 1.9, One)
	tests = appendFloatTests(tests, 2.0, 10.0, Other)
	runTests(t, "fr", tests)
}

func TestGerman(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, Other)
	runTests(t, "de", tests)
}

func TestIcelandic(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{2, Other},
		{11, Other},
		{21, One},
		{111, Other},
		{"0.0", Other},
		{"0.1", One},
		{"2.0", Other},
	}
	runTests(t, "is", tests)
}

func TestIndonesian(t *testing.T) {
	tests := appendIntTests(nil, 0, 10, Other)
	tests = appendFloatTests(tests, 0, 10, Other)
	runTests(t, "id", tests)
}

func TestItalian(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, Other)
	runTests(t, "it", tests)
}

func TestKorean(t *testing.T) {
	tests := appendIntTests(nil, 0, 10, Other)
	tests = appendFloatTests(tests, 0, 10, Other)
	runTests(t, "ko", tests)
}

func TestLatvian(t *testing.T) {
	tests := []pluralTest{
		{0, Zero},
		{"0", Zero},
		{"0.1", One},
		{1, One},
		{"1", One},
		{onePlusEpsilon, One},
		{"10.0", Zero},
		{"10.1", One},
		{"10.2", Other},
		{21, One},
	}
	tests = appendFloatTests(tests, 0.2, 0.9, Other)
	tests = appendFloatTests(tests, 1.2, 1.9, Other)
	tests = appendIntTests(tests, 2, 9, Other)
	tests = appendIntTests(tests, 10, 20, Zero)
	tests = appendIntTests(tests, 22, 29, Other)
	runTests(t, "lv", tests)
}

func TestJapanese(t *testing.T) {
	tests := appendIntTests(nil, 0, 10, Other)
	tests = appendFloatTests(tests, 0, 10, Other)
	runTests(t, "ja", tests)
}

func TestLithuanian(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{2, Few},
		{3, Few},
		{9, Few},
		{10, Other},
		{11, Other},
		{"0.1", Many},
		{"0.7", Many},
		{"1.0", One},
		{onePlusEpsilon, Many},
		{"2.0", Few},
		{"10.0", Other},
	}
	runTests(t, "lt", tests)
}

func TestMalay(t *testing.T) {
	tests := appendIntTests(nil, 0, 10, Other)
	tests = appendFloatTests(tests, 0, 10, Other)
	runTests(t, "ms", tests)
}

func TestPolish(t *testing.T) {
	tests := []pluralTest{
		{0, Many},
		{1, One},
		{2, Few},
		{3, Few},
		{4, Few},
		{5, Many},
		{19, Many},
		{20, Many},
		{10, Many},
		{11, Many},
		{52, Few},
		{"0.1", Other},
		{"0.7", Other},
		{"1.5", Other},
		{"1.0", Other},
		{onePlusEpsilon, Other},
		{"2.0", Other},
		{"10.0", Other},
	}
	runTests(t, "pl", tests)
}

func TestPortuguese(t *testing.T) {
	tests := []pluralTest{
		{0, One},
		{"0.0", One},
		{1, One},
		{"1.0", One},
		{onePlusEpsilon, One},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0, 1.5, One)
	tests = appendFloatTests(tests, 2, 10.0, Other)
	runTests(t, "pt", tests)
}

func TestMacedonian(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{"1.1", One},
		{"2.1", One},
		{onePlusEpsilon, One},
		{2, Other},
		{"2.2", Other},
		{11, One},
	}
	runTests(t, "mk", tests)
}

func TestRussian(t *testing.T) {
	tests := []pluralTest{
		{0, Many},
		{1, One},
		{2, Few},
		{3, Few},
		{4, Few},
		{5, Many},
		{19, Many},
		{20, Many},
		{21, One},
		{11, Many},
		{52, Few},
		{101, One},
		{"0.1", Other},
		{"0.7", Other},
		{"1.5", Other},
		{"1.0", Other},
		{onePlusEpsilon, Other},
		{"2.0", Other},
		{"10.0", Other},
	}
	runTests(t, "ru", tests)
}

func TestSpanish(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{"1", One},
		{"1.0", One},
		{"1.00", One},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 0.9, Other)
	tests = appendFloatTests(tests, 1.1, 10.0, Other)
	runTests(t, "es", tests)
}

func TestNorweigan(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{"1", One},
		{"1.0", One},
		{"1.00", One},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 0.9, Other)
	tests = appendFloatTests(tests, 1.1, 10.0, Other)
	runTests(t, "no", tests)
}

func TestBulgarian(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{2, Other},
		{3, Other},
		{9, Other},
		{10, Other},
		{11, Other},
		{"0.1", Other},
		{"0.7", Other},
		{"1.0", One},
		{"1.001", Other},
		{onePlusEpsilon, Other},
		{"1.1", Other},
		{"2.0", Other},
		{"10.0", Other},
	}
	runTests(t, "bg", tests)
}

func TestSwedish(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, Other)
	runTests(t, "sv", tests)
}

func TestThai(t *testing.T) {
	tests := appendIntTests(nil, 0, 10, Other)
	tests = appendFloatTests(tests, 0, 10, Other)
	runTests(t, "th", tests)
}

func TestVietnamese(t *testing.T) {
	tests := appendIntTests(nil, 0, 10, Other)
	tests = appendFloatTests(tests, 0, 10, Other)
	runTests(t, "vi", tests)
}

func TestTurkish(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{"1", One},
		{"1.0", One},
		{"1.00", One},
		{"1.001", Other},
		{"1.100", Other},
		{"1.101", Other},
		{onePlusEpsilon, Other},
		{2, Other},
		{"0.7", Other},
		{"2.0", Other},
	}
	runTests(t, "tr", tests)
}

func TestUkrainian(t *testing.T) {
	tests := []pluralTest{
		{0, Many},
		{1, One},
		{2, Few},
		{3, Few},
		{4, Few},
		{5, Many},
		{19, Many},
		{20, Many},
		{21, One},
		{11, Many},
		{52, Few},
		{101, One},
		{"0.1", Other},
		{"0.7", Other},
		{"1.5", Other},
		{"1.0", Other},
		{onePlusEpsilon, Other},
		{"2.0", Other},
		{"10.0", Other},
	}
	runTests(t, "uk", tests)
}

func TestCroatian(t *testing.T) {
	tests := makeCroatianBosnianSerbianTests()
	runTests(t, "hr", tests)
}

func TestBosnian(t *testing.T) {
	tests := makeCroatianBosnianSerbianTests()
	runTests(t, "bs", tests)
}

func TestSerbian(t *testing.T) {
	tests := makeCroatianBosnianSerbianTests()
	runTests(t, "sr", tests)
}

func makeCroatianBosnianSerbianTests() []pluralTest {
	return []pluralTest{
		{1, One},
		{"0.1", One},
		{21, One},
		{101, One},
		{1001, One},
		{51, One},
		{"1.1", One},
		{"5.1", One},
		{"100.1", One},
		{"1000.1", One},
		{2, Few},
		{"0.2", Few},
		{22, Few},
		{"1.2", Few},
		{24, Few},
		{"2.4", Few},
		{102, Few},
		{"100.2", Few},
		{1002, Few},
		{"1000.2", Few},
		{5, Other},
		{"0.5", Other},
		{0, Other},
		{100, Other},
		{19, Other},
		{"0.0", Other},
		{"100.0", Other},
		{"1000.0", Other},
	}
}

func TestTigrinya(t *testing.T) {
	tests := []pluralTest{
		{0, One},
		{1, One},
	}
	tests = appendIntTests(tests, 2, 10, Other)
	tests = appendFloatTests(tests, 1.1, 10.0, Other)
	runTests(t, "ti", tests)
}

func appendIntTests(tests []pluralTest, from, to int, p Plural) []pluralTest {
	for i := from; i <= to; i++ {
		tests = append(tests, pluralTest{i, p})
	}
	return tests
}

func appendFloatTests(tests []pluralTest, from, to float64, p Plural) []pluralTest {
	stride := 0.1
	format := "%.1f"
	for f := from; f < to; f += stride {
		tests = append(tests, pluralTest{fmt.Sprintf(format, f), p})
	}
	tests = append(tests, pluralTest{fmt.Sprintf(format, to), p})
	return tests
}

func runTests(t *testing.T, pluralSpecID string, tests []pluralTest) {
	pluralSpecID = normalizePluralSpecID(pluralSpecID)
	if spec := pluralSpecs[pluralSpecID]; spec != nil {
		for _, test := range tests {
			if plural, err := spec.Plural(test.num); plural != test.plural {
				t.Errorf("%s: PluralCategory(%#v) returned %s, %v; expected %s", pluralSpecID, test.num, plural, err, test.plural)
			}
		}
	} else {
		t.Errorf("could not find plural spec for locale %s", pluralSpecID)
	}

}
