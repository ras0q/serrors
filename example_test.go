package serrors_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/ras0q/serrors"
)

func Example() {
	serrors.SetAutoFrame(true) // [optional] attach caller frame on every [serrors.Wrap] / [serrors.New]
	defer serrors.SetAutoFrame(false)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &_handlerOptsForTest))
	slog.SetDefault(logger)

	if err := doSomething(); err != nil {
		slog.Error("something failed", "error", err)
	}

	// Output: {"time":"0001-01-01T00:00:00Z","level":"ERROR","msg":"something failed","error":{"msg":"open file","filename":"notfound.txt","userID":"001","frame":{"function":"github.com/ras0q/serrors_test.doSomething","file":"example_test.go","line":36},"cause":"open notfound.txt: no such file or directory"}}
}

func doSomething() error {
	filename := "notfound.txt"

	_, err := os.Open(filename)
	if err != nil {
		return serrors.Wrap(
			err,                  // wrapeed error
			"open file",          // user-facing message
			"filename", filename, // slog-compatible key-value attribute (expanded)
			slog.String("userID", "001"), // slog-compatible key-value attribute (with slog.[Type])
			serrors.Frame(),              // [optional] caller frame
		)
	}

	return nil
}

// NOTE: _handlerOptsForTest fixes time and frame paths for stable Example output.
var _handlerOptsForTest = slog.HandlerOptions{
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "time" {
			a.Value = slog.TimeValue(time.Time{})
		}
		if a.Key == "frame" {
			if src, ok := a.Value.Any().(*slog.Source); ok {
				return slog.Any("frame", &slog.Source{
					Function: src.Function,
					File:     filepath.Base(src.File),
					Line:     src.Line,
				})
			}
		}

		return a
	},
}
