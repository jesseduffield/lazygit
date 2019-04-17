package errors

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestStackFormat(t *testing.T) {

	defer func() {
		err := recover()
		if err != 'a' {
			t.Fatal(err)
		}

		e, expected := Errorf("hi"), callers()

		bs := [][]uintptr{e.stack, expected}

		if err := compareStacks(bs[0], bs[1]); err != nil {
			t.Errorf("Stack didn't match")
			t.Errorf(err.Error())
		}

		stack := string(e.Stack())

		if !strings.Contains(stack, "a: b(5)") {
			t.Errorf("Stack trace does not contain source line: 'a: b(5)'")
			t.Errorf(stack)
		}
		if !strings.Contains(stack, "error_test.go:") {
			t.Errorf("Stack trace does not contain file name: 'error_test.go:'")
			t.Errorf(stack)
		}
	}()

	a()
}

func TestSkipWorks(t *testing.T) {

	defer func() {
		err := recover()
		if err != 'a' {
			t.Fatal(err)
		}

		bs := [][]uintptr{Wrap("hi", 2).stack, callersSkip(2)}

		if err := compareStacks(bs[0], bs[1]); err != nil {
			t.Errorf("Stack didn't match")
			t.Errorf(err.Error())
		}
	}()

	a()
}

func TestNew(t *testing.T) {

	err := New("foo")

	if err.Error() != "foo" {
		t.Errorf("Wrong message")
	}

	err = New(fmt.Errorf("foo"))

	if err.Error() != "foo" {
		t.Errorf("Wrong message")
	}

	bs := [][]uintptr{New("foo").stack, callers()}

	if err := compareStacks(bs[0], bs[1]); err != nil {
		t.Errorf("Stack didn't match")
		t.Errorf(err.Error())
	}

	if err.ErrorStack() != err.TypeName()+" "+err.Error()+"\n"+string(err.Stack()) {
		t.Errorf("ErrorStack is in the wrong format")
	}
}

func TestIs(t *testing.T) {

	if Is(nil, io.EOF) {
		t.Errorf("nil is an error")
	}

	if !Is(io.EOF, io.EOF) {
		t.Errorf("io.EOF is not io.EOF")
	}

	if !Is(io.EOF, New(io.EOF)) {
		t.Errorf("io.EOF is not New(io.EOF)")
	}

	if !Is(New(io.EOF), New(io.EOF)) {
		t.Errorf("New(io.EOF) is not New(io.EOF)")
	}

	if Is(io.EOF, fmt.Errorf("io.EOF")) {
		t.Errorf("io.EOF is fmt.Errorf")
	}

}

func TestWrapError(t *testing.T) {

	e := func() error {
		return Wrap("hi", 1)
	}()

	if e.Error() != "hi" {
		t.Errorf("Constructor with a string failed")
	}

	if Wrap(fmt.Errorf("yo"), 0).Error() != "yo" {
		t.Errorf("Constructor with an error failed")
	}

	if Wrap(e, 0) != e {
		t.Errorf("Constructor with an Error failed")
	}

	if Wrap(nil, 0).Error() != "<nil>" {
		t.Errorf("Constructor with nil failed")
	}
}

func TestWrapPrefixError(t *testing.T) {

	e := func() error {
		return WrapPrefix("hi", "prefix", 1)
	}()

	if e.Error() != "prefix: hi" {
		t.Errorf("Constructor with a string failed")
	}

	if WrapPrefix(fmt.Errorf("yo"), "prefix", 0).Error() != "prefix: yo" {
		t.Errorf("Constructor with an error failed")
	}

	prefixed := WrapPrefix(e, "prefix", 0)
	original := e.(*Error)

	if prefixed.Err != original.Err || !reflect.DeepEqual(prefixed.stack, original.stack) || !reflect.DeepEqual(prefixed.frames, original.frames) || prefixed.Error() != "prefix: prefix: hi" {
		t.Errorf("Constructor with an Error failed")
	}

	if original.Error() == prefixed.Error() {
		t.Errorf("WrapPrefix changed the original error")
	}

	if WrapPrefix(nil, "prefix", 0).Error() != "prefix: <nil>" {
		t.Errorf("Constructor with nil failed")
	}

	if !strings.HasSuffix(original.StackFrames()[0].File, "error_test.go") || strings.HasSuffix(original.StackFrames()[1].File, "error_test.go") {
		t.Errorf("Skip failed")
	}
}

