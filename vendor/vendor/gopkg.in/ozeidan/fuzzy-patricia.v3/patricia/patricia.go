// Copyright (c) 2014 The go-patricia AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package patricia

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

//------------------------------------------------------------------------------
// Trie
//------------------------------------------------------------------------------

const (
	defaultMaxPrefixPerNode = 10
)

var (
	maxPrefixPerNode = defaultMaxPrefixPerNode
)

type (
	// Prefix is the type of node prefixes
	Prefix []byte
	// Item is just interface{}
	Item interface{}
	// VisitorFunc is the type of functions passed to visit function
	VisitorFunc func(prefix Prefix, item Item) error
	// FuzzyVisitorFunc additionaly returns how many characters were skipped which can be sorted on
	FuzzyVisitorFunc func(prefix Prefix, item Item, skipped int) error
)

// Trie is a generic patricia trie that allows fast retrieval of items by prefix.
// and other funky stuff.
//
// Trie is not thread-safe.
type Trie struct {
	prefix Prefix
	item   Item
	mask   uint64

	children childList
}

// Public API ------------------------------------------------------------------

// NewTrie constructs a new trie.
func NewTrie() *Trie {
	trie := &Trie{}

	trie.children = newSuperDenseChildList()

	trie.mask = 0
	return trie
}

// SetMaxPrefixPerNode sets the maximum length of a prefix before it is split into two nodes
func SetMaxPrefixPerNode(value int) {
	maxPrefixPerNode = value
}

// Clone makes a copy of an existing trie.
// Items stored in both tries become shared, obviously.
func (trie *Trie) Clone() *Trie {
	return &Trie{
		prefix:   append(Prefix(nil), trie.prefix...),
		item:     trie.item,
		children: trie.children.clone(),
	}
}

// Item returns the item stored in the root of this trie.
func (trie *Trie) Item() Item {
	return trie.item
}

// Insert inserts a new item into the trie using the given prefix. Insert does
// not replace existing items. It returns false if an item was already in place.
func (trie *Trie) Insert(key Prefix, item Item) (inserted bool) {
	return trie.put(key, item, false)
}

// Set works much like Insert, but it always sets the item, possibly replacing
// the item previously inserted.
func (trie *Trie) Set(key Prefix, item Item) {
	trie.put(key, item, true)
}

// Get returns the item located at key.
//
// This method is a bit dangerous, because Get can as well end up in an internal
// node that is not really representing any user-defined value. So when nil is
// a valid value being used, it is not possible to tell if the value was inserted
// into the tree by the user or not. A possible workaround for this is not to use
// nil interface as a valid value, even using zero value of any type is enough
// to prevent this bad behaviour.
func (trie *Trie) Get(key Prefix) (item Item) {
	_, node, found, leftover := trie.findSubtree(key)
	if !found || len(leftover) != 0 {
		return nil
	}
	return node.item
}

// Match returns what Get(prefix) != nil would return. The same warning as for
// Get applies here as well.
func (trie *Trie) Match(prefix Prefix) (matchedExactly bool) {
	return trie.Get(prefix) != nil
}

// MatchSubtree returns true when there is a subtree representing extensions
// to key, that is if there are any keys in the tree which have key as prefix.
func (trie *Trie) MatchSubtree(key Prefix) (matched bool) {
	_, _, matched, _ = trie.findSubtree(key)
	return
}

// Visit calls visitor on every node containing a non-nil item
// in alphabetical order.
//
// If an error is returned from visitor, the function stops visiting the tree
// and returns that error, unless it is a special error - SkipSubtree. In that
// case Visit skips the subtree represented by the current node and continues
// elsewhere.
func (trie *Trie) Visit(visitor VisitorFunc) error {
	return trie.walk(nil, visitor)
}

func (trie *Trie) size() int {
	n := 0

	err := trie.walk(nil, func(prefix Prefix, item Item) error {
		n++
		return nil
	})

	if err != nil {
		panic(err)
	}

	return n
}

func (trie *Trie) total() int {
	return 1 + trie.children.total()
}

// VisitSubtree works much like Visit, but it only visits nodes matching prefix.
func (trie *Trie) VisitSubtree(prefix Prefix, visitor VisitorFunc) error {
	// Nil prefix not allowed.
	if prefix == nil {
		panic(ErrNilPrefix)
	}

	// Empty trie must be handled explicitly.
	if trie.prefix == nil {
		return nil
	}

	// Locate the relevant subtree.
	_, root, found, leftover := trie.findSubtree(prefix)
	if !found {
		return nil
	}
	prefix = append(prefix, leftover...)

	// Visit it.
	return root.walk(prefix, visitor)
}

