package jsonparser

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

// Errors
var (
	KeyPathNotFoundError       = errors.New("Key path not found")
	UnknownValueTypeError      = errors.New("Unknown value type")
	MalformedJsonError         = errors.New("Malformed JSON error")
	MalformedStringError       = errors.New("Value is string, but can't find closing '\"' symbol")
	MalformedArrayError        = errors.New("Value is array, but can't find closing ']' symbol")
	MalformedObjectError       = errors.New("Value looks like object, but can't find closing '}' symbol")
	MalformedValueError        = errors.New("Value looks like Number/Boolean/None, but can't find its end: ',' or '}' symbol")
	OverflowIntegerError       = errors.New("Value is number, but overflowed while parsing")
	MalformedStringEscapeError = errors.New("Encountered an invalid escape sequence in a string")
)

// How much stack space to allocate for unescaping JSON strings; if a string longer
// than this needs to be escaped, it will result in a heap allocation
const unescapeStackBufSize = 64

func tokenEnd(data []byte) int {
	for i, c := range data {
		switch c {
		case ' ', '\n', '\r', '\t', ',', '}', ']':
			return i
		}
	}

	return len(data)
}

func findTokenStart(data []byte, token byte) int {
	for i := len(data) - 1; i >= 0; i-- {
		switch data[i] {
		case token:
			return i
		case '[', '{':
			return 0
		}
	}

	return 0
}

func findKeyStart(data []byte, key string) (int, error) {
	i := 0
	ln := len(data)
	if ln > 0 && (data[0] == '{' || data[0] == '[') {
		i = 1
	}
	var stackbuf [unescapeStackBufSize]byte // stack-allocated array for allocation-free unescaping of small strings

	if ku, err := Unescape(StringToBytes(key), stackbuf[:]); err == nil {
		key = bytesToString(&ku)
	}

	for i < ln {
		switch data[i] {
		case '"':
			i++
			keyBegin := i

			strEnd, keyEscaped := stringEnd(data[i:])
			if strEnd == -1 {
				break
			}
			i += strEnd
			keyEnd := i - 1

			valueOffset := nextToken(data[i:])
			if valueOffset == -1 {
				break
			}

			i += valueOffset

			// if string is a key, and key level match
			k := data[keyBegin:keyEnd]
			// for unescape: if there are no escape sequences, this is cheap; if there are, it is a
			// bit more expensive, but causes no allocations unless len(key) > unescapeStackBufSize
			if keyEscaped {
				if ku, err := Unescape(k, stackbuf[:]); err != nil {
					break
				} else {
					k = ku
				}
			}

			if data[i] == ':' && len(key) == len(k) && bytesToString(&k) == key {
				return keyBegin - 1, nil
			}

		case '[':
			end := blockEnd(data[i:], data[i], ']')
			if end != -1 {
				i = i + end
			}
		case '{':
			end := blockEnd(data[i:], data[i], '}')
			if end != -1 {
				i = i + end
			}
		}
		i++
	}

	return -1, KeyPathNotFoundError
}

func tokenStart(data []byte) int {
	for i := len(data) - 1; i >= 0; i-- {
		switch data[i] {
		case '\n', '\r', '\t', ',', '{', '[':
			return i
		}
	}

	return 0
}

// Find position of next character which is not whitespace
func nextToken(data []byte) int {
	for i, c := range data {
		switch c {
		case ' ', '\n', '\r', '\t':
			continue
		default:
			return i
		}
	}

	return -1
}

// Find position of last character which is not whitespace
func lastToken(data []byte) int {
	for i := len(data) - 1; i >= 0; i-- {
		switch data[i] {
		case ' ', '\n', '\r', '\t':
			continue
		default:
			return i
		}
	}

	return -1
}

// Tries to find the end of string
// Support if string contains escaped quote symbols.
func stringEnd(data []byte) (int, bool) {
	escaped := false
	for i, c := range data {
		if c == '"' {
			if !escaped {
				return i + 1, false
			} else {
				j := i - 1
				for {
					if j < 0 || data[j] != '\\' {
						return i + 1, true // even number of backslashes
					}
					j--
					if j < 0 || data[j] != '\\' {
						break // odd number of backslashes
					}
					j--

				}
			}
		} else if c == '\\' {
			escaped = true
		}
	}

	return -1, escaped
}

