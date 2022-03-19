package list

func Some[T any](slice []T, test func(T) bool) bool {
	for _, value := range slice {
		if test(value) {
			return true
		}
	}

	return false
}

func Every[T any](slice []T, test func(T) bool) bool {
	for _, value := range slice {
		if !test(value) {
			return false
		}
	}

	return true
}

func Map[T any, V any](slice []T, f func(T) V) []V {
	result := make([]V, len(slice))
	for i, value := range slice {
		result[i] = f(value)
	}

	return result
}

func MapInPlace[T any](slice []T, f func(T) T) {
	for i, value := range slice {
		slice[i] = f(value)
	}
}

func Filter[T any](slice []T, test func(T) bool) []T {
	result := make([]T, 0)
	for _, element := range slice {
		if test(element) {
			result = append(result, element)
		}
	}
	return result
}

func FilterInPlace[T any](slice []T, test func(T) bool) []T {
	newLength := 0
	for _, element := range slice {
		if test(element) {
			slice[newLength] = element
			newLength++
		}
	}

	return slice[:newLength]
}

func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i := range slice {
		result[i] = slice[len(slice)-1-i]
	}
	return result
}

func ReverseInPlace[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}
