package trie

import (
	//"fmt"
	"testing"
)

func Test(t *testing.T) {
	x := New()
	x.Insert(1, 100)
	if x.Len() != 1 {
		t.Errorf("expected len 1")
	}
	if x.Get(1).(int) != 100 {
		t.Errorf("expected to get 100 for 1")
	}
	x.Remove(1)
	if x.Len() != 0 {
		t.Errorf("expected len 0")
	}
	x.Insert(2, 200)
	x.Insert(1, 100)
	vs := make([]int, 0)
	x.Do(func(k, v interface{}) bool {
		vs = append(vs, k.(int))
		return true
	})
	if len(vs) != 2 || vs[0] != 1 || vs[1] != 2 {
		t.Errorf("expected in order traversal")
	}
}