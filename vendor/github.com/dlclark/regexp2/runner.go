package regexp2

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/dlclark/regexp2/syntax"
)

type runner struct {
	re   *Regexp
	code *syntax.Code

	runtextstart int // starting point for search

	runtext    []rune // text to search
	runtextpos int    // current position in text
	runtextend int

	// The backtracking stack.  Opcodes use this to store data regarding
	// what they have matched and where to backtrack to.  Each "frame" on
	// the stack takes the form of [CodePosition Data1 Data2...], where
	// CodePosition is the position of the current opcode and
	// the data values are all optional.  The CodePosition can be negative, and
	// these values (also called "back2") are used by the BranchMark family of opcodes
	// to indicate whether they are backtracking after a successful or failed
	// match.
	// When we backtrack, we pop the CodePosition off the stack, set the current
	// instruction pointer to that code position, and mark the opcode
	// with a backtracking flag ("Back").  Each opcode then knows how to
	// handle its own data.
	runtrack    []int
	runtrackpos int

	// This stack is used to track text positions across different opcodes.
	// For example, in /(a*b)+/, the parentheses result in a SetMark/CaptureMark
	// pair. SetMark records the text position before we match a*b.  Then
	// CaptureMark uses that position to figure out where the capture starts.
	// Opcodes which push onto this stack are always paired with other opcodes
	// which will pop the value from it later.  A successful match should mean
	// that this stack is empty.
	runstack    []int
	runstackpos int

	// The crawl stack is used to keep track of captures.  Every time a group
	// has a capture, we push its group number onto the runcrawl stack.  In
	// the case of a balanced match, we push BOTH groups onto the stack.
	runcrawl    []int
	runcrawlpos int

	runtrackcount int // count of states that may do backtracking

	runmatch *Match // result object

	ignoreTimeout       bool
	timeout             time.Duration // timeout in milliseconds (needed for actual)
	timeoutChecksToSkip int
	timeoutAt           time.Time

	operator        syntax.InstOp
	codepos         int
	rightToLeft     bool
	caseInsensitive bool
}

// run searches for matches and can continue from the previous match
//
// quick is usually false, but can be true to not return matches, just put it in caches
// textstart is -1 to start at the "beginning" (depending on Right-To-Left), otherwise an index in input
// input is the string to search for our regex pattern
func (re *Regexp) run(quick bool, textstart int, input []rune) (*Match, error) {

	// get a cached runner
	runner := re.getRunner()
	defer re.putRunner(runner)

	if textstart < 0 {
		if re.RightToLeft() {
			textstart = len(input)
		} else {
			textstart = 0
		}
	}

	return runner.scan(input, textstart, quick, re.MatchTimeout)
}

// Scans the string to find the first match. Uses the Match object
// both to feed text in and as a place to store matches that come out.
//
// All the action is in the Go() method. Our
// responsibility is to load up the class members before
// calling Go.
//
// The optimizer can compute a set of candidate starting characters,
// and we could use a separate method Skip() that will quickly scan past
// any characters that we know can't match.
func (r *runner) scan(rt []rune, textstart int, quick bool, timeout time.Duration) (*Match, error) {
	r.timeout = timeout
	r.ignoreTimeout = (time.Duration(math.MaxInt64) == timeout)
	r.runtextstart = textstart
	r.runtext = rt
	r.runtextend = len(rt)

	stoppos := r.runtextend
	bump := 1

	if r.re.RightToLeft() {
		bump = -1
		stoppos = 0
	}

	r.runtextpos = textstart
	initted := false

	r.startTimeoutWatch()
	for {
		if r.re.Debug() {
			//fmt.Printf("\nSearch content: %v\n", string(r.runtext))
			fmt.Printf("\nSearch range: from 0 to %v\n", r.runtextend)
			fmt.Printf("Firstchar search starting at %v stopping at %v\n", r.runtextpos, stoppos)
		}

		if r.findFirstChar() {
			if err := r.checkTimeout(); err != nil {
				return nil, err
			}

			if !initted {
				r.initMatch()
				initted = true
			}

			if r.re.Debug() {
				fmt.Printf("Executing engine starting at %v\n\n", r.runtextpos)
			}

			if err := r.execute(); err != nil {
				return nil, err
			}

			if r.runmatch.matchcount[0] > 0 {
				// We'll return a match even if it touches a previous empty match
				return r.tidyMatch(quick), nil
			}

			// reset state for another go
			r.runtrackpos = len(r.runtrack)
			r.runstackpos = len(r.runstack)
			r.runcrawlpos = len(r.runcrawl)
		}

		// failure!

		if r.runtextpos == stoppos {
			r.tidyMatch(true)
			return nil, nil
		}

		// Recognize leading []* and various anchors, and bump on failure accordingly

		// r.bump by one and start again

		r.runtextpos += bump
	}
	// We never get here
}

