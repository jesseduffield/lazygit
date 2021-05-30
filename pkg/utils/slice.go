package utils

// IncludesString if the list contains the string
func IncludesString(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// IncludesInt if the list contains the Int
func IncludesInt(list []int, a int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

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

// UnionInt returns the union of two int arrays
func UnionInt(a, b []int) []int {
	m := make(map[int]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			// this does not mutate the original a slice
			// though it does mutate the backing array I believe
			// but that doesn't matter because if you later want to append to the
			// original a it must see that the backing array has been changed
			// and create a new one
			a = append(a, item)
		}
	}
	return a
}

// DifferenceInt returns the difference of two int arrays
func DifferenceInt(a, b []int) []int {
	result := []int{}
	m := make(map[int]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			result = append(result, item)
		}
	}
	return result
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
