# serrors

Structured errors that work well with `slog`. An error just implements `slog.LogValuer`, so an ordinary handler prints its message, attributes, and cause.

This follows the idea in Go’s [Error handling and Go](https://go.dev/blog/error-handling-and-go): a simple message for callers and users (`Error()`), richer detail for operators (`slog` attributes on the same error).

```go
package main

import (
	// ...
	"log/slog"
	"github.com/ras0q/serrors"
)

func main() {
	if err := doSomething(); err != nil {
		// developer-facing error
		// Output: 2000/01/01 00:00:00 ERROR doSomething error.msg="open file" error.filename=notfound.txt error.userID=100 error.cause="open notfound.txt: no such file or directory"
		slog.Error("doSomething", "error", err)

		// user-facing error
		// Output: doSomething: open file: open notfound.txt: no such file or directory
		fmt.Printf("doSomething: %+v\n", err)
	}
}

func doSomething() error {
	filename := "notfound.txt"
	userID := 100

	_, err := os.Open(filename)
	if err != nil {
		return serrors.Wrap(
			err, "open file",
			// add key-value attributes (slog-compatible!)
			"filename", filename,
			slog.Int("userID", userID),
			// ...
		)
	}

	return nil
}
```

## Frame

To record where an error was wrapped, pass `serrors.Frame()` to `Wrap` or `New`. It adds a `frame` attribute (function, file, line) in logs only; it is not part of `err.Error()`.

For every error in the process, call `serrors.SetAutoFrame(true)` once (for example from `main`). `New` and `Wrap` then attach a frame automatically unless you already passed `Frame()` on that call.
