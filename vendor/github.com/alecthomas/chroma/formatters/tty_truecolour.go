package formatters

import (
	"fmt"
	"io"

	"github.com/alecthomas/chroma"
)

// TTY16m is a true-colour terminal formatter.
var TTY16m = Register("terminal16m", chroma.FormatterFunc(trueColourFormatter))

func trueColourFormatter(w io.Writer, style *chroma.Style, it chroma.Iterator) error {
	style = clearBackground(style)
	for token := it(); token != chroma.EOF; token = it() {
		entry := style.Get(token.Type)
		if !entry.IsZero() {
			out := ""
			if entry.Bold == chroma.Yes {
				out += "\033[1m"
			}
			if entry.Underline == chroma.Yes {
				out += "\033[4m"
			}
			if entry.Italic == chroma.Yes {
				out += "\033[3m"
			}
			if entry.Colour.IsSet() {
				out += fmt.Sprintf("\033[38;2;%d;%d;%dm", entry.Colour.Red(), entry.Colour.Green(), entry.Colour.Blue())
			}
			if entry.Background.IsSet() {
				out += fmt.Sprintf("\033[48;2;%d;%d;%dm", entry.Background.Red(), entry.Background.Green(), entry.Background.Blue())
			}
			fmt.Fprint(w, out)
		}
		fmt.Fprint(w, token.Value)
		if !entry.IsZero() {
			fmt.Fprint(w, "\033[0m")
		}
	}
	return nil
}
