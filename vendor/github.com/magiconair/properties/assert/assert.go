// Copyright 2018 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package assert provides helper functions for testing.
package assert

import (
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// skip defines the default call depth
const skip = 2

// Equal asserts that got and want are equal as defined by
// reflect.DeepEqual. The test fails with msg if they are not equal.
func Equal(t *testing.T, got, want interface{}, msg ...string) {
	if x := equal(2, got, want, msg...); x != "" {
		fmt.Println(x)
		t.Fail()
	}
}

func equal(skip int, got, want interface{}, msg ...string) string {
	if !reflect.DeepEqual(got, want) {
		return fail(skip, "got %v want %v %s", got, want, strings.Join(msg, " "))
	}
	return ""
}

// Panic asserts that function fn() panics.
// It assumes that recover() either returns a string or
// an error and fails if the message does not match
// the regular expression in 'matches'.
func Panic(t *testing.T, fn func(), matches string) {
	if x := doesPanic(2, fn, matches); x != "" {
		fmt.Println(x)
		t.Fail()
	}
}

func doesPanic(skip int, fn func(), expr string) (err string) {
	defer func() {
		r := recover()
		if r == nil {
			err = fail(skip, "did not panic")
			return
		}
		var v string
		switch r.(type) {
		case error:
			v = r.(error).Error()
		case string:
			v = r.(string)
		}
		err = matches(skip, v, expr)
	}()
	fn()
	return ""
}

// Matches asserts that a value matches a given regular expression.
func Matches(t *testing.T, value, expr string) {
	if x := matches(2, value, expr); x != "" {
		fmt.Println(x)
		t.Fail()
	}
}

func matches(skip int, value, expr string) string {
	ok, err := regexp.MatchString(expr, value)
	if err != nil {
		return fail(skip, "invalid pattern %q. %s", expr, err)
	}
	if !ok {
		return fail(skip, "got %s which does not match %s", value, expr)
	}
	return ""
}

func fail(skip int, format string, args ...interface{}) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("\t%s:%d: %s\n", filepath.Base(file), line, fmt.Sprintf(format, args...))
}
