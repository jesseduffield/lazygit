package str

//import "testing"
import "fmt"

//import "strings"

func ExampleBetween() {
	eg(1, Between("<a>foo</a>", "<a>", "</a>"))
	eg(2, Between("<a>foo</a></a>", "<a>", "</a>"))
	eg(3, Between("<a><a>foo</a></a>", "<a>", "</a>"))
	eg(4, Between("<a><a>foo</a></a>", "<a>", "</a>"))
	eg(5, Between("<a>foo", "<a>", "</a>"))
	eg(6, Between("Some strings } are very {weird}, dont you think?", "{", "}"))
	eg(7, Between("This is ateststring", "", "test"))
	eg(8, Between("This is ateststring", "test", ""))
	// Output:
	// 1: foo
	// 2: foo
	// 3: <a>foo
	// 4: <a>foo
	// 5:
	// 6: weird
	// 7: This is a
	// 8: string
}

func ExampleBetweenF() {
	eg(1, Pipe("abc", BetweenF("a", "c")))
	// Output:
	// 1: b
}

func ExampleCamelize() {
	eg(1, Camelize("data_rate"))
	eg(2, Camelize("background-color"))
	eg(3, Camelize("-moz-something"))
	eg(4, Camelize("_car_speed_"))
	eg(5, Camelize("yes_we_can"))
	// Output:
	// 1: dataRate
	// 2: backgroundColor
	// 3: MozSomething
	// 4: CarSpeed
	// 5: yesWeCan
}

func ExampleCapitalize() {
	eg(1, Capitalize("abc"))
	eg(2, Capitalize("ABC"))
	// Output:
	// 1: Abc
	// 2: Abc
}

func ExampleCharAt() {
	eg(1, CharAt("abc", 1))
	eg(2, CharAt("", -1))
	eg(3, CharAt("", 0))
	eg(4, CharAt("", 10))
	eg(5, CharAt("abc", -1))
	eg(6, CharAt("abc", 10))
	// Output:
	// 1: b
	// 2:
	// 3:
	// 4:
	// 5:
	// 6:
}

func ExampleCharAtF() {
	eg(1, Pipe("abc", CharAtF(1)))
	// Output:
	// 1: b
}

func ExampleChompLeft() {
	eg(1, ChompLeft("foobar", "foo"))
	eg(2, ChompLeft("foobar", "bar"))
	eg(3, ChompLeft("", "foo"))
	eg(4, ChompLeft("", ""))
	eg(5, ChompLeft("foo", ""))
	// Output:
	// 1: bar
	// 2: foobar
	// 3:
	// 4:
	// 5: foo
}

func ExampleChompLeftF() {
	eg(1, Pipe("abc", ChompLeftF("ab")))
	// Output:
	// 1: c
}

func ExampleChompRight() {
	eg(1, ChompRight("foobar", "foo"))
	eg(2, ChompRight("foobar", "bar"))
	eg(3, ChompRight("", "foo"))
	eg(4, ChompRight("", ""))
	// Output:
	// 1: foobar
	// 2: foo
	// 3:
	// 4:
}

func ExampleChompRightF() {
	eg(1, Pipe("abc", ChompRightF("bc")))
	// Output:
	// 1: a
}

func ExampleClassify() {
	eg(1, Classify("data_rate"))
	eg(2, Classify("background-color"))
	eg(3, Classify("-moz-something"))
	eg(4, Classify("_car_speed_"))
	eg(5, Classify("yes_we_can"))
	// Output:
	// 1: DataRate
	// 2: BackgroundColor
	// 3: MozSomething
	// 4: CarSpeed
	// 5: YesWeCan
}

func ExampleClean() {
	eg(1, Clean("clean"))
	eg(2, Clean(""))
	eg(3, Clean(" please\t    clean \t \n  me "))
	// Output:
	// 1: clean
	// 2:
	// 3: please clean me
}

