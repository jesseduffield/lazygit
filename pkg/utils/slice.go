package utils

// NextIndex returns the index of the element that comes after the given number
func NextIndex(numbers []int, currentNumber int) int {
	for index, number := range numbers {
		if number > currentNumber {
			return index
		}
	}
	return len(numbers) - 1
}

// PrevIndex returns the index that comes before the given number, cycling if we reach the end
func PrevIndex(numbers []int, currentNumber int) int {
	end := len(numbers) - 1
	for i := end; i >= 0; i-- {
		if numbers[i] < currentNumber {
			return i
		}
	}
	return 0
}

// NextIntInCycle returns the next int in a slice, returning to the first index if we've reached the end
func NextIntInCycle(sl []int, current int) int {
	for i, val := range sl {
		if val == current {
			if i == len(sl)-1 {
				return sl[0]
			}
			return sl[i+1]
		}
	}
	return sl[0]
}

// PrevIntInCycle returns the prev int in a slice, returning to the first index if we've reached the end
func PrevIntInCycle(sl []int, current int) int {
	for i, val := range sl {
		if val == current {
			if i > 0 {
				return sl[i-1]
			}
			return sl[len(sl)-1]
		}
	}
	return sl[len(sl)-1]
}

func StringArraysOverlap(strArrA []string, strArrB []string) bool {
	for _, first := range strArrA {
		for _, second := range strArrB {
			if first == second {
				return true
			}
		}
	}

	return false
}

func Limit(values []string, limit int) []string {
	if len(values) > limit {
		return values[:limit]
	}
	return values
}

func LimitStr(value string, limit int) string {
	n := 0
	for i := range value {
		if n >= limit {
			return value[:i]
		}
		n++
	}
	return value
}