// Find end of the data structure, array or object.
// For array openSym and closeSym will be '[' and ']', for object '{' and '}'
func blockEnd(data []byte, openSym byte, closeSym byte) int {
	level := 0
	i := 0
	ln := len(data)

	for i < ln {
		switch data[i] {
		case '"': // If inside string, skip it
			se, _ := stringEnd(data[i+1:])
			if se == -1 {
				return -1
			}
			i += se
		case openSym: // If open symbol, increase level
			level++
		case closeSym: // If close symbol, increase level
			level--

			// If we have returned to the original level, we're done
			if level == 0 {
				return i + 1
			}
		}
		i++
	}

	return -1
}

func searchKeys(data []byte, keys ...string) int {
	keyLevel := 0
	level := 0
	i := 0
	ln := len(data)
	lk := len(keys)
	lastMatched := true

	if lk == 0 {
		return 0
	}

	var stackbuf [unescapeStackBufSize]byte // stack-allocated array for allocation-free unescaping of small strings

	for i < ln {
		switch data[i] {
		case '"':
			i++
			keyBegin := i

			strEnd, keyEscaped := stringEnd(data[i:])
			if strEnd == -1 {
				return -1
			}
			i += strEnd
			keyEnd := i - 1

			valueOffset := nextToken(data[i:])
			if valueOffset == -1 {
				return -1
			}

			i += valueOffset

			// if string is a key
			if data[i] == ':' {
				if level < 1 {
					return -1
				}

				key := data[keyBegin:keyEnd]

				// for unescape: if there are no escape sequences, this is cheap; if there are, it is a
				// bit more expensive, but causes no allocations unless len(key) > unescapeStackBufSize
				var keyUnesc []byte
				if !keyEscaped {
					keyUnesc = key
				} else if ku, err := Unescape(key, stackbuf[:]); err != nil {
					return -1
				} else {
					keyUnesc = ku
				}

				if level <= len(keys) {
					if equalStr(&keyUnesc, keys[level-1]) {
						lastMatched = true

						// if key level match
						if keyLevel == level-1 {
							keyLevel++
							// If we found all keys in path
							if keyLevel == lk {
								return i + 1
							}
						}
					} else {
						lastMatched = false
					}
				} else {
					return -1
				}
			} else {
				i--
			}
		case '{':

			// in case parent key is matched then only we will increase the level otherwise can directly
			// can move to the end of this block
			if !lastMatched {
				end := blockEnd(data[i:], '{', '}')
				if end == -1 {
					return -1
				}
				i += end - 1
			} else {
				level++
			}
		case '}':
			level--
			if level == keyLevel {
				keyLevel--
			}
		case '[':
			// If we want to get array element by index
			if keyLevel == level && keys[level][0] == '[' {
				var keyLen = len(keys[level])
				if keyLen < 3 || keys[level][0] != '[' || keys[level][keyLen-1] != ']' {
					return -1
				}
				aIdx, err := strconv.Atoi(keys[level][1 : keyLen-1])
				if err != nil {
					return -1
				}
				var curIdx int
				var valueFound []byte
				var valueOffset int
				var curI = i
				ArrayEach(data[i:], func(value []byte, dataType ValueType, offset int, err error) {
					if curIdx == aIdx {
						valueFound = value
						valueOffset = offset
						if dataType == String {
							valueOffset = valueOffset - 2
							valueFound = data[curI+valueOffset : curI+valueOffset+len(value)+2]
						}
					}
					curIdx += 1
				})

				if valueFound == nil {
					return -1
				} else {
					subIndex := searchKeys(valueFound, keys[level+1:]...)
					if subIndex < 0 {
						return -1
					}
					return i + valueOffset + subIndex
				}
			} else {
				// Do not search for keys inside arrays
				if arraySkip := blockEnd(data[i:], '[', ']'); arraySkip == -1 {
					return -1
				} else {
					i += arraySkip - 1
				}
			}
		case ':': // If encountered, JSON data is malformed
			return -1
		}

		i++
	}

	return -1
}

