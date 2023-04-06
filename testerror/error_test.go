package testerror

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrorF(t *testing.T) {
	err := Newf("test %s", "error")
	assert.Equal(t, "test error", err.Error())
}

func TestNewAsError(t *testing.T) {
	f := func(e error) string {
		return e.Error()
	}
	e := New("test error")
	assert.Equal(t, "test error", f(e))
}

func TestNewErrorNill(t *testing.T) {
	f := func() *Error {
		return nil
	}()
	assert.Nil(t, f)
}

func TestErrorContains(t *testing.T) {
	e := New("test error")
	e.ErrorContains(t, "test")
}

func TestErrorContainsFail(t *testing.T) {
	e := New("test error")
	var t1 testing.T
	e.ErrorContains(&t1, "fail")
	assert.True(t, t1.Failed())
}

func TestErrorNotContains(t *testing.T) {
	e := New("test error")
	e.ErrorNotContains(t, "notcontained")
}

func TestErrorNotContainsFail(t *testing.T) {
	e := New("test error")
	var t1 testing.T
	e.ErrorNotContains(&t1, "test")
	assert.True(t, t1.Failed())
}

func TestErrorIsNil(t *testing.T) {
	var e *Error = nil
	e.ErrorIsNil(t)
}

func TestErrorIsNilFail(t *testing.T) {
	e := New("test error")
	var t1 testing.T
	e.ErrorIsNil(&t1)
	assert.True(t, t1.Failed())
}

func TestErrorIsNilFatal(t *testing.T) {
	var e *Error = nil
	e.ErrorIsNilFatal(t)
}

func TestErrorNotNilFatal(t *testing.T) {
	var e *Error = New("test error")
	e.ErrorNotNilFatal(t)
}

func TestErrorAsError(t *testing.T) {
	var e *Error
	assert.Nil(t, e)
	assert.NoError(
		t,
		e.AsError(),
	)
}

func TestErrorAsErrorNotNil(t *testing.T) {
	var e *Error = New("test error")
	assert.Error(
		t,
		e.AsError(),
	)
}

func TestErrorNotNill(t *testing.T) {
	e := New("test error")
	e.ErrorNotNil(t)
}

func TestErrorNotNillFail(t *testing.T) {
	var e *Error = nil
	var t1 testing.T
	e.ErrorNotNil(&t1)
	assert.True(t, t1.Failed())
}
