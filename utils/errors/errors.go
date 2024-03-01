package errors

import "strings"

// TraceableError can help trace the origin of an error.
type TraceableError interface {
	error
	From(message string) TraceableError
	Last() string
}

// TError is a struct that implements the TraceableError interface.
type TError struct {
	Message []string
}

// Implement the error interface.
func (e *TError) Error() string {
	return strings.Join(e.Message, " > ")
}

func New(message string) *TError {
	return &TError{[]string{message}}
}

func (e *TError) From(message string) TraceableError {
	nmsg := append([]string{message}, e.Message...)
	return &TError{nmsg}
}

func (e *TError) Last() string {
	return e.Message[0]
}