func sameTree(p1, p2 []string) bool {
	minLen := len(p1)
	if len(p2) < minLen {
		minLen = len(p2)
	}

	for pi_1, p_1 := range p1[:minLen] {
		if p2[pi_1] != p_1 {
			return false
		}
	}

	return true
}

func EachKey(data []byte, cb func(int, []byte, ValueType, error), paths ...[]string) int {
	var x struct{}
	pathFlags := make([]bool, len(paths))
	var level, pathsMatched, i int
	ln := len(data)

	var maxPath int
	for _, p := range paths {
		if len(p) > maxPath {
			maxPath = len(p)
		}
	}

	pathsBuf := make([]string, maxPath)

	for i < ln {
		switch data[i] {
		case '"':
			i++
			keyBegin := i

			strEnd, keyEscaped := stringEnd(data[i:])
			if strEnd == -1 {
				return -1
			}
			i += strEnd

			keyEnd := i - 1

			valueOffset := nextToken(data[i:])
			if valueOffset == -1 {
				return -1
			}

			i += valueOffset

			// if string is a key, and key level match
			if data[i] == ':' {
				match := -1
				key := data[keyBegin:keyEnd]

				// for unescape: if there are no escape sequences, this is cheap; if there are, it is a
				// bit more expensive, but causes no allocations unless len(key) > unescapeStackBufSize
				var keyUnesc []byte
				if !keyEscaped {
					keyUnesc = key
				} else {
					var stackbuf [unescapeStackBufSize]byte
					if ku, err := Unescape(key, stackbuf[:]); err != nil {
						return -1
					} else {
						keyUnesc = ku
					}
				}

				if maxPath >= level {
					if level < 1 {
						cb(-1, nil, Unknown, MalformedJsonError)
						return -1
					}

					pathsBuf[level-1] = bytesToString(&keyUnesc)
					for pi, p := range paths {
						if len(p) != level || pathFlags[pi] || !equalStr(&keyUnesc, p[level-1]) || !sameTree(p, pathsBuf[:level]) {
							continue
						}

						match = pi

						pathsMatched++
						pathFlags[pi] = true

						v, dt, _, e := Get(data[i+1:])
						cb(pi, v, dt, e)

						if pathsMatched == len(paths) {
							break
						}
					}
					if pathsMatched == len(paths) {
						return i
					}
				}

				if match == -1 {
					tokenOffset := nextToken(data[i+1:])
					i += tokenOffset

					if data[i] == '{' {
						blockSkip := blockEnd(data[i:], '{', '}')
						i += blockSkip + 1
					}
				}

				if i < ln {
					switch data[i] {
					case '{', '}', '[', '"':
						i--
					}
				}
			} else {
				i--
			}
		case '{':
			level++
		case '}':
			level--
		case '[':
			var ok bool
			arrIdxFlags := make(map[int]struct{})
			pIdxFlags := make([]bool, len(paths))

			if level < 0 {
				cb(-1, nil, Unknown, MalformedJsonError)
				return -1
			}

			for pi, p := range paths {
				if len(p) < level+1 || pathFlags[pi] || p[level][0] != '[' || !sameTree(p, pathsBuf[:level]) {
					continue
				}
				if len(p[level]) >= 2 {
					aIdx, _ := strconv.Atoi(p[level][1 : len(p[level])-1])
					arrIdxFlags[aIdx] = x
					pIdxFlags[pi] = true
				}
			}

			if len(arrIdxFlags) > 0 {
				level++

				var curIdx int
				arrOff, _ := ArrayEach(data[i:], func(value []byte, dataType ValueType, offset int, err error) {
					if _, ok = arrIdxFlags[curIdx]; ok {
						for pi, p := range paths {
							if pIdxFlags[pi] {
								aIdx, _ := strconv.Atoi(p[level-1][1 : len(p[level-1])-1])

								if curIdx == aIdx {
									of := searchKeys(value, p[level:]...)

									pathsMatched++
									pathFlags[pi] = true

									if of != -1 {
										v, dt, _, e := Get(value[of:])
										cb(pi, v, dt, e)
									}
								}
							}
						}
					}

					curIdx += 1
				})

				if pathsMatched == len(paths) {
					return i
				}

				i += arrOff - 1
			} else {
				// Do not search for keys inside arrays
				if arraySkip := blockEnd(data[i:], '[', ']'); arraySkip == -1 {
					return -1
				} else {
					i += arraySkip - 1
				}
			}
		case ']':
			level--
		}

		i++
	}

	return -1
}

