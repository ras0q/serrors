package serrors_test

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/ras0q/serrors"
)

func Test_Handler(t *testing.T) {
	var buf bytes.Buffer
	base := slog.NewTextHandler(&buf, nil)
	wrapped := serrors.NewHandler(base)
	logger := slog.New(wrapped)

	err := serrors.Wrap(errors.New("orig"), "wrap", slog.String("k", "v"))
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	logger.ErrorContext(t.Context(), "message", "error", err)

	out := buf.String()
	if !strings.Contains(out, "k=v") {
		t.Fatalf("expected attribute in output; got: %s", out)
	}
}