func ExampleDasherize() {
	eg(1, Dasherize("dataRate"))
	eg(2, Dasherize("CarSpeed"))
	eg(3, Dasherize("yesWeCan"))
	eg(4, Dasherize(""))
	eg(5, Dasherize("ABC"))
	// Output:
	// 1: data-rate
	// 2: -car-speed
	// 3: yes-we-can
	// 4:
	// 5: -a-b-c
}

func ExampleDecodeHTMLEntities() {
	eg(1, DecodeHTMLEntities("Ken Thompson &amp; Dennis Ritchie"))
	eg(2, DecodeHTMLEntities("3 &lt; 4"))
	eg(3, DecodeHTMLEntities("http:&#47;&#47;"))
	// Output:
	// 1: Ken Thompson & Dennis Ritchie
	// 2: 3 < 4
	// 3: http://
}

func ExampleEnsurePrefix() {
	eg(1, EnsurePrefix("foobar", "foo"))
	eg(2, EnsurePrefix("bar", "foo"))
	eg(3, EnsurePrefix("", ""))
	eg(4, EnsurePrefix("foo", ""))
	eg(5, EnsurePrefix("", "foo"))
	// Output:
	// 1: foobar
	// 2: foobar
	// 3:
	// 4: foo
	// 5: foo
}

func ExampleEnsurePrefixF() {
	eg(1, Pipe("dir", EnsurePrefixF("./")))
	// Output:
	// 1: ./dir
}

func ExampleEnsureSuffix() {
	eg(1, EnsureSuffix("foobar", "bar"))
	eg(2, EnsureSuffix("foo", "bar"))
	eg(3, EnsureSuffix("", ""))
	eg(4, EnsureSuffix("foo", ""))
	eg(5, EnsureSuffix("", "bar"))
	// Output:
	// 1: foobar
	// 2: foobar
	// 3:
	// 4: foo
	// 5: bar
}

func ExampleHumanize() {
	eg(1, Humanize("the_humanize_string_method"))
	eg(2, Humanize("ThehumanizeStringMethod"))
	eg(3, Humanize("the humanize string method"))
	// Output:
	// 1: The humanize string method
	// 2: Thehumanize string method
	// 3: The humanize string method
}

func ExampleIif() {
	eg(1, Iif(true, "T", "F"))
	eg(2, Iif(false, "T", "F"))
	// Output:
	// 1: T
	// 2: F
}

func ExampleIndexOf() {
	eg(1, IndexOf("abcdef", "a", 0))
	eg(2, IndexOf("abcdef", "a", 3))
	eg(3, IndexOf("abcdef", "a", -2))
	eg(4, IndexOf("abcdef", "a", 10))
	eg(5, IndexOf("", "a", 0))
	eg(6, IndexOf("abcdef", "", 2))
	eg(7, IndexOf("abcdef", "", 1000))
	// Output:
	// 1: 0
	// 2: -1
	// 3: -1
	// 4: -1
	// 5: -1
	// 6: 2
	// 7: 6
}

func ExampleIsAlpha() {
	eg(1, IsAlpha("afaf"))
	eg(2, IsAlpha("FJslfjkasfs"))
	eg(3, IsAlpha("áéúóúÁÉÍÓÚãõÃÕàèìòùÀÈÌÒÙâêîôûÂÊÎÔÛäëïöüÄËÏÖÜçÇ"))
	eg(4, IsAlpha("adflj43faljsdf"))
	eg(5, IsAlpha("33"))
	eg(6, IsAlpha("TT....TTTafafetstYY"))
	eg(7, IsAlpha("-áéúóúÁÉÍÓÚãõÃÕàèìòùÀÈÌÒÙâêîôûÂÊÎÔÛäëïöüÄËÏÖÜçÇ"))
	// Output:
	// 1: true
	// 2: true
	// 3: true
	// 4: false
	// 5: false
	// 6: false
	// 7: false
}

