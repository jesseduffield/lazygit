// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
)

func printMap(txt string, m *treemap.Map) {
	fmt.Print(txt, " { ")
	m.Each(func(key interface{}, value interface{}) {
		fmt.Print(key, ":", value, " ")
	})
	fmt.Println("}")
}

// EunumerableWithKeyExample to demonstrate basic usage of EunumerableWithKey
func main() {
	m := treemap.NewWithStringComparator()
	m.Put("g", 7)
	m.Put("f", 6)
	m.Put("e", 5)
	m.Put("d", 4)
	m.Put("c", 3)
	m.Put("b", 2)
	m.Put("a", 1)
	printMap("Initial", m) // { a:1 b:2 c:3 d:4 e:5 f:6 g:7 }

	even := m.Select(func(key interface{}, value interface{}) bool {
		return value.(int)%2 == 0
	})
	printMap("Elements with even values", even) // { b:2 d:4 f:6 }

	foundKey, foundValue := m.Find(func(key interface{}, value interface{}) bool {
		return value.(int)%2 == 0 && value.(int)%3 == 0
	})
	if foundKey != nil {
		fmt.Println("Element with value divisible by 2 and 3 found is", foundValue, "with key", foundKey) // value: 6, index: 4
	}

	square := m.Map(func(key interface{}, value interface{}) (interface{}, interface{}) {
		return key.(string) + key.(string), value.(int) * value.(int)
	})
	printMap("Elements' values squared and letters duplicated", square) // { aa:1 bb:4 cc:9 dd:16 ee:25 ff:36 gg:49 }

	bigger := m.Any(func(key interface{}, value interface{}) bool {
		return value.(int) > 5
	})
	fmt.Println("Map contains element whose value is bigger than 5 is", bigger) // true

	positive := m.All(func(key interface{}, value interface{}) bool {
		return value.(int) > 0
	})
	fmt.Println("All map's elements have positive values is", positive) // true

	evenNumbersSquared := m.Select(func(key interface{}, value interface{}) bool {
		return value.(int)%2 == 0
	}).Map(func(key interface{}, value interface{}) (interface{}, interface{}) {
		return key, value.(int) * value.(int)
	})
	printMap("Chaining", evenNumbersSquared) // { b:4 d:16 f:36 }
}
