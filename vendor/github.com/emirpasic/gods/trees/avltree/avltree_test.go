// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package avltree

import (
	"fmt"
	"testing"
)

func TestAVLTreePut(t *testing.T) {
	tree := NewWithIntComparator()
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite

	if actualValue := tree.Size(); actualValue != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}
	if actualValue, expectedValue := fmt.Sprintf("%d%d%d%d%d%d%d", tree.Keys()...), "1234567"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s%s%s%s%s%s%s", tree.Values()...), "abcdefg"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	tests1 := [][]interface{}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, nil, false},
	}

	for _, test := range tests1 {
		// retrievals
		actualValue, actualFound := tree.Get(test[0])
		if actualValue != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}
}

func TestAVLTreeRemove(t *testing.T) {
	tree := NewWithIntComparator()
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite

	tree.Remove(5)
	tree.Remove(6)
	tree.Remove(7)
	tree.Remove(8)
	tree.Remove(5)

	if actualValue, expectedValue := fmt.Sprintf("%d%d%d%d", tree.Keys()...), "1234"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s%s%s%s", tree.Values()...), "abcd"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s%s%s%s", tree.Values()...), "abcd"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue := tree.Size(); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}

	tests2 := [][]interface{}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, nil, false},
		{6, nil, false},
		{7, nil, false},
		{8, nil, false},
	}

	for _, test := range tests2 {
		actualValue, actualFound := tree.Get(test[0])
		if actualValue != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}

	tree.Remove(1)
	tree.Remove(4)
	tree.Remove(2)
	tree.Remove(3)
	tree.Remove(2)
	tree.Remove(2)

	if actualValue, expectedValue := fmt.Sprintf("%s", tree.Keys()), "[]"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s", tree.Values()), "[]"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if empty, size := tree.Empty(), tree.Size(); empty != true || size != -0 {
		t.Errorf("Got %v expected %v", empty, true)
	}

}