func (r *runner) execute() error {

	r.goTo(0)

	for {

		if r.re.Debug() {
			r.dumpState()
		}

		if err := r.checkTimeout(); err != nil {
			return err
		}

		switch r.operator {
		case syntax.Stop:
			return nil

		case syntax.Nothing:
			break

		case syntax.Goto:
			r.goTo(r.operand(0))
			continue

		case syntax.Testref:
			if !r.runmatch.isMatched(r.operand(0)) {
				break
			}
			r.advance(1)
			continue

		case syntax.Lazybranch:
			r.trackPush1(r.textPos())
			r.advance(1)
			continue

		case syntax.Lazybranch | syntax.Back:
			r.trackPop()
			r.textto(r.trackPeek())
			r.goTo(r.operand(0))
			continue

		case syntax.Setmark:
			r.stackPush(r.textPos())
			r.trackPush()
			r.advance(0)
			continue

		case syntax.Nullmark:
			r.stackPush(-1)
			r.trackPush()
			r.advance(0)
			continue

		case syntax.Setmark | syntax.Back, syntax.Nullmark | syntax.Back:
			r.stackPop()
			break

		case syntax.Getmark:
			r.stackPop()
			r.trackPush1(r.stackPeek())
			r.textto(r.stackPeek())
			r.advance(0)
			continue

		case syntax.Getmark | syntax.Back:
			r.trackPop()
			r.stackPush(r.trackPeek())
			break

		case syntax.Capturemark:
			if r.operand(1) != -1 && !r.runmatch.isMatched(r.operand(1)) {
				break
			}
			r.stackPop()
			if r.operand(1) != -1 {
				r.transferCapture(r.operand(0), r.operand(1), r.stackPeek(), r.textPos())
			} else {
				r.capture(r.operand(0), r.stackPeek(), r.textPos())
			}
			r.trackPush1(r.stackPeek())

			r.advance(2)

			continue

		case syntax.Capturemark | syntax.Back:
			r.trackPop()
			r.stackPush(r.trackPeek())
			r.uncapture()
			if r.operand(0) != -1 && r.operand(1) != -1 {
				r.uncapture()
			}

			break

		case syntax.Branchmark:
			r.stackPop()

			matched := r.textPos() - r.stackPeek()

			if matched != 0 { // Nonempty match -> loop now
				r.trackPush2(r.stackPeek(), r.textPos()) // Save old mark, textpos
				r.stackPush(r.textPos())                 // Make new mark
				r.goTo(r.operand(0))                     // Loop
			} else { // Empty match -> straight now
				r.trackPushNeg1(r.stackPeek()) // Save old mark
				r.advance(1)                   // Straight
			}
			continue

		case syntax.Branchmark | syntax.Back:
			r.trackPopN(2)
			r.stackPop()
			r.textto(r.trackPeekN(1))      // Recall position
			r.trackPushNeg1(r.trackPeek()) // Save old mark
			r.advance(1)                   // Straight
			continue

		case syntax.Branchmark | syntax.Back2:
			r.trackPop()
			r.stackPush(r.trackPeek()) // Recall old mark
			break                      // Backtrack

		case syntax.Lazybranchmark:
			{
				// We hit this the first time through a lazy loop and after each
				// successful match of the inner expression.  It simply continues
				// on and doesn't loop.
				r.stackPop()

				oldMarkPos := r.stackPeek()

				if r.textPos() != oldMarkPos { // Nonempty match -> try to loop again by going to 'back' state
					if oldMarkPos != -1 {
						r.trackPush2(oldMarkPos, r.textPos()) // Save old mark, textpos
					} else {
						r.trackPush2(r.textPos(), r.textPos())
					}
				} else {
					// The inner expression found an empty match, so we'll go directly to 'back2' if we
					// backtrack.  In this case, we need to push something on the stack, since back2 pops.
					// However, in the case of ()+? or similar, this empty match may be legitimate, so push the text
					// position associated with that empty match.
					r.stackPush(oldMarkPos)

					r.trackPushNeg1(r.stackPeek()) // Save old mark
				}
				r.advance(1)
				continue
			}

		case syntax.Lazybranchmark | syntax.Back:

			// After the first time, Lazybranchmark | syntax.Back occurs
			// with each iteration of the loop, and therefore with every attempted
			// match of the inner expression.  We'll try to match the inner expression,
			// then go back to Lazybranchmark if successful.  If the inner expression
			// fails, we go to Lazybranchmark | syntax.Back2

			r.trackPopN(2)
			pos := r.trackPeekN(1)
			r.trackPushNeg1(r.trackPeek()) // Save old mark
			r.stackPush(pos)               // Make new mark
			r.textto(pos)                  // Recall position
			r.goTo(r.operand(0))           // Loop
			continue

		case syntax.Lazybranchmark | syntax.Back2:
			// The lazy loop has failed.  We'll do a true backtrack and
			// start over before the lazy loop.
			r.stackPop()
			r.trackPop()
			r.stackPush(r.trackPeek()) // Recall old mark
			break

		case syntax.Setcount:
			r.stackPush2(r.textPos(), r.operand(0))
			r.trackPush()
			r.advance(1)
			continue

		case syntax.Nullcount:
			r.stackPush2(-1, r.operand(0))
			r.trackPush()
			r.advance(1)
			continue

		case syntax.Setcount | syntax.Back:
			r.stackPopN(2)
			break

		case syntax.Nullcount | syntax.Back:
			r.stackPopN(2)
			break

		case syntax.Branchcount:
			// r.stackPush:
			//  0: Mark
			//  1: Count

			r.stackPopN(2)
			mark := r.stackPeek()
			count := r.stackPeekN(1)
			matched := r.textPos() - mark

			if count >= r.operand(1) || (matched == 0 && count >= 0) { // Max loops or empty match -> straight now
				r.trackPushNeg2(mark, count) // Save old mark, count
				r.advance(2)                 // Straight
			} else { // Nonempty match -> count+loop now
				r.trackPush1(mark)                 // remember mark
				r.stackPush2(r.textPos(), count+1) // Make new mark, incr count
				r.goTo(r.operand(0))               // Loop
			}
			continue

		case syntax.Branchcount | syntax.Back:
			// r.trackPush:
			//  0: Previous mark
			// r.stackPush:
			//  0: Mark (= current pos, discarded)
			//  1: Count
			r.trackPop()
			r.stackPopN(2)
			if r.stackPeekN(1) > 0 { // Positive -> can go straight
				r.textto(r.stackPeek())                           // Zap to mark
				r.trackPushNeg2(r.trackPeek(), r.stackPeekN(1)-1) // Save old mark, old count
				r.advance(2)                                      // Straight
				continue
			}
			r.stackPush2(r.trackPeek(), r.stackPeekN(1)-1) // recall old mark, old count
			break

		case syntax.Branchcount | syntax.Back2:
			// r.trackPush:
			//  0: Previous mark
			//  1: Previous count
			r.trackPopN(2)
			r.stackPush2(r.trackPeek(), r.trackPeekN(1)) // Recall old mark, old count
			break                                        // Backtrack

		case syntax.Lazybranchcount:
			// r.stackPush:
			//  0: Mark
			//  1: Count

			r.stackPopN(2)
			mark := r.stackPeek()
			count := r.stackPeekN(1)

			if count < 0 { // Negative count -> loop now
				r.trackPushNeg1(mark)              // Save old mark
				r.stackPush2(r.textPos(), count+1) // Make new mark, incr count
				r.goTo(r.operand(0))               // Loop
			} else { // Nonneg count -> straight now
				r.trackPush3(mark, count, r.textPos()) // Save mark, count, position
				r.advance(2)                           // Straight
			}
			continue

		case syntax.Lazybranchcount | syntax.Back:
			// r.trackPush:
			//  0: Mark
			//  1: Count
			//  2: r.textPos

			r.trackPopN(3)
			mark := r.trackPeek()
			textpos := r.trackPeekN(2)

			if r.trackPeekN(1) < r.operand(1) && textpos != mark { // Under limit and not empty match -> loop
				r.textto(textpos)                        // Recall position
				r.stackPush2(textpos, r.trackPeekN(1)+1) // Make new mark, incr count
				r.trackPushNeg1(mark)                    // Save old mark
				r.goTo(r.operand(0))                     // Loop
				continue
			} else { // Max loops or empty match -> backtrack
				r.stackPush2(r.trackPeek(), r.trackPeekN(1)) // Recall old mark, count
				break                                        // backtrack
			}

		case syntax.Lazybranchcount | syntax.Back2:
			// r.trackPush:
			//  0: Previous mark
			// r.stackPush:
			//  0: Mark (== current pos, discarded)
			//  1: Count
			r.trackPop()
			r.stackPopN(2)
			r.stackPush2(r.trackPeek(), r.stackPeekN(1)-1) // Recall old mark, count
			break                                          // Backtrack

		case syntax.Setjump:
			r.stackPush2(r.trackpos(), r.crawlpos())
			r.trackPush()
			r.advance(0)
			continue

		case syntax.Setjump | syntax.Back:
			r.stackPopN(2)
			break

		case syntax.Backjump:
			// r.stackPush:
			//  0: Saved trackpos
			//  1: r.crawlpos
			r.stackPopN(2)
			r.trackto(r.stackPeek())

			for r.crawlpos() != r.stackPeekN(1) {
				r.uncapture()
			}

			break

		case syntax.Forejump:
			// r.stackPush:
			//  0: Saved trackpos
			//  1: r.crawlpos
			r.stackPopN(2)
			r.trackto(r.stackPeek())
			r.trackPush1(r.stackPeekN(1))
			r.advance(0)
			continue

		case syntax.Forejump | syntax.Back:
			// r.trackPush:
			//  0: r.crawlpos
			r.trackPop()

			for r.crawlpos() != r.trackPeek() {
				r.uncapture()
			}

			break

		case syntax.Bol:
			if r.leftchars() > 0 && r.charAt(r.textPos()-1) != '\n' {
				break
			}
			r.advance(0)
			continue

		case syntax.Eol:
			if r.rightchars() > 0 && r.charAt(r.textPos()) != '\n' {
				break
			}
			r.advance(0)
			continue

		case syntax.Boundary:
			if !r.isBoundary(r.textPos(), 0, r.runtextend) {
				break
			}
			r.advance(0)
			continue

		case syntax.Nonboundary:
			if r.isBoundary(r.textPos(), 0, r.runtextend) {
				break
			}
			r.advance(0)
			continue

		case syntax.ECMABoundary:
			if !r.isECMABoundary(r.textPos(), 0, r.runtextend) {
				break
			}
			r.advance(0)
			continue

		case syntax.NonECMABoundary:
			if r.isECMABoundary(r.textPos(), 0, r.runtextend) {
				break
			}
			r.advance(0)
			continue

		case syntax.Beginning:
			if r.leftchars() > 0 {
				break
			}
			r.advance(0)
			continue

		case syntax.Start:
			if r.textPos() != r.textstart() {
				break
			}
			r.advance(0)
			continue

		case syntax.EndZ:
			rchars := r.rightchars()
			if rchars > 1 {
				break
			}
			// RE2 and EcmaScript define $ as "asserts position at the end of the string"
			// PCRE/.NET adds "or before the line terminator right at the end of the string (if any)"
			if (r.re.options & (RE2 | ECMAScript)) != 0 {
				// RE2/Ecmascript mode
				if rchars > 0 {
					break
				}
			} else if rchars == 1 && r.charAt(r.textPos()) != '\n' {
				// "regular" mode
				break
			}

			r.advance(0)
			continue

		case syntax.End:
			if r.rightchars() > 0 {
				break
			}
			r.advance(0)
			continue

		case syntax.One:
			if r.forwardchars() < 1 || r.forwardcharnext() != rune(r.operand(0)) {
				break
			}

			r.advance(1)
			continue

		case syntax.Notone:
			if r.forwardchars() < 1 || r.forwardcharnext() == rune(r.operand(0)) {
				break
			}

			r.advance(1)
			continue

		case syntax.Set:

			if r.forwardchars() < 1 || !r.code.Sets[r.operand(0)].CharIn(r.forwardcharnext()) {
				break
			}

			r.advance(1)
			continue

		case syntax.Multi:
			if !r.runematch(r.code.Strings[r.operand(0)]) {
				break
			}

			r.advance(1)
			continue

		case syntax.Ref:

			capnum := r.operand(0)

			if r.runmatch.isMatched(capnum) {
				if !r.refmatch(r.runmatch.matchIndex(capnum), r.runmatch.matchLength(capnum)) {
					break
				}
			} else {
				if (r.re.options & ECMAScript) == 0 {
					break
				}
			}

			r.advance(1)
			continue

		case syntax.Onerep:

			c := r.operand(1)

			if r.forwardchars() < c {
				break
			}

			ch := rune(r.operand(0))

			for c > 0 {
				if r.forwardcharnext() != ch {
					goto BreakBackward
				}
				c--
			}

			r.advance(2)
			continue

		case syntax.Notonerep:

			c := r.operand(1)

			if r.forwardchars() < c {
				break
			}
			ch := rune(r.operand(0))

			for c > 0 {
				if r.forwardcharnext() == ch {
					goto BreakBackward
				}
				c--
			}

			r.advance(2)
			continue

		case syntax.Setrep:

			c := r.operand(1)

			if r.forwardchars() < c {
				break
			}

			set := r.code.Sets[r.operand(0)]

			for c > 0 {
				if !set.CharIn(r.forwardcharnext()) {
					goto BreakBackward
				}
				c--
			}

			r.advance(2)
			continue

		case syntax.Oneloop:

			c := r.operand(1)

			if c > r.forwardchars() {
				c = r.forwardchars()
			}

			ch := rune(r.operand(0))
			i := c

			for ; i > 0; i-- {
				if r.forwardcharnext() != ch {
					r.backwardnext()
					break
				}
			}

			if c > i {
				r.trackPush2(c-i-1, r.textPos()-r.bump())
			}

			r.advance(2)
			continue

		case syntax.Notoneloop:

			c := r.operand(1)

			if c > r.forwardchars() {
				c = r.forwardchars()
			}

			ch := rune(r.operand(0))
			i := c

			for ; i > 0; i-- {
				if r.forwardcharnext() == ch {
					r.backwardnext()
					break
				}
			}

			if c > i {
				r.trackPush2(c-i-1, r.textPos()-r.bump())
			}

			r.advance(2)
			continue

		case syntax.Setloop:

			c := r.operand(1)

			if c > r.forwardchars() {
				c = r.forwardchars()
			}

			set := r.code.Sets[r.operand(0)]
			i := c

			for ; i > 0; i-- {
				if !set.CharIn(r.forwardcharnext()) {
					r.backwardnext()
					break
				}
			}

			if c > i {
				r.trackPush2(c-i-1, r.textPos()-r.bump())
			}

			r.advance(2)
			continue

		case syntax.Oneloop | syntax.Back, syntax.Notoneloop | syntax.Back:

			r.trackPopN(2)
			i := r.trackPeek()
			pos := r.trackPeekN(1)

			r.textto(pos)

			if i > 0 {
				r.trackPush2(i-1, pos-r.bump())
			}

			r.advance(2)
			continue

		case syntax.Setloop | syntax.Back:

			r.trackPopN(2)
			i := r.trackPeek()
			pos := r.trackPeekN(1)

			r.textto(pos)

			if i > 0 {
				r.trackPush2(i-1, pos-r.bump())
			}

			r.advance(2)
			continue

		case syntax.Onelazy, syntax.Notonelazy:

			c := r.operand(1)

			if c > r.forwardchars() {
				c = r.forwardchars()
			}

			if c > 0 {
				r.trackPush2(c-1, r.textPos())
			}

			r.advance(2)
			continue

		case syntax.Setlazy:

			c := r.operand(1)

			if c > r.forwardchars() {
				c = r.forwardchars()
			}

			if c > 0 {
				r.trackPush2(c-1, r.textPos())
			}

			r.advance(2)
			continue

		case syntax.Onelazy | syntax.Back:

			r.trackPopN(2)
			pos := r.trackPeekN(1)
			r.textto(pos)

			if r.forwardcharnext() != rune(r.operand(0)) {
				break
			}

			i := r.trackPeek()

			if i > 0 {
				r.trackPush2(i-1, pos+r.bump())
			}

			r.advance(2)
			continue

		case syntax.Notonelazy | syntax.Back:

			r.trackPopN(2)
			pos := r.trackPeekN(1)
			r.textto(pos)

			if r.forwardcharnext() == rune(r.operand(0)) {
				break
			}

			i := r.trackPeek()

			if i > 0 {
				r.trackPush2(i-1, pos+r.bump())
			}

			r.advance(2)
			continue

		case syntax.Setlazy | syntax.Back:

			r.trackPopN(2)
			pos := r.trackPeekN(1)
			r.textto(pos)

			if !r.code.Sets[r.operand(0)].CharIn(r.forwardcharnext()) {
				break
			}

			i := r.trackPeek()

			if i > 0 {
				r.trackPush2(i-1, pos+r.bump())
			}

			r.advance(2)
			continue

		default:
			return errors.New("unknown state in regex runner")
		}

	BreakBackward:
		;

		// "break Backward" comes here:
		r.backtrack()
	}
}

