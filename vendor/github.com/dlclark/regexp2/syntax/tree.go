package syntax

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
)

type RegexTree struct {
	root       *regexNode
	caps       map[int]int
	capnumlist []int
	captop     int
	Capnames   map[string]int
	Caplist    []string
	options    RegexOptions
}

// It is built into a parsed tree for a regular expression.

// Implementation notes:
//
// Since the node tree is a temporary data structure only used
// during compilation of the regexp to integer codes, it's
// designed for clarity and convenience rather than
// space efficiency.
//
// RegexNodes are built into a tree, linked by the n.children list.
// Each node also has a n.parent and n.ichild member indicating
// its parent and which child # it is in its parent's list.
//
// RegexNodes come in as many types as there are constructs in
// a regular expression, for example, "concatenate", "alternate",
// "one", "rept", "group". There are also node types for basic
// peephole optimizations, e.g., "onerep", "notsetrep", etc.
//
// Because perl 5 allows "lookback" groups that scan backwards,
// each node also gets a "direction". Normally the value of
// boolean n.backward = false.
//
// During parsing, top-level nodes are also stacked onto a parse
// stack (a stack of trees). For this purpose we have a n.next
// pointer. [Note that to save a few bytes, we could overload the
// n.parent pointer instead.]
//
// On the parse stack, each tree has a "role" - basically, the
// nonterminal in the grammar that the parser has currently
// assigned to the tree. That code is stored in n.role.
//
// Finally, some of the different kinds of nodes have data.
// Two integers (for the looping constructs) are stored in
// n.operands, an an object (either a string or a set)
// is stored in n.data
type regexNode struct {
	t        nodeType
	children []*regexNode
	str      []rune
	set      *CharSet
	ch       rune
	m        int
	n        int
	options  RegexOptions
	next     *regexNode
}

type nodeType int32

const (
	// The following are leaves, and correspond to primitive operations

	ntOnerep      nodeType = 0  // lef,back char,min,max    a {n}
	ntNotonerep            = 1  // lef,back char,min,max    .{n}
	ntSetrep               = 2  // lef,back set,min,max     [\d]{n}
	ntOneloop              = 3  // lef,back char,min,max    a {,n}
	ntNotoneloop           = 4  // lef,back char,min,max    .{,n}
	ntSetloop              = 5  // lef,back set,min,max     [\d]{,n}
	ntOnelazy              = 6  // lef,back char,min,max    a {,n}?
	ntNotonelazy           = 7  // lef,back char,min,max    .{,n}?
	ntSetlazy              = 8  // lef,back set,min,max     [\d]{,n}?
	ntOne                  = 9  // lef      char            a
	ntNotone               = 10 // lef      char            [^a]
	ntSet                  = 11 // lef      set             [a-z\s]  \w \s \d
	ntMulti                = 12 // lef      string          abcd
	ntRef                  = 13 // lef      group           \#
	ntBol                  = 14 //                          ^
	ntEol                  = 15 //                          $
	ntBoundary             = 16 //                          \b
	ntNonboundary          = 17 //                          \B
	ntBeginning            = 18 //                          \A
	ntStart                = 19 //                          \G
	ntEndZ                 = 20 //                          \Z
	ntEnd                  = 21 //                          \Z

	// Interior nodes do not correspond to primitive operations, but
	// control structures compositing other operations

	// Concat and alternate take n children, and can run forward or backwards

	ntNothing     = 22 //          []
	ntEmpty       = 23 //          ()
	ntAlternate   = 24 //          a|b
	ntConcatenate = 25 //          ab
	ntLoop        = 26 // m,x      * + ? {,}
	ntLazyloop    = 27 // m,x      *? +? ?? {,}?
	ntCapture     = 28 // n        ()
	ntGroup       = 29 //          (?:)
	ntRequire     = 30 //          (?=) (?<=)
	ntPrevent     = 31 //          (?!) (?<!)
	ntGreedy      = 32 //          (?>) (?<)
	ntTestref     = 33 //          (?(n) | )
	ntTestgroup   = 34 //          (?(...) | )

	ntECMABoundary    = 41 //                          \b
	ntNonECMABoundary = 42 //                          \B
)

func newRegexNode(t nodeType, opt RegexOptions) *regexNode {
	return &regexNode{
		t:       t,
		options: opt,
	}
}

func newRegexNodeCh(t nodeType, opt RegexOptions, ch rune) *regexNode {
	return &regexNode{
		t:       t,
		options: opt,
		ch:      ch,
	}
}

func newRegexNodeStr(t nodeType, opt RegexOptions, str []rune) *regexNode {
	return &regexNode{
		t:       t,
		options: opt,
		str:     str,
	}
}

