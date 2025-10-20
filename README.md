# serrors

Structured errors that work well with `slog.Error(...)`.

```go
package main

import (
    // ...
    "log/slog"
	"github.com/ras0q/serrors"
)

func main() {
	base := slog.NewJSONHandler(os.Stderr, nil)
	wrapped := serrors.NewHandler(base)
	logger := slog.New(wrapped)

	// Set the global default logger for the process
	slog.SetDefault(logger)

	// run your application...
    ctx := context.Background()
	if err := doSomething(ctx); err != nil {
        slog.ErrorContext(ctx, "doSomething", "error", err)
    }
}

func doSomething(ctx context.Context) error {
    filename := "notfound.txt"

    _, err := os.Open(filename)
    if err != nil {
        return serrors.Wrap(
            err, "open file",
            // add key-value attributes (slog-compatible!)
            "filename", filename,
            slog.String("userID", "001")
            // ...
        )
    }

    return nil
}
```
