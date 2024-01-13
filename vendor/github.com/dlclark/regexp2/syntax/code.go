package syntax

import (
	"bytes"
	"fmt"
	"math"
)

// similar to prog.go in the go regex package...also with comment 'may not belong in this package'

// File provides operator constants for use by the Builder and the Machine.

// Implementation notes:
//
// Regexps are built into RegexCodes, which contain an operation array,
// a string table, and some constants.
//
// Each operation is one of the codes below, followed by the integer
// operands specified for each op.
//
// Strings and sets are indices into a string table.

type InstOp int

const (
	// 					    lef/back operands        description

	Onerep    InstOp = 0 // lef,back char,min,max    a {n}
	Notonerep        = 1 // lef,back char,min,max    .{n}
	Setrep           = 2 // lef,back set,min,max     [\d]{n}

	Oneloop    = 3 // lef,back char,min,max    a {,n}
	Notoneloop = 4 // lef,back char,min,max    .{,n}
	Setloop    = 5 // lef,back set,min,max     [\d]{,n}

	Onelazy    = 6 // lef,back char,min,max    a {,n}?
	Notonelazy = 7 // lef,back char,min,max    .{,n}?
	Setlazy    = 8 // lef,back set,min,max     [\d]{,n}?

	One    = 9  // lef      char            a
	Notone = 10 // lef      char            [^a]
	Set    = 11 // lef      set             [a-z\s]  \w \s \d

	Multi = 12 // lef      string          abcd
	Ref   = 13 // lef      group           \#

	Bol         = 14 //                          ^
	Eol         = 15 //                          $
	Boundary    = 16 //                          \b
	Nonboundary = 17 //                          \B
	Beginning   = 18 //                          \A
	Start       = 19 //                          \G
	EndZ        = 20 //                          \Z
	End         = 21 //                          \Z

	Nothing = 22 //                          Reject!

	// Primitive control structures

	Lazybranch      = 23 // back     jump            straight first
	Branchmark      = 24 // back     jump            branch first for loop
	Lazybranchmark  = 25 // back     jump            straight first for loop
	Nullcount       = 26 // back     val             set counter, null mark
	Setcount        = 27 // back     val             set counter, make mark
	Branchcount     = 28 // back     jump,limit      branch++ if zero<=c<limit
	Lazybranchcount = 29 // back     jump,limit      same, but straight first
	Nullmark        = 30 // back                     save position
	Setmark         = 31 // back                     save position
	Capturemark     = 32 // back     group           define group
	Getmark         = 33 // back                     recall position
	Setjump         = 34 // back                     save backtrack state
	Backjump        = 35 //                          zap back to saved state
	Forejump        = 36 //                          zap backtracking state
	Testref         = 37 //                          backtrack if ref undefined
	Goto            = 38 //          jump            just go

	Prune = 39 //                          prune it baby
	Stop  = 40 //                          done!

	ECMABoundary    = 41 //                          \b
	NonECMABoundary = 42 //                          \B

	// Modifiers for alternate modes

	Mask  = 63  // Mask to get unmodified ordinary operator
	Rtl   = 64  // bit to indicate that we're reverse scanning.
	Back  = 128 // bit to indicate that we're backtracking.
	Back2 = 256 // bit to indicate that we're backtracking on a second branch.
	Ci    = 512 // bit to indicate that we're case-insensitive.
)

type Code struct {
	Codes       []int       // the code
	Strings     [][]rune    // string table
	Sets        []*CharSet  //character set table
	TrackCount  int         // how many instructions use backtracking
	Caps        map[int]int // mapping of user group numbers -> impl group slots
	Capsize     int         // number of impl group slots
	FcPrefix    *Prefix     // the set of candidate first characters (may be null)
	BmPrefix    *BmPrefix   // the fixed prefix string as a Boyer-Moore machine (may be null)
	Anchors     AnchorLoc   // the set of zero-length start anchors (RegexFCD.Bol, etc)
	RightToLeft bool        // true if right to left
}

func opcodeBacktracks(op InstOp) bool {
	op &= Mask

	switch op {
	case Oneloop, Notoneloop, Setloop, Onelazy, Notonelazy, Setlazy, Lazybranch, Branchmark, Lazybranchmark,
		Nullcount, Setcount, Branchcount, Lazybranchcount, Setmark, Capturemark, Getmark, Setjump, Backjump,
		Forejump, Goto:
		return true

	default:
		return false
	}
}

func opcodeSize(op InstOp) int {
	op &= Mask

	switch op {
	case Nothing, Bol, Eol, Boundary, Nonboundary, ECMABoundary, NonECMABoundary, Beginning, Start, EndZ,
		End, Nullmark, Setmark, Getmark, Setjump, Backjump, Forejump, Stop:
		return 1

	case One, Notone, Multi, Ref, Testref, Goto, Nullcount, Setcount, Lazybranch, Branchmark, Lazybranchmark,
		Prune, Set:
		return 2

	case Capturemark, Branchcount, Lazybranchcount, Onerep, Notonerep, Oneloop, Notoneloop, Onelazy, Notonelazy,
		Setlazy, Setrep, Setloop:
		return 3

	default:
		panic(fmt.Errorf("Unexpected op code: %v", op))
	}
}

