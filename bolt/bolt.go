package bolt

// The following applies only to this file (bolt.go).
// The MIT License (MIT)

// Copyright (c) 2013 Ben Johnson

// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
//https://github.com/boltdb/bolt/blob/51f99c862475898df9773747d3accd05a7ca33c1/bolt_test.go

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// Assert fails the test if the condition is false. Is fatal.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		decorateAndLog(msg, v...)
		tb.FailNow()
	}
}

// Check fails the test if the condition is false. Is not fatal.
func Check(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		decorateAndLog(msg, v...)
		tb.Fail()
	}
}

// fails the test if an err is not nil. Is fatal upon failure.
func AssertOk(tb testing.TB, err error) {
	if err != nil {
		decorateAndLog("unexpected error: %s", err.Error())
		tb.FailNow()
	}
}

// fails the test if an err is not nil. Is not fatal.
func CheckOk(tb testing.TB, err error) {
	if err != nil {
		decorateAndLog("unexpected error: %s", err.Error())
		tb.Fail()
	}
}

//fails the test if exp is not equal to act. Is fatal upon failure
func AssertEquals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		decorateAndLog("\n\texp: %#v\n\tgot: %#v", exp, act)
		tb.FailNow()
	}
}

//fails the test if exp is not equal to act. Is not fatal
func CheckEquals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		decorateAndLog("\n\texp: %#v\n\tgot: %#v", exp, act)
		tb.Fail()
	}
}

// Prefixes the message with the file name and line number at the point of the original error
// and prints it to standard output.
func decorateAndLog(msg string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(2)

	if ok {
		format := fmt.Sprintf("\033[7;1;31m%s:%d: %s\033[0;39m\n\n", filepath.Base(file), line, msg)
		fmt.Printf(format, v...)
	} else {
		format := fmt.Sprintf("\033[7;1;31m ???:1: %s\033[0;39m\n\n", msg)
		fmt.Printf(format, v...)
	}
}
