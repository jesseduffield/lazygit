package litter

import (
	"fmt"
	"reflect"
	"sort"
)

// mapReusedPointers takes a structure, and recursively maps all pointers mentioned in the tree,
// detecting circular references, and providing a list of all pointers that was referenced at
// least twice by the provided structure.
func mapReusedPointers(v reflect.Value) ptrmap {
	pm := &pointerVisitor{}
	pm.consider(v)
	return pm.reused
}

// A map of pointers.
type ptrinfo struct {
	id     int
	parent *ptrmap
}

func (p *ptrinfo) label() string {
	if p.id == -1 {
		p.id = p.parent.count
		p.parent.count++
	}
	return fmt.Sprintf("p%d", p.id)
}

type ptrkey struct {
	p uintptr
	t reflect.Type
}

func ptrkeyFor(v reflect.Value) (k ptrkey) {
	k.p = v.Pointer()
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.IsValid() {
		k.t = v.Type()
	}
	return
}

type ptrmap struct {
	m     map[ptrkey]*ptrinfo
	count int
}

// Returns true if contains a pointer.
func (pm *ptrmap) contains(v reflect.Value) bool {
	if pm.m != nil {
		_, ok := pm.m[ptrkeyFor(v)]
		return ok
	}
	return false
}

// Gets a pointer.
func (pm *ptrmap) get(v reflect.Value) (*ptrinfo, bool) {
	if pm.m != nil {
		p, ok := pm.m[ptrkeyFor(v)]
		return p, ok
	}
	return nil, false
}

// Removes a pointer.
func (pm *ptrmap) remove(v reflect.Value) {
	if pm.m != nil {
		delete(pm.m, ptrkeyFor(v))
	}
}

// Adds a pointer.
func (pm *ptrmap) add(p reflect.Value) bool {
	if pm.contains(p) {
		return false
	}
	pm.put(p)
	return true
}

// Adds a pointer (slow path).
func (pm *ptrmap) put(v reflect.Value) {
	if pm.m == nil {
		pm.m = make(map[ptrkey]*ptrinfo, 31)
	}

	key := ptrkeyFor(v)
	if _, ok := pm.m[key]; !ok {
		pm.m[key] = &ptrinfo{id: -1, parent: pm}
	}
}

type pointerVisitor struct {
	pointers ptrmap
	reused   ptrmap
}

// Recursively consider v and each of its children, updating the map according to the
// semantics of MapReusedPointers
func (pv *pointerVisitor) consider(v reflect.Value) {
	if v.Kind() == reflect.Invalid {
		return
	}
	if isPointerValue(v) { // pointer is 0 for unexported fields
		if pv.tryAddPointer(v) {
			// No use descending inside this value, since it have been seen before and all its descendants
			// have been considered
			return
		}
	}

	// Now descend into any children of this value
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		numEntries := v.Len()
		for i := 0; i < numEntries; i++ {
			pv.consider(v.Index(i))
		}

	case reflect.Interface:
		pv.consider(v.Elem())

	case reflect.Ptr:
		pv.consider(v.Elem())

	case reflect.Map:
		keys := v.MapKeys()
		sort.Sort(mapKeySorter{
			keys:    keys,
			options: &Config,
		})
		for _, key := range keys {
			pv.consider(v.MapIndex(key))
		}

	case reflect.Struct:
		numFields := v.NumField()
		for i := 0; i < numFields; i++ {
			pv.consider(v.Field(i))
		}
	}
}

// addPointer to the pointerMap, update reusedPointers. Returns true if pointer was reused
func (pv *pointerVisitor) tryAddPointer(v reflect.Value) bool {
	// Is this allready known to be reused?
	if pv.reused.contains(v) {
		return true
	}

	// Have we seen it once before?
	if pv.pointers.contains(v) {
		// Add it to the register of pointers we have seen more than once
		pv.reused.add(v)
		return true
	}

	// This pointer was new to us
	pv.pointers.add(v)
	return false
}
