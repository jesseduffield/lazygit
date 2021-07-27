package terminfo

//go:generate go run gen.go

// BoolCapName returns the bool capability name.
func BoolCapName(i int) string {
	return boolCapNames[2*i]
}

// BoolCapNameShort returns the short bool capability name.
func BoolCapNameShort(i int) string {
	return boolCapNames[2*i+1]
}

// NumCapName returns the num capability name.
func NumCapName(i int) string {
	return numCapNames[2*i]
}

// NumCapNameShort returns the short num capability name.
func NumCapNameShort(i int) string {
	return numCapNames[2*i+1]
}

// StringCapName returns the string capability name.
func StringCapName(i int) string {
	return stringCapNames[2*i]
}

// StringCapNameShort returns the short string capability name.
func StringCapNameShort(i int) string {
	return stringCapNames[2*i+1]
}
