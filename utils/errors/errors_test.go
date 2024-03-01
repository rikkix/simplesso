package errors_test

import (
	"testing"

	"github.com/rikkix/simplesso/utils/errors"
)

func TestError(t *testing.T) {
	msg1 := "error 1"
	msg2 := "error 2"

	var err errors.TraceableError
	err = errors.New(msg1)
	if err.Error() != msg1 {
		t.Errorf("Error() failed. Expected: %s, got: %s", msg1, err.Error())
	}

	err = err.From(msg2)
	if err.Error() != msg2 + " > " + msg1 {
		t.Errorf("From() failed. Expected: %s, got: %s", msg2 + " > " + msg1, err.Error())
	}
}