func eg(index int, example interface{}) {
	output := fmt.Sprintf("%d: %v", index, example)
	fmt.Printf("%s\n", Clean(output))
}

func ExampleIsAlphaNumeric() {
	eg(1, IsAlphaNumeric("afaf35353afaf"))
	eg(2, IsAlphaNumeric("FFFF99fff"))
	eg(3, IsAlphaNumeric("99"))
	eg(4, IsAlphaNumeric("afff"))
	eg(5, IsAlphaNumeric("Infinity"))
	eg(6, IsAlphaNumeric("áéúóúÁÉÍÓÚãõÃÕàèìòùÀÈÌÒÙâêîôûÂÊÎÔÛäëïöüÄËÏÖÜçÇ1234567890"))
	eg(7, IsAlphaNumeric("-Infinity"))
	eg(8, IsAlphaNumeric("-33"))
	eg(9, IsAlphaNumeric("aaff.."))
	eg(10, IsAlphaNumeric(".áéúóúÁÉÍÓÚãõÃÕàèìòùÀÈÌÒÙâêîôûÂÊÎÔÛäëïöüÄËÏÖÜçÇ1234567890"))
	// Output:
	// 1: true
	// 2: true
	// 3: true
	// 4: true
	// 5: true
	// 6: true
	// 7: false
	// 8: false
	// 9: false
	// 10: false
}

func ExampleIsEmpty() {
	eg(1, IsEmpty(" "))
	eg(2, IsEmpty("\t\t\t   "))
	eg(3, IsEmpty("\t\n "))
	eg(4, IsEmpty("hi"))
	// Output:
	// 1: true
	// 2: true
	// 3: true
	// 4: false
}

func ExampleIsLower() {
	eg(1, IsLower("a"))
	eg(2, IsLower("A"))
	eg(3, IsLower("abc"))
	eg(4, IsLower("aBc"))
	eg(5, IsLower("áéúóúãõàèìòùâêîôûäëïöüç"))
	eg(6, IsLower("hi jp"))
	eg(7, IsLower("ÁÉÍÓÚÃÕÀÈÌÒÙÂÊÎÔÛÄËÏÖÜÇ"))
	eg(8, IsLower("áéúóúãõàèìòùâêîôûäëïöüçÁ"))
	eg(9, IsLower("áéúóúãõàèìòùâêîôû äëïöüç"))
	// Output:
	// 1: true
	// 2: false
	// 3: true
	// 4: false
	// 5: true
	// 6: false
	// 7: false
	// 8: false
	// 9: false
}

func ExampleIsNumeric() {
	eg(1, IsNumeric("3"))
	eg(2, IsNumeric("34.22"))
	eg(3, IsNumeric("-22.33"))
	eg(4, IsNumeric("NaN"))
	eg(5, IsNumeric("Infinity"))
	eg(6, IsNumeric("-Infinity"))
	eg(7, IsNumeric("JP"))
	eg(8, IsNumeric("-5"))
	eg(9, IsNumeric("00099242424"))
	// Output:
	// 1: true
	// 2: false
	// 3: false
	// 4: false
	// 5: false
	// 6: false
	// 7: false
	// 8: false
	// 9: true
}

func ExampleIsUpper() {
	eg(1, IsUpper("a"))
	eg(2, IsUpper("A"))
	eg(3, IsUpper("ABC"))
	eg(4, IsUpper("aBc"))
	eg(5, IsUpper("áéúóúãõàèìòùâêîôûäëïöüç"))
	eg(6, IsUpper("HI JP"))
	eg(7, IsUpper("ÁÉÍÓÚÃÕÀÈÌÒÙÂÊÎÔÛÄËÏÖÜÇ"))
	eg(8, IsUpper("áéúóúãõàèìòùâêîôûäëïöüçÁ"))
	eg(9, IsUpper("ÁÉÍÓÚÃÕÀÈÌÒÙÂÊÎ ÔÛÄËÏÖÜÇ"))
	// Output:
	// 1: false
	// 2: true
	// 3: true
	// 4: false
	// 5: false
	// 6: false
	// 7: true
	// 8: false
	// 9: false
}

