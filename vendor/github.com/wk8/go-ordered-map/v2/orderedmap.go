// Package orderedmap implements an ordered map, i.e. a map that also keeps track of
// the order in which keys were inserted.
//
// All operations are constant-time.
//
// Github repo: https://github.com/wk8/go-ordered-map
//
package orderedmap

import (
	"fmt"

	list "github.com/bahlo/generic-list-go"
)

type Pair[K comparable, V any] struct {
	Key   K
	Value V

	element *list.Element[*Pair[K, V]]
}

type OrderedMap[K comparable, V any] struct {
	pairs map[K]*Pair[K, V]
	list  *list.List[*Pair[K, V]]
}

type initConfig[K comparable, V any] struct {
	capacity    int
	initialData []Pair[K, V]
}

type InitOption[K comparable, V any] func(config *initConfig[K, V])

// WithCapacity allows giving a capacity hint for the map, akin to the standard make(map[K]V, capacity).
func WithCapacity[K comparable, V any](capacity int) InitOption[K, V] {
	return func(c *initConfig[K, V]) {
		c.capacity = capacity
	}
}

// WithInitialData allows passing in initial data for the map.
func WithInitialData[K comparable, V any](initialData ...Pair[K, V]) InitOption[K, V] {
	return func(c *initConfig[K, V]) {
		c.initialData = initialData
		if c.capacity < len(initialData) {
			c.capacity = len(initialData)
		}
	}
}

// New creates a new OrderedMap.
// options can either be one or several InitOption[K, V], or a single integer,
// which is then interpreted as a capacity hint, à la make(map[K]V, capacity).
func New[K comparable, V any](options ...any) *OrderedMap[K, V] { //nolint:varnamelen
	orderedMap := &OrderedMap[K, V]{}

	var config initConfig[K, V]
	for _, untypedOption := range options {
		switch option := untypedOption.(type) {
		case int:
			if len(options) != 1 {
				invalidOption()
			}
			config.capacity = option

		case InitOption[K, V]:
			option(&config)

		default:
			invalidOption()
		}
	}

	orderedMap.initialize(config.capacity)
	orderedMap.AddPairs(config.initialData...)

	return orderedMap
}

const invalidOptionMessage = `when using orderedmap.New[K,V]() with options, either provide one or several InitOption[K, V]; or a single integer which is then interpreted as a capacity hint, à la make(map[K]V, capacity).` //nolint:lll

func invalidOption() { panic(invalidOptionMessage) }

func (om *OrderedMap[K, V]) initialize(capacity int) {
	om.pairs = make(map[K]*Pair[K, V], capacity)
	om.list = list.New[*Pair[K, V]]()
}

// Get looks for the given key, and returns the value associated with it,
// or V's nil value if not found. The boolean it returns says whether the key is present in the map.
func (om *OrderedMap[K, V]) Get(key K) (val V, present bool) {
	if pair, present := om.pairs[key]; present {
		return pair.Value, true
	}

	return
}

// Load is an alias for Get, mostly to present an API similar to `sync.Map`'s.
func (om *OrderedMap[K, V]) Load(key K) (V, bool) {
	return om.Get(key)
}

// Value returns the value associated with the given key or the zero value.
func (om *OrderedMap[K, V]) Value(key K) (val V) {
	if pair, present := om.pairs[key]; present {
		val = pair.Value
	}
	return
}

// GetPair looks for the given key, and returns the pair associated with it,
// or nil if not found. The Pair struct can then be used to iterate over the ordered map
// from that point, either forward or backward.
func (om *OrderedMap[K, V]) GetPair(key K) *Pair[K, V] {
	return om.pairs[key]
}

// Set sets the key-value pair, and returns what `Get` would have returned
// on that key prior to the call to `Set`.
func (om *OrderedMap[K, V]) Set(key K, value V) (val V, present bool) {
	if pair, present := om.pairs[key]; present {
		oldValue := pair.Value
		pair.Value = value
		return oldValue, true
	}

	pair := &Pair[K, V]{
		Key:   key,
		Value: value,
	}
	pair.element = om.list.PushBack(pair)
	om.pairs[key] = pair

	return
}

// AddPairs allows setting multiple pairs at a time. It's equivalent to calling
// Set on each pair sequentially.
func (om *OrderedMap[K, V]) AddPairs(pairs ...Pair[K, V]) {
	for _, pair := range pairs {
		om.Set(pair.Key, pair.Value)
	}
}

// Store is an alias for Set, mostly to present an API similar to `sync.Map`'s.
func (om *OrderedMap[K, V]) Store(key K, value V) (V, bool) {
	return om.Set(key, value)
}

// Delete removes the key-value pair, and returns what `Get` would have returned
// on that key prior to the call to `Delete`.
func (om *OrderedMap[K, V]) Delete(key K) (val V, present bool) {
	if pair, present := om.pairs[key]; present {
		om.list.Remove(pair.element)
		delete(om.pairs, key)
		return pair.Value, true
	}
	return
}

