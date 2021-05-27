package gocui

import "time"

func (v *View) loaderLines() [][]cell {
	duplicate := make([][]cell, len(v.lines))
	for i := range v.lines {
		if i < len(v.lines)-1 {
			duplicate[i] = make([]cell, len(v.lines[i]))
			copy(duplicate[i], v.lines[i])
		} else {
			duplicate[i] = make([]cell, len(v.lines[i])+2)
			copy(duplicate[i], v.lines[i])
			duplicate[i][len(duplicate[i])-2] = cell{chr: ' '}
			duplicate[i][len(duplicate[i])-1] = Loader()
		}
	}

	return duplicate
}

// Loader can show a loading animation
func Loader() cell {
	characters := "|/-\\"
	now := time.Now()
	nanos := now.UnixNano()
	index := nanos / 50000000 % int64(len(characters))
	str := characters[index : index+1]
	chr := []rune(str)[0]
	return cell{
		chr: chr,
	}
}
