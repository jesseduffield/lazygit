// Package svg contains an SVG formatter.
package svg

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"

	"github.com/alecthomas/chroma"
)

// Option sets an option of the SVG formatter.
type Option func(f *Formatter)

// FontFamily sets the font-family.
func FontFamily(fontFamily string) Option { return func(f *Formatter) { f.fontFamily = fontFamily } }

// EmbedFontFile embeds given font file
func EmbedFontFile(fontFamily string, fileName string) (option Option, err error) {
	var format FontFormat
	switch path.Ext(fileName) {
	case ".woff":
		format = WOFF
	case ".woff2":
		format = WOFF2
	case ".ttf":
		format = TRUETYPE
	default:
		return nil, errors.New("unexpected font file suffix")
	}

	var content []byte
	if content, err = ioutil.ReadFile(fileName); err == nil {
		option = EmbedFont(fontFamily, base64.StdEncoding.EncodeToString(content), format)
	}
	return
}

// EmbedFont embeds given base64 encoded font
func EmbedFont(fontFamily string, font string, format FontFormat) Option {
	return func(f *Formatter) { f.fontFamily = fontFamily; f.embeddedFont = font; f.fontFormat = format }
}

// New SVG formatter.
func New(options ...Option) *Formatter {
	f := &Formatter{fontFamily: "Consolas, Monaco, Lucida Console, Liberation Mono, DejaVu Sans Mono, Bitstream Vera Sans Mono, Courier New, monospace"}
	for _, option := range options {
		option(f)
	}
	return f
}

// Formatter that generates SVG.
type Formatter struct {
	fontFamily   string
	embeddedFont string
	fontFormat   FontFormat
}

func (f *Formatter) Format(w io.Writer, style *chroma.Style, iterator chroma.Iterator) (err error) {
	f.writeSVG(w, style, iterator.Tokens())
	return err
}

var svgEscaper = strings.NewReplacer(
	`&`, "&amp;",
	`<`, "&lt;",
	`>`, "&gt;",
	`"`, "&quot;",
	` `, "&#160;",
	`	`, "&#160;&#160;&#160;&#160;",
)

// EscapeString escapes special characters.
func escapeString(s string) string {
	return svgEscaper.Replace(s)
}

func (f *Formatter) writeSVG(w io.Writer, style *chroma.Style, tokens []chroma.Token) { // nolint: gocyclo
	svgStyles := f.styleToSVG(style)
	lines := chroma.SplitTokensIntoLines(tokens)

	fmt.Fprint(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	fmt.Fprint(w, "<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.0//EN\" \"http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd\">\n")
	fmt.Fprintf(w, "<svg width=\"%dpx\" height=\"%dpx\" xmlns=\"http://www.w3.org/2000/svg\">\n", 8*maxLineWidth(lines), 10+int(16.8*float64(len(lines)+1)))

	if f.embeddedFont != "" {
		f.writeFontStyle(w)
	}

	fmt.Fprintf(w, "<rect width=\"100%%\" height=\"100%%\" fill=\"%s\"/>\n", style.Get(chroma.Background).Background.String())
	fmt.Fprintf(w, "<g font-family=\"%s\" font-size=\"14px\" fill=\"%s\">\n", f.fontFamily, style.Get(chroma.Text).Colour.String())

	f.writeTokenBackgrounds(w, lines, style)

	for index, tokens := range lines {
		fmt.Fprintf(w, "<text x=\"0\" y=\"%fem\" xml:space=\"preserve\">", 1.2*float64(index+1))

		for _, token := range tokens {
			text := escapeString(token.String())
			attr := f.styleAttr(svgStyles, token.Type)
			if attr != "" {
				text = fmt.Sprintf("<tspan %s>%s</tspan>", attr, text)
			}
			fmt.Fprint(w, text)
		}
		fmt.Fprint(w, "</text>")
	}

	fmt.Fprint(w, "\n</g>\n")
	fmt.Fprint(w, "</svg>\n")
}

func maxLineWidth(lines [][]chroma.Token) int {
	maxWidth := 0
	for _, tokens := range lines {
		length := 0
		for _, token := range tokens {
			length += len(strings.ReplaceAll(token.String(), `	`, "    "))
		}
		if length > maxWidth {
			maxWidth = length
		}
	}
	return maxWidth
}

// There is no background attribute for text in SVG so simply calculate the position and text
// of tokens with a background color that differs from the default and add a rectangle for each before
// adding the token.
func (f *Formatter) writeTokenBackgrounds(w io.Writer, lines [][]chroma.Token, style *chroma.Style) {
	for index, tokens := range lines {
		lineLength := 0
		for _, token := range tokens {
			length := len(strings.ReplaceAll(token.String(), `	`, "    "))
			tokenBackground := style.Get(token.Type).Background
			if tokenBackground.IsSet() && tokenBackground != style.Get(chroma.Background).Background {
				fmt.Fprintf(w, "<rect id=\"%s\" x=\"%dch\" y=\"%fem\" width=\"%dch\" height=\"1.2em\" fill=\"%s\" />\n", escapeString(token.String()), lineLength, 1.2*float64(index)+0.25, length, style.Get(token.Type).Background.String())
			}
			lineLength += length
		}
	}
}

type FontFormat int

// https://transfonter.org/formats
const (
	WOFF FontFormat = iota
	WOFF2
	TRUETYPE
)

var fontFormats = [...]string{
	"woff",
	"woff2",
	"truetype",
}

func (f *Formatter) writeFontStyle(w io.Writer) {
	fmt.Fprintf(w, `<style>
@font-face {
	font-family: '%s';
	src: url(data:application/x-font-%s;charset=utf-8;base64,%s) format('%s');'
	font-weight: normal;
	font-style: normal;
}
</style>`, f.fontFamily, fontFormats[f.fontFormat], f.embeddedFont, fontFormats[f.fontFormat])
}

func (f *Formatter) styleAttr(styles map[chroma.TokenType]string, tt chroma.TokenType) string {
	if _, ok := styles[tt]; !ok {
		tt = tt.SubCategory()
		if _, ok := styles[tt]; !ok {
			tt = tt.Category()
			if _, ok := styles[tt]; !ok {
				return ""
			}
		}
	}
	return styles[tt]
}

func (f *Formatter) styleToSVG(style *chroma.Style) map[chroma.TokenType]string {
	converted := map[chroma.TokenType]string{}
	bg := style.Get(chroma.Background)
	// Convert the style.
	for t := range chroma.StandardTypes {
		entry := style.Get(t)
		if t != chroma.Background {
			entry = entry.Sub(bg)
		}
		if entry.IsZero() {
			continue
		}
		converted[t] = StyleEntryToSVG(entry)
	}
	return converted
}

// StyleEntryToSVG converts a chroma.StyleEntry to SVG attributes.
func StyleEntryToSVG(e chroma.StyleEntry) string {
	var styles []string

	if e.Colour.IsSet() {
		styles = append(styles, "fill=\""+e.Colour.String()+"\"")
	}
	if e.Bold == chroma.Yes {
		styles = append(styles, "font-weight=\"bold\"")
	}
	if e.Italic == chroma.Yes {
		styles = append(styles, "font-style=\"italic\"")
	}
	if e.Underline == chroma.Yes {
		styles = append(styles, "text-decoration=\"underline\"")
	}
	return strings.Join(styles, " ")
}