// Len returns the length of the ordered map.
func (om *OrderedMap[K, V]) Len() int {
	if om == nil || om.pairs == nil {
		return 0
	}
	return len(om.pairs)
}

// Oldest returns a pointer to the oldest pair. It's meant to be used to iterate on the ordered map's
// pairs from the oldest to the newest, e.g.:
// for pair := orderedMap.Oldest(); pair != nil; pair = pair.Next() { fmt.Printf("%v => %v\n", pair.Key, pair.Value) }
func (om *OrderedMap[K, V]) Oldest() *Pair[K, V] {
	if om == nil || om.list == nil {
		return nil
	}
	return listElementToPair(om.list.Front())
}

// Newest returns a pointer to the newest pair. It's meant to be used to iterate on the ordered map's
// pairs from the newest to the oldest, e.g.:
// for pair := orderedMap.Oldest(); pair != nil; pair = pair.Next() { fmt.Printf("%v => %v\n", pair.Key, pair.Value) }
func (om *OrderedMap[K, V]) Newest() *Pair[K, V] {
	if om == nil || om.list == nil {
		return nil
	}
	return listElementToPair(om.list.Back())
}

// Next returns a pointer to the next pair.
func (p *Pair[K, V]) Next() *Pair[K, V] {
	return listElementToPair(p.element.Next())
}

// Prev returns a pointer to the previous pair.
func (p *Pair[K, V]) Prev() *Pair[K, V] {
	return listElementToPair(p.element.Prev())
}

func listElementToPair[K comparable, V any](element *list.Element[*Pair[K, V]]) *Pair[K, V] {
	if element == nil {
		return nil
	}
	return element.Value
}

// KeyNotFoundError may be returned by functions in this package when they're called with keys that are not present
// in the map.
type KeyNotFoundError[K comparable] struct {
	MissingKey K
}

func (e *KeyNotFoundError[K]) Error() string {
	return fmt.Sprintf("missing key: %v", e.MissingKey)
}

// MoveAfter moves the value associated with key to its new position after the one associated with markKey.
// Returns an error iff key or markKey are not present in the map. If an error is returned,
// it will be a KeyNotFoundError.
func (om *OrderedMap[K, V]) MoveAfter(key, markKey K) error {
	elements, err := om.getElements(key, markKey)
	if err != nil {
		return err
	}
	om.list.MoveAfter(elements[0], elements[1])
	return nil
}

// MoveBefore moves the value associated with key to its new position before the one associated with markKey.
// Returns an error iff key or markKey are not present in the map. If an error is returned,
// it will be a KeyNotFoundError.
func (om *OrderedMap[K, V]) MoveBefore(key, markKey K) error {
	elements, err := om.getElements(key, markKey)
	if err != nil {
		return err
	}
	om.list.MoveBefore(elements[0], elements[1])
	return nil
}

func (om *OrderedMap[K, V]) getElements(keys ...K) ([]*list.Element[*Pair[K, V]], error) {
	elements := make([]*list.Element[*Pair[K, V]], len(keys))
	for i, k := range keys {
		pair, present := om.pairs[k]
		if !present {
			return nil, &KeyNotFoundError[K]{k}
		}
		elements[i] = pair.element
	}
	return elements, nil
}

// MoveToBack moves the value associated with key to the back of the ordered map,
// i.e. makes it the newest pair in the map.
// Returns an error iff key is not present in the map. If an error is returned,
// it will be a KeyNotFoundError.
func (om *OrderedMap[K, V]) MoveToBack(key K) error {
	_, err := om.GetAndMoveToBack(key)
	return err
}

// MoveToFront moves the value associated with key to the front of the ordered map,
// i.e. makes it the oldest pair in the map.
// Returns an error iff key is not present in the map. If an error is returned,
// it will be a KeyNotFoundError.
func (om *OrderedMap[K, V]) MoveToFront(key K) error {
	_, err := om.GetAndMoveToFront(key)
	return err
}

// GetAndMoveToBack combines Get and MoveToBack in the same call. If an error is returned,
// it will be a KeyNotFoundError.
func (om *OrderedMap[K, V]) GetAndMoveToBack(key K) (val V, err error) {
	if pair, present := om.pairs[key]; present {
		val = pair.Value
		om.list.MoveToBack(pair.element)
	} else {
		err = &KeyNotFoundError[K]{key}
	}

	return
}

// GetAndMoveToFront combines Get and MoveToFront in the same call. If an error is returned,
// it will be a KeyNotFoundError.
func (om *OrderedMap[K, V]) GetAndMoveToFront(key K) (val V, err error) {
	if pair, present := om.pairs[key]; present {
		val = pair.Value
		om.list.MoveToFront(pair.element)
	} else {
		err = &KeyNotFoundError[K]{key}
	}

	return
}
