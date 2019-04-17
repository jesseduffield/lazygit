package tst

import (
	//"fmt"
	"math/rand"
	"testing"
)

func randomString() string {
	n := 3 + rand.Intn(10)
	bs := make([]byte, n)
	for i := 0; i<n; i++ {
		bs[i] = byte(97 + rand.Intn(25))
	}
	return string(bs)
}

func Test(t *testing.T) {
	tree := New()
	tree.Insert("test", 1)
	if tree.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	if !tree.Has("test") {
		t.Errorf("expecting to find key=test")
	}
	
	tree.Insert("testing", 2)
	tree.Insert("abcd", 0)
		
	found := false
	tree.Do(func(key string, val interface{})bool {
		if key == "test" && val.(int) == 1 {
			found = true
		}
		return true
	})
	if !found {
		t.Errorf("expecting iterator to find test")
	}
	
	tree.Remove("testing")
	tree.Remove("abcd")
	
	v := tree.Remove("test")
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has("test") {
		t.Errorf("expecting not to find key=test")
	}
	if v.(int) != 1 {
		t.Errorf("expecting value=1")
	}
}

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()
	strs := make([]string, b.N)
	for i := 0; i<b.N; i++ {
		strs[i] = randomString()
	}
	b.StartTimer()
	
	tree := New()
	for i, str := range strs {
		tree.Insert(str, i)
	}
}

func BenchmarkMapInsert(b *testing.B) {
	b.StopTimer()
	strs := make([]string, b.N)
	for i := 0; i<b.N; i++ {
		strs[i] = randomString()
	}
	b.StartTimer()
	
	m := make(map[string]int)
	for i, str := range strs {
		m[str] = i
	}
}