package lo

// Keys creates a slice of the map keys.
// Play: https://go.dev/play/p/Uu11fHASqrU
func Keys[K comparable, V any](in ...map[K]V) []K {
	size := 0
	for i := range in {
		size += len(in[i])
	}
	result := make([]K, 0, size)

	for i := range in {
		for k := range in[i] {
			result = append(result, k)
		}
	}

	return result
}

// UniqKeys creates a slice of unique keys in the map.
// Play: https://go.dev/play/p/TPKAb6ILdHk
func UniqKeys[K comparable, V any](in ...map[K]V) []K {
	size := 0
	for i := range in {
		size += len(in[i])
	}

	seen := make(map[K]struct{}, size)
	result := make([]K, 0)

	for i := range in {
		for k := range in[i] {
			if _, exists := seen[k]; exists {
				continue
			}
			seen[k] = struct{}{}
			result = append(result, k)
		}
	}

	return result
}

// HasKey returns whether the given key exists.
// Play: https://go.dev/play/p/aVwubIvECqS
func HasKey[K comparable, V any](in map[K]V, key K) bool {
	_, ok := in[key]
	return ok
}

// Values creates a slice of the map values.
// Play: https://go.dev/play/p/nnRTQkzQfF6
func Values[K comparable, V any](in ...map[K]V) []V {
	size := 0
	for i := range in {
		size += len(in[i])
	}
	result := make([]V, 0, size)

	for i := range in {
		for _, v := range in[i] {
			result = append(result, v)
		}
	}

	return result
}