// Data types available in valid JSON data.
type ValueType int

const (
	NotExist = ValueType(iota)
	String
	Number
	Object
	Array
	Boolean
	Null
	Unknown
)

func (vt ValueType) String() string {
	switch vt {
	case NotExist:
		return "non-existent"
	case String:
		return "string"
	case Number:
		return "number"
	case Object:
		return "object"
	case Array:
		return "array"
	case Boolean:
		return "boolean"
	case Null:
		return "null"
	default:
		return "unknown"
	}
}

var (
	trueLiteral  = []byte("true")
	falseLiteral = []byte("false")
	nullLiteral  = []byte("null")
)

func createInsertComponent(keys []string, setValue []byte, comma, object bool) []byte {
	isIndex := string(keys[0][0]) == "["
	offset := 0
	lk := calcAllocateSpace(keys, setValue, comma, object)
	buffer := make([]byte, lk, lk)
	if comma {
		offset += WriteToBuffer(buffer[offset:], ",")
	}
	if isIndex && !comma {
		offset += WriteToBuffer(buffer[offset:], "[")
	} else {
		if object {
			offset += WriteToBuffer(buffer[offset:], "{")
		}
		if !isIndex {
			offset += WriteToBuffer(buffer[offset:], "\"")
			offset += WriteToBuffer(buffer[offset:], keys[0])
			offset += WriteToBuffer(buffer[offset:], "\":")
		}
	}

	for i := 1; i < len(keys); i++ {
		if string(keys[i][0]) == "[" {
			offset += WriteToBuffer(buffer[offset:], "[")
		} else {
			offset += WriteToBuffer(buffer[offset:], "{\"")
			offset += WriteToBuffer(buffer[offset:], keys[i])
			offset += WriteToBuffer(buffer[offset:], "\":")
		}
	}
	offset += WriteToBuffer(buffer[offset:], string(setValue))
	for i := len(keys) - 1; i > 0; i-- {
		if string(keys[i][0]) == "[" {
			offset += WriteToBuffer(buffer[offset:], "]")
		} else {
			offset += WriteToBuffer(buffer[offset:], "}")
		}
	}
	if isIndex && !comma {
		offset += WriteToBuffer(buffer[offset:], "]")
	}
	if object && !isIndex {
		offset += WriteToBuffer(buffer[offset:], "}")
	}
	return buffer
}

func calcAllocateSpace(keys []string, setValue []byte, comma, object bool) int {
	isIndex := string(keys[0][0]) == "["
	lk := 0
	if comma {
		// ,
		lk += 1
	}
	if isIndex && !comma {
		// []
		lk += 2
	} else {
		if object {
			// {
			lk += 1
		}
		if !isIndex {
			// "keys[0]"
			lk += len(keys[0]) + 3
		}
	}


	lk += len(setValue)
	for i := 1; i < len(keys); i++ {
		if string(keys[i][0]) == "[" {
			// []
			lk += 2
		} else {
			// {"keys[i]":setValue}
			lk += len(keys[i]) + 5
		}
	}

	if object && !isIndex {
		// }
		lk += 1
	}

	return lk
}

func WriteToBuffer(buffer []byte, str string) int {
	copy(buffer, str)
	return len(str)
}

