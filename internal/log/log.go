package log

import (
	"fmt"

	"go-media-stream/internal/customerrors"
)

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func (l *Logger) Error(err error) {
	fmt.Println(customerrors.NewWithError(err))
}
