package jsonparser

func FuzzParseString(data []byte) int {
	r, err := ParseString(data)
	if err != nil || r == "" {
		return 0
	}
	return 1
}

func FuzzEachKey(data []byte) int {
	paths := [][]string{
		{"name"},
		{"order"},
		{"nested", "a"},
		{"nested", "b"},
		{"nested2", "a"},
		{"nested", "nested3", "b"},
		{"arr", "[1]", "b"},
		{"arrInt", "[3]"},
		{"arrInt", "[5]"},
		{"nested"},
		{"arr", "["},
		{"a\n", "b\n"},
	}
	EachKey(data, func(idx int, value []byte, vt ValueType, err error) {}, paths...)
	return 1
}

func FuzzDelete(data []byte) int {
	Delete(data, "test")
	return 1
}

func FuzzSet(data []byte) int {
	_, err := Set(data, []byte(`"new value"`), "test")
	if err != nil {
		return 0
	}
	return 1
}

func FuzzObjectEach(data []byte) int {
	_ = ObjectEach(data, func(key, value []byte, valueType ValueType, off int) error {
		return nil
	})
	return 1
}

func FuzzParseFloat(data []byte) int {
	_, err := ParseFloat(data)
	if err != nil {
		return 0
	}
	return 1
}

func FuzzParseInt(data []byte) int {
	_, err := ParseInt(data)
	if err != nil {
		return 0
	}
	return 1
}

func FuzzParseBool(data []byte) int {
	_, err := ParseBoolean(data)
	if err != nil {
		return 0
	}
	return 1
}

func FuzzTokenStart(data []byte) int {
	_ = tokenStart(data)
	return 1
}

func FuzzGetString(data []byte) int {
	_, err := GetString(data, "test")
	if err != nil {
		return 0
	}
	return 1
}

func FuzzGetFloat(data []byte) int {
	_, err := GetFloat(data, "test")
	if err != nil {
		return 0
	}
	return 1
}

func FuzzGetInt(data []byte) int {
	_, err := GetInt(data, "test")
	if err != nil {
		return 0
	}
	return 1
}

func FuzzGetBoolean(data []byte) int {
	_, err := GetBoolean(data, "test")
	if err != nil {
		return 0
	}
	return 1
}

func FuzzGetUnsafeString(data []byte) int {
	_, err := GetUnsafeString(data, "test")
	if err != nil {
		return 0
	}
	return 1
}