func newRegexNodeSet(t nodeType, opt RegexOptions, set *CharSet) *regexNode {
	return &regexNode{
		t:       t,
		options: opt,
		set:     set,
	}
}

func newRegexNodeM(t nodeType, opt RegexOptions, m int) *regexNode {
	return &regexNode{
		t:       t,
		options: opt,
		m:       m,
	}
}
func newRegexNodeMN(t nodeType, opt RegexOptions, m, n int) *regexNode {
	return &regexNode{
		t:       t,
		options: opt,
		m:       m,
		n:       n,
	}
}

func (n *regexNode) writeStrToBuf(buf *bytes.Buffer) {
	for i := 0; i < len(n.str); i++ {
		buf.WriteRune(n.str[i])
	}
}

func (n *regexNode) addChild(child *regexNode) {
	reduced := child.reduce()
	n.children = append(n.children, reduced)
	reduced.next = n
}

func (n *regexNode) insertChildren(afterIndex int, nodes []*regexNode) {
	newChildren := make([]*regexNode, 0, len(n.children)+len(nodes))
	n.children = append(append(append(newChildren, n.children[:afterIndex]...), nodes...), n.children[afterIndex:]...)
}

// removes children including the start but not the end index
func (n *regexNode) removeChildren(startIndex, endIndex int) {
	n.children = append(n.children[:startIndex], n.children[endIndex:]...)
}

// Pass type as OneLazy or OneLoop
func (n *regexNode) makeRep(t nodeType, min, max int) {
	n.t += (t - ntOne)
	n.m = min
	n.n = max
}

func (n *regexNode) reduce() *regexNode {
	switch n.t {
	case ntAlternate:
		return n.reduceAlternation()

	case ntConcatenate:
		return n.reduceConcatenation()

	case ntLoop, ntLazyloop:
		return n.reduceRep()

	case ntGroup:
		return n.reduceGroup()

	case ntSet, ntSetloop:
		return n.reduceSet()

	default:
		return n
	}
}

// Basic optimization. Single-letter alternations can be replaced
// by faster set specifications, and nested alternations with no
// intervening operators can be flattened:
//
// a|b|c|def|g|h -> [a-c]|def|[gh]
// apple|(?:orange|pear)|grape -> apple|orange|pear|grape
func (n *regexNode) reduceAlternation() *regexNode {
	if len(n.children) == 0 {
		return newRegexNode(ntNothing, n.options)
	}

	wasLastSet := false
	lastNodeCannotMerge := false
	var optionsLast RegexOptions
	var i, j int

	for i, j = 0, 0; i < len(n.children); i, j = i+1, j+1 {
		at := n.children[i]

		if j < i {
			n.children[j] = at
		}

		for {
			if at.t == ntAlternate {
				for k := 0; k < len(at.children); k++ {
					at.children[k].next = n
				}
				n.insertChildren(i+1, at.children)

				j--
			} else if at.t == ntSet || at.t == ntOne {
				// Cannot merge sets if L or I options differ, or if either are negated.
				optionsAt := at.options & (RightToLeft | IgnoreCase)

				if at.t == ntSet {
					if !wasLastSet || optionsLast != optionsAt || lastNodeCannotMerge || !at.set.IsMergeable() {
						wasLastSet = true
						lastNodeCannotMerge = !at.set.IsMergeable()
						optionsLast = optionsAt
						break
					}
				} else if !wasLastSet || optionsLast != optionsAt || lastNodeCannotMerge {
					wasLastSet = true
					lastNodeCannotMerge = false
					optionsLast = optionsAt
					break
				}

				// The last node was a Set or a One, we're a Set or One and our options are the same.
				// Merge the two nodes.
				j--
				prev := n.children[j]

				var prevCharClass *CharSet
				if prev.t == ntOne {
					prevCharClass = &CharSet{}
					prevCharClass.addChar(prev.ch)
				} else {
					prevCharClass = prev.set
				}

				if at.t == ntOne {
					prevCharClass.addChar(at.ch)
				} else {
					prevCharClass.addSet(*at.set)
				}

				prev.t = ntSet
				prev.set = prevCharClass
			} else if at.t == ntNothing {
				j--
			} else {
				wasLastSet = false
				lastNodeCannotMerge = false
			}
			break
		}
	}

	if j < i {
		n.removeChildren(j, i)
	}

	return n.stripEnation(ntNothing)
}