/*

Del - Receives existing data structure, path to delete.

Returns:
`data` - return modified data

*/
func Delete(data []byte, keys ...string) []byte {
	lk := len(keys)
	if lk == 0 {
		return data[:0]
	}

	array := false
	if len(keys[lk-1]) > 0 && string(keys[lk-1][0]) == "[" {
		array = true
	}

	var startOffset, keyOffset int
	endOffset := len(data)
	var err error
	if !array {
		if len(keys) > 1 {
			_, _, startOffset, endOffset, err = internalGet(data, keys[:lk-1]...)
			if err == KeyPathNotFoundError {
				// problem parsing the data
				return data
			}
		}

		keyOffset, err = findKeyStart(data[startOffset:endOffset], keys[lk-1])
		if err == KeyPathNotFoundError {
			// problem parsing the data
			return data
		}
		keyOffset += startOffset
		_, _, _, subEndOffset, _ := internalGet(data[startOffset:endOffset], keys[lk-1])
		endOffset = startOffset + subEndOffset
		tokEnd := tokenEnd(data[endOffset:])
		tokStart := findTokenStart(data[:keyOffset], ","[0])

		if data[endOffset+tokEnd] == ","[0] {
			endOffset += tokEnd + 1
		} else if data[endOffset+tokEnd] == " "[0] && len(data) > endOffset+tokEnd+1 && data[endOffset+tokEnd+1] == ","[0] {
			endOffset += tokEnd + 2
		} else if data[endOffset+tokEnd] == "}"[0] && data[tokStart] == ","[0] {
			keyOffset = tokStart
		}
	} else {
		_, _, keyOffset, endOffset, err = internalGet(data, keys...)
		if err == KeyPathNotFoundError {
			// problem parsing the data
			return data
		}

		tokEnd := tokenEnd(data[endOffset:])
		tokStart := findTokenStart(data[:keyOffset], ","[0])

		if data[endOffset+tokEnd] == ","[0] {
			endOffset += tokEnd + 1
		} else if data[endOffset+tokEnd] == "]"[0] && data[tokStart] == ","[0] {
			keyOffset = tokStart
		}
	}

	// We need to remove remaining trailing comma if we delete las element in the object
	prevTok := lastToken(data[:keyOffset])
	remainedValue := data[endOffset:]

	var newOffset int
	if nextToken(remainedValue) > -1 && remainedValue[nextToken(remainedValue)] == '}' && data[prevTok] == ',' {
		newOffset = prevTok
	} else {
		newOffset = prevTok + 1
	}

	// We have to make a copy here if we don't want to mangle the original data, because byte slices are
	// accessed by reference and not by value
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	data = append(dataCopy[:newOffset], dataCopy[endOffset:]...)

	return data
}

/*

Set - Receives existing data structure, path to set, and data to set at that key.

Returns:
`value` - modified byte array
`err` - On any parsing error

*/
func Set(data []byte, setValue []byte, keys ...string) (value []byte, err error) {
	// ensure keys are set
	if len(keys) == 0 {
		return nil, KeyPathNotFoundError
	}

	_, _, startOffset, endOffset, err := internalGet(data, keys...)
	if err != nil {
		if err != KeyPathNotFoundError {
			// problem parsing the data
			return nil, err
		}
		// full path doesnt exist
		// does any subpath exist?
		var depth int
		for i := range keys {
			_, _, start, end, sErr := internalGet(data, keys[:i+1]...)
			if sErr != nil {
				break
			} else {
				endOffset = end
				startOffset = start
				depth++
			}
		}
		comma := true
		object := false
		if endOffset == -1 {
			firstToken := nextToken(data)
			// We can't set a top-level key if data isn't an object
			if firstToken < 0 || data[firstToken] != '{' {
				return nil, KeyPathNotFoundError
			}
			// Don't need a comma if the input is an empty object
			secondToken := firstToken + 1 + nextToken(data[firstToken+1:])
			if data[secondToken] == '}' {
				comma = false
			}
			// Set the top level key at the end (accounting for any trailing whitespace)
			// This assumes last token is valid like '}', could check and return error
			endOffset = lastToken(data)
		}
		depthOffset := endOffset
		if depth != 0 {
			// if subpath is a non-empty object, add to it
			// or if subpath is a non-empty array, add to it
			if (data[startOffset] == '{' && data[startOffset+1+nextToken(data[startOffset+1:])] != '}') ||
				(data[startOffset] == '[' && data[startOffset+1+nextToken(data[startOffset+1:])] == '{') && keys[depth:][0][0] == 91 {
				depthOffset--
				startOffset = depthOffset
				// otherwise, over-write it with a new object
			} else {
				comma = false
				object = true
			}
		} else {
			startOffset = depthOffset
		}
		value = append(data[:startOffset], append(createInsertComponent(keys[depth:], setValue, comma, object), data[depthOffset:]...)...)
	} else {
		// path currently exists
		startComponent := data[:startOffset]
		endComponent := data[endOffset:]

		value = make([]byte, len(startComponent)+len(endComponent)+len(setValue))
		newEndOffset := startOffset + len(setValue)
		copy(value[0:startOffset], startComponent)
		copy(value[startOffset:newEndOffset], setValue)
		copy(value[newEndOffset:], endComponent)
	}
	return value, nil
}