// increase the size of stack and track storage
func (r *runner) ensureStorage() {
	if r.runstackpos < r.runtrackcount*4 {
		doubleIntSlice(&r.runstack, &r.runstackpos)
	}
	if r.runtrackpos < r.runtrackcount*4 {
		doubleIntSlice(&r.runtrack, &r.runtrackpos)
	}
}

func doubleIntSlice(s *[]int, pos *int) {
	oldLen := len(*s)
	newS := make([]int, oldLen*2)

	copy(newS[oldLen:], *s)
	*pos += oldLen
	*s = newS
}

// Save a number on the longjump unrolling stack
func (r *runner) crawl(i int) {
	if r.runcrawlpos == 0 {
		doubleIntSlice(&r.runcrawl, &r.runcrawlpos)
	}
	r.runcrawlpos--
	r.runcrawl[r.runcrawlpos] = i
}

// Remove a number from the longjump unrolling stack
func (r *runner) popcrawl() int {
	val := r.runcrawl[r.runcrawlpos]
	r.runcrawlpos++
	return val
}

// Get the height of the stack
func (r *runner) crawlpos() int {
	return len(r.runcrawl) - r.runcrawlpos
}

func (r *runner) advance(i int) {
	r.codepos += (i + 1)
	r.setOperator(r.code.Codes[r.codepos])
}

