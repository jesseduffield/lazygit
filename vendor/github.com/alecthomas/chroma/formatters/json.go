package formatters

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/alecthomas/chroma"
)

// JSON formatter outputs the raw token structures as JSON.
var JSON = Register("json", chroma.FormatterFunc(func(w io.Writer, s *chroma.Style, it chroma.Iterator) error {
	fmt.Fprintln(w, "[")
	i := 0
	for t := it(); t != chroma.EOF; t = it() {
		if i > 0 {
			fmt.Fprintln(w, ",")
		}
		i++
		bytes, err := json.Marshal(t)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, "  "+string(bytes)); err != nil {
			return err
		}
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "]")
	return nil
}))