func ExampleLeft() {
	eg(1, Left("abcdef", 0))
	eg(2, Left("abcdef", 1))
	eg(3, Left("abcdef", 4))
	eg(4, Left("abcdef", -2))
	// Output:
	// 1:
	// 2: a
	// 3: abcd
	// 4: ef
}

func ExampleLeftOf() {
	eg(1, LeftOf("abcdef", "def"))
	eg(2, LeftOf("abcdef", "abc"))
	eg(3, LeftOf("abcdef", ""))
	eg(4, LeftOf("", "abc"))
	eg(5, LeftOf("abcdef", "xyz"))
	// Output:
	// 1: abc
	// 2:
	// 3: abcdef
	// 4:
	// 5:
}

func ExampleLines() {
	eg(1, Lines("a\r\nb\nc\r\n"))
	eg(2, Lines("a\r\nb\nc\r\nd"))
	// Output:
	// 1: [a b c ]
	// 2: [a b c d]
}

func ExampleMatch() {
	eg(1, Match("foobar", `^fo.*r$`))
	eg(2, Match("foobar", `^fo.*x$`))
	eg(3, Match("", `^fo.*x$`))
	// Output:
	// 1: true
	// 2: false
	// 3: false
}

func ExamplePad() {
	eg(1, Pad("hello", "x", 5))
	eg(2, Pad("hello", "x", 10))
	eg(3, Pad("hello", "x", 11))
	eg(4, Pad("hello", "x", 6))
	eg(5, Pad("hello", "x", 1))
	// Output:
	// 1: hello
	// 2: xxxhelloxx
	// 3: xxxhelloxxx
	// 4: xhello
	// 5: hello
}

func ExamplePadLeft() {
	eg(1, PadLeft("hello", "x", 5))
	eg(2, PadLeft("hello", "x", 10))
	eg(3, PadLeft("hello", "x", 11))
	eg(4, PadLeft("hello", "x", 6))
	eg(5, PadLeft("hello", "x", 1))
	// Output:
	// 1: hello
	// 2: xxxxxhello
	// 3: xxxxxxhello
	// 4: xhello
	// 5: hello
}

func ExamplePadRight() {
	eg(1, PadRight("hello", "x", 5))
	eg(2, PadRight("hello", "x", 10))
	eg(3, PadRight("hello", "x", 11))
	eg(4, PadRight("hello", "x", 6))
	eg(5, PadRight("hello", "x", 1))
	// Output:
	// 1: hello
	// 2: helloxxxxx
	// 3: helloxxxxxx
	// 4: hellox
	// 5: hello
}

func ExamplePipe() {
	eg(1, Pipe("\nabcdef   \n", Clean, BetweenF("a", "f"), ChompLeftF("bc")))
	// Output:
	// 1: de
}

func ExampleReplaceF() {
	eg(1, Pipe("abcdefab", ReplaceF("ab", "x", -1)))
	eg(2, Pipe("abcdefab", ReplaceF("ab", "x", 1)))
	eg(3, Pipe("abcdefab", ReplaceF("ab", "x", 0)))
	// Output:
	// 1: xcdefx
	// 2: xcdefab
	// 3: abcdefab
}

func ExampleReplacePattern() {
	eg(1, ReplacePattern("aabbcc", `a`, "x"))
	// Output:
	// 1: xxbbcc
}

func ExampleReplacePatternF() {
	eg(1, Pipe("aabbcc", ReplacePatternF(`a`, "x")))
	// Output:
	// 1: xxbbcc
}

func ExampleReverse() {
	eg(1, Reverse("abc"))
	eg(2, Reverse("中文"))
	// Output:
	// 1: cba
	// 2: 文中
}

