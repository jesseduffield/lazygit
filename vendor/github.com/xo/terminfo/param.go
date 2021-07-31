package terminfo

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
)

// parametizer represents the a scan state for a parameterized string.
type parametizer struct {
	// z is the string to parameterize
	z []byte

	// pos is the current position in s.
	pos int

	// nest is the current nest level.
	nest int

	// s is the variable stack.
	s stack

	// skipElse keeps the state of skipping else.
	skipElse bool

	// buf is the result buffer.
	buf *bytes.Buffer

	// params are the parameters to interpolate.
	params [9]interface{}

	// vars are dynamic variables.
	vars [26]interface{}
}

// staticVars are the static, global variables.
var staticVars = struct {
	vars [26]interface{}
	sync.Mutex
}{}

var parametizerPool = sync.Pool{
	New: func() interface{} {
		p := new(parametizer)
		p.buf = bytes.NewBuffer(make([]byte, 0, 45))
		return p
	},
}

// newParametizer returns a new initialized parametizer from the pool.
func newParametizer(z []byte) *parametizer {
	p := parametizerPool.Get().(*parametizer)
	p.z = z

	return p
}

// reset resets the parametizer.
func (p *parametizer) reset() {
	p.pos, p.nest = 0, 0

	p.s.reset()
	p.buf.Reset()

	p.params, p.vars = [9]interface{}{}, [26]interface{}{}

	parametizerPool.Put(p)
}

// stateFn represents the state of the scanner as a function that returns the
// next state.
type stateFn func() stateFn

// exec executes the parameterizer, interpolating the supplied parameters.
func (p *parametizer) exec() string {
	for state := p.scanTextFn; state != nil; {
		state = state()
	}
	return p.buf.String()
}

// peek returns the next byte.
func (p *parametizer) peek() (byte, error) {
	if p.pos >= len(p.z) {
		return 0, io.EOF
	}
	return p.z[p.pos], nil
}

// writeFrom writes the characters from ppos to pos to the buffer.
func (p *parametizer) writeFrom(ppos int) {
	if p.pos > ppos {
		// append remaining characters.
		p.buf.Write(p.z[ppos:p.pos])
	}
}

func (p *parametizer) scanTextFn() stateFn {
	ppos := p.pos
	for {
		ch, err := p.peek()
		if err != nil {
			p.writeFrom(ppos)
			return nil
		}

		if ch == '%' {
			p.writeFrom(ppos)
			p.pos++
			return p.scanCodeFn
		}

		p.pos++
	}
}

func (p *parametizer) scanCodeFn() stateFn {
	ch, err := p.peek()
	if err != nil {
		return nil
	}

	switch ch {
	case '%':
		p.buf.WriteByte('%')

	case ':':
		// this character is used to avoid interpreting "%-" and "%+" as operators.
		// the next character is where the format really begins.
		p.pos++
		_, err = p.peek()
		if err != nil {
			return nil
		}
		return p.scanFormatFn

	case '#', ' ', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
		return p.scanFormatFn

	case 'o':
		p.buf.WriteString(strconv.FormatInt(int64(p.s.popInt()), 8))

	case 'd':
		p.buf.WriteString(strconv.Itoa(p.s.popInt()))

	case 'x':
		p.buf.WriteString(strconv.FormatInt(int64(p.s.popInt()), 16))

	case 'X':
		p.buf.WriteString(strings.ToUpper(strconv.FormatInt(int64(p.s.popInt()), 16)))

	case 's':
		p.buf.WriteString(p.s.popString())

	case 'c':
		p.buf.WriteByte(p.s.popByte())

	case 'p':
		p.pos++
		return p.pushParamFn

	case 'P':
		p.pos++
		return p.setDsVarFn

	case 'g':
		p.pos++
		return p.getDsVarFn

	case '\'':
		p.pos++
		ch, err = p.peek()
		if err != nil {
			return nil
		}

		p.s.push(ch)

		// skip the '\''
		p.pos++

	case '{':
		p.pos++
		return p.pushIntfn

	case 'l':
		p.s.push(len(p.s.popString()))

	case '+':
		bi, ai := p.s.popInt(), p.s.popInt()
		p.s.push(ai + bi)

	case '-':
		bi, ai := p.s.popInt(), p.s.popInt()
		p.s.push(ai - bi)

	case '*':
		bi, ai := p.s.popInt(), p.s.popInt()
		p.s.push(ai * bi)

	case '/':
		bi, ai := p.s.popInt(), p.s.popInt()
		if bi != 0 {
			p.s.push(ai / bi)
		} else {
			p.s.push(0)
		}

	case 'm':
		bi, ai := p.s.popInt(), p.s.popInt()
		if bi != 0 {
			p.s.push(ai % bi)
		} else {
			p.s.push(0)
		}

	case '&':
		bi, ai := p.s.popInt(), p.s.popInt()
		p.s.push(ai & bi)

	case '|':
		bi, ai := p.s.popInt(), p.s.popInt()
		p.s.push(ai | bi)

	case '^':
		bi, ai := p.s.popInt(), p.s.popInt()
		p.s.push(ai ^ bi)

	case '=':
		bi, ai := p.s.popInt(), p.s.popInt()
		p.s.push(ai == bi)

	case '>':
		bi, ai := p.s.popInt(), p.s.popInt()
		p.s.push(ai > bi)

	case '<':
		bi, ai := p.s.popInt(), p.s.popInt()
		p.s.push(ai < bi)

	case 'A':
		bi, ai := p.s.popBool(), p.s.popBool()
		p.s.push(ai && bi)

	case 'O':
		bi, ai := p.s.popBool(), p.s.popBool()
		p.s.push(ai || bi)

	case '!':
		p.s.push(!p.s.popBool())

	case '~':
		p.s.push(^p.s.popInt())

	case 'i':
		for i := range p.params[:2] {
			if n, ok := p.params[i].(int); ok {
				p.params[i] = n + 1
			}
		}

	case '?', ';':

	case 't':
		return p.scanThenFn

	case 'e':
		p.skipElse = true
		return p.skipTextFn
	}

	p.pos++

	return p.scanTextFn
}