// Basic optimization. Adjacent strings can be concatenated.
//
// (?:abc)(?:def) -> abcdef
func (n *regexNode) reduceConcatenation() *regexNode {
	// Eliminate empties and concat adjacent strings/chars

	var optionsLast RegexOptions
	var optionsAt RegexOptions
	var i, j int

	if len(n.children) == 0 {
		return newRegexNode(ntEmpty, n.options)
	}

	wasLastString := false

	for i, j = 0, 0; i < len(n.children); i, j = i+1, j+1 {
		var at, prev *regexNode

		at = n.children[i]

		if j < i {
			n.children[j] = at
		}

		if at.t == ntConcatenate &&
			((at.options & RightToLeft) == (n.options & RightToLeft)) {
			for k := 0; k < len(at.children); k++ {
				at.children[k].next = n
			}

			//insert at.children at i+1 index in n.children
			n.insertChildren(i+1, at.children)

			j--
		} else if at.t == ntMulti || at.t == ntOne {
			// Cannot merge strings if L or I options differ
			optionsAt = at.options & (RightToLeft | IgnoreCase)

			if !wasLastString || optionsLast != optionsAt {
				wasLastString = true
				optionsLast = optionsAt
				continue
			}

			j--
			prev = n.children[j]

			if prev.t == ntOne {
				prev.t = ntMulti
				prev.str = []rune{prev.ch}
			}

			if (optionsAt & RightToLeft) == 0 {
				if at.t == ntOne {
					prev.str = append(prev.str, at.ch)
				} else {
					prev.str = append(prev.str, at.str...)
				}
			} else {
				if at.t == ntOne {
					// insert at the front by expanding our slice, copying the data over, and then setting the value
					prev.str = append(prev.str, 0)
					copy(prev.str[1:], prev.str)
					prev.str[0] = at.ch
				} else {
					//insert at the front...this one we'll make a new slice and copy both into it
					merge := make([]rune, len(prev.str)+len(at.str))
					copy(merge, at.str)
					copy(merge[len(at.str):], prev.str)
					prev.str = merge
				}
			}
		} else if at.t == ntEmpty {
			j--
		} else {
			wasLastString = false
		}
	}

	if j < i {
		// remove indices j through i from the children
		n.removeChildren(j, i)
	}

	return n.stripEnation(ntEmpty)
}

// Nested repeaters just get multiplied with each other if they're not
// too lumpy
func (n *regexNode) reduceRep() *regexNode {

	u := n
	t := n.t
	min := n.m
	max := n.n

	for {
		if len(u.children) == 0 {
			break
		}

		child := u.children[0]

		// multiply reps of the same type only
		if child.t != t {
			childType := child.t

			if !(childType >= ntOneloop && childType <= ntSetloop && t == ntLoop ||
				childType >= ntOnelazy && childType <= ntSetlazy && t == ntLazyloop) {
				break
			}
		}

		// child can be too lumpy to blur, e.g., (a {100,105}) {3} or (a {2,})?
		// [but things like (a {2,})+ are not too lumpy...]
		if u.m == 0 && child.m > 1 || child.n < child.m*2 {
			break
		}

		u = child
		if u.m > 0 {
			if (math.MaxInt32-1)/u.m < min {
				u.m = math.MaxInt32
			} else {
				u.m = u.m * min
			}
		}
		if u.n > 0 {
			if (math.MaxInt32-1)/u.n < max {
				u.n = math.MaxInt32
			} else {
				u.n = u.n * max
			}
		}
	}

	if math.MaxInt32 == min {
		return newRegexNode(ntNothing, n.options)
	}
	return u

}

// Simple optimization. If a concatenation or alternation has only
// one child strip out the intermediate node. If it has zero children,
// turn it into an empty.
func (n *regexNode) stripEnation(emptyType nodeType) *regexNode {
	switch len(n.children) {
	case 0:
		return newRegexNode(emptyType, n.options)
	case 1:
		return n.children[0]
	default:
		return n
	}
}

func (n *regexNode) reduceGroup() *regexNode {
	u := n

	for u.t == ntGroup {
		u = u.children[0]
	}

	return u
}

// Simple optimization. If a set is a singleton, an inverse singleton,
// or empty, it's transformed accordingly.
func (n *regexNode) reduceSet() *regexNode {
	// Extract empty-set, one and not-one case as special

	if n.set == nil {
		n.t = ntNothing
	} else if n.set.IsSingleton() {
		n.ch = n.set.SingletonChar()
		n.set = nil
		n.t += (ntOne - ntSet)
	} else if n.set.IsSingletonInverse() {
		n.ch = n.set.SingletonChar()
		n.set = nil
		n.t += (ntNotone - ntSet)
	}

	return n
}

