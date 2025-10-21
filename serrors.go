package serrors

import (
	"context"
	"errors"
	"log/slog"
	"strings"
)

type LogAttrsProvider interface {
	LogAttrs() []any
}

type LogAttrsProviderError interface {
	error
	LogAttrsProvider
}

type Handler struct {
	slog.Handler
}

var _ slog.Handler = (*Handler)(nil)

func NewHandler(h slog.Handler) *Handler {
	return &Handler{h}
}

// Handle implements slog.Handler.
func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	record.Attrs(func(a slog.Attr) bool {
		if a.Key != "error" {
			return true
		}

		v := a.Value.Any()
		if v == nil {
			return true
		}

		err, ok := v.(LogAttrsProviderError)
		if !ok {
			return true
		}

		record.Add(err.LogAttrs()...)

		return true
	})

	return h.Handler.Handle(ctx, record)
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