func (r *runner) goTo(newpos int) {
	// when branching backward or in place, ensure storage
	if newpos <= r.codepos {
		r.ensureStorage()
	}

	r.setOperator(r.code.Codes[newpos])
	r.codepos = newpos
}

func (r *runner) textto(newpos int) {
	r.runtextpos = newpos
}

func (r *runner) trackto(newpos int) {
	r.runtrackpos = len(r.runtrack) - newpos
}

func (r *runner) textstart() int {
	return r.runtextstart
}

func (r *runner) textPos() int {
	return r.runtextpos
}

// push onto the backtracking stack
func (r *runner) trackpos() int {
	return len(r.runtrack) - r.runtrackpos
}

func (r *runner) trackPush() {
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = r.codepos
}

func (r *runner) trackPush1(I1 int) {
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I1
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = r.codepos
}

func (r *runner) trackPush2(I1, I2 int) {
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I1
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I2
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = r.codepos
}

func (r *runner) trackPush3(I1, I2, I3 int) {
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I1
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I2
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I3
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = r.codepos
}

func (r *runner) trackPushNeg1(I1 int) {
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I1
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = -r.codepos
}

func (r *runner) trackPushNeg2(I1, I2 int) {
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I1
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = I2
	r.runtrackpos--
	r.runtrack[r.runtrackpos] = -r.codepos
}

