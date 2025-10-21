package serrors_test

import (
	"log/slog"
	"os"
	"time"

	"github.com/ras0q/serrors"
)

func Example() {
	base := slog.NewJSONHandler(os.Stdout, &handlerOptsForTest)
	wrapped := serrors.NewHandler(base)
	logger := slog.New(wrapped)

	slog.SetDefault(logger)

	if err := doSomething(); err != nil {
		slog.Error("something failed", "error", err)
	}

	// Output: {"time":"0001-01-01T00:00:00Z","level":"ERROR","msg":"something failed","error":"open file: open notfound.txt: no such file or directory","filename":"notfound.txt","userID":"001"}
}

func doSomething() error {
	filename := "notfound.txt"

	_, err := os.Open(filename)
	if err != nil {
		return serrors.Wrap(
			err, "open file",
			// add key-value attributes (compatible to slog!)
			"filename", filename,
			slog.String("userID", "001"),
			// ...
		)
	}

	return nil
}

// NOTE: handlerOptsForTest is a utility to fix the time attribute during testing.
var handlerOptsForTest = slog.HandlerOptions{
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "time" {
			a.Value = slog.TimeValue(time.Time{})
		}

		return a
	},
}
