package trie

import (
	"fmt"
)

type (
	Trie struct {
		root *node
		size int
	}
	node struct {
		key interface{}
		value interface{}
		next [256]*node
	}
	iterator struct {
		step int
		node *node
		prev *iterator
	}
)

func toBytes(obj interface{}) []byte {
	switch o := obj.(type) {
	case []byte:
		return o
	case string:
		return []byte(o)
	}
	return []byte(fmt.Sprint(obj))
}

func New() *Trie {
	return &Trie{nil,0}
}
func (this *Trie) Do(handler func(interface{},interface{})bool) {
	if this.size > 0 {
		this.root.do(handler)
	}
}
func (this *Trie) Get(key interface{}) interface{} {
	if this.size == 0 {
		return nil
	}
	
	bs := toBytes(key)
	cur := this.root
	for i := 0; i < len(bs); i++ {
		if cur.next[bs[i]] != nil {
			cur = cur.next[bs[i]]
		} else {
			return nil
		}
	}
	return cur.value
}
func (this *Trie) Has(key interface{}) bool {
	return this.Get(key) != nil
}
func (this *Trie) Init() {
	this.root = nil
	this.size = 0
}
func (this *Trie) Insert(key interface{}, value interface{}) {
	if this.size == 0 {
		this.root = newNode()
	}
	
	bs := toBytes(key)
	cur := this.root
	for i := 0; i < len(bs); i++ {
		if cur.next[bs[i]] != nil {
			cur = cur.next[bs[i]]
		} else {
			cur.next[bs[i]] = newNode()
			cur = cur.next[bs[i]]
		}
	}	
	if cur.key == nil {
		this.size++
	}
	cur.key = key
	cur.value = value
}
func (this *Trie) Len() int {
	return this.size
}
func (this *Trie) Remove(key interface{}) interface{} {
	if this.size == 0 {
		return nil
	}
	bs := toBytes(key)
	cur := this.root
	
	for i := 0; i < len(bs); i++ {
		if cur.next[bs[i]] != nil {
			cur = cur.next[bs[i]]
		} else {
			return nil
		}
	}
	
	// TODO: cleanup dead nodes
	
	val := cur.value
	
	if cur.value != nil {
		this.size--
		cur.value = nil
		cur.key = nil
	}
	return val
}
func (this *Trie) String() string {
	str := "{"
	i := 0
	this.Do(func(k, v interface{}) bool {
		if i > 0 {
			str += ", "
		}
		str += fmt.Sprint(k, ":", v)
		i++
		return true
	})
	str += "}"
	return str
}

func newNode() *node {
	var next [256]*node
	return &node{nil,nil,next}
}
func (this *node) do(handler func(interface{}, interface{}) bool) bool {
	for i := 0; i < 256; i++ {
		if this.next[i] != nil {
			if this.next[i].key != nil {
				if !handler(this.next[i].key, this.next[i].value) {
					return false
				}
			}
			if !this.next[i].do(handler) {
				return false
			}
		}
	}
	return true
}