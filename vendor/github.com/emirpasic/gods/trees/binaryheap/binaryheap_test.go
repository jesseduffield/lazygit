// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binaryheap

import (
	"math/rand"
	"testing"
)

func TestBinaryHeapPush(t *testing.T) {
	heap := NewWithIntComparator()

	if actualValue := heap.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	heap.Push(3) // [3]
	heap.Push(2) // [2,3]
	heap.Push(1) // [1,3,2](2 swapped with 1, hence last)

	if actualValue := heap.Values(); actualValue[0].(int) != 1 || actualValue[1].(int) != 3 || actualValue[2].(int) != 2 {
		t.Errorf("Got %v expected %v", actualValue, "[1,2,3]")
	}
	if actualValue := heap.Empty(); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := heap.Size(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := heap.Peek(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestBinaryHeapPushBulk(t *testing.T) {
	heap := NewWithIntComparator()

	heap.Push(15, 20, 3, 1, 2)

	if actualValue := heap.Values(); actualValue[0].(int) != 1 || actualValue[1].(int) != 2 || actualValue[2].(int) != 3 {
		t.Errorf("Got %v expected %v", actualValue, "[1,2,3]")
	}
	if actualValue, ok := heap.Pop(); actualValue != 1 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestBinaryHeapPop(t *testing.T) {
	heap := NewWithIntComparator()

	if actualValue := heap.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}

	heap.Push(3) // [3]
	heap.Push(2) // [2,3]
	heap.Push(1) // [1,3,2](2 swapped with 1, hence last)
	heap.Pop()   // [3,2]

	if actualValue, ok := heap.Peek(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := heap.Pop(); actualValue != 2 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	if actualValue, ok := heap.Pop(); actualValue != 3 || !ok {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	if actualValue, ok := heap.Pop(); actualValue != nil || ok {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue := heap.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := heap.Values(); len(actualValue) != 0 {
		t.Errorf("Got %v expected %v", actualValue, "[]")
	}
}

func TestBinaryHeapRandom(t *testing.T) {
	heap := NewWithIntComparator()

	rand.Seed(3)
	for i := 0; i < 10000; i++ {
		r := int(rand.Int31n(30))
		heap.Push(r)
	}

	prev, _ := heap.Pop()
	for !heap.Empty() {
		curr, _ := heap.Pop()
		if prev.(int) > curr.(int) {
			t.Errorf("Heap property invalidated. prev: %v current: %v", prev, curr)
		}
		prev = curr
	}
}

func TestBinaryHeapIteratorOnEmpty(t *testing.T) {
	heap := NewWithIntComparator()
	it := heap.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty heap")
	}
}

func TestBinaryHeapIteratorNext(t *testing.T) {
	heap := NewWithIntComparator()
	heap.Push(3) // [3]
	heap.Push(2) // [2,3]
	heap.Push(1) // [1,3,2](2 swapped with 1, hence last)

	it := heap.Iterator()
	count := 0
	for it.Next() {
		count++
		index := it.Index()
		value := it.Value()
		switch index {
		case 0:
			if actualValue, expectedValue := value, 1; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value, 3; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value, 2; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
		if actualValue, expectedValue := index, count-1; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}
	if actualValue, expectedValue := count, 3; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBinaryHeapIteratorPrev(t *testing.T) {
	heap := NewWithIntComparator()
	heap.Push(3) // [3]
	heap.Push(2) // [2,3]
	heap.Push(1) // [1,3,2](2 swapped with 1, hence last)

	it := heap.Iterator()
	for it.Next() {
	}
	count := 0
	for it.Prev() {
		count++
		index := it.Index()
		value := it.Value()
		switch index {
		case 0:
			if actualValue, expectedValue := value, 1; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value, 3; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value, 2; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
		if actualValue, expectedValue := index, 3-count; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}
	if actualValue, expectedValue := count, 3; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestBinaryHeapIteratorBegin(t *testing.T) {
	heap := NewWithIntComparator()
	it := heap.Iterator()
	it.Begin()
	heap.Push(2)
	heap.Push(3)
	heap.Push(1)
	for it.Next() {
	}
	it.Begin()
	it.Next()
	if index, value := it.Index(), it.Value(); index != 0 || value != 1 {
		t.Errorf("Got %v,%v expected %v,%v", index, value, 0, 1)
	}
}

func TestListIteratorEnd(t *testing.T) {
	heap := NewWithIntComparator()
	it := heap.Iterator()

	if index := it.Index(); index != -1 {
		t.Errorf("Got %v expected %v", index, -1)
	}

	it.End()
	if index := it.Index(); index != 0 {
		t.Errorf("Got %v expected %v", index, 0)
	}

	heap.Push(3) // [3]
	heap.Push(2) // [2,3]
	heap.Push(1) // [1,3,2](2 swapped with 1, hence last)
	it.End()
	if index := it.Index(); index != heap.Size() {
		t.Errorf("Got %v expected %v", index, heap.Size())
	}

	it.Prev()
	if index, value := it.Index(), it.Value(); index != heap.Size()-1 || value != 2 {
		t.Errorf("Got %v,%v expected %v,%v", index, value, heap.Size()-1, 2)
	}
}

func TestStackIteratorFirst(t *testing.T) {
	heap := NewWithIntComparator()
	it := heap.Iterator()
	if actualValue, expectedValue := it.First(), false; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	heap.Push(3) // [3]
	heap.Push(2) // [2,3]
	heap.Push(1) // [1,3,2](2 swapped with 1, hence last)
	if actualValue, expectedValue := it.First(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if index, value := it.Index(), it.Value(); index != 0 || value != 1 {
		t.Errorf("Got %v,%v expected %v,%v", index, value, 0, 1)
	}
}

func TestBinaryHeapIteratorLast(t *testing.T) {
	tree := NewWithIntComparator()
	it := tree.Iterator()
	if actualValue, expectedValue := it.Last(), false; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	tree.Push(2)
	tree.Push(3)
	tree.Push(1) // [1,3,2](2 swapped with 1, hence last)
	if actualValue, expectedValue := it.Last(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if index, value := it.Index(), it.Value(); index != 2 || value != 2 {
		t.Errorf("Got %v,%v expected %v,%v", index, value, 2, 2)
	}
}

func TestBinaryHeapSerialization(t *testing.T) {
	heap := NewWithStringComparator()

	heap.Push("c") // ["c"]
	heap.Push("b") // ["b","c"]
	heap.Push("a") // ["a","c","b"]("b" swapped with "a", hence last)

	var err error
	assert := func() {
		if actualValue := heap.Values(); actualValue[0].(string) != "a" || actualValue[1].(string) != "c" || actualValue[2].(string) != "b" {
			t.Errorf("Got %v expected %v", actualValue, "[1,3,2]")
		}
		if actualValue := heap.Size(); actualValue != 3 {
			t.Errorf("Got %v expected %v", actualValue, 3)
		}
		if actualValue, ok := heap.Peek(); actualValue != "a" || !ok {
			t.Errorf("Got %v expected %v", actualValue, "a")
		}
		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	json, err := heap.ToJSON()
	assert()

	err = heap.FromJSON(json)
	assert()
}

func benchmarkPush(b *testing.B, heap *Heap, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			heap.Push(n)
		}
	}
}

func benchmarkPop(b *testing.B, heap *Heap, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			heap.Pop()
		}
	}
}

func BenchmarkBinaryHeapPop100(b *testing.B) {
	b.StopTimer()
	size := 100
	heap := NewWithIntComparator()
	for n := 0; n < size; n++ {
		heap.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, heap, size)
}

func BenchmarkBinaryHeapPop1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	heap := NewWithIntComparator()
	for n := 0; n < size; n++ {
		heap.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, heap, size)
}

func BenchmarkBinaryHeapPop10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	heap := NewWithIntComparator()
	for n := 0; n < size; n++ {
		heap.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, heap, size)
}

func BenchmarkBinaryHeapPop100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	heap := NewWithIntComparator()
	for n := 0; n < size; n++ {
		heap.Push(n)
	}
	b.StartTimer()
	benchmarkPop(b, heap, size)
}

func BenchmarkBinaryHeapPush100(b *testing.B) {
	b.StopTimer()
	size := 100
	heap := NewWithIntComparator()
	b.StartTimer()
	benchmarkPush(b, heap, size)
}

func BenchmarkBinaryHeapPush1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	heap := NewWithIntComparator()
	for n := 0; n < size; n++ {
		heap.Push(n)
	}
	b.StartTimer()
	benchmarkPush(b, heap, size)
}

func BenchmarkBinaryHeapPush10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	heap := NewWithIntComparator()
	for n := 0; n < size; n++ {
		heap.Push(n)
	}
	b.StartTimer()
	benchmarkPush(b, heap, size)
}

func BenchmarkBinaryHeapPush100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	heap := NewWithIntComparator()
	for n := 0; n < size; n++ {
		heap.Push(n)
	}
	b.StartTimer()
	benchmarkPush(b, heap, size)
}
