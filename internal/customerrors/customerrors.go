package customerrors

import (
	"fmt"
	"runtime"
)

type mainError struct {
	msg  string
	err  error
	file string
	line int
	fn   string
}

func (e *mainError) Error() string {
	message := "\033[31m"
	message += fmt.Sprintf("Message: %s\n", e.msg)
	message += fmt.Sprintf("Error: %v\n", e.err)
	message += fmt.Sprintf("File %s %d\n", e.file, e.line)
	message += fmt.Sprintf("Function %s\n", e.fn)
	message += "\033[0m"
	return message
}

func New(msg string) error {
	e := &mainError{msg: msg}
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return e
	}
	e.line = line
	e.file = file
	e.fn = runtime.FuncForPC(pc).Name()
	return e
}

func NewWithError(err error) error {
	e := &mainError{msg: "", err: err}
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return e
	}
	e.line = line
	e.file = file
	e.fn = runtime.FuncForPC(pc).Name()
	return e
}

func (e *mainError) Unwrap() error {
	return e.err
}