// UniqValues creates a slice of unique values in the map.
// Play: https://go.dev/play/p/nf6bXMh7rM3
func UniqValues[K, V comparable](in ...map[K]V) []V {
	size := 0
	for i := range in {
		size += len(in[i])
	}

	seen := make(map[V]struct{}, size)
	result := make([]V, 0)

	for i := range in {
		for _, v := range in[i] {
			if _, exists := seen[v]; exists {
				continue
			}
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

// ValueOr returns the value of the given key or the fallback value if the key is not present.
// Play: https://go.dev/play/p/bAq9mHErB4V
func ValueOr[K comparable, V any](in map[K]V, key K, fallback V) V {
	if v, ok := in[key]; ok {
		return v
	}
	return fallback
}

// PickBy returns same map type filtered by given predicate.
// Play: https://go.dev/play/p/kdg8GR_QMmf
func PickBy[K comparable, V any, Map ~map[K]V](in Map, predicate func(key K, value V) bool) Map {
	r := Map{}
	for k, v := range in {
		if predicate(k, v) {
			r[k] = v
		}
	}
	return r
}

// PickByErr returns same map type filtered by given predicate.
// It returns the first error returned by the predicate.
func PickByErr[K comparable, V any, Map ~map[K]V](in Map, predicate func(key K, value V) (bool, error)) (Map, error) {
	r := Map{}
	for k, v := range in {
		ok, err := predicate(k, v)
		if err != nil {
			return nil, err
		}
		if ok {
			r[k] = v
		}
	}
	return r, nil
}

// PickByKeys returns same map type filtered by given keys.
// Play: https://go.dev/play/p/R1imbuci9qU
func PickByKeys[K comparable, V any, Map ~map[K]V](in Map, keys []K) Map {
	r := Map{}
	for i := range keys {
		if v, ok := in[keys[i]]; ok {
			r[keys[i]] = v
		}
	}
	return r
}

// PickByValues returns same map type filtered by given values.
// Play: https://go.dev/play/p/1zdzSvbfsJc
func PickByValues[K, V comparable, Map ~map[K]V](in Map, values []V) Map {
	r := Map{}

	seen := Keyify(values)
	for k, v := range in {
		if _, ok := seen[v]; ok {
			r[k] = v
		}
	}
	return r
}

// OmitBy returns same map type filtered by given predicate.
// Play: https://go.dev/play/p/EtBsR43bdsd
func OmitBy[K comparable, V any, Map ~map[K]V](in Map, predicate func(key K, value V) bool) Map {
	r := Map{}
	for k, v := range in {
		if !predicate(k, v) {
			r[k] = v
		}
	}
	return r
}

// OmitByErr returns same map type filtered by given predicate.
// It returns the first error returned by the predicate.
func OmitByErr[K comparable, V any, Map ~map[K]V](in Map, predicate func(key K, value V) (bool, error)) (Map, error) {
	r := Map{}
	for k, v := range in {
		ok, err := predicate(k, v)
		if err != nil {
			return nil, err
		}
		if !ok {
			r[k] = v
		}
	}
	return r, nil
}

// OmitByKeys returns same map type filtered by given keys.
// Play: https://go.dev/play/p/t1QjCrs-ysk
func OmitByKeys[K comparable, V any, Map ~map[K]V](in Map, keys []K) Map {
	r := Map{}
	for k, v := range in {
		r[k] = v
	}
	for i := range keys {
		delete(r, keys[i])
	}
	return r
}

// OmitByValues returns same map type filtered by given values.
// Play: https://go.dev/play/p/9UYZi-hrs8j
func OmitByValues[K, V comparable, Map ~map[K]V](in Map, values []V) Map {
	r := Map{}

	seen := Keyify(values)
	for k, v := range in {
		if _, ok := seen[v]; !ok {
			r[k] = v
		}
	}

	return r
}

// Entries transforms a map into a slice of key/value pairs.
// Play: https://go.dev/play/p/_t4Xe34-Nl5
func Entries[K comparable, V any](in map[K]V) []Entry[K, V] {
	entries := make([]Entry[K, V], 0, len(in))

	for k, v := range in {
		entries = append(entries, Entry[K, V]{
			Key:   k,
			Value: v,
		})
	}

	return entries
}

// ToPairs transforms a map into a slice of key/value pairs.
// Alias of Entries().
// Play: https://go.dev/play/p/3Dhgx46gawJ
func ToPairs[K comparable, V any](in map[K]V) []Entry[K, V] {
	return Entries(in)
}

// FromEntries transforms a slice of key/value pairs into a map.
// Play: https://go.dev/play/p/oIr5KHFGCEN
func FromEntries[K comparable, V any](entries []Entry[K, V]) map[K]V {
	out := make(map[K]V, len(entries))

	for i := range entries {
		out[entries[i].Key] = entries[i].Value
	}

	return out
}

// FromPairs transforms a slice of key/value pairs into a map.
// Alias of FromEntries().
// Play: https://go.dev/play/p/oIr5KHFGCEN
func FromPairs[K comparable, V any](entries []Entry[K, V]) map[K]V {
	return FromEntries(entries)
}

// Invert creates a map composed of the inverted keys and values. If map
// contains duplicate values, subsequent values overwrite property assignments
// of previous values.
// Play: https://go.dev/play/p/rFQ4rak6iA1
func Invert[K, V comparable](in map[K]V) map[V]K {
	out := make(map[V]K, len(in))

	for k, v := range in {
		out[v] = k
	}

	return out
}

// Assign merges multiple maps from left to right.
// Play: https://go.dev/play/p/VhwfJOyxf5o
func Assign[K comparable, V any, Map ~map[K]V](maps ...Map) Map {
	count := 0
	for i := range maps {
		count += len(maps[i])
	}

	out := make(Map, count)
	for i := range maps {
		for k, v := range maps[i] {
			out[k] = v
		}
	}

	return out
}

// ChunkEntries splits a map into a slice of elements in groups of length equal to its size. If the map cannot be split evenly,
// the final chunk will contain the remaining elements.
// Play: https://go.dev/play/p/X_YQL6mmoD-
func ChunkEntries[K comparable, V any](m map[K]V, size int) []map[K]V {
	if size <= 0 {
		panic("lo.ChunkEntries: size must be greater than 0")
	}

	count := len(m)
	if count == 0 {
		return []map[K]V{}
	}

	result := make([]map[K]V, 0, ((count-1)/size)+1)

	for k, v := range m {
		if len(result) == 0 || len(result[len(result)-1]) == size {
			result = append(result, make(map[K]V, size))
		}

		result[len(result)-1][k] = v
	}

	return result
}

// MapKeys manipulates map keys and transforms it to a map of another type.
// Play: https://go.dev/play/p/9_4WPIqOetJ
func MapKeys[K comparable, V any, R comparable](in map[K]V, iteratee func(value V, key K) R) map[R]V {
	result := make(map[R]V, len(in))

	for k, v := range in {
		result[iteratee(v, k)] = v
	}

	return result
}

// MapKeysErr manipulates map keys and transforms it to a map of another type.
// It returns the first error returned by the iteratee.
func MapKeysErr[K comparable, V any, R comparable](in map[K]V, iteratee func(value V, key K) (R, error)) (map[R]V, error) {
	result := make(map[R]V, len(in))

	for k, v := range in {
		r, err := iteratee(v, k)
		if err != nil {
			return nil, err
		}
		result[r] = v
	}

	return result, nil
}

// MapValues manipulates map values and transforms it to a map of another type.
// Play: https://go.dev/play/p/T_8xAfvcf0W
func MapValues[K comparable, V, R any](in map[K]V, iteratee func(value V, key K) R) map[K]R {
	result := make(map[K]R, len(in))

	for k, v := range in {
		result[k] = iteratee(v, k)
	}

	return result
}

// MapValuesErr manipulates map values and transforms it to a map of another type.
// It returns the first error returned by the iteratee.
func MapValuesErr[K comparable, V, R any](in map[K]V, iteratee func(value V, key K) (R, error)) (map[K]R, error) {
	result := make(map[K]R, len(in))

	for k, v := range in {
		r, err := iteratee(v, k)
		if err != nil {
			return nil, err
		}
		result[k] = r
	}

	return result, nil
}

// MapEntries manipulates map entries and transforms it to a map of another type.
// Play: https://go.dev/play/p/VuvNQzxKimT
func MapEntries[K1 comparable, V1 any, K2 comparable, V2 any](in map[K1]V1, iteratee func(key K1, value V1) (K2, V2)) map[K2]V2 {
	result := make(map[K2]V2, len(in))

	for k1 := range in {
		k2, v2 := iteratee(k1, in[k1])
		result[k2] = v2
	}

	return result
}

// MapEntriesErr manipulates map entries and transforms it to a map of another type.
// It returns the first error returned by the iteratee.
func MapEntriesErr[K1 comparable, V1 any, K2 comparable, V2 any](in map[K1]V1, iteratee func(key K1, value V1) (K2, V2, error)) (map[K2]V2, error) {
	result := make(map[K2]V2, len(in))

	for k1 := range in {
		k2, v2, err := iteratee(k1, in[k1])
		if err != nil {
			return nil, err
		}
		result[k2] = v2
	}

	return result, nil
}

// MapToSlice transforms a map into a slice based on specified iteratee.
// Play: https://go.dev/play/p/ZuiCZpDt6LD
func MapToSlice[K comparable, V, R any](in map[K]V, iteratee func(key K, value V) R) []R {
	result := make([]R, 0, len(in))

	for k, v := range in {
		result = append(result, iteratee(k, v))
	}

	return result
}

// MapToSliceErr transforms a map into a slice based on specified iteratee.
// It returns the first error returned by the iteratee.
func MapToSliceErr[K comparable, V, R any](in map[K]V, iteratee func(key K, value V) (R, error)) ([]R, error) {
	result := make([]R, 0, len(in))

	for k, v := range in {
		r, err := iteratee(k, v)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
}

// FilterMapToSlice transforms a map into a slice based on specified iteratee.
// The iteratee returns a value and a boolean. If the boolean is true, the value is added to the result slice.
// If the boolean is false, the value is not added to the result slice.
// The order of the keys in the input map is not specified and the order of the keys in the output slice is not guaranteed.
// Play: https://go.dev/play/p/jgsD_Kil9pV
func FilterMapToSlice[K comparable, V, R any](in map[K]V, iteratee func(key K, value V) (R, bool)) []R {
	result := make([]R, 0, len(in))

	for k, v := range in {
		if v, ok := iteratee(k, v); ok {
			result = append(result, v)
		}
	}

	return result
}

// FilterMapToSliceErr transforms a map into a slice based on specified iteratee.
// The iteratee returns a value, a boolean, and an error. If the boolean is true, the value is added to the result slice.
// If the boolean is false, the value is not added to the result slice.
// If an error is returned, iteration stops immediately and returns the error.
// The order of the keys in the input map is not specified and the order of the keys in the output slice is not guaranteed.
func FilterMapToSliceErr[K comparable, V, R any](in map[K]V, iteratee func(key K, value V) (R, bool, error)) ([]R, error) {
	result := make([]R, 0, len(in))

	for k, v := range in {
		r, ok, err := iteratee(k, v)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, r)
		}
	}

	return result, nil
}

// FilterKeys transforms a map into a slice based on predicate returns true for specific elements.
// It is a mix of lo.Filter() and lo.Keys().
// Play: https://go.dev/play/p/OFlKXlPrBAe
func FilterKeys[K comparable, V any](in map[K]V, predicate func(key K, value V) bool) []K {
	result := make([]K, 0)

	for k, v := range in {
		if predicate(k, v) {
			result = append(result, k)
		}
	}

	return result
}

// FilterValues transforms a map into a slice based on predicate returns true for specific elements.
// It is a mix of lo.Filter() and lo.Values().
// Play: https://go.dev/play/p/YVD5r_h-LX-
func FilterValues[K comparable, V any](in map[K]V, predicate func(key K, value V) bool) []V {
	result := make([]V, 0)

	for k, v := range in {
		if predicate(k, v) {
			result = append(result, v)
		}
	}

	return result
}

// FilterKeysErr transforms a map into a slice of keys based on predicate that can return an error.
// It is a mix of lo.Filter() and lo.Keys() with error handling.
// If the predicate returns true, the key is added to the result slice.
// If the predicate returns an error, iteration stops immediately and returns the error.
// The order of the keys in the input map is not specified.
func FilterKeysErr[K comparable, V any](in map[K]V, predicate func(key K, value V) (bool, error)) ([]K, error) {
	result := make([]K, 0)

	for k, v := range in {
		ok, err := predicate(k, v)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, k)
		}
	}

	return result, nil
}

// FilterValuesErr transforms a map into a slice of values based on predicate that can return an error.
// It is a mix of lo.Filter() and lo.Values() with error handling.
// If the predicate returns true, the value is added to the result slice.
// If the predicate returns an error, iteration stops immediately and returns the error.
// The order of the keys in the input map is not specified.
func FilterValuesErr[K comparable, V any](in map[K]V, predicate func(key K, value V) (bool, error)) ([]V, error) {
	result := make([]V, 0)

	for k, v := range in {
		ok, err := predicate(k, v)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, v)
		}
	}

	return result, nil
}
