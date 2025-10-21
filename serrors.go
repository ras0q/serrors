package serrors

import (
	"errors"
	"strings"
)

type LogAttrsProvider interface {
	LogAttrs() []any
}

type LogAttrsProviderError interface {
	error
	LogAttrsProvider
}

type structuredError struct {
	Cause error
	Msg   string
	Args  []any
}

var _ LogAttrsProviderError = structuredError{}
var _ LogAttrsProviderError = (*structuredError)(nil)

func (s structuredError) Error() string {
	var buf strings.Builder
	buf.WriteString(s.Msg)

	if s.Cause != nil {
		buf.WriteString(": ")
		buf.WriteString(s.Cause.Error())
	}

	return buf.String()
}

func (s structuredError) Unwrap() error {
	return s.Cause
}

// LogAttrs implements LogAttrsProvider
func (s structuredError) LogAttrs() []any {
	args := s.Args

	for err := s.Cause; err != nil; err = errors.Unwrap(err) {
		provider, ok := err.(LogAttrsProviderError)
		if !ok {
			break
		}

		args = append(args, provider.LogAttrs()...)
	}

	return args
}

func New(msg string, args ...any) error {
	return Wrap(nil, msg, args...)
}

func Wrap(cause error, msg string, args ...any) error {
	return &structuredError{
		Cause: cause,
		Msg:   msg,
		Args:  args,
	}
}
