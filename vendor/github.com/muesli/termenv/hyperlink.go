package termenv

import (
	"fmt"
)

// Hyperlink creates a hyperlink using OSC8.
func Hyperlink(link, name string) string {
	return output.Hyperlink(link, name)
}

// Hyperlink creates a hyperlink using OSC8.
func (o *Output) Hyperlink(link, name string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", link, name)
}