type potentialSubtree struct {
	idx     int
	skipped int
	prefix  Prefix
	node    *Trie
}

// VisitFuzzy visits every node that is succesfully matched via fuzzy matching
func (trie *Trie) VisitFuzzy(partial Prefix, caseInsensitive bool, visitor FuzzyVisitorFunc) error {
	if len(partial) == 0 {
		return trie.VisitPrefixes(partial, caseInsensitive, func(prefix Prefix, item Item) error {
			return visitor(prefix, item, 0)
		})
	}

	var (
		m   uint64
		cmp uint64
		i   int
		p   potentialSubtree
	)

	potential := []potentialSubtree{potentialSubtree{node: trie, prefix: Prefix(""), idx: 0}}
	for l := len(potential); l > 0; l = len(potential) {
		i = l - 1
		p = potential[i]

		potential = potential[:i]
		m = makePrefixMask(partial[p.idx:])

		if caseInsensitive {
			cmp = caseInsensitiveMask(p.node.mask)
		} else {
			cmp = p.node.mask
		}

		if (cmp & m) != m {
			continue
		}

		matchCount, skipped := fuzzyMatchCount(p.node.prefix,
			partial[p.idx:], p.idx, caseInsensitive)
		p.idx += matchCount
		if p.idx != 0 {
			p.skipped += skipped
		}

		if p.idx == len(partial) {
			fullPrefix := append(p.prefix, p.node.prefix...)

			err := p.node.walk(Prefix(""), func(prefix Prefix, item Item) error {
				key := make([]byte, len(fullPrefix), len(fullPrefix)+len(prefix))
				copy(key, fullPrefix)
				key = append(key, prefix...)

				err := visitor(key, item, p.skipped)
				if err != nil {
					return err
				}

				return nil
			})
			if err != nil {
				return err
			}

			continue
		}

		for _, c := range p.node.children.getChildren() {
			if c != nil {
				newPrefix := make(Prefix, len(p.prefix), len(p.prefix)+len(p.node.prefix))
				copy(newPrefix, p.prefix)
				newPrefix = append(newPrefix, p.node.prefix...)
				potential = append(potential, potentialSubtree{
					node:    c,
					prefix:  newPrefix,
					idx:     p.idx,
					skipped: p.skipped,
				})
			} else {
				fmt.Println("warning, child isn il")
			}
		}
	}

	return nil
}

func fuzzyMatchCount(prefix, query Prefix, idx int, caseInsensitive bool) (count, skipped int) {
	for i := 0; i < len(prefix); i++ {
		var match bool

		if caseInsensitive {
			match = matchCaseInsensitive(prefix[i], query[count])
		} else {
			match = prefix[i] == query[count]
		}

		if !match {
			if count+idx > 0 {
				skipped++
			}
			continue
		}

		count++
		if count >= len(query) {
			return
		}
	}
	return
}

// VisitSubstring takes a substring and visits all the nodes that whos prefix contains this substring
func (trie *Trie) VisitSubstring(substring Prefix, caseInsensitive bool, visitor VisitorFunc) error {
	if len(substring) == 0 {
		return trie.VisitSubtree(substring, visitor)
	}

	var (
		m            uint64
		cmp          uint64
		i            int
		p            potentialSubtree
		suffixLen    int
		maxSuffixLen = len(substring) - 1
	)

	potential := []potentialSubtree{potentialSubtree{node: trie, prefix: nil}}
	for l := len(potential); l > 0; l = len(potential) {
		i = l - 1
		p = potential[i]

		potential = potential[:i]

		if len(p.prefix) < maxSuffixLen {
			suffixLen = len(p.prefix)
		} else {
			suffixLen = maxSuffixLen
		}

		searchBytes := append(p.prefix[len(p.prefix)-suffixLen:], p.node.prefix...)

		contains := false

		if caseInsensitive {
			contains = bytes.Contains(bytes.ToUpper(searchBytes), bytes.ToUpper(substring))
		} else {
			contains = bytes.Contains(searchBytes, substring)
		}

		if contains {
			fullPrefix := append(p.prefix, p.node.prefix...)
			err := p.node.walk(Prefix(""), func(prefix Prefix, item Item) error {
				key := make([]byte, len(fullPrefix), len(fullPrefix)+len(prefix))
				copy(key, fullPrefix)
				key = append(key, prefix...)
				copy(key, append(fullPrefix, prefix...))

				err := visitor(key, item)
				if err != nil {
					return err
				}

				return nil
			})
			if err != nil {
				return err
			}
		}

		newPrefix := make(Prefix, len(p.prefix), len(p.prefix)+len(p.node.prefix))
		copy(newPrefix, p.prefix)
		newPrefix = append(newPrefix, p.node.prefix...)

		overLap := overlapLength(newPrefix, substring, caseInsensitive)
		m = makePrefixMask(substring[overLap:])

		for _, c := range p.node.children.getChildren() {
			if caseInsensitive {
				cmp = caseInsensitiveMask(c.mask)
			} else {
				cmp = c.mask
			}
			if c != nil && (cmp&m == m) {
				potential = append(potential, potentialSubtree{
					node:   c,
					prefix: newPrefix,
				})
			}
		}
	}

	return nil
}