func ExampleRight() {
	eg(1, Right("abcdef", 0))
	eg(2, Right("abcdef", 1))
	eg(3, Right("abcdef", 4))
	eg(4, Right("abcdef", -2))
	// Output:
	// 1:
	// 2: f
	// 3: cdef
	// 4: ab
}

func ExampleRightOf() {
	eg(1, RightOf("abcdef", "abc"))
	eg(2, RightOf("abcdef", "def"))
	eg(3, RightOf("abcdef", ""))
	eg(4, RightOf("", "abc"))
	eg(5, RightOf("abcdef", "xyz"))
	// Output:
	// 1: def
	// 2:
	// 3: abcdef
	// 4:
	// 5:
}

func ExampleRightF() {
	eg(1, Pipe("abcdef", RightF(3)))
	// Output:
	// 1: def
}

func ExampleSliceContains() {
	eg(1, SliceContains([]string{"foo", "bar"}, "foo"))
	eg(2, SliceContains(nil, "foo"))
	eg(3, SliceContains([]string{"foo", "bar"}, "bah"))
	eg(4, SliceContains([]string{"foo", "bar"}, ""))
	// Output:
	// 1: true
	// 2: false
	// 3: false
	// 4: false
}

func ExampleSliceIndexOf() {
	eg(1, SliceIndexOf([]string{"foo", "bar"}, "foo"))
	eg(2, SliceIndexOf(nil, "foo"))
	eg(3, SliceIndexOf([]string{"foo", "bar"}, "bah"))
	eg(4, SliceIndexOf([]string{"foo", "bar"}, ""))
	eg(5, SliceIndexOf([]string{"foo", "bar"}, "bar"))
	// Output:
	// 1: 0
	// 2: -1
	// 3: -1
	// 4: -1
	// 5: 1
}

func ExampleSlugify() {
	eg(1, Slugify("foo bar"))
	eg(2, Slugify("foo/bar bah"))
	eg(3, Slugify("foo-bar--bah"))
	// Output:
	// 1: foo-bar
	// 2: foobar-bah
	// 3: foo-bar-bah
}

func ExampleStripPunctuation() {
	eg(1, StripPunctuation("My, st[ring] *full* of %punct)"))
	// Output:
	// 1: My string full of punct
}

func ExampleStripTags() {
	eg(1, StripTags("<p>just <b>some</b> text</p>"))
	eg(2, StripTags("<p>just <b>some</b> text</p>", "p"))
	eg(3, StripTags("<a><p>just <b>some</b> text</p></a>", "a", "p"))
	eg(4, StripTags("<a><p>just <b>some</b> text</p></a>", "b"))
	// Output:
	// 1: just some text
	// 2: just <b>some</b> text
	// 3: just <b>some</b> text
	// 4: <a><p>just some text</p></a>
}

func ExampleSubstr() {
	eg(1, Substr("abcdef", 2, -1))
	eg(2, Substr("abcdef", 2, 0))
	eg(3, Substr("abcdef", 2, 1))
	eg(4, Substr("abcdef", 2, 3))
	eg(5, Substr("abcdef", 2, 4))
	eg(6, Substr("abcdef", 2, 100))
	eg(7, Substr("abcdef", 0, 1))
	// Output:
	// 1:
	// 2:
	// 3: c
	// 4: cde
	// 5: cdef
	// 6: cdef
	// 7: a
}

