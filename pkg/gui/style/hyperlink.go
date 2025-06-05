package style

import "fmt"

// Render the given text as an OSC 8 hyperlink
func PrintHyperlink(text string, link string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", link, text)
}

// Render a link where the text is the same as a link
func PrintSimpleHyperlink(link string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", link, link)
}