func overlapLength(prefix, query Prefix, caseInsensitive bool) int {
	startLength := len(query) - 1
	if len(prefix) < startLength {
		startLength = len(prefix)
	}
	for i := startLength; i > 0; i-- {
		suffix := prefix[len(prefix)-i:]
		queryPrefix := query[:i]
		if caseInsensitive {
			if bytes.EqualFold(suffix, queryPrefix) {
				return i
			}
		} else if bytes.Equal(suffix, queryPrefix) {
			return i
		}
	}

	return 0
}

// VisitPrefixes visits only nodes that represent prefixes of key.
// To say the obvious, returning SkipSubtree from visitor makes no sense here.
func (trie *Trie) VisitPrefixes(key Prefix, caseInsensitive bool, visitor VisitorFunc) error {
	// Nil key not allowed.
	if key == nil {
		panic(ErrNilPrefix)
	}

	// Empty trie must be handled explicitly.
	if trie.prefix == nil {
		return nil
	}

	// Walk the path matching key prefixes.
	node := trie
	prefix := key
	offset := 0
	for {
		// Compute what part of prefix matches.
		common := node.longestCommonPrefixLength(key, caseInsensitive)
		key = key[common:]
		offset += common

		// Partial match means that there is no subtree matching prefix.
		if common < len(node.prefix) {
			return nil
		}

		// Call the visitor.
		if item := node.item; item != nil {
			if err := visitor(prefix[:offset], item); err != nil {
				return err
			}
		}

		if len(key) == 0 {
			// This node represents key, we are finished.
			return nil
		}

		// There is some key suffix left, move to the children.
		child := node.children.next(key[0])
		if child == nil {
			// There is nowhere to continue, return.
			return nil
		}

		node = child
	}
}

// Delete deletes the item represented by the given prefix.
//
// True is returned if the matching node was found and deleted.
func (trie *Trie) Delete(key Prefix) (deleted bool) {
	// Nil prefix not allowed.
	if key == nil {
		panic(ErrNilPrefix)
	}

	// Empty trie must be handled explicitly.
	if trie.prefix == nil {
		return false
	}

	// Find the relevant node.
	path, found, _ := trie.findSubtreePath(key)
	if !found {
		return false
	}

	node := path[len(path)-1]
	var parent *Trie
	if len(path) != 1 {
		parent = path[len(path)-2]
	}

	// If the item is already set to nil, there is nothing to do.
	if node.item == nil {
		return false
	}

	// Delete the item.
	node.item = nil

	// Initialise i before goto.
	// Will be used later in a loop.
	i := len(path) - 1

	// In case there are some child nodes, we cannot drop the whole subtree.
	// We can try to compact nodes, though.
	if node.children.length() != 0 {
		goto Compact
	}

	// In case we are at the root, just reset it and we are done.
	if parent == nil {
		node.reset()
		return true
	}

	// We can drop a subtree.
	// Find the first ancestor that has its value set or it has 2 or more child nodes.
	// That will be the node where to drop the subtree at.
	for ; i >= 0; i-- {
		if current := path[i]; current.item != nil || current.children.length() >= 2 {
			break
		}
	}

	// Handle the case when there is no such node.
	// In other words, we can reset the whole tree.
	if i == -1 {
		path[0].reset()
		return true
	}

	// We can just remove the subtree here.
	node = path[i]
	if i == 0 {
		parent = nil
	} else {
		parent = path[i-1]
	}
	// i+1 is always a valid index since i is never pointing to the last node.
	// The loop above skips at least the last node since we are sure that the item
	// is set to nil and it has no children, othewise we would be compacting instead.
	node.children.remove(path[i+1].prefix[0])

	// lastly, the bitmasks of all of the parent nodes have to be updated again, since
	// a child node of all of them has bin removed
	for ; i >= 0; i-- {
		n := path[i]
		n.mask = n.children.combinedMask()
	}

Compact:
	// The node is set to the first non-empty ancestor,
	// so try to compact since that might be possible now.
	if compacted := node.compact(); compacted != node {
		if parent == nil {
			*node = *compacted
		} else {
			parent.children.replace(node.prefix[0], compacted)
			*parent = *parent.compact()
		}
	}

	return true
}

