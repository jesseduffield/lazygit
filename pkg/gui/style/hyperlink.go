package style

import "fmt"

// Render the given text as an OSC 8 hyperlink
func PrintHyperlink(text string, link string, underline bool) string {
	result := fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", link, text)
	if underline {
		return AttrUnderline.Sprint(result)
	}
	return result
}

// Render a link where the text is the same as a link
func PrintSimpleHyperlink(link string, underline bool) string {
	result := fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", link, link)
	if underline {
		return AttrUnderline.Sprint(result)
	}
	return result
}
