package style

import (
	"fmt"
	"strings"
)

// Render the given text as an OSC 8 hyperlink
func PrintHyperlink(text string, link string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", link, text)
}

// Render a link where the text is the same as a link
func PrintSimpleHyperlink(link string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", link, link)
}

func UnderlineLinks[T string | []byte](text T) string {
	result := ""
	remaining := string(text)
	for {
		linkStart := strings.Index(remaining, "https://")
		if linkStart == -1 {
			break
		}

		linkEnd := strings.IndexAny(remaining[linkStart:], " \n>")
		if linkEnd == -1 {
			linkEnd = len(remaining)
		} else {
			linkEnd += linkStart
		}
		underlinedLink := PrintSimpleHyperlink(remaining[linkStart:linkEnd])
		result += remaining[:linkStart] + underlinedLink
		remaining = remaining[linkEnd:]
	}
	return result + remaining
}