// DeleteSubtree finds the subtree exactly matching prefix and deletes it.
//
// True is returned if the subtree was found and deleted.
func (trie *Trie) DeleteSubtree(prefix Prefix) (deleted bool) {
	// Nil prefix not allowed.
	if prefix == nil {
		panic(ErrNilPrefix)
	}

	// Empty trie must be handled explicitly.
	if trie.prefix == nil {
		return false
	}

	// Locate the relevant subtree.
	parent, root, found, _ := trie.findSubtree(prefix)
	path, _, _ := trie.findSubtreePath(prefix)
	if !found {
		return false
	}

	// If we are in the root of the trie, reset the trie.
	if parent == nil {
		root.reset()
		return true
	}

	// Otherwise remove the root node from its parent.
	parent.children.remove(root.prefix[0])

	// update masks
	parent.mask = parent.children.combinedMask()
	for i := len(path) - 1; i >= 0; i-- {
		n := path[i]
		n.mask = n.children.combinedMask()
	}

	return true
}

// Internal helper methods -----------------------------------------------------

func (trie *Trie) empty() bool {
	return trie.item == nil && trie.children.length() == 0
}

func (trie *Trie) reset() {
	trie.prefix = nil
	trie.children = newSuperDenseChildList()
}

func makePrefixMask(key Prefix) uint64 {
	var mask uint64
	for _, b := range key {
		if b >= '0' && b <= '9' {
			// 0-9 bits: 0-9
			b -= 48
		} else if b >= 'A' && b <= 'Z' {
			// A-Z bits: 10-35
			b -= 55
		} else if b >= 'a' && b <= 'z' {
			// a-z bits: 36-61
			b -= 61
		} else if b == '.' {
			b = 62
		} else if b == '-' {
			b = 63
		} else {
			continue
		}
		mask |= uint64(1) << uint64(b)
	}
	return mask
}

const upperBits = 0xFFFFFFC00
const lowerBits = 0x3FFFFFF000000000

func caseInsensitiveMask(mask uint64) uint64 {
	mask |= (mask & upperBits) << uint64(26)
	mask |= (mask & lowerBits) >> uint64(26)
	return mask
}

var charmap = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz.-"

func (trie *Trie) put(key Prefix, item Item, replace bool) (inserted bool) {
	// Nil prefix not allowed.
	if key == nil {
		panic(ErrNilPrefix)
	}

	var (
		common int
		node   = trie
		child  *Trie
		mask   uint64
	)

	mask = makePrefixMask(key)

	if node.prefix == nil {
		node.mask |= mask
		if len(key) <= maxPrefixPerNode {
			node.prefix = key
			goto InsertItem
		}
		node.prefix = key[:maxPrefixPerNode]
		key = key[maxPrefixPerNode:]
		mask = makePrefixMask(key)
		goto AppendChild
	}

	for {
		// Compute the longest common prefix length.
		common = node.longestCommonPrefixLength(key, false)
		key = key[common:]

		// Only a part matches, split.
		if common < len(node.prefix) {
			goto SplitPrefix
		}

		// common == len(node.prefix) since never (common > len(node.prefix))
		// common == len(former key) <-> 0 == len(key)
		// -> former key == node.prefix
		if len(key) == 0 {
			goto InsertItem
		}

		node.mask |= mask
		// Check children for matching prefix.
		child = node.children.next(key[0])
		if child == nil {
			goto AppendChild
		}
		node = child
	}

SplitPrefix:
	// Split the prefix if necessary.
	child = new(Trie)
	*child = *node
	*node = *NewTrie()
	node.prefix = child.prefix[:common]
	child.prefix = child.prefix[common:]
	child = child.compact()
	node.children = node.children.add(child)
	node.mask = child.mask
	node.mask |= mask
	mask = makePrefixMask(key)

AppendChild:
	// Keep appending children until whole prefix is inserted.
	// This loop starts with empty node.prefix that needs to be filled.
	for len(key) != 0 {
		child := NewTrie()
		child.mask = mask
		if len(key) <= maxPrefixPerNode {
			child.prefix = key
			node.children = node.children.add(child)
			node = child
			goto InsertItem
		} else {
			child.prefix = key[:maxPrefixPerNode]
			key = key[maxPrefixPerNode:]
			mask = makePrefixMask(key)
			node.children = node.children.add(child)
			node = child
		}
	}

InsertItem:
	// Try to insert the item if possible.
	if replace || node.item == nil {
		node.item = item
		return true
	}
	return false
}