var codeStr = []string{
	"Onerep", "Notonerep", "Setrep",
	"Oneloop", "Notoneloop", "Setloop",
	"Onelazy", "Notonelazy", "Setlazy",
	"One", "Notone", "Set",
	"Multi", "Ref",
	"Bol", "Eol", "Boundary", "Nonboundary", "Beginning", "Start", "EndZ", "End",
	"Nothing",
	"Lazybranch", "Branchmark", "Lazybranchmark",
	"Nullcount", "Setcount", "Branchcount", "Lazybranchcount",
	"Nullmark", "Setmark", "Capturemark", "Getmark",
	"Setjump", "Backjump", "Forejump", "Testref", "Goto",
	"Prune", "Stop",
	"ECMABoundary", "NonECMABoundary",
}

func operatorDescription(op InstOp) string {
	desc := codeStr[op&Mask]
	if (op & Ci) != 0 {
		desc += "-Ci"
	}
	if (op & Rtl) != 0 {
		desc += "-Rtl"
	}
	if (op & Back) != 0 {
		desc += "-Back"
	}
	if (op & Back2) != 0 {
		desc += "-Back2"
	}

	return desc
}

// OpcodeDescription is a humman readable string of the specific offset
func (c *Code) OpcodeDescription(offset int) string {
	buf := &bytes.Buffer{}

	op := InstOp(c.Codes[offset])
	fmt.Fprintf(buf, "%06d ", offset)

	if opcodeBacktracks(op & Mask) {
		buf.WriteString("*")
	} else {
		buf.WriteString(" ")
	}
	buf.WriteString(operatorDescription(op))
	buf.WriteString("(")
	op &= Mask

	switch op {
	case One, Notone, Onerep, Notonerep, Oneloop, Notoneloop, Onelazy, Notonelazy:
		buf.WriteString("Ch = ")
		buf.WriteString(CharDescription(rune(c.Codes[offset+1])))

	case Set, Setrep, Setloop, Setlazy:
		buf.WriteString("Set = ")
		buf.WriteString(c.Sets[c.Codes[offset+1]].String())

	case Multi:
		fmt.Fprintf(buf, "String = %s", string(c.Strings[c.Codes[offset+1]]))

	case Ref, Testref:
		fmt.Fprintf(buf, "Index = %d", c.Codes[offset+1])

	case Capturemark:
		fmt.Fprintf(buf, "Index = %d", c.Codes[offset+1])
		if c.Codes[offset+2] != -1 {
			fmt.Fprintf(buf, ", Unindex = %d", c.Codes[offset+2])
		}

	case Nullcount, Setcount:
		fmt.Fprintf(buf, "Value = %d", c.Codes[offset+1])

	case Goto, Lazybranch, Branchmark, Lazybranchmark, Branchcount, Lazybranchcount:
		fmt.Fprintf(buf, "Addr = %d", c.Codes[offset+1])
	}

	switch op {
	case Onerep, Notonerep, Oneloop, Notoneloop, Onelazy, Notonelazy, Setrep, Setloop, Setlazy:
		buf.WriteString(", Rep = ")
		if c.Codes[offset+2] == math.MaxInt32 {
			buf.WriteString("inf")
		} else {
			fmt.Fprintf(buf, "%d", c.Codes[offset+2])
		}

	case Branchcount, Lazybranchcount:
		buf.WriteString(", Limit = ")
		if c.Codes[offset+2] == math.MaxInt32 {
			buf.WriteString("inf")
		} else {
			fmt.Fprintf(buf, "%d", c.Codes[offset+2])
		}

	}

	buf.WriteString(")")

	return buf.String()
}

func (c *Code) Dump() string {
	buf := &bytes.Buffer{}

	if c.RightToLeft {
		fmt.Fprintln(buf, "Direction:  right-to-left")
	} else {
		fmt.Fprintln(buf, "Direction:  left-to-right")
	}
	if c.FcPrefix == nil {
		fmt.Fprintln(buf, "Firstchars: n/a")
	} else {
		fmt.Fprintf(buf, "Firstchars: %v\n", c.FcPrefix.PrefixSet.String())
	}

	if c.BmPrefix == nil {
		fmt.Fprintln(buf, "Prefix:     n/a")
	} else {
		fmt.Fprintf(buf, "Prefix:     %v\n", Escape(c.BmPrefix.String()))
	}

	fmt.Fprintf(buf, "Anchors:    %v\n", c.Anchors)
	fmt.Fprintln(buf)

	if c.BmPrefix != nil {
		fmt.Fprintln(buf, "BoyerMoore:")
		fmt.Fprintln(buf, c.BmPrefix.Dump("    "))
	}
	for i := 0; i < len(c.Codes); i += opcodeSize(InstOp(c.Codes[i])) {
		fmt.Fprintln(buf, c.OpcodeDescription(i))
	}

	return buf.String()
}
