package syntax

import (
	"bytes"
	"fmt"
	"math"
	"os"
)

func Write(tree *RegexTree) (*Code, error) {
	w := writer{
		intStack:   make([]int, 0, 32),
		emitted:    make([]int, 2),
		stringhash: make(map[string]int),
		sethash:    make(map[string]int),
	}

	code, err := w.codeFromTree(tree)

	if tree.options&Debug > 0 && code != nil {
		os.Stdout.WriteString(code.Dump())
		os.Stdout.WriteString("\n")
	}

	return code, err
}

type writer struct {
	emitted []int

	intStack    []int
	curpos      int
	stringhash  map[string]int
	stringtable [][]rune
	sethash     map[string]int
	settable    []*CharSet
	counting    bool
	count       int
	trackcount  int
	caps        map[int]int
}

const (
	beforeChild nodeType = 64
	afterChild           = 128
	//MaxPrefixSize is the largest number of runes we'll use for a BoyerMoyer prefix
	MaxPrefixSize = 50
)

// The top level RegexCode generator. It does a depth-first walk
// through the tree and calls EmitFragment to emits code before
// and after each child of an interior node, and at each leaf.
//
// It runs two passes, first to count the size of the generated
// code, and second to generate the code.
//
// We should time it against the alternative, which is
// to just generate the code and grow the array as we go.
func (w *writer) codeFromTree(tree *RegexTree) (*Code, error) {
	var (
		curNode  *regexNode
		curChild int
		capsize  int
	)
	// construct sparse capnum mapping if some numbers are unused

	if tree.capnumlist == nil || tree.captop == len(tree.capnumlist) {
		capsize = tree.captop
		w.caps = nil
	} else {
		capsize = len(tree.capnumlist)
		w.caps = tree.caps
		for i := 0; i < len(tree.capnumlist); i++ {
			w.caps[tree.capnumlist[i]] = i
		}
	}

	w.counting = true

	for {
		if !w.counting {
			w.emitted = make([]int, w.count)
		}

		curNode = tree.root
		curChild = 0

		w.emit1(Lazybranch, 0)

		for {
			if len(curNode.children) == 0 {
				w.emitFragment(curNode.t, curNode, 0)
			} else if curChild < len(curNode.children) {
				w.emitFragment(curNode.t|beforeChild, curNode, curChild)

				curNode = curNode.children[curChild]

				w.pushInt(curChild)
				curChild = 0
				continue
			}

			if w.emptyStack() {
				break
			}

			curChild = w.popInt()
			curNode = curNode.next

			w.emitFragment(curNode.t|afterChild, curNode, curChild)
			curChild++
		}

		w.patchJump(0, w.curPos())
		w.emit(Stop)

		if !w.counting {
			break
		}

		w.counting = false
	}

	fcPrefix := getFirstCharsPrefix(tree)
	prefix := getPrefix(tree)
	rtl := (tree.options & RightToLeft) != 0

	var bmPrefix *BmPrefix
	//TODO: benchmark string prefixes
	if prefix != nil && len(prefix.PrefixStr) > 0 && MaxPrefixSize > 0 {
		if len(prefix.PrefixStr) > MaxPrefixSize {
			// limit prefix changes to 10k
			prefix.PrefixStr = prefix.PrefixStr[:MaxPrefixSize]
		}
		bmPrefix = newBmPrefix(prefix.PrefixStr, prefix.CaseInsensitive, rtl)
	} else {
		bmPrefix = nil
	}

	return &Code{
		Codes:       w.emitted,
		Strings:     w.stringtable,
		Sets:        w.settable,
		TrackCount:  w.trackcount,
		Caps:        w.caps,
		Capsize:     capsize,
		FcPrefix:    fcPrefix,
		BmPrefix:    bmPrefix,
		Anchors:     getAnchors(tree),
		RightToLeft: rtl,
	}, nil
}