func getType(data []byte, offset int) ([]byte, ValueType, int, error) {
	var dataType ValueType
	endOffset := offset

	// if string value
	if data[offset] == '"' {
		dataType = String
		if idx, _ := stringEnd(data[offset+1:]); idx != -1 {
			endOffset += idx + 1
		} else {
			return nil, dataType, offset, MalformedStringError
		}
	} else if data[offset] == '[' { // if array value
		dataType = Array
		// break label, for stopping nested loops
		endOffset = blockEnd(data[offset:], '[', ']')

		if endOffset == -1 {
			return nil, dataType, offset, MalformedArrayError
		}

		endOffset += offset
	} else if data[offset] == '{' { // if object value
		dataType = Object
		// break label, for stopping nested loops
		endOffset = blockEnd(data[offset:], '{', '}')

		if endOffset == -1 {
			return nil, dataType, offset, MalformedObjectError
		}

		endOffset += offset
	} else {
		// Number, Boolean or None
		end := tokenEnd(data[endOffset:])

		if end == -1 {
			return nil, dataType, offset, MalformedValueError
		}

		value := data[offset : endOffset+end]

		switch data[offset] {
		case 't', 'f': // true or false
			if bytes.Equal(value, trueLiteral) || bytes.Equal(value, falseLiteral) {
				dataType = Boolean
			} else {
				return nil, Unknown, offset, UnknownValueTypeError
			}
		case 'u', 'n': // undefined or null
			if bytes.Equal(value, nullLiteral) {
				dataType = Null
			} else {
				return nil, Unknown, offset, UnknownValueTypeError
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			dataType = Number
		default:
			return nil, Unknown, offset, UnknownValueTypeError
		}

		endOffset += end
	}
	return data[offset:endOffset], dataType, endOffset, nil
}

/*
Get - Receives data structure, and key path to extract value from.

Returns:
`value` - Pointer to original data structure containing key value, or just empty slice if nothing found or error
`dataType` -    Can be: `NotExist`, `String`, `Number`, `Object`, `Array`, `Boolean` or `Null`
`offset` - Offset from provided data structure where key value ends. Used mostly internally, for example for `ArrayEach` helper.
`err` - If key not found or any other parsing issue it should return error. If key not found it also sets `dataType` to `NotExist`

Accept multiple keys to specify path to JSON value (in case of quering nested structures).
If no keys provided it will try to extract closest JSON value (simple ones or object/array), useful for reading streams or arrays, see `ArrayEach` implementation.
*/
func Get(data []byte, keys ...string) (value []byte, dataType ValueType, offset int, err error) {
	a, b, _, d, e := internalGet(data, keys...)
	return a, b, d, e
}

func internalGet(data []byte, keys ...string) (value []byte, dataType ValueType, offset, endOffset int, err error) {
	if len(keys) > 0 {
		if offset = searchKeys(data, keys...); offset == -1 {
			return nil, NotExist, -1, -1, KeyPathNotFoundError
		}
	}

	// Go to closest value
	nO := nextToken(data[offset:])
	if nO == -1 {
		return nil, NotExist, offset, -1, MalformedJsonError
	}

	offset += nO
	value, dataType, endOffset, err = getType(data, offset)
	if err != nil {
		return value, dataType, offset, endOffset, err
	}

	// Strip quotes from string values
	if dataType == String {
		value = value[1 : len(value)-1]
	}

	return value[:len(value):len(value)], dataType, offset, endOffset, nil
}

// ArrayEach is used when iterating arrays, accepts a callback function with the same return arguments as `Get`.
func ArrayEach(data []byte, cb func(value []byte, dataType ValueType, offset int, err error), keys ...string) (offset int, err error) {
	if len(data) == 0 {
		return -1, MalformedObjectError
	}

	nT := nextToken(data)
	if nT == -1 {
		return -1, MalformedJsonError
	}

	offset = nT + 1

	if len(keys) > 0 {
		if offset = searchKeys(data, keys...); offset == -1 {
			return offset, KeyPathNotFoundError
		}

		// Go to closest value
		nO := nextToken(data[offset:])
		if nO == -1 {
			return offset, MalformedJsonError
		}

		offset += nO

		if data[offset] != '[' {
			return offset, MalformedArrayError
		}

		offset++
	}

	nO := nextToken(data[offset:])
	if nO == -1 {
		return offset, MalformedJsonError
	}

	offset += nO

	if data[offset] == ']' {
		return offset, nil
	}

	for true {
		v, t, o, e := Get(data[offset:])

		if e != nil {
			return offset, e
		}

		if o == 0 {
			break
		}

		if t != NotExist {
			cb(v, t, offset+o-len(v), e)
		}

		if e != nil {
			break
		}

		offset += o

		skipToToken := nextToken(data[offset:])
		if skipToToken == -1 {
			return offset, MalformedArrayError
		}
		offset += skipToToken

		if data[offset] == ']' {
			break
		}

		if data[offset] != ',' {
			return offset, MalformedArrayError
		}

		offset++
	}

	return offset, nil
}

// ObjectEach iterates over the key-value pairs of a JSON object, invoking a given callback for each such entry
func ObjectEach(data []byte, callback func(key []byte, value []byte, dataType ValueType, offset int) error, keys ...string) (err error) {
	offset := 0

	// Descend to the desired key, if requested
	if len(keys) > 0 {
		if off := searchKeys(data, keys...); off == -1 {
			return KeyPathNotFoundError
		} else {
			offset = off
		}
	}

	// Validate and skip past opening brace
	if off := nextToken(data[offset:]); off == -1 {
		return MalformedObjectError
	} else if offset += off; data[offset] != '{' {
		return MalformedObjectError
	} else {
		offset++
	}

	// Skip to the first token inside the object, or stop if we find the ending brace
	if off := nextToken(data[offset:]); off == -1 {
		return MalformedJsonError
	} else if offset += off; data[offset] == '}' {
		return nil
	}

	// Loop pre-condition: data[offset] points to what should be either the next entry's key, or the closing brace (if it's anything else, the JSON is malformed)
	for offset < len(data) {
		// Step 1: find the next key
		var key []byte

		// Check what the the next token is: start of string, end of object, or something else (error)
		switch data[offset] {
		case '"':
			offset++ // accept as string and skip opening quote
		case '}':
			return nil // we found the end of the object; stop and return success
		default:
			return MalformedObjectError
		}

		// Find the end of the key string
		var keyEscaped bool
		if off, esc := stringEnd(data[offset:]); off == -1 {
			return MalformedJsonError
		} else {
			key, keyEscaped = data[offset:offset+off-1], esc
			offset += off
		}

		// Unescape the string if needed
		if keyEscaped {
			var stackbuf [unescapeStackBufSize]byte // stack-allocated array for allocation-free unescaping of small strings
			if keyUnescaped, err := Unescape(key, stackbuf[:]); err != nil {
				return MalformedStringEscapeError
			} else {
				key = keyUnescaped
			}
		}

		// Step 2: skip the colon
		if off := nextToken(data[offset:]); off == -1 {
			return MalformedJsonError
		} else if offset += off; data[offset] != ':' {
			return MalformedJsonError
		} else {
			offset++
		}

		// Step 3: find the associated value, then invoke the callback
		if value, valueType, off, err := Get(data[offset:]); err != nil {
			return err
		} else if err := callback(key, value, valueType, offset+off); err != nil { // Invoke the callback here!
			return err
		} else {
			offset += off
		}

		// Step 4: skip over the next comma to the following token, or stop if we hit the ending brace
		if off := nextToken(data[offset:]); off == -1 {
			return MalformedArrayError
		} else {
			offset += off
			switch data[offset] {
			case '}':
				return nil // Stop if we hit the close brace
			case ',':
				offset++ // Ignore the comma
			default:
				return MalformedObjectError
			}
		}

		// Skip to the next token after the comma
		if off := nextToken(data[offset:]); off == -1 {
			return MalformedArrayError
		} else {
			offset += off
		}
	}

	return MalformedObjectError // we shouldn't get here; it's expected that we will return via finding the ending brace
}

// GetUnsafeString returns the value retrieved by `Get`, use creates string without memory allocation by mapping string to slice memory. It does not handle escape symbols.
func GetUnsafeString(data []byte, keys ...string) (val string, err error) {
	v, _, _, e := Get(data, keys...)

	if e != nil {
		return "", e
	}

	return bytesToString(&v), nil
}

// GetString returns the value retrieved by `Get`, cast to a string if possible, trying to properly handle escape and utf8 symbols
// If key data type do not match, it will return an error.
func GetString(data []byte, keys ...string) (val string, err error) {
	v, t, _, e := Get(data, keys...)

	if e != nil {
		return "", e
	}

	if t != String {
		return "", fmt.Errorf("Value is not a string: %s", string(v))
	}

	// If no escapes return raw content
	if bytes.IndexByte(v, '\\') == -1 {
		return string(v), nil
	}

	return ParseString(v)
}

// GetFloat returns the value retrieved by `Get`, cast to a float64 if possible.
// The offset is the same as in `Get`.
// If key data type do not match, it will return an error.
func GetFloat(data []byte, keys ...string) (val float64, err error) {
	v, t, _, e := Get(data, keys...)

	if e != nil {
		return 0, e
	}

	if t != Number {
		return 0, fmt.Errorf("Value is not a number: %s", string(v))
	}

	return ParseFloat(v)
}

// GetInt returns the value retrieved by `Get`, cast to a int64 if possible.
// If key data type do not match, it will return an error.
func GetInt(data []byte, keys ...string) (val int64, err error) {
	v, t, _, e := Get(data, keys...)

	if e != nil {
		return 0, e
	}

	if t != Number {
		return 0, fmt.Errorf("Value is not a number: %s", string(v))
	}

	return ParseInt(v)
}

// GetBoolean returns the value retrieved by `Get`, cast to a bool if possible.
// The offset is the same as in `Get`.
// If key data type do not match, it will return error.
func GetBoolean(data []byte, keys ...string) (val bool, err error) {
	v, t, _, e := Get(data, keys...)

	if e != nil {
		return false, e
	}

	if t != Boolean {
		return false, fmt.Errorf("Value is not a boolean: %s", string(v))
	}

	return ParseBoolean(v)
}

// ParseBoolean parses a Boolean ValueType into a Go bool (not particularly useful, but here for completeness)
func ParseBoolean(b []byte) (bool, error) {
	switch {
	case bytes.Equal(b, trueLiteral):
		return true, nil
	case bytes.Equal(b, falseLiteral):
		return false, nil
	default:
		return false, MalformedValueError
	}
}

// ParseString parses a String ValueType into a Go string (the main parsing work is unescaping the JSON string)
func ParseString(b []byte) (string, error) {
	var stackbuf [unescapeStackBufSize]byte // stack-allocated array for allocation-free unescaping of small strings
	if bU, err := Unescape(b, stackbuf[:]); err != nil {
		return "", MalformedValueError
	} else {
		return string(bU), nil
	}
}

// ParseNumber parses a Number ValueType into a Go float64
func ParseFloat(b []byte) (float64, error) {
	if v, err := parseFloat(&b); err != nil {
		return 0, MalformedValueError
	} else {
		return v, nil
	}
}

// ParseInt parses a Number ValueType into a Go int64
func ParseInt(b []byte) (int64, error) {
	if v, ok, overflow := parseInt(b); !ok {
		if overflow {
			return 0, OverflowIntegerError
		}
		return 0, MalformedValueError
	} else {
		return v, nil
	}
}