func ExampleTemplateWithDelimiters() {
	eg(1, TemplateWithDelimiters("Hello {{name}} at {{date-year}}", map[string]interface{}{"name": "foo", "date-year": 2014}, "{{", "}}"))
	eg(2, TemplateWithDelimiters("Hello #{name} at #{date-year}", map[string]interface{}{"name": "foo", "date-year": 2014}, "#{", "}"))
	eg(3, TemplateWithDelimiters("Hello (name) at (date-year)", map[string]interface{}{"name": "foo", "date-year": 2014}, "(", ")"))
	eg(4, TemplateWithDelimiters("Hello [name] at [date-year]", map[string]interface{}{"name": "foo", "date-year": 2014}, "[", "]"))
	eg(5, TemplateWithDelimiters("Hello *name* at *date-year*", map[string]interface{}{"name": "foo", "date-year": 2014}, "*", "*"))
	eg(6, TemplateWithDelimiters("Hello $name$ at $date-year$", map[string]interface{}{"name": "foo", "date-year": 2014}, "$", "$"))
	// Output:
	// 1: Hello foo at 2014
	// 2: Hello foo at 2014
	// 3: Hello foo at 2014
	// 4: Hello foo at 2014
	// 5: Hello foo at 2014
	// 6: Hello foo at 2014
}

func ExampleTemplate() {
	eg(1, Template("Hello {{name}} at {{date-year}}", map[string]interface{}{"name": "foo", "date-year": 2014}))
	eg(2, Template("Hello {{name}}", map[string]interface{}{"name": ""}))
	SetTemplateDelimiters("{", "}")
	eg(3, Template("Hello {name} at {date-year}", map[string]interface{}{"name": "foo", "date-year": 2014}))
	// Output:
	// 1: Hello foo at 2014
	// 2: Hello
	// 3: Hello foo at 2014
}

func ExampleToArgv() {
	eg(1, QuoteItems(ToArgv(`GO_ENV=test gosu --watch foo@release "some quoted string 'inside'"`)))
	eg(2, QuoteItems(ToArgv(`gosu foo\ bar`)))
	eg(3, QuoteItems(ToArgv(`gosu --test="some arg" -w -s a=123`)))
	// Output:
	// 1: ["GO_ENV=test" "gosu" "--watch" "foo@release" "some quoted string 'inside'"]
	// 2: ["gosu" "foo bar"]
	// 3: ["gosu" "--test=some arg" "-w" "-s" "a=123"]
}

func ExampleToBool() {
	eg(1, ToBool("true"))
	eg(2, ToBool("yes"))
	eg(3, ToBool("1"))
	eg(4, ToBool("on"))
	eg(5, ToBool("false"))
	eg(6, ToBool("no"))
	eg(7, ToBool("0"))
	eg(8, ToBool("off"))
	eg(9, ToBool(""))
	eg(10, ToBool("?"))
	// Output:
	// 1: true
	// 2: true
	// 3: true
	// 4: true
	// 5: false
	// 6: false
	// 7: false
	// 8: false
	// 9: false
	// 10: false
}

func ExampleToBoolOr() {
	eg(1, ToBoolOr("foo", true))
	eg(2, ToBoolOr("foo", false))
	eg(3, ToBoolOr("true", false))
	eg(4, ToBoolOr("", true))
	// Output:
	// 1: true
	// 2: false
	// 3: true
	// 4: true
}

func ExampleToIntOr() {
	eg(1, ToIntOr("foo", 0))
	eg(2, ToIntOr("", 1))
	eg(3, ToIntOr("100", 0))
	eg(4, ToIntOr("-1", 1))
	// Output:
	// 1: 0
	// 2: 1
	// 3: 100
	// 4: -1
}

func ExampleUnderscore() {
	eg(1, Underscore("fooBar"))
	eg(2, Underscore("FooBar"))
	eg(3, Underscore(""))
	eg(4, Underscore("x"))
	// Output:
	// 1: foo_bar
	// 2: _foo_bar
	// 3:
	// 4: x
}

func ExampleWrapHTML() {
	eg(1, WrapHTML("foo", "span", nil))
	eg(2, WrapHTML("foo", "", nil))
	eg(3, WrapHTML("foo", "", map[string]string{"class": "bar"}))
	// Output:
	// 1: <span>foo</span>
	// 2: <div>foo</div>
	// 3: <div class="bar">foo</div>
}

func ExampleWrapHTMLF() {
	eg(1, Pipe("foo", WrapHTMLF("div", nil)))
	// Output:
	// 1: <div>foo</div>
}