func (p *parametizer) scanFormatFn() stateFn {
	// the character was already read, so no need to check the error.
	ch, _ := p.peek()

	// 6 should be the maximum length of a format string, for example "%:-9.9d".
	f := []byte{'%', ch, 0, 0, 0, 0}

	var err error

	for {
		p.pos++
		ch, err = p.peek()
		if err != nil {
			return nil
		}

		f = append(f, ch)
		switch ch {
		case 'o', 'd', 'x', 'X':
			fmt.Fprintf(p.buf, string(f), p.s.popInt())
			break

		case 's':
			fmt.Fprintf(p.buf, string(f), p.s.popString())
			break

		case 'c':
			fmt.Fprintf(p.buf, string(f), p.s.popByte())
			break
		}
	}

	p.pos++

	return p.scanTextFn
}

func (p *parametizer) pushParamFn() stateFn {
	ch, err := p.peek()
	if err != nil {
		return nil
	}

	if ai := int(ch - '1'); ai >= 0 && ai < len(p.params) {
		p.s.push(p.params[ai])
	} else {
		p.s.push(0)
	}

	// skip the '}'
	p.pos++

	return p.scanTextFn
}

func (p *parametizer) setDsVarFn() stateFn {
	ch, err := p.peek()
	if err != nil {
		return nil
	}

	if ch >= 'A' && ch <= 'Z' {
		staticVars.Lock()
		staticVars.vars[int(ch-'A')] = p.s.pop()
		staticVars.Unlock()
	} else if ch >= 'a' && ch <= 'z' {
		p.vars[int(ch-'a')] = p.s.pop()
	}

	p.pos++
	return p.scanTextFn
}

func (p *parametizer) getDsVarFn() stateFn {
	ch, err := p.peek()
	if err != nil {
		return nil
	}

	var a byte
	if ch >= 'A' && ch <= 'Z' {
		a = 'A'
	} else if ch >= 'a' && ch <= 'z' {
		a = 'a'
	}

	staticVars.Lock()
	p.s.push(staticVars.vars[int(ch-a)])
	staticVars.Unlock()

	p.pos++

	return p.scanTextFn
}

func (p *parametizer) pushIntfn() stateFn {
	var ai int
	for {
		ch, err := p.peek()
		if err != nil {
			return nil
		}

		p.pos++
		if ch < '0' || ch > '9' {
			p.s.push(ai)
			return p.scanTextFn
		}

		ai = (ai * 10) + int(ch-'0')
	}
}

func (p *parametizer) scanThenFn() stateFn {
	p.pos++

	if p.s.popBool() {
		return p.scanTextFn
	}

	p.skipElse = false

	return p.skipTextFn
}

func (p *parametizer) skipTextFn() stateFn {
	for {
		ch, err := p.peek()
		if err != nil {
			return nil
		}

		p.pos++
		if ch == '%' {
			break
		}
	}

	if p.skipElse {
		return p.skipElseFn
	}

	return p.skipThenFn
}

func (p *parametizer) skipThenFn() stateFn {
	ch, err := p.peek()
	if err != nil {
		return nil
	}

	p.pos++
	switch ch {
	case ';':
		if p.nest == 0 {
			return p.scanTextFn
		}
		p.nest--

	case '?':
		p.nest++

	case 'e':
		if p.nest == 0 {
			return p.scanTextFn
		}
	}

	return p.skipTextFn
}

func (p *parametizer) skipElseFn() stateFn {
	ch, err := p.peek()
	if err != nil {
		return nil
	}

	p.pos++
	switch ch {
	case ';':
		if p.nest == 0 {
			return p.scanTextFn
		}
		p.nest--

	case '?':
		p.nest++
	}

	return p.skipTextFn
}

// Printf evaluates a parameterized terminfo value z, interpolating params.
func Printf(z []byte, params ...interface{}) string {
	p := newParametizer(z)
	defer p.reset()

	// make sure we always have 9 parameters -- makes it easier
	// later to skip checks and its faster
	for i := 0; i < len(p.params) && i < len(params); i++ {
		p.params[i] = params[i]
	}

	return p.exec()
}

// Fprintf evaluates a parameterized terminfo value z, interpolating params and
// writing to w.
func Fprintf(w io.Writer, z []byte, params ...interface{}) {
	w.Write([]byte(Printf(z, params...)))
}
