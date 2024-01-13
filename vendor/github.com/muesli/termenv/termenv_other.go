//go:build js || plan9 || aix
// +build js plan9 aix

package termenv

func colorProfile() Profile {
	return ANSI256
}

func (o Output) foregroundColor() Color {
	// default gray
	return ANSIColor(7)
}

func (o Output) backgroundColor() Color {
	// default black
	return ANSIColor(0)
}
