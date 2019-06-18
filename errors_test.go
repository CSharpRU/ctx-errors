package errors_test

import (
	errors "github.com/eyudkin/ctx-errors"
	"testing"
)

func Test_Value(t *testing.T) {
	err := occurErrPackageBar()
	if err == nil {
		t.Fatalf("error expected, nil received")
	}

	// We're able to check error code from package foo...
	if errors.Value(err, ErrorCodeKeyFoo) != Error1Foo {
		t.Errorf("foo error 1 expected")
	}
	// And the same way we're able to check error code from package bar.
	if errors.Value(err, ErrorCodeKeyBar) != ErrorABar {
		t.Errorf("bar error A expected")
	}
}

// Package foo.

type errKeyFoo int8

const (
	ErrorCodeKeyFoo errKeyFoo = iota
	Error1Foo
	Error2Foo
	Error3Foo
	ErrAddDataFoo
)

func occurErrPackageFoo() error {
	return errors.New("foo err").
		WithValue(ErrorCodeKeyFoo, Error1Foo).
		WithValue(ErrAddDataFoo, "additional data from foo")
}

// Package bar
type errKeyBar int8

const (
	ErrorCodeKeyBar errKeyBar = iota
	ErrorABar
	ErrorBBar
	ErrorCBar
	ErrAddDataBar
)

// Function which calls occurErrPackageFoo and adds its own context for its error.
func occurErrPackageBar() error {
	if err := occurErrPackageFoo(); err != nil {
		return errors.Wrap(err, "bar err").WithValue(ErrorCodeKeyBar, ErrorABar).
			WithValue(ErrAddDataBar, "additional data from bar")
	}

	return nil
}
