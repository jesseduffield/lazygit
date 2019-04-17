package skip

import (
	//"fmt"
	"testing"
)

func Test(t *testing.T) {
	sl := New(func(a,b interface{})bool {
		return a.(int) < b.(int)
	})
	sl.Insert(1, 100)
	if sl.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	sl.Insert(1, 1000)
	if sl.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	if sl.Get(1).(int) != 1000 {
		t.Errorf("expecting sl[1] == 1000")
	}
	sl.Remove(1)
	if sl.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	
	sl.Insert(2, 200)
	sl.Insert(1, 100)
	vs := make([]int, 0)
	sl.Do(func(k, v interface{}) bool {
		vs = append(vs, k.(int))
		return true
	})
	if len(vs) != 2 || vs[0] != 1 || vs[1] != 2 {
		t.Errorf("expecting sorted iteration of all keys")
	}
}
