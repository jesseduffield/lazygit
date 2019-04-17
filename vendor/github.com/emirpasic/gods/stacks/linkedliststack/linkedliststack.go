// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package linkedliststack implements a stack backed by a singly-linked list.
//
// Structure is not thread safe.
//
// Reference:https://en.wikipedia.org/wiki/Stack_%28abstract_data_type%29#Linked_list
package linkedliststack

import (
	"fmt"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/emirpasic/gods/stacks"
	"strings"
)

func assertStackImplementation() {
	var _ stacks.Stack = (*Stack)(nil)
}

// Stack holds elements in a singly-linked-list
type Stack struct {
	list *singlylinkedlist.List
}

// New nnstantiates a new empty stack
func New() *Stack {
	return &Stack{list: &singlylinkedlist.List{}}
}

// Push adds a value onto the top of the stack
func (stack *Stack) Push(value interface{}) {
	stack.list.Prepend(value)
}

// Pop removes top element on stack and returns it, or nil if stack is empty.
// Second return parameter is true, unless the stack was empty and there was nothing to pop.
func (stack *Stack) Pop() (value interface{}, ok bool) {
	value, ok = stack.list.Get(0)
	stack.list.Remove(0)
	return
}

// Peek returns top element on the stack without removing it, or nil if stack is empty.
// Second return parameter is true, unless the stack was empty and there was nothing to peek.
func (stack *Stack) Peek() (value interface{}, ok bool) {
	return stack.list.Get(0)
}

// Empty returns true if stack does not contain any elements.
func (stack *Stack) Empty() bool {
	return stack.list.Empty()
}

// Size returns number of elements within the stack.
func (stack *Stack) Size() int {
	return stack.list.Size()
}

// Clear removes all elements from the stack.
func (stack *Stack) Clear() {
	stack.list.Clear()
}

// Values returns all elements in the stack (LIFO order).
func (stack *Stack) Values() []interface{} {
	return stack.list.Values()
}

// String returns a string representation of container
func (stack *Stack) String() string {
	str := "LinkedListStack\n"
	values := []string{}
	for _, value := range stack.list.Values() {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

// Check that the index is within bounds of the list
func (stack *Stack) withinRange(index int) bool {
	return index >= 0 && index < stack.list.Size()
}
