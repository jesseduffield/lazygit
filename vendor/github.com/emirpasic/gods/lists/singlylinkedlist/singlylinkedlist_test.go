// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package singlylinkedlist

import (
	"fmt"
	"testing"

	"github.com/emirpasic/gods/utils"
)

func TestListNew(t *testing.T) {
	list1 := New()

	if actualValue := list1.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	list2 := New(1, "b")

	if actualValue := list2.Size(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	if actualValue, ok := list2.Get(0); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}

	if actualValue, ok := list2.Get(1); actualValue != "b" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "b")
	}

	if actualValue, ok := list2.Get(2); actualValue != nil || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
}

func TestListAdd(t *testing.T) {
	list := New()
	list.Add("a")
	list.Add("b", "c")
	if actualValue := list.Empty(); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := list.Size(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := list.Get(2); actualValue != "c" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "c")
	}
}

func TestListRemove(t *testing.T) {
	list := New()
	list.Add("a")
	list.Add("b", "c")
	list.Remove(2)
	if actualValue, ok := list.Get(2); actualValue != nil || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	list.Remove(1)
	list.Remove(0)
	list.Remove(0) // no effect
	if actualValue := list.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := list.Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestListGet(t *testing.T) {
	list := New()
	list.Add("a")
	list.Add("b", "c")
	if actualValue, ok := list.Get(0); actualValue != "a" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "a")
	}
	if actualValue, ok := list.Get(1); actualValue != "b" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "b")
	}
	if actualValue, ok := list.Get(2); actualValue != "c" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "c")
	}
	if actualValue, ok := list.Get(3); actualValue != nil || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	list.Remove(0)
	if actualValue, ok := list.Get(0); actualValue != "b" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "b")
	}
}

func TestListSwap(t *testing.T) {
	list := New()
	list.Add("a")
	list.Add("b", "c")
	list.Swap(0, 1)
	if actualValue, ok := list.Get(0); actualValue != "b" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "c")
	}
}

func TestListSort(t *testing.T) {
	list := New()
	list.Sort(utils.StringComparator)
	list.Add("e", "f", "g", "a", "b", "c", "d")
	list.Sort(utils.StringComparator)
	for i := 1; i < list.Size(); i++ {
		a, _ := list.Get(i - 1)
		b, _ := list.Get(i)
		if a.(string) > b.(string) {
			t.Errorf("Not sorted! %s > %s", a, b)
		}
	}
}