func ExampleErrorf(x int) (int, error) {
	if x%2 == 1 {
		return 0, Errorf("can only halve even numbers, got %d", x)
	}
	return x / 2, nil
}

func ExampleWrapError() (error, error) {
	// Wrap io.EOF with the current stack-trace and return it
	return nil, Wrap(io.EOF, 0)
}

func ExampleWrapError_skip() {
	defer func() {
		if err := recover(); err != nil {
			// skip 1 frame (the deferred function) and then return the wrapped err
			err = Wrap(err, 1)
		}
	}()
}

func ExampleIs(reader io.Reader, buff []byte) {
	_, err := reader.Read(buff)
	if Is(err, io.EOF) {
		return
	}
}

func ExampleNew(UnexpectedEOF error) error {
	// calling New attaches the current stacktrace to the existing UnexpectedEOF error
	return New(UnexpectedEOF)
}

func ExampleWrap() error {

	if err := recover(); err != nil {
		return Wrap(err, 1)
	}

	return a()
}

func ExampleError_Error(err error) {
	fmt.Println(err.Error())
}

func ExampleError_ErrorStack(err error) {
	fmt.Println(err.(*Error).ErrorStack())
}

func ExampleError_Stack(err *Error) {
	fmt.Println(err.Stack())
}

func ExampleError_TypeName(err *Error) {
	fmt.Println(err.TypeName(), err.Error())
}

func ExampleError_StackFrames(err *Error) {
	for _, frame := range err.StackFrames() {
		fmt.Println(frame.File, frame.LineNumber, frame.Package, frame.Name)
	}
}

func a() error {
	b(5)
	return nil
}

func b(i int) {
	c()
}

func c() {
	panic('a')
}

// compareStacks will compare a stack created using the errors package (actual)
// with a reference stack created with the callers function (expected). The
// first entry is compared inexact since the actual and expected stacks cannot
// be created at the exact same program counter position so the first entry
// will always differ somewhat. Returns nil if the stacks are equal enough and
// an error containing a detailed error message otherwise.
func compareStacks(actual, expected []uintptr) error {
	if len(actual) != len(expected) {
		return stackCompareError("Stacks does not have equal length", actual, expected)
	}
	for i, pc := range actual {
		if i == 0 {
			firstEntryDiff := (int)(expected[i]) - (int)(pc)
			if firstEntryDiff < -27 || firstEntryDiff > 27 {
				return stackCompareError(fmt.Sprintf("First entry PC diff to large (%d)", firstEntryDiff), actual, expected)
			}
		} else if pc != expected[i] {
			return stackCompareError(fmt.Sprintf("Stacks does not match entry %d (and maybe others)", i), actual, expected)
		}
	}
	return nil
}

func stackCompareError(msg string, actual, expected []uintptr) error {
	return fmt.Errorf("%s\nActual stack trace:\n%s\nExpected stack trace:\n%s", msg, readableStackTrace(actual), readableStackTrace(expected))
}

func callers() []uintptr {
	return callersSkip(1)
}

func callersSkip(skip int) []uintptr {
	callers := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(skip+2, callers[:])
	return callers[:length]
}

func readableStackTrace(callers []uintptr) string {
	var result bytes.Buffer
	frames := callersToFrames(callers)
	for _, frame := range frames {
		result.WriteString(fmt.Sprintf("%s:%d (%#x)\n\t%s\n", frame.File, frame.Line, frame.PC, frame.Function))
	}
	return result.String()
}

func callersToFrames(callers []uintptr) []runtime.Frame {
	frames := make([]runtime.Frame, 0, len(callers))
	framesPtr := runtime.CallersFrames(callers)
	for {
		frame, more := framesPtr.Next()
		frames = append(frames, frame)
		if !more {
			return frames
		}
	}
}
