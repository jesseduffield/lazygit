# regexp2 - full featured regular expressions for Go
Regexp2 is a feature-rich RegExp engine for Go.  It doesn't have constant time guarantees like the built-in `regexp` package, but it allows backtracking and is compatible with Perl5 and .NET.  You'll likely be better off with the RE2 engine from the `regexp` package and should only use this if you need to write very complex patterns or require compatibility with .NET.

## Basis of the engine
The engine is ported from the .NET framework's System.Text.RegularExpressions.Regex engine.  That engine was open sourced in 2015 under the MIT license.  There are some fundamental differences between .NET strings and Go strings that required a bit of borrowing from the Go framework regex engine as well.  I cleaned up a couple of the dirtier bits during the port (regexcharclass.cs was terrible), but the parse tree, code emmitted, and therefore patterns matched should be identical.

## Installing
This is a go-gettable library, so install is easy:

    go get github.com/dlclark/regexp2/...

## Usage
Usage is similar to the Go `regexp` package.  Just like in `regexp`, you start by converting a regex into a state machine via the `Compile` or `MustCompile` methods.  They ultimately do the same thing, but `MustCompile` will panic if the regex is invalid.  You can then use the provided `Regexp` struct to find matches repeatedly.  A `Regexp` struct is safe to use across goroutines.

```go
re := regexp2.MustCompile(`Your pattern`, 0)
if isMatch, _ := re.MatchString(`Something to match`); isMatch {
    //do something
}
```

The only error that the `*Match*` methods *should* return is a Timeout if you set the `re.MatchTimeout` field.  Any other error is a bug in the `regexp2` package.  If you need more details about capture groups in a match then use the `FindStringMatch` method, like so:

```go
if m, _ := re.FindStringMatch(`Something to match`); m != nil {
    // the whole match is always group 0
    fmt.Printf("Group 0: %v\n", m.String())

    // you can get all the groups too
    gps := m.Groups()

    // a group can be captured multiple times, so each cap is separately addressable
    fmt.Printf("Group 1, first capture", gps[1].Captures[0].String())
    fmt.Printf("Group 1, second capture", gps[1].Captures[1].String())
}
```

Group 0 is embedded in the Match.  Group 0 is an automatically-assigned group that encompasses the whole pattern.  This means that `m.String()` is the same as `m.Group.String()` and `m.Groups()[0].String()`

The __last__ capture is embedded in each group, so `g.String()` will return the same thing as `g.Capture.String()` and  `g.Captures[len(g.Captures)-1].String()`.

## Compare `regexp` and `regexp2`
| Category | regexp | regexp2 |
| --- | --- | --- |
| Catastrophic backtracking possible | no, constant execution time guarantees | yes, if your pattern is at risk you can use the `re.MatchTimeout` field |
| Python-style capture groups `(?P<name>re)` | yes | no (yes in RE2 compat mode) |
| .NET-style capture groups `(?<name>re)` or `(?'name're)` | no | yes |
| comments `(?#comment)` | no | yes |
| branch numbering reset `(?\|a\|b)` | no | no |
| possessive match `(?>re)` | no | yes |
| positive lookahead `(?=re)` | no | yes |
| negative lookahead `(?!re)` | no | yes |
| positive lookbehind `(?<=re)` | no | yes |
| negative lookbehind `(?<!re)` | no | yes |
| back reference `\1` | no | yes |
| named back reference `\k'name'` | no | yes |
| named ascii character class `[[:foo:]]`| yes | no (yes in RE2 compat mode) |
| conditionals `(?(expr)yes\|no)` | no | yes |

## RE2 compatibility mode
The default behavior of `regexp2` is to match the .NET regexp engine, however the `RE2` option is provided to change the parsing to increase compatibility with RE2.  Using the `RE2` option when compiling a regexp will not take away any features, but will change the following behaviors:
* add support for named ascii character classes (e.g. `[[:foo:]]`)
* add support for python-style capture groups (e.g. `(P<name>re)`)
* change singleline behavior for `$` to only match end of string (like RE2) (see [#24](https://github.com/dlclark/regexp2/issues/24))
 
```go
re := regexp2.MustCompile(`Your RE2-compatible pattern`, regexp2.RE2)
if isMatch, _ := re.MatchString(`Something to match`); isMatch {
    //do something
}
```

This feature is a work in progress and I'm open to ideas for more things to put here (maybe more relaxed character escaping rules?).


## Library features that I'm still working on
- Regex split

## Potential bugs
I've run a battery of tests against regexp2 from various sources and found the debug output matches the .NET engine, but .NET and Go handle strings very differently.  I've attempted to handle these differences, but most of my testing deals with basic ASCII with a little bit of multi-byte Unicode.  There's a chance that there are bugs in the string handling related to character sets with supplementary Unicode chars.  Right-to-Left support is coded, but not well tested either.

## Find a bug?
I'm open to new issues and pull requests with tests if you find something odd!
