package serrors

import (
	"log/slog"
	"runtime"
)

// Frame returns a [slog.Attr] with the source location where it is evaluated.
// Use it as an argument to [New] or [Wrap] to record that call site in logs.
func Frame() slog.Attr {
	var pcs [1]uintptr
	if runtime.Callers(2, pcs[:]) == 0 {
		return slog.Group("frame")
	}

	frame, _ := runtime.CallersFrames(pcs[:]).Next()

	return slog.Any("frame", &slog.Source{
		Function: frame.Function,
		File:     frame.File,
		Line:     frame.Line,
	})
}
