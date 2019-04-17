package skip

import (
	"fmt"
	"math/rand"
	"time"
)

type (
	node struct {
		next []*node
		key interface{}
		value interface{}
	}
	SkipList struct {
		root *node
		size int
		less func(interface{},interface{})bool
		gen *rand.Rand
		probability float64
	}
)
// Create a new skip list
func New(less func(interface{},interface{})bool) *SkipList {
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := &node{make([]*node, 0),nil,nil}
	return &SkipList{n, 0, less, gen, 0.75}
}
func (this *SkipList) Do(f func(interface{}, interface{})bool) {
	if this.size == 0 {
		return
	}
	cur := this.root.next[0]
	for cur != nil {
		if !f(cur.key, cur.value) {
			break
		}
		cur = cur.next[0]
	}
}
// Get an item from the skip list
func (this *SkipList) Get(key interface{}) interface{} {
	if this.size == 0 {
		return nil
	}
	
	cur := this.root
	// Start at the top
	for i := len(cur.next)-1; i >= 0; i-- {
		for this.less(cur.next[i].key, key) {
			cur = cur.next[i]
		}
	}
	cur = cur.next[0]
	
	if this.equals(cur.key, key) {
		return cur.value
	}
	
	return nil
}
// Insert a new item into the skip list
func (this *SkipList) Insert(key interface{}, value interface{}) {
	prev := this.getPrevious(key)
	
	// Already in the list so just update the value
	if len(prev) > 0 && prev[0].next[0] != nil && this.equals(prev[0].next[0].key, key) {
		prev[0].next[0].value = value
		return
	}

	h := len(this.root.next)	
	nh := this.pickHeight()
	n := &node{make([]*node, nh),key,value}
	
	// Higher than anything seen before, so tack it on top
	if nh > h {
		this.root.next = append(this.root.next, n)
	}	
	
	// Update the previous nodes
	for i := 0; i < h && i < nh; i++ {
		n.next[i] = prev[i].next[i]
		prev[i].next[i] = n
	}
	
	this.size++
}
// Get the length of the skip list
func (this *SkipList) Len() int {
	return this.size
}
// Remove an item from the skip list
func (this *SkipList) Remove(key interface{}) interface{} {
	prev := this.getPrevious(key)
	if len(prev) == 0 {
		return nil
	}
	cur := prev[0].next[0]
	
	// If we found it
	if cur != nil && this.equals(key, cur.key) {
		// Change all the linked lists
		for i := 0; i < len(prev); i++ {
			if prev[i] != nil && prev[i].next[i] != nil {
				prev[i].next[i] = cur.next[i]
			}
		}
		
		// Kill off the upper links if they're nil
		for i := len(this.root.next)-1; i>=0; i-- {
			if this.root.next[i] == nil {
				this.root.next = this.root.next[:i]
			} else {
				break
			}
		}
		
		this.size--
		
		return cur.value
	}
	
	return nil
}
// String representation of the list
func (this *SkipList) String() string {	
	str := "{"
	if len(this.root.next) > 0 {
		cur := this.root.next[0]
		for cur != nil {
			str += fmt.Sprint(cur.key)
			str += ":"
			str += fmt.Sprint(cur.value)
			str += " "
			cur = cur.next[0]
		}
	}
	str += "}"
	
	return str
}
// Get a vertical list of nodes of all the things that occur
//  immediately before "key"
func (this *SkipList) getPrevious(key interface{}) []*node {
	cur := this.root
	h := len(cur.next)
	nodes := make([]*node, h)
	for i := h-1; i >= 0; i-- {
		for cur.next[i] != nil && this.less(cur.next[i].key, key) {
			cur = cur.next[i]
		}
		nodes[i] = cur
	}
	return nodes
}
// Defines an equals method in terms of "less"
func (this *SkipList) equals(a, b interface{}) bool {
	return !this.less(a,b) && !this.less(b,a)
}
// Pick a random height
func (this *SkipList) pickHeight() int {
	h := 1
	for this.gen.Float64() > this.probability {
		h++
	}
	if h > len(this.root.next) {
		return h + 1
	}
	return h
}
