package bolt

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"testing"
)

type metaTester struct {
	testing.TB
	failed bool
	fatal  bool
}

func (t *metaTester) Fail() {
	t.failed = true
}

func (t *metaTester) FailNow() {
	t.failed = true
	t.fatal = true
}

func testerFormat(msg string, v ...interface{}) string {
	format := fmt.Sprintf("\033[7;1;31m%s\033[0;39m\n\n", msg)
	return fmt.Sprintf(format, v...)
}

func collectOutput(f func()) string {
	old := os.Stdout     // keep backup of the real stdout
	r, w, _ := os.Pipe() // hook up stdout to the write end of the pipe
	os.Stdout = w
	defer func() { os.Stdout = old }() // just in case our visitor panics or whatever.
	f()                                // may write to stdout
	w.Close()
	buf, _ := ioutil.ReadAll(r)
	r.Close()
	return string(buf)
}

func TestCheckFormatOrder(t *testing.T) {
	tt := &metaTester{}
	_, _, lineNo, _ := runtime.Caller(0) // this must be just before before the test call
	theOutPut := collectOutput(func() { Check(tt, false, "check val6=%[3]d, val5=%[2]d, val4=%[1]d", 4, 5, 6) })

	if theOutPut == "" {
		t.Error(testerFormat("expected output, got nothing."))
	}
	expected := testerFormat("bolt_test.go:%d: check val6=6, val5=5, val4=4", lineNo+1)

	if expected != theOutPut {
		t.Error(testerFormat("expected output:\n%s\ngot:\n%s", expected, theOutPut))
	}
}

func TestCheckPositive(t *testing.T) {
	var someVar = 11
	tt := &metaTester{}
	theOutPut := collectOutput(func() { Check(tt, true, "this is an error %d", someVar) })

	if tt.failed || tt.fatal {
		t.Error(testerFormat("did not expect fail or fatal"))
	}
	if theOutPut != "" {
		t.Error(testerFormat("unexpected output: %s", theOutPut))
	}
}

func TestCheckNegative(t *testing.T) {
	var someVar = 11
	tt := &metaTester{}
	_, _, lineNo, _ := runtime.Caller(0) // this must be just before before the test call
	theOutPut := collectOutput(func() { Check(tt, false, "check someVar is %d", someVar) })

	if !tt.failed || tt.fatal {
		t.Error(testerFormat("expected fail but not fatal"))
	}
	if theOutPut == "" {
		t.Error(testerFormat("expected output, got nothing."))
	}
	expected := testerFormat("bolt_test.go:%d: check someVar is 11", lineNo+1)
	if theOutPut != expected {
		t.Error(testerFormat("expected output:\n%s\ngot:\n%s", expected, theOutPut))
	}
}

func TestAssertPositive(t *testing.T) {
	var someVar = 11
	tt := &metaTester{}
	theOutPut := collectOutput(func() { Assert(tt, true, "this is an error %d", someVar) })

	if tt.failed || tt.fatal {
		t.Error(testerFormat("did not expect fail or fatal"))
	}
	if theOutPut != "" {
		t.Error(testerFormat("unexpected output: %s", theOutPut))
	}
}

func TestAssertNegative(t *testing.T) {
	var someVar = 11
	tt := &metaTester{}
	_, _, lineNo, _ := runtime.Caller(0) // this must be just before before the test call
	theOutPut := collectOutput(func() { Assert(tt, false, "Assert someVar is %d", someVar) })

	if !tt.failed || !tt.fatal {
		t.Error(testerFormat("expected fail and fatal to be true"))
	}
	if theOutPut == "" {
		t.Error(testerFormat("expected output, got nothing."))

	}
	expected := testerFormat("bolt_test.go:%d: Assert someVar is 11", lineNo+1)
	if theOutPut != expected {
		t.Error(testerFormat("expected output:\n%s\ngot:\n%s", expected, theOutPut))
	}
}

func TestCheckOkPositive(t *testing.T) {
	tt := &metaTester{}
	theOutput := collectOutput(func() { CheckOk(tt, nil) }) // should do nothing.

	if tt.failed || tt.fatal {
		t.Error(testerFormat("did not expect fail or fatal"))
	}
	if theOutput != "" {
		t.Error(testerFormat("unexpected output: %s", theOutput))
	}
}

