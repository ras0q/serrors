package serrors

import (
	"strings"
)

type SError struct {
	Cause error
	Msg   string
	Args  []any
}

func (s SError) Error() string {
	var buf strings.Builder
	buf.WriteString(s.Msg)

	if s.Cause != nil {
		buf.WriteString(": ")
		buf.WriteString(s.Msg)
	}

	return buf.String()
}

func (s SError) Unwrap() error {
	return s.Cause
}

func New(msg string, args ...any) error {
	return Wrap(nil, msg, args...)
}

func Wrap(cause error, msg string, args ...any) error {
	return &SError{
		Cause: cause,
		Msg:   msg,
		Args:  args,
	}
}