func (r *runner) backtrack() {
	newpos := r.runtrack[r.runtrackpos]
	r.runtrackpos++

	if r.re.Debug() {
		if newpos < 0 {
			fmt.Printf("       Backtracking (back2) to code position %v\n", -newpos)
		} else {
			fmt.Printf("       Backtracking to code position %v\n", newpos)
		}
	}

	if newpos < 0 {
		newpos = -newpos
		r.setOperator(r.code.Codes[newpos] | syntax.Back2)
	} else {
		r.setOperator(r.code.Codes[newpos] | syntax.Back)
	}

	// When branching backward, ensure storage
	if newpos < r.codepos {
		r.ensureStorage()
	}

	r.codepos = newpos
}

func (r *runner) setOperator(op int) {
	r.caseInsensitive = (0 != (op & syntax.Ci))
	r.rightToLeft = (0 != (op & syntax.Rtl))
	r.operator = syntax.InstOp(op & ^(syntax.Rtl | syntax.Ci))
}

func (r *runner) trackPop() {
	r.runtrackpos++
}

// pop framesize items from the backtracking stack
func (r *runner) trackPopN(framesize int) {
	r.runtrackpos += framesize
}

// Technically we are actually peeking at items already popped.  So if you want to
// get and pop the top item from the stack, you do
// r.trackPop();
// r.trackPeek();
func (r *runner) trackPeek() int {
	return r.runtrack[r.runtrackpos-1]
}