func TestCheckOkNegative(t *testing.T) {
	tt := &metaTester{}
	err := errors.New("this was not entirely unexpected, actually.")
	_, _, lineNo, _ := runtime.Caller(0) // this must be just before before the test call
	theOutPut := collectOutput(func() { CheckOk(tt, err) })

	if !tt.failed || tt.fatal {
		t.Error(testerFormat("expected fail but not fatal"))
	}
	if theOutPut == "" {
		t.Error(testerFormat("expected output, got nothing."))
	}
	expected := testerFormat("bolt_test.go:%d: unexpected error: this was not entirely unexpected, actually.", lineNo+1)
	if theOutPut != expected {
		t.Error(testerFormat("expected output:\n%s\ngot:\n%s", expected, theOutPut))
	}
}

func TestAssertOkPositive(t *testing.T) {
	tt := &metaTester{}
	theOutPut := collectOutput(func() { AssertOk(tt, nil) }) // AssertOk is the FUT

	if tt.failed || tt.fatal {
		t.Error(testerFormat("did not expect fail or fatal"))
	}
	if theOutPut != "" {
		t.Error(testerFormat("unexpected output: %s", theOutPut))
	}
}

func TestAssertOkNegative(t *testing.T) {
	tt := &metaTester{}
	err := errors.New("this was not entirely unexpected, actually.")
	_, _, lineNo, _ := runtime.Caller(0)                     // this must be just before before the test call
	theOutPut := collectOutput(func() { AssertOk(tt, err) }) // AssertOk is the actual FUT

	if !tt.failed || !tt.fatal {
		t.Error(testerFormat("expected fail and fatal to be true"))
	}
	if theOutPut == "" {
		t.Error(testerFormat("expected output, got nothing."))
	}
	expected := testerFormat("bolt_test.go:%d: unexpected error: this was not entirely unexpected, actually.", lineNo+1)
	if theOutPut != expected {
		t.Error(testerFormat("expected output:\n%s\ngot:\n%s", expected, theOutPut))
	}
}

func TestCheckEqualsPositive(t *testing.T) {
	tt := &metaTester{}
	var someVar1 = 1
	var somevar2 = someVar1
	theOutPut := collectOutput(func() { CheckEquals(tt, someVar1, somevar2) }) // no message

	if tt.failed || tt.fatal {
		t.Error(testerFormat("did not expect fail or fatal"))
	}
	if theOutPut != "" {
		t.Error(testerFormat("unexpected output: %s", theOutPut))
	}
}

func TestCheckEqualsNegative(t *testing.T) {
	tt := &metaTester{}
	var someVar1 = 11
	var somevar2 = 12
	_, _, lineNo, _ := runtime.Caller(0)                                       // this must be just before before the test call
	theOutPut := collectOutput(func() { CheckEquals(tt, someVar1, somevar2) }) // some message

	if !tt.failed || tt.fatal {
		t.Error(testerFormat("expected fail but not fatal"))
	}
	if theOutPut == "" {
		t.Error(testerFormat("expected output, got nothing."))
	}
	expected := testerFormat("bolt_test.go:%d: \n\texp: 11\n\tgot: 12", lineNo+1)
	if theOutPut != expected {
		t.Error(testerFormat("expected output was:\n%s\ngot:\n%s", expected, theOutPut))
	}
}

func TestAssertEqualsPositive(t *testing.T) {
	tt := &metaTester{}
	var someVar1 = 11
	var somevar2 = someVar1
	theOutPut := collectOutput(func() { AssertEquals(tt, someVar1, somevar2) }) // no fatal message

	if tt.failed || tt.fatal {
		t.Error(testerFormat("did not expect fail or fatal"))
	}
	if theOutPut != "" {
		t.Error(testerFormat("unexpected output: %s", theOutPut))
	}
}

func TestAssertEqualsNegative(t *testing.T) {
	tt := &metaTester{}
	var someVar1 = 11
	var somevar2 = 12
	_, _, lineNo, _ := runtime.Caller(0)                                        // this must be just before before the test call
	theOutPut := collectOutput(func() { AssertEquals(tt, someVar1, somevar2) }) // no fatal message

	if !tt.failed || !tt.fatal {
		t.Error(testerFormat("expected fail and fatal to be true"))
	}
	if theOutPut == "" {
		t.Error(testerFormat("expected output, got nothing."))
	}
	expected := testerFormat("bolt_test.go:%d: \n\texp: 11\n\tgot: 12", lineNo+1)
	if theOutPut != expected {
		t.Error(testerFormat("expected output was:\n%s\ngot:\n%s", expected, theOutPut))
	}
}
