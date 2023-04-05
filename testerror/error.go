package testerror

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// Error is a simple error type that allows us to chain checks using methods.
// However due to Go's way of handling interface types, when a nil *Error value is used as
// an error type (e.g. when passed into a func accepting an error) the underlying concrete
// type is *Error and it will not pass the usual error != nil check.
// Instead use the AsError() method to get a regular error type.
// Or use the reflect package `val := reflect.ValueOf(myCheckError); val.IsNil()`.
type Error struct {
	msg string
}

func New(msg string) *Error {
	return &Error{
		msg: msg,
	}
}

func Newf(format string, args ...any) *Error {
	return &Error{
		msg: fmt.Sprintf(format, args...),
	}
}

// Implement Error interface
func (e *Error) Error() string {
	return e.msg
}

// AsError returns a regular error type that can be used in the usual way.
// This fixes some issues when comparing nil types, which can fail as the underlying types are different.
// e.g. comparing the a nil error interface type to a nil *Error type from this package can fail.
func (e *Error) AsError() error {
	if e == nil {
		return nil
	}
	return errors.New(e.msg)
}

func (e *Error) ErrorIsNil(t *testing.T) {
	if e != nil {
		t.Errorf(e.msg)
	}
}

func (e *Error) ErrorIsNilFatal(t *testing.T) {
	if e != nil {
		t.Fatalf(e.msg)
	}
}

func (e *Error) ErrorNotNil(t *testing.T) {
	if e == nil {
		t.Errorf("error is nil")
	}
}

func (e *Error) ErrorNotNilFatal(t *testing.T) {
	if e == nil {
		t.Fatalf("error is nil")
	}
}

func (e *Error) ErrorContains(t *testing.T, substr string) {
	if !strings.Contains(e.msg, substr) {
		t.Errorf("error '%s' does not contain substring '%s'", e.msg, substr)
	}
}

func (e *Error) ErrorNotContains(t *testing.T, substr string) {
	if strings.Contains(e.msg, substr) {
		t.Errorf("error '%s' does contain substring '%s'", e.msg, substr)
	}
}