// The main RegexCode generator. It does a depth-first walk
// through the tree and calls EmitFragment to emits code before
// and after each child of an interior node, and at each leaf.
func (w *writer) emitFragment(nodetype nodeType, node *regexNode, curIndex int) error {
	bits := InstOp(0)

	if nodetype <= ntRef {
		if (node.options & RightToLeft) != 0 {
			bits |= Rtl
		}
		if (node.options & IgnoreCase) != 0 {
			bits |= Ci
		}
	}
	ntBits := nodeType(bits)

	switch nodetype {
	case ntConcatenate | beforeChild, ntConcatenate | afterChild, ntEmpty:
		break

	case ntAlternate | beforeChild:
		if curIndex < len(node.children)-1 {
			w.pushInt(w.curPos())
			w.emit1(Lazybranch, 0)
		}

	case ntAlternate | afterChild:
		if curIndex < len(node.children)-1 {
			lbPos := w.popInt()
			w.pushInt(w.curPos())
			w.emit1(Goto, 0)
			w.patchJump(lbPos, w.curPos())
		} else {
			for i := 0; i < curIndex; i++ {
				w.patchJump(w.popInt(), w.curPos())
			}
		}
		break

	case ntTestref | beforeChild:
		if curIndex == 0 {
			w.emit(Setjump)
			w.pushInt(w.curPos())
			w.emit1(Lazybranch, 0)
			w.emit1(Testref, w.mapCapnum(node.m))
			w.emit(Forejump)
		}

	case ntTestref | afterChild:
		if curIndex == 0 {
			branchpos := w.popInt()
			w.pushInt(w.curPos())
			w.emit1(Goto, 0)
			w.patchJump(branchpos, w.curPos())
			w.emit(Forejump)
			if len(node.children) <= 1 {
				w.patchJump(w.popInt(), w.curPos())
			}
		} else if curIndex == 1 {
			w.patchJump(w.popInt(), w.curPos())
		}

	case ntTestgroup | beforeChild:
		if curIndex == 0 {
			w.emit(Setjump)
			w.emit(Setmark)
			w.pushInt(w.curPos())
			w.emit1(Lazybranch, 0)
		}

	case ntTestgroup | afterChild:
		if curIndex == 0 {
			w.emit(Getmark)
			w.emit(Forejump)
		} else if curIndex == 1 {
			Branchpos := w.popInt()
			w.pushInt(w.curPos())
			w.emit1(Goto, 0)
			w.patchJump(Branchpos, w.curPos())
			w.emit(Getmark)
			w.emit(Forejump)
			if len(node.children) <= 2 {
				w.patchJump(w.popInt(), w.curPos())
			}
		} else if curIndex == 2 {
			w.patchJump(w.popInt(), w.curPos())
		}

	case ntLoop | beforeChild, ntLazyloop | beforeChild:

		if node.n < math.MaxInt32 || node.m > 1 {
			if node.m == 0 {
				w.emit1(Nullcount, 0)
			} else {
				w.emit1(Setcount, 1-node.m)
			}
		} else if node.m == 0 {
			w.emit(Nullmark)
		} else {
			w.emit(Setmark)
		}

		if node.m == 0 {
			w.pushInt(w.curPos())
			w.emit1(Goto, 0)
		}
		w.pushInt(w.curPos())

	case ntLoop | afterChild, ntLazyloop | afterChild:

		startJumpPos := w.curPos()
		lazy := (nodetype - (ntLoop | afterChild))

		if node.n < math.MaxInt32 || node.m > 1 {
			if node.n == math.MaxInt32 {
				w.emit2(InstOp(Branchcount+lazy), w.popInt(), math.MaxInt32)
			} else {
				w.emit2(InstOp(Branchcount+lazy), w.popInt(), node.n-node.m)
			}
		} else {
			w.emit1(InstOp(Branchmark+lazy), w.popInt())
		}

		if node.m == 0 {
			w.patchJump(w.popInt(), startJumpPos)
		}

	case ntGroup | beforeChild, ntGroup | afterChild:

	case ntCapture | beforeChild:
		w.emit(Setmark)

	case ntCapture | afterChild:
		w.emit2(Capturemark, w.mapCapnum(node.m), w.mapCapnum(node.n))

	case ntRequire | beforeChild:
		// NOTE: the following line causes lookahead/lookbehind to be
		// NON-BACKTRACKING. It can be commented out with (*)
		w.emit(Setjump)

		w.emit(Setmark)

	case ntRequire | afterChild:
		w.emit(Getmark)

		// NOTE: the following line causes lookahead/lookbehind to be
		// NON-BACKTRACKING. It can be commented out with (*)
		w.emit(Forejump)

	case ntPrevent | beforeChild:
		w.emit(Setjump)
		w.pushInt(w.curPos())
		w.emit1(Lazybranch, 0)

	case ntPrevent | afterChild:
		w.emit(Backjump)
		w.patchJump(w.popInt(), w.curPos())
		w.emit(Forejump)

	case ntGreedy | beforeChild:
		w.emit(Setjump)

	case ntGreedy | afterChild:
		w.emit(Forejump)

	case ntOne, ntNotone:
		w.emit1(InstOp(node.t|ntBits), int(node.ch))

	case ntNotoneloop, ntNotonelazy, ntOneloop, ntOnelazy:
		if node.m > 0 {
			if node.t == ntOneloop || node.t == ntOnelazy {
				w.emit2(Onerep|bits, int(node.ch), node.m)
			} else {
				w.emit2(Notonerep|bits, int(node.ch), node.m)
			}
		}
		if node.n > node.m {
			if node.n == math.MaxInt32 {
				w.emit2(InstOp(node.t|ntBits), int(node.ch), math.MaxInt32)
			} else {
				w.emit2(InstOp(node.t|ntBits), int(node.ch), node.n-node.m)
			}
		}

	case ntSetloop, ntSetlazy:
		if node.m > 0 {
			w.emit2(Setrep|bits, w.setCode(node.set), node.m)
		}
		if node.n > node.m {
			if node.n == math.MaxInt32 {
				w.emit2(InstOp(node.t|ntBits), w.setCode(node.set), math.MaxInt32)
			} else {
				w.emit2(InstOp(node.t|ntBits), w.setCode(node.set), node.n-node.m)
			}
		}

	case ntMulti:
		w.emit1(InstOp(node.t|ntBits), w.stringCode(node.str))

	case ntSet:
		w.emit1(InstOp(node.t|ntBits), w.setCode(node.set))

	case ntRef:
		w.emit1(InstOp(node.t|ntBits), w.mapCapnum(node.m))

	case ntNothing, ntBol, ntEol, ntBoundary, ntNonboundary, ntECMABoundary, ntNonECMABoundary, ntBeginning, ntStart, ntEndZ, ntEnd:
		w.emit(InstOp(node.t))

	default:
		return fmt.Errorf("unexpected opcode in regular expression generation: %v", nodetype)
	}

	return nil
}

