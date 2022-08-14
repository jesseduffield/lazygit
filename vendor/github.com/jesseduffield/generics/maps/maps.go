package maps

func Keys[Key comparable, Value any](m map[Key]Value) []Key {
	keys := make([]Key, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func Values[Key comparable, Value any](m map[Key]Value) []Value {
	values := make([]Value, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

func TransformValues[Key comparable, Value any, NewValue any](
	m map[Key]Value, fn func(Value) NewValue,
) map[Key]NewValue {
	output := make(map[Key]NewValue)
	for key, value := range m {
		output[key] = fn(value)
	}
	return output
}

func TransformKeys[Key comparable, Value any, NewKey comparable](m map[Key]Value, fn func(Key) NewKey) map[NewKey]Value {
	output := make(map[NewKey]Value)
	for key, value := range m {
		output[fn(key)] = value
	}
	return output
}

func MapToSlice[Key comparable, Value any, Mapped any](m map[Key]Value, f func(Key, Value) Mapped) []Mapped {
	output := make([]Mapped, 0, len(m))
	for key, value := range m {
		output = append(output, f(key, value))
	}
	return output
}

func Filter[Key comparable, Value any](m map[Key]Value, f func(Key, Value) bool) map[Key]Value {
	output := map[Key]Value{}
	for key, value := range m {
		if f(key, value) {
			output[key] = value
		}
	}
	return output
}
