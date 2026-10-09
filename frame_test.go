package serrors_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/ras0q/serrors"
)

func TestFrame(t *testing.T) {
	t.Parallel()

	err := newErrorWithFrame()

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	logger.Error("event", "error", err)

	var logLine struct {
		Error struct {
			Frame struct {
				Function string `json:"function"`
				File     string `json:"file"`
				Line     int    `json:"line"`
			} `json:"frame"`
		} `json:"error"`
	}
	if err := json.Unmarshal(buf.Bytes(), &logLine); err != nil {
		t.Fatalf("unmarshal log: %v\nlog: %s", err, buf.String())
	}

	if logLine.Error.Frame.Line == 0 {
		t.Fatalf("frame line is zero: %+v", logLine.Error.Frame)
	}
	if !strings.HasSuffix(logLine.Error.Frame.File, "frame_test.go") {
		t.Fatalf("unexpected file: %q", logLine.Error.Frame.File)
	}
	if !strings.Contains(logLine.Error.Frame.Function, "newErrorWithFrame") {
		t.Fatalf("unexpected function: %q", logLine.Error.Frame.Function)
	}
}

func newErrorWithFrame() error {
	return serrors.New("wrapped", serrors.Frame())
}

func TestAutoFrame(t *testing.T) {
	serrors.SetAutoFrame(true)
	t.Cleanup(func() { serrors.SetAutoFrame(false) })

	err := serrors.New("wrapped")

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	logger.Error("event", "error", err)

	var logLine struct {
		Error struct {
			Frame struct {
				Function string `json:"function"`
				File     string `json:"file"`
				Line     int    `json:"line"`
			} `json:"frame"`
		} `json:"error"`
	}
	if err := json.Unmarshal(buf.Bytes(), &logLine); err != nil {
		t.Fatalf("unmarshal log: %v\nlog: %s", err, buf.String())
	}

	if logLine.Error.Frame.Line == 0 {
		t.Fatalf("frame line is zero: %+v", logLine.Error.Frame)
	}
	if !strings.Contains(logLine.Error.Frame.Function, "TestAutoFrame") {
		t.Fatalf("unexpected function: %q", logLine.Error.Frame.Function)
	}
}

func TestAutoFrame_skipsWhenFramePresent(t *testing.T) {
	serrors.SetAutoFrame(true)
	t.Cleanup(func() { serrors.SetAutoFrame(false) })

	err := serrors.New("wrapped", serrors.Frame())

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	logger.Error("event", "error", err)

	var logLine struct {
		Error struct {
			Frame struct {
				Function string `json:"function"`
			} `json:"frame"`
		} `json:"error"`
	}
	if err := json.Unmarshal(buf.Bytes(), &logLine); err != nil {
		t.Fatalf("unmarshal log: %v\nlog: %s", err, buf.String())
	}

	if !strings.Contains(logLine.Error.Frame.Function, "TestAutoFrame_skipsWhenFramePresent") {
		t.Fatalf("unexpected function: %q", logLine.Error.Frame.Function)
	}
}