func (trie *Trie) compact() *Trie {
	// Only a node with a single child can be compacted.
	if trie.children.length() != 1 {
		return trie
	}

	child := trie.children.head()

	// If any item is set, we cannot compact since we want to retain
	// the ability to do searching by key. This makes compaction less usable,
	// but that simply cannot be avoided.
	if trie.item != nil || child.item != nil {
		return trie
	}

	// Make sure the combined prefixes fit into a single node.
	if len(trie.prefix)+len(child.prefix) > maxPrefixPerNode {
		return trie
	}

	// Concatenate the prefixes, move the items.
	child.prefix = append(trie.prefix, child.prefix...)
	child.mask = trie.mask
	if trie.item != nil {
		child.item = trie.item
	}

	return child
}

func (trie *Trie) findSubtree(prefix Prefix) (parent *Trie, root *Trie, found bool, leftover Prefix) {
	// Find the subtree matching prefix.
	root = trie
	for {
		// Compute what part of prefix matches.
		common := root.longestCommonPrefixLength(prefix, false)
		prefix = prefix[common:]

		// We used up the whole prefix, subtree found.
		if len(prefix) == 0 {
			found = true
			leftover = root.prefix[common:]
			return
		}

		// Partial match means that there is no subtree matching prefix.
		if common < len(root.prefix) {
			leftover = root.prefix[common:]
			return
		}

		// There is some prefix left, move to the children.
		child := root.children.next(prefix[0])
		if child == nil {
			// There is nowhere to continue, there is no subtree matching prefix.
			return
		}

		parent = root
		root = child
	}
}

func (trie *Trie) findSubtreePath(prefix Prefix) (path []*Trie, found bool, leftover Prefix) {
	// Find the subtree matching prefix.
	root := trie
	var subtreePath []*Trie
	for {
		// Append the current root to the path.
		subtreePath = append(subtreePath, root)

		// Compute what part of prefix matches.
		common := root.longestCommonPrefixLength(prefix, false)
		prefix = prefix[common:]

		// We used up the whole prefix, subtree found.
		if len(prefix) == 0 {
			path = subtreePath
			found = true
			leftover = root.prefix[common:]
			return
		}

		// Partial match means that there is no subtree matching prefix.
		if common < len(root.prefix) {
			leftover = root.prefix[common:]
			return
		}

		// There is some prefix left, move to the children.
		child := root.children.next(prefix[0])
		if child == nil {
			// There is nowhere to continue, there is no subtree matching prefix.
			return
		}

		root = child
	}
}

func (trie *Trie) walk(actualRootPrefix Prefix, visitor VisitorFunc) error {
	var prefix Prefix
	// Allocate a bit more space for prefix at the beginning.
	if actualRootPrefix == nil {
		prefix = make(Prefix, 32+len(trie.prefix))
		copy(prefix, trie.prefix)
		prefix = prefix[:len(trie.prefix)]
	} else {
		prefix = make(Prefix, 32+len(actualRootPrefix))
		copy(prefix, actualRootPrefix)
		prefix = prefix[:len(actualRootPrefix)]
	}

	// Visit the root first. Not that this works for empty trie as well since
	// in that case item == nil && len(children) == 0.
	if trie.item != nil {
		if err := visitor(prefix, trie.item); err != nil {
			if err == SkipSubtree {
				return nil
			}
			return err
		}
	}

	// Then continue to the children.
	return trie.children.walk(&prefix, visitor)
}

func (trie *Trie) longestCommonPrefixLength(prefix Prefix, caseInsensitive bool) (i int) {
	for ; i < len(prefix) && i < len(trie.prefix); i++ {
		p := prefix[i]
		t := trie.prefix[i]

		if caseInsensitive {
			if !(matchCaseInsensitive(t, p)) {
				break
			}
		} else {
			if p != t {
				break
			}
		}
	}
	return
}

func matchCaseInsensitive(a byte, b byte) bool {
	return a == b+32 || b == a+32 || a == b
}

func (trie *Trie) dump() string {
	writer := &bytes.Buffer{}
	trie.print(writer, 0)
	return writer.String()
}

func (trie *Trie) print(writer io.Writer, indent int) {
	fmt.Fprintf(writer, "%s%s %v\n", strings.Repeat(" ", indent), string(trie.prefix), trie.item)
	trie.children.print(writer, indent+2)
}

// Errors ----------------------------------------------------------------------

var (
	SkipSubtree  = errors.New("Skip this subtree")
	ErrNilPrefix = errors.New("Nil prefix passed into a method call")
)
