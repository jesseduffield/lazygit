package queue

import (
	"testing"
)

func Test(t *testing.T) {
	q := New()
	
	if q.Len() != 0 {
		t.Errorf("Length should be 0")		
	}
	
	q.Enqueue(1)
	
	if q.Len() != 1 {
		t.Errorf("Length should be 1")
	}
	
	if q.Peek().(int) != 1 {
		t.Errorf("Enqueued value should be 1")
	}
	
	v := q.Dequeue()
	
	if v.(int) != 1 {
		t.Errorf("Dequeued value should be 1")
	}
	
	if q.Peek() != nil || q.Dequeue() != nil {
		t.Errorf("Empty queue should have no values")
	}
	
	q.Enqueue(1)	
	q.Enqueue(2)
	
	if q.Len() != 2 {
		t.Errorf("Length should be 2")
	}
	
	if q.Peek().(int) != 1 {
		t.Errorf("First value should be 1")
	}
	
	q.Dequeue()
	
	if q.Peek().(int) != 2 {
		t.Errorf("Next value should be 2")
	}
}