func TestAVLTreeLeftAndRight(t *testing.T) {
	tree := NewWithIntComparator()

	if actualValue := tree.Left(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue := tree.Right(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	tree.Put(1, "a")
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x") // overwrite
	tree.Put(2, "b")

	if actualValue, expectedValue := fmt.Sprintf("%d", tree.Left().Key), "1"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s", tree.Left().Value), "x"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := fmt.Sprintf("%d", tree.Right().Key), "7"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s", tree.Right().Value), "g"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestAVLTreeCeilingAndFloor(t *testing.T) {
	tree := NewWithIntComparator()

	if node, found := tree.Floor(0); node != nil || found {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}
	if node, found := tree.Ceiling(0); node != nil || found {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")

	if node, found := tree.Floor(4); node.Key != 4 || !found {
		t.Errorf("Got %v expected %v", node.Key, 4)
	}
	if node, found := tree.Floor(0); node != nil || found {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	if node, found := tree.Ceiling(4); node.Key != 4 || !found {
		t.Errorf("Got %v expected %v", node.Key, 4)
	}
	if node, found := tree.Ceiling(8); node != nil || found {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}
}

func TestAVLTreeIteratorNextOnEmpty(t *testing.T) {
	tree := NewWithIntComparator()
	it := tree.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty tree")
	}
}

func TestAVLTreeIteratorPrevOnEmpty(t *testing.T) {
	tree := NewWithIntComparator()
	it := tree.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty tree")
	}
}

func TestAVLTreeIterator1Next(t *testing.T) {
	tree := NewWithIntComparator()
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite
	// │   ┌── 7
	// └── 6
	//     │   ┌── 5
	//     └── 4
	//         │   ┌── 3
	//         └── 2
	//             └── 1
	it := tree.Iterator()
	count := 0
	for it.Next() {
		count++
		key := it.Key()
		switch key {
		case count:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestAVLTreeIterator1Prev(t *testing.T) {
	tree := NewWithIntComparator()
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite
	// │   ┌── 7
	// └── 6
	//     │   ┌── 5
	//     └── 4
	//         │   ┌── 3
	//         └── 2
	//             └── 1
	it := tree.Iterator()
	for it.Next() {
	}
	countDown := tree.size
	for it.Prev() {
		key := it.Key()
		switch key {
		case countDown:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
		countDown--
	}
	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestAVLTreeIterator2Next(t *testing.T) {
	tree := NewWithIntComparator()
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")
	it := tree.Iterator()
	count := 0
	for it.Next() {
		count++
		key := it.Key()
		switch key {
		case count:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestAVLTreeIterator2Prev(t *testing.T) {
	tree := NewWithIntComparator()
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")
	it := tree.Iterator()
	for it.Next() {
	}
	countDown := tree.size
	for it.Prev() {
		key := it.Key()
		switch key {
		case countDown:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
		countDown--
	}
	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestAVLTreeIterator3Next(t *testing.T) {
	tree := NewWithIntComparator()
	tree.Put(1, "a")
	it := tree.Iterator()
	count := 0
	for it.Next() {
		count++
		key := it.Key()
		switch key {
		case count:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestAVLTreeIterator3Prev(t *testing.T) {
	tree := NewWithIntComparator()
	tree.Put(1, "a")
	it := tree.Iterator()
	for it.Next() {
	}
	countDown := tree.size
	for it.Prev() {
		key := it.Key()
		switch key {
		case countDown:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
		countDown--
	}
	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestAVLTreeIterator4Next(t *testing.T) {
	tree := NewWithIntComparator()
	tree.Put(13, 5)
	tree.Put(8, 3)
	tree.Put(17, 7)
	tree.Put(1, 1)
	tree.Put(11, 4)
	tree.Put(15, 6)
	tree.Put(25, 9)
	tree.Put(6, 2)
	tree.Put(22, 8)
	tree.Put(27, 10)
	// │           ┌── 27
	// │       ┌── 25
	// │       │   └── 22
	// │   ┌── 17
	// │   │   └── 15
	// └── 13
	//     │   ┌── 11
	//     └── 8
	//         │   ┌── 6
	//         └── 1
	it := tree.Iterator()
	count := 0
	for it.Next() {
		count++
		value := it.Value()
		switch value {
		case count:
			if actualValue, expectedValue := value, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := value, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestAVLTreeIterator4Prev(t *testing.T) {
	tree := NewWithIntComparator()
	tree.Put(13, 5)
	tree.Put(8, 3)
	tree.Put(17, 7)
	tree.Put(1, 1)
	tree.Put(11, 4)
	tree.Put(15, 6)
	tree.Put(25, 9)
	tree.Put(6, 2)
	tree.Put(22, 8)
	tree.Put(27, 10)
	// │           ┌── 27
	// │       ┌── 25
	// │       │   └── 22
	// │   ┌── 17
	// │   │   └── 15
	// └── 13
	//     │   ┌── 11
	//     └── 8
	//         │   ┌── 6
	//         └── 1
	it := tree.Iterator()
	count := tree.Size()
	for it.Next() {
	}
	for it.Prev() {
		value := it.Value()
		switch value {
		case count:
			if actualValue, expectedValue := value, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := value, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
		count--
	}
	if actualValue, expectedValue := count, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestAVLTreeIteratorBegin(t *testing.T) {
	tree := NewWithIntComparator()
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")
	it := tree.Iterator()

	if it.Key() != nil {
		t.Errorf("Got %v expected %v", it.Key(), nil)
	}

	it.Begin()

	if it.Key() != nil {
		t.Errorf("Got %v expected %v", it.Key(), nil)
	}

	for it.Next() {
	}

	it.Begin()

	if it.Key() != nil {
		t.Errorf("Got %v expected %v", it.Key(), nil)
	}

	it.Next()
	if key, value := it.Key(), it.Value(); key != 1 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 1, "a")
	}
}

func TestAVLTreeIteratorEnd(t *testing.T) {
	tree := NewWithIntComparator()
	it := tree.Iterator()

	if it.Key() != nil {
		t.Errorf("Got %v expected %v", it.Key(), nil)
	}

	it.End()
	if it.Key() != nil {
		t.Errorf("Got %v expected %v", it.Key(), nil)
	}

	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")
	it.End()
	if it.Key() != nil {
		t.Errorf("Got %v expected %v", it.Key(), nil)
	}

	it.Prev()
	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
	}
}

func TestAVLTreeIteratorFirst(t *testing.T) {
	tree := NewWithIntComparator()
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")
	it := tree.Iterator()
	if actualValue, expectedValue := it.First(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if key, value := it.Key(), it.Value(); key != 1 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 1, "a")
	}
}

func TestAVLTreeIteratorLast(t *testing.T) {
	tree := NewWithIntComparator()
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")
	it := tree.Iterator()
	if actualValue, expectedValue := it.Last(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
	}
}

func TestAVLTreeSerialization(t *testing.T) {
	tree := NewWithStringComparator()
	tree.Put("c", "3")
	tree.Put("b", "2")
	tree.Put("a", "1")

	var err error
	assert := func() {
		if actualValue, expectedValue := tree.Size(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if actualValue := tree.Keys(); actualValue[0].(string) != "a" || actualValue[1].(string) != "b" || actualValue[2].(string) != "c" {
			t.Errorf("Got %v expected %v", actualValue, "[a,b,c]")
		}
		if actualValue := tree.Values(); actualValue[0].(string) != "1" || actualValue[1].(string) != "2" || actualValue[2].(string) != "3" {
			t.Errorf("Got %v expected %v", actualValue, "[1,2,3]")
		}
		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	json, err := tree.ToJSON()
	assert()

	err = tree.FromJSON(json)
	assert()
}

func benchmarkGet(b *testing.B, tree *Tree, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Get(n)
		}
	}
}

func benchmarkPut(b *testing.B, tree *Tree, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Put(n, struct{}{})
		}
	}
}

func benchmarkRemove(b *testing.B, tree *Tree, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Remove(n)
		}
	}
}

func BenchmarkAVLTreeGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkAVLTreeGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkAVLTreeGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkAVLTreeGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkGet(b, tree, size)
}

func BenchmarkAVLTreePut100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewWithIntComparator()
	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkAVLTreePut1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkAVLTreePut10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkAVLTreePut100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkPut(b, tree, size)
}

func BenchmarkAVLTreeRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, tree, size)
}

func BenchmarkAVLTreeRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, tree, size)
}

func BenchmarkAVLTreeRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, tree, size)
}

func BenchmarkAVLTreeRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkRemove(b, tree, size)
}