func TestListClear(t *testing.T) {
	list := New()
	list.Add("e", "f", "g", "a", "b", "c", "d")
	list.Clear()
	if actualValue := list.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := list.Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestListContains(t *testing.T) {
	list := New()
	list.Add("a")
	list.Add("b", "c")
	if actualValue := list.Contains("a"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := list.Contains("a", "b", "c"); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := list.Contains("a", "b", "c", "d"); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	list.Clear()
	if actualValue := list.Contains("a"); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := list.Contains("a", "b", "c"); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
}

func TestListValues(t *testing.T) {
	list := New()
	list.Add("a")
	list.Add("b", "c")
	if actualValue, expectedValue := fmt.Sprintf("%s%s%s", list.Values()...), "abc"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestListIndexOf(t *testing.T) {
	list := New()

	expectedIndex := -1
	if index := list.IndexOf("a"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	list.Add("a")
	list.Add("b", "c")

	expectedIndex = 0
	if index := list.IndexOf("a"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	expectedIndex = 1
	if index := list.IndexOf("b"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	expectedIndex = 2
	if index := list.IndexOf("c"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}
}

func TestListInsert(t *testing.T) {
	list := New()
	list.Insert(0, "b", "c")
	list.Insert(0, "a")
	list.Insert(10, "x") // ignore
	if actualValue := list.Size(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	list.Insert(3, "d") // append
	if actualValue := list.Size(); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, 4)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s%s%s%s", list.Values()...), "abcd"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestListSet(t *testing.T) {
	list := New()
	list.Set(0, "a")
	list.Set(1, "b")
	if actualValue := list.Size(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	list.Set(2, "c") // append
	if actualValue := list.Size(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	list.Set(4, "d")  // ignore
	list.Set(1, "bb") // update
	if actualValue := list.Size(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s%s%s", list.Values()...), "abbc"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestListEach(t *testing.T) {
	list := New()
	list.Add("a", "b", "c")
	list.Each(func(index int, value interface{}) {
		switch index {
		case 0:
			if actualValue, expectedValue := value, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
	})
}

func TestListMap(t *testing.T) {
	list := New()
	list.Add("a", "b", "c")
	mappedList := list.Map(func(index int, value interface{}) interface{} {
		return "mapped: " + value.(string)
	})
	if actualValue, _ := mappedList.Get(0); actualValue != "mapped: a" {
		t.Errorf("Got %v expected %v", actualValue, "mapped: a")
	}
	if actualValue, _ := mappedList.Get(1); actualValue != "mapped: b" {
		t.Errorf("Got %v expected %v", actualValue, "mapped: b")
	}
	if actualValue, _ := mappedList.Get(2); actualValue != "mapped: c" {
		t.Errorf("Got %v expected %v", actualValue, "mapped: c")
	}
	if mappedList.Size() != 3 {
		t.Errorf("Got %v expected %v", mappedList.Size(), 3)
	}
}

func TestListSelect(t *testing.T) {
	list := New()
	list.Add("a", "b", "c")
	selectedList := list.Select(func(index int, value interface{}) bool {
		return value.(string) >= "a" && value.(string) <= "b"
	})
	if actualValue, _ := selectedList.Get(0); actualValue != "a" {
		t.Errorf("Got %v expected %v", actualValue, "value: a")
	}
	if actualValue, _ := selectedList.Get(1); actualValue != "b" {
		t.Errorf("Got %v expected %v", actualValue, "value: b")
	}
	if selectedList.Size() != 2 {
		t.Errorf("Got %v expected %v", selectedList.Size(), 3)
	}
}

func TestListAny(t *testing.T) {
	list := New()
	list.Add("a", "b", "c")
	any := list.Any(func(index int, value interface{}) bool {
		return value.(string) == "c"
	})
	if any != true {
		t.Errorf("Got %v expected %v", any, true)
	}
	any = list.Any(func(index int, value interface{}) bool {
		return value.(string) == "x"
	})
	if any != false {
		t.Errorf("Got %v expected %v", any, false)
	}
}
func TestListAll(t *testing.T) {
	list := New()
	list.Add("a", "b", "c")
	all := list.All(func(index int, value interface{}) bool {
		return value.(string) >= "a" && value.(string) <= "c"
	})
	if all != true {
		t.Errorf("Got %v expected %v", all, true)
	}
	all = list.All(func(index int, value interface{}) bool {
		return value.(string) >= "a" && value.(string) <= "b"
	})
	if all != false {
		t.Errorf("Got %v expected %v", all, false)
	}
}
func TestListFind(t *testing.T) {
	list := New()
	list.Add("a", "b", "c")
	foundIndex, foundValue := list.Find(func(index int, value interface{}) bool {
		return value.(string) == "c"
	})
	if foundValue != "c" || foundIndex != 2 {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, foundIndex, "c", 2)
	}
	foundIndex, foundValue = list.Find(func(index int, value interface{}) bool {
		return value.(string) == "x"
	})
	if foundValue != nil || foundIndex != -1 {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, foundIndex, nil, nil)
	}
}
func TestListChaining(t *testing.T) {
	list := New()
	list.Add("a", "b", "c")
	chainedList := list.Select(func(index int, value interface{}) bool {
		return value.(string) > "a"
	}).Map(func(index int, value interface{}) interface{} {
		return value.(string) + value.(string)
	})
	if chainedList.Size() != 2 {
		t.Errorf("Got %v expected %v", chainedList.Size(), 2)
	}
	if actualValue, ok := chainedList.Get(0); actualValue != "bb" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "b")
	}
	if actualValue, ok := chainedList.Get(1); actualValue != "cc" || !ok {
		t.Errorf("Got %v expected %v", actualValue, "c")
	}
}

func TestListIteratorNextOnEmpty(t *testing.T) {
	list := New()
	it := list.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestListIteratorNext(t *testing.T) {
	list := New()
	list.Add("a", "b", "c")
	it := list.Iterator()
	count := 0
	for it.Next() {
		count++
		index := it.Index()
		value := it.Value()
		switch index {
		case 0:
			if actualValue, expectedValue := value, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
	}
	if actualValue, expectedValue := count, 3; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestListIteratorBegin(t *testing.T) {
	list := New()
	it := list.Iterator()
	it.Begin()
	list.Add("a", "b", "c")
	for it.Next() {
	}
	it.Begin()
	it.Next()
	if index, value := it.Index(), it.Value(); index != 0 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", index, value, 0, "a")
	}
}

func TestListIteratorFirst(t *testing.T) {
	list := New()
	it := list.Iterator()
	if actualValue, expectedValue := it.First(), false; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	list.Add("a", "b", "c")
	if actualValue, expectedValue := it.First(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if index, value := it.Index(), it.Value(); index != 0 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", index, value, 0, "a")
	}
}

func TestListSerialization(t *testing.T) {
	list := New()
	list.Add("a", "b", "c")

	var err error
	assert := func() {
		if actualValue, expectedValue := fmt.Sprintf("%s%s%s", list.Values()...), "abc"; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if actualValue, expectedValue := list.Size(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	json, err := list.ToJSON()
	assert()

	err = list.FromJSON(json)
	assert()
}

func benchmarkGet(b *testing.B, list *List, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Get(n)
		}
	}
}

func benchmarkAdd(b *testing.B, list *List, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Add(n)
		}
	}
}

func benchmarkRemove(b *testing.B, list *List, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			list.Remove(n)
		}
	}
}

func BenchmarkSinglyLinkedListGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := New()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkGet(b, list, size)
}

func BenchmarkSinglyLinkedListGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := New()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkGet(b, list, size)
}

func BenchmarkSinglyLinkedListGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := New()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkGet(b, list, size)
}

func BenchmarkSinglyLinkedListGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	list := New()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkGet(b, list, size)
}

func BenchmarkSinglyLinkedListAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := New()
	b.StartTimer()
	benchmarkAdd(b, list, size)
}

func BenchmarkSinglyLinkedListAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := New()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkAdd(b, list, size)
}

func BenchmarkSinglyLinkedListAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := New()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkAdd(b, list, size)
}

func BenchmarkSinglyLinkedListAdd100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	list := New()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkAdd(b, list, size)
}

func BenchmarkSinglyLinkedListRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	list := New()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, list, size)
}

func BenchmarkSinglyLinkedListRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	list := New()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, list, size)
}

func BenchmarkSinglyLinkedListRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	list := New()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, list, size)
}

func BenchmarkSinglyLinkedListRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	list := New()
	for n := 0; n < size; n++ {
		list.Add(n)
	}
	b.StartTimer()
	benchmarkRemove(b, list, size)
}