func (n *regexNode) reverseLeft() *regexNode {
	if n.options&RightToLeft != 0 && n.t == ntConcatenate && len(n.children) > 0 {
		//reverse children order
		for left, right := 0, len(n.children)-1; left < right; left, right = left+1, right-1 {
			n.children[left], n.children[right] = n.children[right], n.children[left]
		}
	}

	return n
}

func (n *regexNode) makeQuantifier(lazy bool, min, max int) *regexNode {
	if min == 0 && max == 0 {
		return newRegexNode(ntEmpty, n.options)
	}

	if min == 1 && max == 1 {
		return n
	}

	switch n.t {
	case ntOne, ntNotone, ntSet:
		if lazy {
			n.makeRep(Onelazy, min, max)
		} else {
			n.makeRep(Oneloop, min, max)
		}
		return n

	default:
		var t nodeType
		if lazy {
			t = ntLazyloop
		} else {
			t = ntLoop
		}
		result := newRegexNodeMN(t, n.options, min, max)
		result.addChild(n)
		return result
	}
}

// debug functions

var typeStr = []string{
	"Onerep", "Notonerep", "Setrep",
	"Oneloop", "Notoneloop", "Setloop",
	"Onelazy", "Notonelazy", "Setlazy",
	"One", "Notone", "Set",
	"Multi", "Ref",
	"Bol", "Eol", "Boundary", "Nonboundary",
	"Beginning", "Start", "EndZ", "End",
	"Nothing", "Empty",
	"Alternate", "Concatenate",
	"Loop", "Lazyloop",
	"Capture", "Group", "Require", "Prevent", "Greedy",
	"Testref", "Testgroup",
	"Unknown", "Unknown", "Unknown",
	"Unknown", "Unknown", "Unknown",
	"ECMABoundary", "NonECMABoundary",
}

func (n *regexNode) description() string {
	buf := &bytes.Buffer{}

	buf.WriteString(typeStr[n.t])

	if (n.options & ExplicitCapture) != 0 {
		buf.WriteString("-C")
	}
	if (n.options & IgnoreCase) != 0 {
		buf.WriteString("-I")
	}
	if (n.options & RightToLeft) != 0 {
		buf.WriteString("-L")
	}
	if (n.options & Multiline) != 0 {
		buf.WriteString("-M")
	}
	if (n.options & Singleline) != 0 {
		buf.WriteString("-S")
	}
	if (n.options & IgnorePatternWhitespace) != 0 {
		buf.WriteString("-X")
	}
	if (n.options & ECMAScript) != 0 {
		buf.WriteString("-E")
	}

	switch n.t {
	case ntOneloop, ntNotoneloop, ntOnelazy, ntNotonelazy, ntOne, ntNotone:
		buf.WriteString("(Ch = " + CharDescription(n.ch) + ")")
		break
	case ntCapture:
		buf.WriteString("(index = " + strconv.Itoa(n.m) + ", unindex = " + strconv.Itoa(n.n) + ")")
		break
	case ntRef, ntTestref:
		buf.WriteString("(index = " + strconv.Itoa(n.m) + ")")
		break
	case ntMulti:
		fmt.Fprintf(buf, "(String = %s)", string(n.str))
		break
	case ntSet, ntSetloop, ntSetlazy:
		buf.WriteString("(Set = " + n.set.String() + ")")
		break
	}

	switch n.t {
	case ntOneloop, ntNotoneloop, ntOnelazy, ntNotonelazy, ntSetloop, ntSetlazy, ntLoop, ntLazyloop:
		buf.WriteString("(Min = ")
		buf.WriteString(strconv.Itoa(n.m))
		buf.WriteString(", Max = ")
		if n.n == math.MaxInt32 {
			buf.WriteString("inf")
		} else {
			buf.WriteString(strconv.Itoa(n.n))
		}
		buf.WriteString(")")

		break
	}

	return buf.String()
}

var padSpace = []byte("                                ")

func (t *RegexTree) Dump() string {
	return t.root.dump()
}

func (n *regexNode) dump() string {
	var stack []int
	CurNode := n
	CurChild := 0

	buf := bytes.NewBufferString(CurNode.description())
	buf.WriteRune('\n')

	for {
		if CurNode.children != nil && CurChild < len(CurNode.children) {
			stack = append(stack, CurChild+1)
			CurNode = CurNode.children[CurChild]
			CurChild = 0

			Depth := len(stack)
			if Depth > 32 {
				Depth = 32
			}
			buf.Write(padSpace[:Depth])
			buf.WriteString(CurNode.description())
			buf.WriteRune('\n')
		} else {
			if len(stack) == 0 {
				break
			}

			CurChild = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			CurNode = CurNode.next
		}
	}
	return buf.String()
}
