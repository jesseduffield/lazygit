package set

import (
	"testing"
)

func Test(t *testing.T) {
	s := New()
	
	s.Insert(5)
	
	if s.Len() != 1 {
		t.Errorf("Length should be 1")
	}
	
	if !s.Has(5) { 
		t.Errorf("Membership test failed")
	}
	
	s.Remove(5)
	
	if s.Len() != 0 {
		t.Errorf("Length should be 0")
	}
	
	if s.Has(5) {
		t.Errorf("The set should be empty")
	}
	
	// Difference
	s1 := New(1,2,3,4,5,6)
	s2 := New(4,5,6)
	s3 := s1.Difference(s2)
	
	if s3.Len() != 3 {
		t.Errorf("Length should be 3")
	}
	
	if !(s3.Has(1) && s3.Has(2) && s3.Has(3)) {
		t.Errorf("Set should only contain 1, 2, 3")
	}
	
	// Intersection
	s3 = s1.Intersection(s2)
	if s3.Len() != 3 {
		t.Errorf("Length should be 3 after intersection")
	}
	
	if !(s3.Has(4) && s3.Has(5) && s3.Has(6)) {
		t.Errorf("Set should contain 4, 5, 6")
	}
	
	// Union
	s4 := New(7,8,9)
	s3 = s2.Union(s4)
	
	if s3.Len() != 6 {
		t.Errorf("Length should be 6 after union")
	}
	
	if !(s3.Has(7)) {
		t.Errorf("Set should contain 4, 5, 6, 7, 8, 9")
	}
	
	// Subset
	if !s1.SubsetOf(s1) {
		t.Errorf("set should be a subset of itself")
	}
	// Proper Subset
	if s1.ProperSubsetOf(s1) {
		t.Errorf("set should not be a subset of itself")
	}
	
}