// get the ith element down on the backtracking stack
func (r *runner) trackPeekN(i int) int {
	return r.runtrack[r.runtrackpos-i-1]
}

// Push onto the grouping stack
func (r *runner) stackPush(I1 int) {
	r.runstackpos--
	r.runstack[r.runstackpos] = I1
}

func (r *runner) stackPush2(I1, I2 int) {
	r.runstackpos--
	r.runstack[r.runstackpos] = I1
	r.runstackpos--
	r.runstack[r.runstackpos] = I2
}

func (r *runner) stackPop() {
	r.runstackpos++
}

// pop framesize items from the grouping stack
func (r *runner) stackPopN(framesize int) {
	r.runstackpos += framesize
}

// Technically we are actually peeking at items already popped.  So if you want to
// get and pop the top item from the stack, you do
// r.stackPop();
// r.stackPeek();
func (r *runner) stackPeek() int {
	return r.runstack[r.runstackpos-1]
}

// get the ith element down on the grouping stack
func (r *runner) stackPeekN(i int) int {
	return r.runstack[r.runstackpos-i-1]
}

func (r *runner) operand(i int) int {
	return r.code.Codes[r.codepos+i+1]
}

func (r *runner) leftchars() int {
	return r.runtextpos
}

func (r *runner) rightchars() int {
	return r.runtextend - r.runtextpos
}

func (r *runner) bump() int {
	if r.rightToLeft {
		return -1
	}
	return 1
}

func (r *runner) forwardchars() int {
	if r.rightToLeft {
		return r.runtextpos
	}
	return r.runtextend - r.runtextpos
}

func (r *runner) forwardcharnext() rune {
	var ch rune
	if r.rightToLeft {
		r.runtextpos--
		ch = r.runtext[r.runtextpos]
	} else {
		ch = r.runtext[r.runtextpos]
		r.runtextpos++
	}

	if r.caseInsensitive {
		return unicode.ToLower(ch)
	}
	return ch
}

func (r *runner) runematch(str []rune) bool {
	var pos int

	c := len(str)
	if !r.rightToLeft {
		if r.runtextend-r.runtextpos < c {
			return false
		}

		pos = r.runtextpos + c
	} else {
		if r.runtextpos-0 < c {
			return false
		}

		pos = r.runtextpos
	}

	if !r.caseInsensitive {
		for c != 0 {
			c--
			pos--
			if str[c] != r.runtext[pos] {
				return false
			}
		}
	} else {
		for c != 0 {
			c--
			pos--
			if str[c] != unicode.ToLower(r.runtext[pos]) {
				return false
			}
		}
	}

	if !r.rightToLeft {
		pos += len(str)
	}

	r.runtextpos = pos

	return true
}

func (r *runner) refmatch(index, len int) bool {
	var c, pos, cmpos int

	if !r.rightToLeft {
		if r.runtextend-r.runtextpos < len {
			return false
		}

		pos = r.runtextpos + len
	} else {
		if r.runtextpos-0 < len {
			return false
		}

		pos = r.runtextpos
	}
	cmpos = index + len

	c = len

	if !r.caseInsensitive {
		for c != 0 {
			c--
			cmpos--
			pos--
			if r.runtext[cmpos] != r.runtext[pos] {
				return false
			}

		}
	} else {
		for c != 0 {
			c--
			cmpos--
			pos--

			if unicode.ToLower(r.runtext[cmpos]) != unicode.ToLower(r.runtext[pos]) {
				return false
			}
		}
	}

	if !r.rightToLeft {
		pos += len
	}

	r.runtextpos = pos

	return true
}

func (r *runner) backwardnext() {
	if r.rightToLeft {
		r.runtextpos++
	} else {
		r.runtextpos--
	}
}

func (r *runner) charAt(j int) rune {
	return r.runtext[j]
}

