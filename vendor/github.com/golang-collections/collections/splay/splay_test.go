package splay

import (
	//"fmt"
	"testing"
)

func Test(t *testing.T) {
	tree := New(func(a,b interface{})bool {
		return a.(string) < b.(string)
	})
	
	tree.Insert("d", 4)
	tree.Insert("b", 2)
	tree.Insert("a", 1)
	tree.Insert("c", 3)
	
	if tree.Len() != 4 {
		t.Errorf("expecting len 4")
	}

	tree.Remove("b")	
	
	if tree.Len() != 3 {
		t.Errorf("expecting len 3")
	}
}