// To avoid recursion, we use a simple integer stack.
// This is the push.
func (w *writer) pushInt(i int) {
	w.intStack = append(w.intStack, i)
}

// Returns true if the stack is empty.
func (w *writer) emptyStack() bool {
	return len(w.intStack) == 0
}

// This is the pop.
func (w *writer) popInt() int {
	//get our item
	idx := len(w.intStack) - 1
	i := w.intStack[idx]
	//trim our slice
	w.intStack = w.intStack[:idx]
	return i
}

// Returns the current position in the emitted code.
func (w *writer) curPos() int {
	return w.curpos
}

// Fixes up a jump instruction at the specified offset
// so that it jumps to the specified jumpDest.
func (w *writer) patchJump(offset, jumpDest int) {
	w.emitted[offset+1] = jumpDest
}

// Returns an index in the set table for a charset
// uses a map to eliminate duplicates.
func (w *writer) setCode(set *CharSet) int {
	if w.counting {
		return 0
	}

	buf := &bytes.Buffer{}

	set.mapHashFill(buf)
	hash := buf.String()
	i, ok := w.sethash[hash]
	if !ok {
		i = len(w.sethash)
		w.sethash[hash] = i
		w.settable = append(w.settable, set)
	}
	return i
}

// Returns an index in the string table for a string.
// uses a map to eliminate duplicates.
func (w *writer) stringCode(str []rune) int {
	if w.counting {
		return 0
	}

	hash := string(str)
	i, ok := w.stringhash[hash]
	if !ok {
		i = len(w.stringhash)
		w.stringhash[hash] = i
		w.stringtable = append(w.stringtable, str)
	}

	return i
}

// When generating code on a regex that uses a sparse set
// of capture slots, we hash them to a dense set of indices
// for an array of capture slots. Instead of doing the hash
// at match time, it's done at compile time, here.
func (w *writer) mapCapnum(capnum int) int {
	if capnum == -1 {
		return -1
	}

	if w.caps != nil {
		return w.caps[capnum]
	}

	return capnum
}

// Emits a zero-argument operation. Note that the emit
// functions all run in two modes: they can emit code, or
// they can just count the size of the code.
func (w *writer) emit(op InstOp) {
	if w.counting {
		w.count++
		if opcodeBacktracks(op) {
			w.trackcount++
		}
		return
	}
	w.emitted[w.curpos] = int(op)
	w.curpos++
}

// Emits a one-argument operation.
func (w *writer) emit1(op InstOp, opd1 int) {
	if w.counting {
		w.count += 2
		if opcodeBacktracks(op) {
			w.trackcount++
		}
		return
	}
	w.emitted[w.curpos] = int(op)
	w.curpos++
	w.emitted[w.curpos] = opd1
	w.curpos++
}

// Emits a two-argument operation.
func (w *writer) emit2(op InstOp, opd1, opd2 int) {
	if w.counting {
		w.count += 3
		if opcodeBacktracks(op) {
			w.trackcount++
		}
		return
	}
	w.emitted[w.curpos] = int(op)
	w.curpos++
	w.emitted[w.curpos] = opd1
	w.curpos++
	w.emitted[w.curpos] = opd2
	w.curpos++
}
