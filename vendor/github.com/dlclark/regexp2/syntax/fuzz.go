// +build gofuzz

package syntax

// Fuzz is the input point for go-fuzz
func Fuzz(data []byte) int {
	sdata := string(data)
	tree, err := Parse(sdata, RegexOptions(0))
	if err != nil {
		return 0
	}

	// translate it to code
	_, err = Write(tree)
	if err != nil {
		panic(err)
	}

	return 1
}
