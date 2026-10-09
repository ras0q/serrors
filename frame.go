package serrors

import (
	"log/slog"
	"runtime"
	"sync/atomic"
)

var autoFrame atomic.Bool

// SetAutoFrame configures whether [New] and [Wrap] attach a frame attribute for
// their call site on every error they create. The setting applies process-wide.
func SetAutoFrame(enabled bool) {
	autoFrame.Store(enabled)
}

// Frame returns a [slog.Attr] with the source location where it is evaluated.
// Use it as an argument to [New] or [Wrap] to record that call site in logs.
func Frame() slog.Attr {
	return frameAttr(3)
}

func frameAttr(skip int) slog.Attr {
	var pcs [1]uintptr
	if runtime.Callers(skip, pcs[:]) == 0 {
		return slog.Group("frame")
	}

	frame, _ := runtime.CallersFrames(pcs[:]).Next()

	return slog.Any("frame", &slog.Source{
		Function: frame.Function,
		File:     frame.File,
		Line:     frame.Line,
	})
}

func attrsContainFrame(attrs []slog.Attr) bool {
	for _, a := range attrs {
		if a.Key == "frame" {
			return true
		}
	}
	return false
}
