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
