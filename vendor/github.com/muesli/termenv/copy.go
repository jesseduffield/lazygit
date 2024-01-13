package termenv

import (
	"github.com/aymanbagabas/go-osc52"
)

// Copy copies text to clipboard using OSC 52 escape sequence.
func (o Output) Copy(str string) {
	out := osc52.NewOutput(o.tty, o.environ.Environ())
	out.Copy(str)
}

// Copy copies text to clipboard using OSC 52 escape sequence.
func Copy(str string) {
	output.Copy(str)
}
