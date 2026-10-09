package serrors

import (
	"log/slog"
	"strings"
)

type structuredError struct {
	Cause error
	Msg   string
	Attrs []slog.Attr
}

var (
	_ error          = structuredError{}
	_ slog.LogValuer = structuredError{}
)

// Error implements [error].
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

// LogValue implements [slog.LogValuer].
func (s structuredError) LogValue() slog.Value {
	attrs := make([]slog.Attr, 0, 1+len(s.Attrs)+1)
	attrs = append(attrs, slog.String("msg", s.Msg))
	attrs = append(attrs, s.Attrs...)
	if s.Cause != nil {
		attrs = append(attrs, slog.Any("cause", s.Cause))
	}

	return slog.GroupValue(attrs...)
}

// New returns an error whose message and args are logged via [slog.LogValuer].
func New(msg string, args ...any) error {
	return wrap(nil, msg, args)
}

// Wrap returns an error whose message, args, and cause are logged via [slog.LogValuer].
// args use slog's key-value rules.
func Wrap(cause error, msg string, args ...any) error {
	return wrap(cause, msg, args)
}

func wrap(cause error, msg string, args []any) error {
	var record slog.Record
	record.Add(args...)
	attrs := make([]slog.Attr, 0, record.NumAttrs())
	record.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})
	if autoFrame.Load() && !attrsContainFrame(attrs) {
		attrs = append(attrs, frameAttr(4))
	}

	return &structuredError{
		Cause: cause,
		Msg:   msg,
		Attrs: attrs,
	}
}
