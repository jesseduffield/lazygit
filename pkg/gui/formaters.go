package gui

// formatGit72 after the git best practices where
// the line should newline at max 72 columns. It formats
// similarly to how git does by breaking the row on the ' '
// before the word that's overflowing
func formatGit72(content []rune) ([]rune, error) {
	// Max length is the same as content length, as only the space
	// at the end of a word may be converted into a newline
	cpy := make([]rune, len(content))
	c := 0
	prevSpace := -1
	for i := 0; i < len(content); i++ {
		nextRune := content[i]
		nextIsNewline := nextRune == '\n'
		if nextRune == ' ' {
			prevSpace = i
		}
		if c == 72 && !nextIsNewline {
			cpy[prevSpace] = '\n'
			// set next as newline, as it just has been injected
			nextIsNewline = true
		} else {
			cpy[i] = nextRune
		}
		if nextIsNewline {
			c = 0
		} else {
			c++
		}
	}
	return cpy, nil
}