func (r *runner) findFirstChar() bool {

	if 0 != (r.code.Anchors & (syntax.AnchorBeginning | syntax.AnchorStart | syntax.AnchorEndZ | syntax.AnchorEnd)) {
		if !r.code.RightToLeft {
			if (0 != (r.code.Anchors&syntax.AnchorBeginning) && r.runtextpos > 0) ||
				(0 != (r.code.Anchors&syntax.AnchorStart) && r.runtextpos > r.runtextstart) {
				r.runtextpos = r.runtextend
				return false
			}
			if 0 != (r.code.Anchors&syntax.AnchorEndZ) && r.runtextpos < r.runtextend-1 {
				r.runtextpos = r.runtextend - 1
			} else if 0 != (r.code.Anchors&syntax.AnchorEnd) && r.runtextpos < r.runtextend {
				r.runtextpos = r.runtextend
			}
		} else {
			if (0 != (r.code.Anchors&syntax.AnchorEnd) && r.runtextpos < r.runtextend) ||
				(0 != (r.code.Anchors&syntax.AnchorEndZ) && (r.runtextpos < r.runtextend-1 ||
					(r.runtextpos == r.runtextend-1 && r.charAt(r.runtextpos) != '\n'))) ||
				(0 != (r.code.Anchors&syntax.AnchorStart) && r.runtextpos < r.runtextstart) {
				r.runtextpos = 0
				return false
			}
			if 0 != (r.code.Anchors&syntax.AnchorBeginning) && r.runtextpos > 0 {
				r.runtextpos = 0
			}
		}

		if r.code.BmPrefix != nil {
			return r.code.BmPrefix.IsMatch(r.runtext, r.runtextpos, 0, r.runtextend)
		}

		return true // found a valid start or end anchor
	} else if r.code.BmPrefix != nil {
		r.runtextpos = r.code.BmPrefix.Scan(r.runtext, r.runtextpos, 0, r.runtextend)

		if r.runtextpos == -1 {
			if r.code.RightToLeft {
				r.runtextpos = 0
			} else {
				r.runtextpos = r.runtextend
			}
			return false
		}

		return true
	} else if r.code.FcPrefix == nil {
		return true
	}

	r.rightToLeft = r.code.RightToLeft
	r.caseInsensitive = r.code.FcPrefix.CaseInsensitive

	set := r.code.FcPrefix.PrefixSet
	if set.IsSingleton() {
		ch := set.SingletonChar()
		for i := r.forwardchars(); i > 0; i-- {
			if ch == r.forwardcharnext() {
				r.backwardnext()
				return true
			}
		}
	} else {
		for i := r.forwardchars(); i > 0; i-- {
			n := r.forwardcharnext()
			//fmt.Printf("%v in %v: %v\n", string(n), set.String(), set.CharIn(n))
			if set.CharIn(n) {
				r.backwardnext()
				return true
			}
		}
	}

	return false
}

func (r *runner) initMatch() {
	// Use a hashtable'ed Match object if the capture numbers are sparse

	if r.runmatch == nil {
		if r.re.caps != nil {
			r.runmatch = newMatchSparse(r.re, r.re.caps, r.re.capsize, r.runtext, r.runtextstart)
		} else {
			r.runmatch = newMatch(r.re, r.re.capsize, r.runtext, r.runtextstart)
		}
	} else {
		r.runmatch.reset(r.runtext, r.runtextstart)
	}

	// note we test runcrawl, because it is the last one to be allocated
	// If there is an alloc failure in the middle of the three allocations,
	// we may still return to reuse this instance, and we want to behave
	// as if the allocations didn't occur. (we used to test _trackcount != 0)

	if r.runcrawl != nil {
		r.runtrackpos = len(r.runtrack)
		r.runstackpos = len(r.runstack)
		r.runcrawlpos = len(r.runcrawl)
		return
	}

	r.initTrackCount()

	tracksize := r.runtrackcount * 8
	stacksize := r.runtrackcount * 8

	if tracksize < 32 {
		tracksize = 32
	}
	if stacksize < 16 {
		stacksize = 16
	}

	r.runtrack = make([]int, tracksize)
	r.runtrackpos = tracksize

	r.runstack = make([]int, stacksize)
	r.runstackpos = stacksize

	r.runcrawl = make([]int, 32)
	r.runcrawlpos = 32
}

func (r *runner) tidyMatch(quick bool) *Match {
	if !quick {
		match := r.runmatch

		r.runmatch = nil

		match.tidy(r.runtextpos)
		return match
	} else {
		// send back our match -- it's not leaving the package, so it's safe to not clean it up
		// this reduces allocs for frequent calls to the "IsMatch" bool-only functions
		return r.runmatch
	}
}

// capture captures a subexpression. Note that the
// capnum used here has already been mapped to a non-sparse
// index (by the code generator RegexWriter).
func (r *runner) capture(capnum, start, end int) {
	if end < start {
		T := end
		end = start
		start = T
	}

	r.crawl(capnum)
	r.runmatch.addMatch(capnum, start, end-start)
}

// transferCapture captures a subexpression. Note that the
// capnum used here has already been mapped to a non-sparse
// index (by the code generator RegexWriter).
func (r *runner) transferCapture(capnum, uncapnum, start, end int) {
	var start2, end2 int

	// these are the two intervals that are cancelling each other

	if end < start {
		T := end
		end = start
		start = T
	}

	start2 = r.runmatch.matchIndex(uncapnum)
	end2 = start2 + r.runmatch.matchLength(uncapnum)

	// The new capture gets the innermost defined interval

	if start >= end2 {
		end = start
		start = end2
	} else if end <= start2 {
		start = start2
	} else {
		if end > end2 {
			end = end2
		}
		if start2 > start {
			start = start2
		}
	}

	r.crawl(uncapnum)
	r.runmatch.balanceMatch(uncapnum)

	if capnum != -1 {
		r.crawl(capnum)
		r.runmatch.addMatch(capnum, start, end-start)
	}
}

