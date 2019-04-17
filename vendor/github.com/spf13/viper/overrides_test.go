package viper

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

type layer int

const (
	defaultLayer layer = iota + 1
	overrideLayer
)

func TestNestedOverrides(t *testing.T) {
	assert := assert.New(t)
	var v *Viper

	// Case 0: value overridden by a value
	overrideDefault(assert, "tom", 10, "tom", 20) // "tom" is first given 10 as default value, then overridden by 20
	override(assert, "tom", 10, "tom", 20)        // "tom" is first given value 10, then overridden by 20
	overrideDefault(assert, "tom.age", 10, "tom.age", 20)
	override(assert, "tom.age", 10, "tom.age", 20)
	overrideDefault(assert, "sawyer.tom.age", 10, "sawyer.tom.age", 20)
	override(assert, "sawyer.tom.age", 10, "sawyer.tom.age", 20)

	// Case 1: key:value overridden by a value
	v = overrideDefault(assert, "tom.age", 10, "tom", "boy") // "tom.age" is first given 10 as default value, then "tom" is overridden by "boy"
	assert.Nil(v.Get("tom.age"))                             // "tom.age" should not exist anymore
	v = override(assert, "tom.age", 10, "tom", "boy")
	assert.Nil(v.Get("tom.age"))

	// Case 2: value overridden by a key:value
	overrideDefault(assert, "tom", "boy", "tom.age", 10) // "tom" is first given "boy" as default value, then "tom" is overridden by map{"age":10}
	override(assert, "tom.age", 10, "tom", "boy")

	// Case 3: key:value overridden by a key:value
	v = overrideDefault(assert, "tom.size", 4, "tom.age", 10)
	assert.Equal(4, v.Get("tom.size")) // value should still be reachable
	v = override(assert, "tom.size", 4, "tom.age", 10)
	assert.Equal(4, v.Get("tom.size"))
	deepCheckValue(assert, v, overrideLayer, []string{"tom", "size"}, 4)

	// Case 4: key:value overridden by a map
	v = overrideDefault(assert, "tom.size", 4, "tom", map[string]interface{}{"age": 10}) // "tom.size" is first given "4" as default value, then "tom" is overridden by map{"age":10}
	assert.Equal(4, v.Get("tom.size"))                                                   // "tom.size" should still be reachable
	assert.Equal(10, v.Get("tom.age"))                                                   // new value should be there
	deepCheckValue(assert, v, overrideLayer, []string{"tom", "age"}, 10)                 // new value should be there
	v = override(assert, "tom.size", 4, "tom", map[string]interface{}{"age": 10})
	assert.Nil(v.Get("tom.size"))
	assert.Equal(10, v.Get("tom.age"))
	deepCheckValue(assert, v, overrideLayer, []string{"tom", "age"}, 10)

	// Case 5: array overridden by a value
	overrideDefault(assert, "tom", []int{10, 20}, "tom", 30)
	override(assert, "tom", []int{10, 20}, "tom", 30)
	overrideDefault(assert, "tom.age", []int{10, 20}, "tom.age", 30)
	override(assert, "tom.age", []int{10, 20}, "tom.age", 30)

	// Case 6: array overridden by an array
	overrideDefault(assert, "tom", []int{10, 20}, "tom", []int{30, 40})
	override(assert, "tom", []int{10, 20}, "tom", []int{30, 40})
	overrideDefault(assert, "tom.age", []int{10, 20}, "tom.age", []int{30, 40})
	v = override(assert, "tom.age", []int{10, 20}, "tom.age", []int{30, 40})
	// explicit array merge:
	s, ok := v.Get("tom.age").([]int)
	if assert.True(ok, "tom[\"age\"] is not a slice") {
		v.Set("tom.age", append(s, []int{50, 60}...))
		assert.Equal([]int{30, 40, 50, 60}, v.Get("tom.age"))
		deepCheckValue(assert, v, overrideLayer, []string{"tom", "age"}, []int{30, 40, 50, 60})
	}
}

func overrideDefault(assert *assert.Assertions, firstPath string, firstValue interface{}, secondPath string, secondValue interface{}) *Viper {
	return overrideFromLayer(defaultLayer, assert, firstPath, firstValue, secondPath, secondValue)
}
func override(assert *assert.Assertions, firstPath string, firstValue interface{}, secondPath string, secondValue interface{}) *Viper {
	return overrideFromLayer(overrideLayer, assert, firstPath, firstValue, secondPath, secondValue)
}

// overrideFromLayer performs the sequential override and low-level checks.
//
// First assignment is made on layer l for path firstPath with value firstValue,
// the second one on the override layer (i.e., with the Set() function)
// for path secondPath with value secondValue.
//
// firstPath and secondPath can include an arbitrary number of dots to indicate
// a nested element.
//
// After each assignment, the value is checked, retrieved both by its full path
// and by its key sequence (successive maps).
func overrideFromLayer(l layer, assert *assert.Assertions, firstPath string, firstValue interface{}, secondPath string, secondValue interface{}) *Viper {
	v := New()
	firstKeys := strings.Split(firstPath, v.keyDelim)
	if assert == nil ||
		len(firstKeys) == 0 || len(firstKeys[0]) == 0 {
		return v
	}

	// Set and check first value
	switch l {
	case defaultLayer:
		v.SetDefault(firstPath, firstValue)
	case overrideLayer:
		v.Set(firstPath, firstValue)
	default:
		return v
	}
	assert.Equal(firstValue, v.Get(firstPath))
	deepCheckValue(assert, v, l, firstKeys, firstValue)

	// Override and check new value
	secondKeys := strings.Split(secondPath, v.keyDelim)
	if len(secondKeys) == 0 || len(secondKeys[0]) == 0 {
		return v
	}
	v.Set(secondPath, secondValue)
	assert.Equal(secondValue, v.Get(secondPath))
	deepCheckValue(assert, v, overrideLayer, secondKeys, secondValue)

	return v
}

// deepCheckValue checks that all given keys correspond to a valid path in the
// configuration map of the given layer, and that the final value equals the one given
func deepCheckValue(assert *assert.Assertions, v *Viper, l layer, keys []string, value interface{}) {
	if assert == nil || v == nil ||
		len(keys) == 0 || len(keys[0]) == 0 {
		return
	}

	// init
	var val interface{}
	var ms string
	switch l {
	case defaultLayer:
		val = v.defaults
		ms = "v.defaults"
	case overrideLayer:
		val = v.override
		ms = "v.override"
	}

	// loop through map
	var m map[string]interface{}
	err := false
	for _, k := range keys {
		if val == nil {
			assert.Fail(fmt.Sprintf("%s is not a map[string]interface{}", ms))
			return
		}

		// deep scan of the map to get the final value
		switch val.(type) {
		case map[interface{}]interface{}:
			m = cast.ToStringMap(val)
		case map[string]interface{}:
			m = val.(map[string]interface{})
		default:
			assert.Fail(fmt.Sprintf("%s is not a map[string]interface{}", ms))
			return
		}
		ms = ms + "[\"" + k + "\"]"
		val = m[k]
	}
	if !err {
		assert.Equal(value, val)
	}
}
