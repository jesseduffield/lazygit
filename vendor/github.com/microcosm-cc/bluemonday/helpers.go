// Copyright (c) 2014, David Kitchen <david@buro9.com>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of the organisation (Microcosm) nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package bluemonday

import (
	"encoding/base64"
	"net/url"
	"regexp"
)

// A selection of regular expressions that can be used as .Matching() rules on
// HTML attributes.
var (
	// CellAlign handles the `align` attribute
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/td#attr-align
	CellAlign = regexp.MustCompile(`(?i)^(center|justify|left|right|char)$`)

	// CellVerticalAlign handles the `valign` attribute
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/td#attr-valign
	CellVerticalAlign = regexp.MustCompile(`(?i)^(baseline|bottom|middle|top)$`)

	// Direction handles the `dir` attribute
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/bdo#attr-dir
	Direction = regexp.MustCompile(`(?i)^(rtl|ltr)$`)

	// ImageAlign handles the `align` attribute on the `image` tag
	// http://www.w3.org/MarkUp/Test/Img/imgtest.html
	ImageAlign = regexp.MustCompile(
		`(?i)^(left|right|top|texttop|middle|absmiddle|baseline|bottom|absbottom)$`,
	)

	// Integer describes whole positive integers (including 0) used in places
	// like td.colspan
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/td#attr-colspan
	Integer = regexp.MustCompile(`^[0-9]+$`)

	// ISO8601 according to the W3 group is only a subset of the ISO8601
	// standard: http://www.w3.org/TR/NOTE-datetime
	//
	// Used in places like time.datetime
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/time#attr-datetime
	//
	// Matches patterns:
	//  Year:
	//     YYYY (eg 1997)
	//  Year and month:
	//     YYYY-MM (eg 1997-07)
	//  Complete date:
	//     YYYY-MM-DD (eg 1997-07-16)
	//  Complete date plus hours and minutes:
	//     YYYY-MM-DDThh:mmTZD (eg 1997-07-16T19:20+01:00)
	//  Complete date plus hours, minutes and seconds:
	//     YYYY-MM-DDThh:mm:ssTZD (eg 1997-07-16T19:20:30+01:00)
	//  Complete date plus hours, minutes, seconds and a decimal fraction of a
	//  second
	//      YYYY-MM-DDThh:mm:ss.sTZD (eg 1997-07-16T19:20:30.45+01:00)
	ISO8601 = regexp.MustCompile(
		`^[0-9]{4}(-[0-9]{2}(-[0-9]{2}([ T][0-9]{2}(:[0-9]{2}){1,2}(.[0-9]{1,6})` +
			`?Z?([\+-][0-9]{2}:[0-9]{2})?)?)?)?$`,
	)

	// ListType encapsulates the common value as well as the latest spec
	// values for lists
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/ol#attr-type
	ListType = regexp.MustCompile(`(?i)^(circle|disc|square|a|A|i|I|1)$`)

	// SpaceSeparatedTokens is used in places like `a.rel` and the common attribute
	// `class` which both contain space delimited lists of data tokens
	// http://www.w3.org/TR/html-markup/datatypes.html#common.data.tokens-def
	// Regexp: \p{L} matches unicode letters, \p{N} matches unicode numbers
	SpaceSeparatedTokens = regexp.MustCompile(`^([\s\p{L}\p{N}_-]+)$`)

	// Number is a double value used on HTML5 meter and progress elements
	// http://www.whatwg.org/specs/web-apps/current-work/multipage/the-button-element.html#the-meter-element
	Number = regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?$`)

	// NumberOrPercent is used predominantly as units of measurement in width
	// and height attributes
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/img#attr-height
	NumberOrPercent = regexp.MustCompile(`^[0-9]+[%]?$`)

	// Paragraph of text in an attribute such as *.'title', img.alt, etc
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes#attr-title
	// Note that we are not allowing chars that could close tags like '>'
	Paragraph = regexp.MustCompile(`^[\p{L}\p{N}\s\-_',\[\]!\./\\\(\)]*$`)

	// dataURIImagePrefix is used by AllowDataURIImages to define the acceptable
	// prefix of data URIs that contain common web image formats.
	//
	// This is not exported as it's not useful by itself, and only has value
	// within the AllowDataURIImages func
	dataURIImagePrefix = regexp.MustCompile(
		`^image/(gif|jpeg|png|svg\+xml|webp);base64,`,
	)
)

// AllowStandardURLs is a convenience function that will enable rel="nofollow"
// on "a", "area" and "link" (if you have allowed those elements) and will
// ensure that the URL values are parseable and either relative or belong to the
// "mailto", "http", or "https" schemes
func (p *Policy) AllowStandardURLs() {
	// URLs must be parseable by net/url.Parse()
	p.RequireParseableURLs(true)

	// !url.IsAbs() is permitted
	p.AllowRelativeURLs(true)

	// Most common URL schemes only
	p.AllowURLSchemes("mailto", "http", "https")

	// For linking elements we will add rel="nofollow" if it does not already exist
	// This applies to "a" "area" "link"
	p.RequireNoFollowOnLinks(true)
}

// AllowStandardAttributes will enable "id", "title" and the language specific
// attributes "dir" and "lang" on all elements that are allowed
func (p *Policy) AllowStandardAttributes() {
	// "dir" "lang" are permitted as both language attributes affect charsets
	// and direction of text.
	p.AllowAttrs("dir").Matching(Direction).Globally()
	p.AllowAttrs(
		"lang",
	).Matching(regexp.MustCompile(`[a-zA-Z]{2,20}`)).Globally()

	// "id" is permitted. This is pretty much as some HTML elements require this
	// to work well ("dfn" is an example of a "id" being value)
	// This does create a risk that JavaScript and CSS within your web page
	// might identify the wrong elements. Ensure that you select things
	// accurately
	p.AllowAttrs("id").Matching(
		regexp.MustCompile(`[a-zA-Z0-9\:\-_\.]+`),
	).Globally()

	// "title" is permitted as it improves accessibility.
	p.AllowAttrs("title").Matching(Paragraph).Globally()
}

// AllowStyling presently enables the class attribute globally.
//
// Note: When bluemonday ships a CSS parser and we can safely sanitise that,
// this will also allow sanitized styling of elements via the style attribute.
func (p *Policy) AllowStyling() {

	// "class" is permitted globally
	p.AllowAttrs("class").Matching(SpaceSeparatedTokens).Globally()
}

// AllowImages enables the img element and some popular attributes. It will also
// ensure that URL values are parseable. This helper does not enable data URI
// images, for that you should also use the AllowDataURIImages() helper.
func (p *Policy) AllowImages() {

	// "img" is permitted
	p.AllowAttrs("align").Matching(ImageAlign).OnElements("img")
	p.AllowAttrs("alt").Matching(Paragraph).OnElements("img")
	p.AllowAttrs("height", "width").Matching(NumberOrPercent).OnElements("img")

	// Standard URLs enabled
	p.AllowStandardURLs()
	p.AllowAttrs("src").OnElements("img")
}

// AllowDataURIImages permits the use of inline images defined in RFC2397
// http://tools.ietf.org/html/rfc2397
// http://en.wikipedia.org/wiki/Data_URI_scheme
//
// Images must have a mimetype matching:
//
//	image/gif
//	image/jpeg
//	image/png
//	image/webp
//
// NOTE: There is a potential security risk to allowing data URIs and you should
// only permit them on content you already trust.
// http://palizine.plynt.com/issues/2010Oct/bypass-xss-filters/
// https://capec.mitre.org/data/definitions/244.html
func (p *Policy) AllowDataURIImages() {

	// URLs must be parseable by net/url.Parse()
	p.RequireParseableURLs(true)

	// Supply a function to validate images contained within data URI
	p.AllowURLSchemeWithCustomPolicy(
		"data",
		func(url *url.URL) (allowUrl bool) {
			if url.RawQuery != "" || url.Fragment != "" {
				return false
			}

			matched := dataURIImagePrefix.FindString(url.Opaque)
			if matched == "" {
				return false
			}

			_, err := base64.StdEncoding.DecodeString(url.Opaque[len(matched):])
			if err != nil {
				return false
			}

			return true
		},
	)
}

// AllowLists will enabled ordered and unordered lists, as well as definition
// lists
func (p *Policy) AllowLists() {
	// "ol" "ul" are permitted
	p.AllowAttrs("type").Matching(ListType).OnElements("ol", "ul")

	// "li" is permitted
	p.AllowAttrs("type").Matching(ListType).OnElements("li")
	p.AllowAttrs("value").Matching(Integer).OnElements("li")

	// "dl" "dt" "dd" are permitted
	p.AllowElements("dl", "dt", "dd")
}

// AllowTables will enable a rich set of elements and attributes to describe
// HTML tables
func (p *Policy) AllowTables() {

	// "table" is permitted
	p.AllowAttrs("height", "width").Matching(NumberOrPercent).OnElements("table")
	p.AllowAttrs("summary").Matching(Paragraph).OnElements("table")

	// "caption" is permitted
	p.AllowElements("caption")

	// "col" "colgroup" are permitted
	p.AllowAttrs("align").Matching(CellAlign).OnElements("col", "colgroup")
	p.AllowAttrs("height", "width").Matching(
		NumberOrPercent,
	).OnElements("col", "colgroup")
	p.AllowAttrs("span").Matching(Integer).OnElements("colgroup", "col")
	p.AllowAttrs("valign").Matching(
		CellVerticalAlign,
	).OnElements("col", "colgroup")

	// "thead" "tr" are permitted
	p.AllowAttrs("align").Matching(CellAlign).OnElements("thead", "tr")
	p.AllowAttrs("valign").Matching(CellVerticalAlign).OnElements("thead", "tr")

	// "td" "th" are permitted
	p.AllowAttrs("abbr").Matching(Paragraph).OnElements("td", "th")
	p.AllowAttrs("align").Matching(CellAlign).OnElements("td", "th")
	p.AllowAttrs("colspan", "rowspan").Matching(Integer).OnElements("td", "th")
	p.AllowAttrs("headers").Matching(
		SpaceSeparatedTokens,
	).OnElements("td", "th")
	p.AllowAttrs("height", "width").Matching(
		NumberOrPercent,
	).OnElements("td", "th")
	p.AllowAttrs(
		"scope",
	).Matching(
		regexp.MustCompile(`(?i)(?:row|col)(?:group)?`),
	).OnElements("td", "th")
	p.AllowAttrs("valign").Matching(CellVerticalAlign).OnElements("td", "th")
	p.AllowAttrs("nowrap").Matching(
		regexp.MustCompile(`(?i)|nowrap`),
	).OnElements("td", "th")

	// "tbody" "tfoot"
	p.AllowAttrs("align").Matching(CellAlign).OnElements("tbody", "tfoot")
	p.AllowAttrs("valign").Matching(
		CellVerticalAlign,
	).OnElements("tbody", "tfoot")
}

func (p *Policy) AllowIFrames(vals ...SandboxValue) {
	p.AllowAttrs("sandbox").OnElements("iframe")

	p.RequireSandboxOnIFrame(vals...)
}