// revert the last capture
func (r *runner) uncapture() {
	capnum := r.popcrawl()
	r.runmatch.removeMatch(capnum)
}

//debug

func (r *runner) dumpState() {
	back := ""
	if r.operator&syntax.Back != 0 {
		back = " Back"
	}
	if r.operator&syntax.Back2 != 0 {
		back += " Back2"
	}
	fmt.Printf("Text:  %v\nTrack: %v\nStack: %v\n       %s%s\n\n",
		r.textposDescription(),
		r.stackDescription(r.runtrack, r.runtrackpos),
		r.stackDescription(r.runstack, r.runstackpos),
		r.code.OpcodeDescription(r.codepos),
		back)
}

func (r *runner) stackDescription(a []int, index int) string {
	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "%v/%v", len(a)-index, len(a))
	if buf.Len() < 8 {
		buf.WriteString(strings.Repeat(" ", 8-buf.Len()))
	}

	buf.WriteRune('(')
	for i := index; i < len(a); i++ {
		if i > index {
			buf.WriteRune(' ')
		}

		buf.WriteString(strconv.Itoa(a[i]))
	}

	buf.WriteRune(')')

	return buf.String()
}

func (r *runner) textposDescription() string {
	buf := &bytes.Buffer{}

	buf.WriteString(strconv.Itoa(r.runtextpos))

	if buf.Len() < 8 {
		buf.WriteString(strings.Repeat(" ", 8-buf.Len()))
	}

	if r.runtextpos > 0 {
		buf.WriteString(syntax.CharDescription(r.runtext[r.runtextpos-1]))
	} else {
		buf.WriteRune('^')
	}

	buf.WriteRune('>')

	for i := r.runtextpos; i < r.runtextend; i++ {
		buf.WriteString(syntax.CharDescription(r.runtext[i]))
	}
	if buf.Len() >= 64 {
		buf.Truncate(61)
		buf.WriteString("...")
	} else {
		buf.WriteRune('$')
	}

	return buf.String()
}

// decide whether the pos
// at the specified index is a boundary or not. It's just not worth
// emitting inline code for this logic.
func (r *runner) isBoundary(index, startpos, endpos int) bool {
	return (index > startpos && syntax.IsWordChar(r.runtext[index-1])) !=
		(index < endpos && syntax.IsWordChar(r.runtext[index]))
}

func (r *runner) isECMABoundary(index, startpos, endpos int) bool {
	return (index > startpos && syntax.IsECMAWordChar(r.runtext[index-1])) !=
		(index < endpos && syntax.IsECMAWordChar(r.runtext[index]))
}

// this seems like a comment to justify randomly picking 1000 :-P
// We have determined this value in a series of experiments where x86 retail
// builds (ono-lab-optimized) were run on different pattern/input pairs. Larger values
// of TimeoutCheckFrequency did not tend to increase performance; smaller values
// of TimeoutCheckFrequency tended to slow down the execution.
const timeoutCheckFrequency int = 1000

func (r *runner) startTimeoutWatch() {
	if r.ignoreTimeout {
		return
	}

	r.timeoutChecksToSkip = timeoutCheckFrequency
	r.timeoutAt = time.Now().Add(r.timeout)
}

func (r *runner) checkTimeout() error {
	if r.ignoreTimeout {
		return nil
	}
	r.timeoutChecksToSkip--
	if r.timeoutChecksToSkip != 0 {
		return nil
	}

	r.timeoutChecksToSkip = timeoutCheckFrequency
	return r.doCheckTimeout()
}

func (r *runner) doCheckTimeout() error {
	current := time.Now()

	if current.Before(r.timeoutAt) {
		return nil
	}

	if r.re.Debug() {
		//Debug.WriteLine("")
		//Debug.WriteLine("RegEx match timeout occurred!")
		//Debug.WriteLine("Specified timeout:       " + TimeSpan.FromMilliseconds(_timeout).ToString())
		//Debug.WriteLine("Timeout check frequency: " + TimeoutCheckFrequency)
		//Debug.WriteLine("Search pattern:          " + _runregex._pattern)
		//Debug.WriteLine("Input:                   " + r.runtext)
		//Debug.WriteLine("About to throw RegexMatchTimeoutException.")
	}

	return fmt.Errorf("match timeout after %v on input `%v`", r.timeout, string(r.runtext))
}

func (r *runner) initTrackCount() {
	r.runtrackcount = r.code.TrackCount
}

// getRunner returns a run to use for matching re.
// It uses the re's runner cache if possible, to avoid
// unnecessary allocation.
func (re *Regexp) getRunner() *runner {
	re.muRun.Lock()
	if n := len(re.runner); n > 0 {
		z := re.runner[n-1]
		re.runner = re.runner[:n-1]
		re.muRun.Unlock()
		return z
	}
	re.muRun.Unlock()
	z := &runner{
		re:   re,
		code: re.code,
	}
	return z
}

// putRunner returns a runner to the re's cache.
// There is no attempt to limit the size of the cache, so it will
// grow to the maximum number of simultaneous matches
// run using re.  (The cache empties when re gets garbage collected.)
func (re *Regexp) putRunner(r *runner) {
	re.muRun.Lock()
	re.runner = append(re.runner, r)
	re.muRun.Unlock()
}
