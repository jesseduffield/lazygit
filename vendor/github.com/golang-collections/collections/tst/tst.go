package tst

import (
	"fmt"
)

type (
	Node struct {
		key byte
		value interface{}
		left, middle, right *Node
	}
	NodeIterator struct {
		step int
		node *Node
		prev *NodeIterator
	}
	TernarySearchTree struct {
		length int
		root *Node
	}
)

// Create a new ternary search tree
func New() *TernarySearchTree {
	tree := &TernarySearchTree{}
	tree.Init()
	return tree
}
// Iterate over the collection
func (this *TernarySearchTree) Do(callback func(string, interface{})bool) {
	if this.Len() == 0 {
		return
	}
	bs := []byte{}
	i := &NodeIterator{0,this.root,nil}
	for i != nil {
		switch i.step {
		// Left
		case 0:
			i.step++
			if i.node.left != nil {
				i = &NodeIterator{0,i.node.left,i}
				continue
			}
		// Value
		case 1:
			i.step++
			if i.node.key > 0 {
				bs = append(bs, i.node.key)
			}
			if i.node.value != nil {
				if !callback(string(bs), i.node.value) {
					return
				}
				continue
			}
		// Middle
		case 2:
			i.step++
			if i.node.middle != nil {
				i = &NodeIterator{0,i.node.middle,i}
				continue
			}
		// Right
		case 3:
			if len(bs) > 0 {
				bs = bs[:len(bs)-1]
			}
			i.step++
			if i.node.right != nil {
				i = &NodeIterator{0,i.node.right,i}
				continue
			}
		// Backtrack
		case 4:
			i = i.prev
		}
	}
}
// Get the value at the specified key. Returns nil if not found.
func (this *TernarySearchTree) Get(key string) interface{} {
	if this.length == 0 {
		return nil
	}

	node := this.root
	bs := []byte(key)
	for i := 0; i < len(bs); {
		b := bs[i]
		if b > node.key {
			if node.right == nil {
				return nil
			}
			node = node.right
		} else if (b < node.key) {
			if node.left == nil {
				return nil
			}
			node = node.left
		} else {
			i++
			if i < len(bs) {
				if node.middle == nil {
					return nil
				}
				node = node.middle
			} else {
				break
			}
		}
	}
	return node.value
}
func (this *TernarySearchTree) GetLongestPrefix(key string) interface{} {
	if this.length == 0 {
		return nil
	}

	n := this.root
	v := n.value
	bs := []byte(key)
	for i := 0; i < len(bs); {
		b := bs[i]
		if n.value != nil {
			v = n.value
		}
		if b > n.key {
			if n.right == nil {
				break
			}
			n = n.right
		} else if b < n.key {
			if n.left == nil {
				break
			}
			n = n.left
		} else {
			i++
			if i < len(bs) {
				if n.middle == nil {
					break
				}
				n = n.middle
			} else {
				break
			}
		}
	}
	if n.value != nil {
		v = n.value
	}
	return v
}
// Test to see whether or not the given key is contained in the tree.
func (this *TernarySearchTree) Has(key string) bool {
	return this.Get(key) != nil
}
// Initialize the tree (reset it so that it's empty). New will do this for you.
func (this *TernarySearchTree) Init() {
	this.length = 0
	this.root = nil
}
// Insert a new key value pair into the collection
func (this *TernarySearchTree) Insert(key string, value interface{}) {
	// If the value is nil then remove this key from the collection
	if value == nil {
		this.Remove(key)
		return
	}

	if this.length == 0 {
		this.root = &Node{0,nil,nil,nil,nil}
	}

	t := this.root
	bs := []byte(key)
	for i := 0; i < len(bs); {
		b := bs[i]
		if b > t.key {
			if t.right == nil {
				t.right = &Node{b,nil,nil,nil,nil}
			}
			t = t.right
		} else if b < t.key {
			if t.left == nil {
				t.left = &Node{b,nil,nil,nil,nil}
			}
			t = t.left
		} else {
			i++
			if i < len(bs) {
				if t.middle == nil {
					t.middle = &Node{bs[i],nil,nil,nil,nil}
				}
				t = t.middle
			}
		}
	}

	if t.value == nil {
		this.length++
	}
	t.value = value
}
// Get the number of items stored in the tree
func (this *TernarySearchTree) Len() int {
	return this.length
}
// Remove a key from the collection
func (this *TernarySearchTree) Remove(key string) interface{} {
	if this.length == 0 {
		return nil
	}

	var remove *Node
	var direction int

	t := this.root
	bs := []byte(key)
	for i := 0; i < len(bs); {
		b := bs[i]
		if b > t.key {
			// Not in the collection
			if t.right == nil {
				return nil
			}
			// This is a branch so we have to keep it
			remove = t
			direction = 1
			// Move to the next node
			t = t.right
		} else if b < t.key {
			// Not in the collection
			if t.left == nil {
				return nil
			}
			// This is a branch so we have to keep it
			remove = t
			direction = -1
			// Move to the next node
			t = t.left
		} else {
			i++
			if i < len(bs) {
				// Not in the collection
				if t.middle == nil {
					return nil
				}
				// Has a value so we need to keep at least this much
				if t.value != nil {
					remove = t
					direction = 0
				}
				// Move to the next node
				t = t.middle
			}
		}
	}

	// If this was the only item in the tree, set the root pointer to nil
	if this.length == 1 {
		this.root = nil
	} else {
		if direction == -1 {
			remove.left = nil
		} else if direction == 0 {
			remove.middle = nil
		} else {
			remove.right = nil
		}
	}
	this.length--
	return t.value
}
func (this *TernarySearchTree) String() string {
	if this.length == 0 {
		return "{}"
	}

	return this.root.String()
}
// Dump the tree to a string for easier debugging
func (this *Node) String() string {
	str := "{" + string(this.key)
	if this.value != nil {
		str += ":" + fmt.Sprint(this.value)
	}
	if this.left != nil {
		str += this.left.String()
	} else {
		str += " "
	}
	if this.middle != nil {
		str += this.middle.String()
	} else {
		str += " "
	}
	if this.right != nil {
		str += this.right.String()
	} else {
		str += " "
	}
	str += "}"
	return str
}
