# str

    import "github.com/mgutz/str"

Package str is a comprehensive set of string functions to build more Go
awesomeness. Str complements Go's standard packages and does not duplicate
functionality found in `strings` or `strconv`.

Str is based on plain functions instead of object-based methods, consistent with
Go standard string packages.

    str.Between("<a>foo</a>", "<a>", "</a>") == "foo"

Str supports pipelining instead of chaining

    s := str.Pipe("\nabcdef\n", Clean, BetweenF("a", "f"), ChompLeftF("bc"))

User-defined filters can be added to the pipeline by inserting a function or
closure that returns a function with this signature

    func(string) string

### Index

* [Variables](#variables)
* [func  Between](#func
[godoc](https://godoc.org/github.com/mgutz/str)
between)
* [func  BetweenF](#func--betweenf)
* [func  Camelize](#func--camelize)
* [func  Capitalize](#func--capitalize)
* [func  CharAt](#func--charat)
* [func  CharAtF](#func--charatf)
* [func  ChompLeft](#func--chompleft)
* [func  ChompLeftF](#func--chompleftf)
* [func  ChompRight](#func--chompright)
* [func  ChompRightF](#func--chomprightf)
* [func  Classify](#func--classify)
* [func  ClassifyF](#func--classifyf)
* [func  Clean](#func--clean)
* [func  Dasherize](#func--dasherize)
* [func  DecodeHTMLEntities](#func--decodehtmlentities)
* [func  EnsurePrefix](#func--ensureprefix)
* [func  EnsurePrefixF](#func--ensureprefixf)
* [func  EnsureSuffix](#func--ensuresuffix)
* [func  EnsureSuffixF](#func--ensuresuffixf)
* [func  EscapeHTML](#func--escapehtml)
* [func  Humanize](#func--humanize)
* [func  Iif](#func--iif)
* [func  IndexOf](#func--indexof)
* [func  IsAlpha](#func--isalpha)
* [func  IsAlphaNumeric](#func--isalphanumeric)
* [func  IsEmpty](#func--isempty)
* [func  IsLower](#func--islower)
* [func  IsNumeric](#func--isnumeric)
* [func  IsUpper](#func--isupper)
* [func  Left](#func--left)
* [func  LeftF](#func--leftf)
* [func  LeftOf](#func--leftof)
* [func  Letters](#func--letters)
* [func  Lines](#func--lines)
* [func  Map](#func--map)
* [func  Match](#func--match)
* [func  Pad](#func--pad)
* [func  PadF](#func--padf)
* [func  PadLeft](#func--padleft)
* [func  PadLeftF](#func--padleftf)
* [func  PadRight](#func--padright)
* [func  PadRightF](#func--padrightf)
* [func  Pipe](#func--pipe)
* [func  QuoteItems](#func--quoteitems)
* [func  ReplaceF](#func--replacef)
* [func  ReplacePattern](#func--replacepattern)
* [func  ReplacePatternF](#func--replacepatternf)
* [func  Reverse](#func--reverse)
* [func  Right](#func--right)
* [func  RightF](#func--rightf)
* [func  RightOf](#func--rightof)
* [func  SetTemplateDelimiters](#func--settemplatedelimiters)
* [func  Slice](#func--slice)
* [func  SliceContains](#func--slicecontains)
* [func  SliceF](#func--slicef)
* [func  SliceIndexOf](#func--sliceindexof)
* [func  Slugify](#func--slugify)
* [func  StripPunctuation](#func--strippunctuation)
* [func  StripTags](#func--striptags)
* [func  Substr](#func--substr)
* [func  SubstrF](#func--substrf)
* [func  Template](#func--template)
* [func  TemplateDelimiters](#func--templatedelimiters)
* [func  TemplateWithDelimiters](#func--templatewithdelimiters)
* [func  ToArgv](#func--toargv)
* [func  ToBool](#func--tobool)
* [func  ToBoolOr](#func--toboolor)
* [func  ToFloat32Or](#func--tofloat32or)
* [func  ToFloat64Or](#func--tofloat64or)
* [func  ToIntOr](#func--tointor)
* [func  Underscore](#func--underscore)
* [func  UnescapeHTML](#func--unescapehtml)
* [func  WrapHTML](#func--wraphtml)
* [func  WrapHTMLF](#func--wraphtmlf)


#### Variables

```go
var ToFloatOr = ToFloat64Or
```
ToFloatOr parses as a float64 or returns defaultValue.

```go
var Verbose = false
```
Verbose flag enables console output for those functions that have counterparts
in Go's excellent stadard packages.

#### func  [Between](#between)

```go
func Between(s, left, right string) string
```
Between extracts a string between left and right strings.

#### func  [BetweenF](#betweenf)

```go
func BetweenF(left, right string) func(string) string
```
BetweenF is the filter form for Between.

#### func  [Camelize](#camelize)

```go
func Camelize(s string) string
```
Camelize return new string which removes any underscores or dashes and convert a
string into camel casing.

#### func  [Capitalize](#capitalize)

```go
func Capitalize(s string) string
```
Capitalize uppercases the first char of s and lowercases the rest.

#### func  [CharAt](#charat)

```go
func CharAt(s string, index int) string
```
CharAt returns a string from the character at the specified position.

#### func  [CharAtF](#charatf)

```go
func CharAtF(index int) func(string) string
```
CharAtF is the filter form of CharAt.

#### func  [ChompLeft](#chompleft)

```go
func ChompLeft(s, prefix string) string
```
ChompLeft removes prefix at the start of a string.

#### func  [ChompLeftF](#chompleftf)

```go
func ChompLeftF(prefix string) func(string) string
```
ChompLeftF is the filter form of ChompLeft.

#### func  [ChompRight](#chompright)

```go
func ChompRight(s, suffix string) string
```
ChompRight removes suffix from end of s.

#### func  [ChompRightF](#chomprightf)

```go
func ChompRightF(suffix string) func(string) string
```
ChompRightF is the filter form of ChompRight.

#### func  [Classify](#classify)

```go
func Classify(s string) string
```
Classify returns a camelized string with the first letter upper cased.

#### func  [ClassifyF](#classifyf)

```go
func ClassifyF(s string) func(string) string
```
ClassifyF is the filter form of Classify.

#### func  [Clean](#clean)

```go
func Clean(s string) string
```
Clean compresses all adjacent whitespace to a single space and trims s.

#### func  [Dasherize](#dasherize)

```go
func Dasherize(s string) string
```
Dasherize converts a camel cased string into a string delimited by dashes.

#### func  [DecodeHTMLEntities](#decodehtmlentities)

```go
func DecodeHTMLEntities(s string) string
```
DecodeHTMLEntities decodes HTML entities into their proper string
representation. DecodeHTMLEntities is an alias for html.UnescapeString

#### func  [EnsurePrefix](#ensureprefix)

```go
func EnsurePrefix(s, prefix string) string
```
EnsurePrefix ensures s starts with prefix.

#### func  [EnsurePrefixF](#ensureprefixf)

```go
func EnsurePrefixF(prefix string) func(string) string
```
EnsurePrefixF is the filter form of EnsurePrefix.

#### func  [EnsureSuffix](#ensuresuffix)

```go
func EnsureSuffix(s, suffix string) string
```
EnsureSuffix ensures s ends with suffix.

#### func  [EnsureSuffixF](#ensuresuffixf)

```go
func EnsureSuffixF(suffix string) func(string) string
```
EnsureSuffixF is the filter form of EnsureSuffix.

#### func  [EscapeHTML](#escapehtml)

```go
func EscapeHTML(s string) string
```
EscapeHTML is alias for html.EscapeString.

#### func  [Humanize](#humanize)

```go
func Humanize(s string) string
```
Humanize transforms s into a human friendly form.

#### func  [Iif](#iif)

```go
func Iif(condition bool, truthy string, falsey string) string
```
Iif is short for immediate if. If condition is true return truthy else falsey.

#### func  [IndexOf](#indexof)

```go
func IndexOf(s string, needle string, start int) int
```
IndexOf finds the index of needle in s starting from start.

#### func  [IsAlpha](#isalpha)

```go
func IsAlpha(s string) bool
```
IsAlpha returns true if a string contains only letters from ASCII (a-z,A-Z).
Other letters from other languages are not supported.

#### func  [IsAlphaNumeric](#isalphanumeric)

```go
func IsAlphaNumeric(s string) bool
```
IsAlphaNumeric returns true if a string contains letters and digits.

#### func  [IsEmpty](#isempty)

```go
func IsEmpty(s string) bool
```
IsEmpty returns true if the string is solely composed of whitespace.

#### func  [IsLower](#islower)

```go
func IsLower(s string) bool
```
IsLower returns true if s comprised of all lower case characters.

#### func  [IsNumeric](#isnumeric)

```go
func IsNumeric(s string) bool
```
IsNumeric returns true if a string contains only digits from 0-9. Other digits
not in Latin (such as Arabic) are not currently supported.

#### func  [IsUpper](#isupper)

```go
func IsUpper(s string) bool
```
IsUpper returns true if s contains all upper case chracters.

#### func  [Left](#left)

```go
func Left(s string, n int) string
```
Left returns the left substring of length n.

#### func  [LeftF](#leftf)

```go
func LeftF(n int) func(string) string
```
LeftF is the filter form of Left.

#### func  [LeftOf](#leftof)

```go
func LeftOf(s string, needle string) string
```
LeftOf returns the substring left of needle.

#### func  [Letters](#letters)

```go
func Letters(s string) []string
```
Letters returns an array of runes as strings so it can be indexed into.

#### func  [Lines](#lines)

```go
func Lines(s string) []string
```
Lines convert windows newlines to unix newlines then convert to an Array of
lines.

#### func  [Map](#map)

```go
func Map(arr []string, iterator func(string) string) []string
```
Map maps an array's iitem through an iterator.

#### func  [Match](#match)

```go
func Match(s, pattern string) bool
```
Match returns true if patterns matches the string

#### func  [Pad](#pad)

```go
func Pad(s, c string, n int) string
```
Pad pads string s on both sides with c until it has length of n.

#### func  [PadF](#padf)

```go
func PadF(c string, n int) func(string) string
```
PadF is the filter form of Pad.

#### func  [PadLeft](#padleft)

```go
func PadLeft(s, c string, n int) string
```
PadLeft pads s on left side with c until it has length of n.

#### func  [PadLeftF](#padleftf)

```go
func PadLeftF(c string, n int) func(string) string
```
PadLeftF is the filter form of PadLeft.

#### func  [PadRight](#padright)

```go
func PadRight(s, c string, n int) string
```
PadRight pads s on right side with c until it has length of n.

#### func  [PadRightF](#padrightf)

```go
func PadRightF(c string, n int) func(string) string
```
PadRightF is the filter form of Padright

#### func  [Pipe](#pipe)

```go
func Pipe(s string, funcs ...func(string) string) string
```
Pipe pipes s through one or more string filters.

#### func  [QuoteItems](#quoteitems)

```go
func QuoteItems(arr []string) []string
```
QuoteItems quotes all items in array, mostly for debugging.

#### func  [ReplaceF](#replacef)

```go
func ReplaceF(old, new string, n int) func(string) string
```
ReplaceF is the filter form of strings.Replace.

#### func  [ReplacePattern](#replacepattern)

```go
func ReplacePattern(s, pattern, repl string) string
```
ReplacePattern replaces string with regexp string. ReplacePattern returns a copy
of src, replacing matches of the Regexp with the replacement string repl. Inside
repl, $ signs are interpreted as in Expand, so for instance $1 represents the
text of the first submatch.

#### func  [ReplacePatternF](#replacepatternf)

```go
func ReplacePatternF(pattern, repl string) func(string) string
```
ReplacePatternF is the filter form of ReplaceRegexp.

#### func  [Reverse](#reverse)

```go
func Reverse(s string) string
```
Reverse a string

#### func  [Right](#right)

```go
func Right(s string, n int) string
```
Right returns the right substring of length n.

#### func  [RightF](#rightf)

```go
func RightF(n int) func(string) string
```
RightF is the Filter version of Right.

#### func  [RightOf](#rightof)

```go
func RightOf(s string, prefix string) string
```
RightOf returns the substring to the right of prefix.

#### func  [SetTemplateDelimiters](#settemplatedelimiters)

```go
func SetTemplateDelimiters(opening, closing string)
```
SetTemplateDelimiters sets the delimiters for Template function. Defaults to
"{{" and "}}"

#### func  [Slice](#slice)

```go
func Slice(s string, start, end int) string
```
Slice slices a string. If end is negative then it is the from the end of the
string.

#### func  [SliceContains](#slicecontains)

```go
func SliceContains(slice []string, val string) bool
```
SliceContains determines whether val is an element in slice.

#### func  [SliceF](#slicef)

```go
func SliceF(start, end int) func(string) string
```
SliceF is the filter for Slice.

#### func  [SliceIndexOf](#sliceindexof)

```go
func SliceIndexOf(slice []string, val string) int
```
SliceIndexOf gets the indx of val in slice. Returns -1 if not found.

#### func  [Slugify](#slugify)

```go
func Slugify(s string) string
```
Slugify converts s into a dasherized string suitable for URL segment.

#### func  [StripPunctuation](#strippunctuation)

```go
func StripPunctuation(s string) string
```
StripPunctuation strips puncation from string.

#### func  [StripTags](#striptags)

```go
func StripTags(s string, tags ...string) string
```
StripTags strips all of the html tags or tags specified by the parameters

#### func  [Substr](#substr)

```go
func Substr(s string, index int, n int) string
```
Substr returns a substring of s starting at index of length n.

#### func  [SubstrF](#substrf)

```go
func SubstrF(index, n int) func(string) string
```
SubstrF is the filter form of Substr.

#### func  [Template](#template)

```go
func Template(s string, values map[string]interface{}) string
```
Template is a string template which replaces template placeholders delimited by
"{{" and "}}" with values from map. The global delimiters may be set with
SetTemplateDelimiters.

#### func  [TemplateDelimiters](#templatedelimiters)

```go
func TemplateDelimiters() (opening string, closing string)
```
TemplateDelimiters is the getter for the opening and closing delimiters for
Template.

#### func  [TemplateWithDelimiters](#templatewithdelimiters)

```go
func TemplateWithDelimiters(s string, values map[string]interface{}, opening, closing string) string
```
TemplateWithDelimiters is string template with user-defineable opening and
closing delimiters.

#### func  [ToArgv](#toargv)

```go
func ToArgv(s string) []string
```
ToArgv converts string s into an argv for exec.

#### func  [ToBool](#tobool)

```go
func ToBool(s string) bool
```
ToBool fuzzily converts truthy values.

#### func  [ToBoolOr](#toboolor)

```go
func ToBoolOr(s string, defaultValue bool) bool
```
ToBoolOr parses s as a bool or returns defaultValue.

#### func  [ToFloat32Or](#tofloat32or)

```go
func ToFloat32Or(s string, defaultValue float32) float32
```
ToFloat32Or parses as a float32 or returns defaultValue on error.

#### func  [ToFloat64Or](#tofloat64or)

```go
func ToFloat64Or(s string, defaultValue float64) float64
```
ToFloat64Or parses s as a float64 or returns defaultValue.

#### func  [ToIntOr](#tointor)

```go
func ToIntOr(s string, defaultValue int) int
```
ToIntOr parses s as an int or returns defaultValue.

#### func  [Underscore](#underscore)

```go
func Underscore(s string) string
```
Underscore returns converted camel cased string into a string delimited by
underscores.

#### func  [UnescapeHTML](#unescapehtml)

```go
func UnescapeHTML(s string) string
```
UnescapeHTML is an alias for html.UnescapeString.

#### func  [WrapHTML](#wraphtml)

```go
func WrapHTML(s string, tag string, attrs map[string]string) string
```
WrapHTML wraps s within HTML tag having attributes attrs. Note, WrapHTML does
not escape s value.

#### func  [WrapHTMLF](#wraphtmlf)

```go
func WrapHTMLF(tag string, attrs map[string]string) func(string) string
```
WrapHTMLF is the filter form of WrapHTML.
