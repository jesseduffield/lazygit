package utils

import (
	"testing"
)

func TestThreadSafeMap(t *testing.T) {
	m := NewThreadSafeMap[int, int]()

	m.Set(1, 1)
	m.Set(2, 2)
	m.Set(3, 3)

	if m.Len() != 3 {
		t.Errorf("Expected length to be 3, got %d", m.Len())
	}

	if !m.Has(1) {
		t.Errorf("Expected to have key 1")
	}

	if m.Has(4) {
		t.Errorf("Expected to not have key 4")
	}

	if _, ok := m.Get(1); !ok {
		t.Errorf("Expected to have key 1")
	}

	if _, ok := m.Get(4); ok {
		t.Errorf("Expected to not have key 4")
	}

	m.Delete(1)

	if m.Has(1) {
		t.Errorf("Expected to not have key 1")
	}

	m.Clear()

	if m.Len() != 0 {
		t.Errorf("Expected length to be 0, got %d", m.Len())
	}
}

func TestThreadSafeMapConcurrentReadWrite(t *testing.T) {
	m := NewThreadSafeMap[int, int]()

	go func() {
		for i := 0; i < 10000; i++ {
			m.Set(0, 0)
		}
	}()

	for i := 0; i < 10000; i++ {
		m.Get(0)
	}
}
