package serrors_test

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"
	"testing/synctest"

	"github.com/ras0q/serrors"
)

func Test_Handler(t *testing.T) {
	t.Parallel()

	type testCase struct {
		err      error
		expected string
	}

	run := func(t *testing.T, tc testCase) {
		synctest.Test(t, func(t *testing.T) {
			var buf bytes.Buffer
			base := slog.NewJSONHandler(&buf, nil)
			wrapped := serrors.NewHandler(base)
			logger := slog.New(wrapped)

			logger.ErrorContext(t.Context(), "message", "error", tc.err)

			actual := buf.String()
			if tc.expected != actual {
				t.Fatalf("output mismatch!\nexpected:\n\t%s\nactual:\n\t%s", tc.expected, actual)
			}
		})
	}

	makeJSONLog := func(extra string) string {
		return `{"time":"2000-01-01T09:00:00+09:00","level":"ERROR","msg":"message",` + extra + "}\n"
	}

	testCases := map[string]testCase{
		"no error": {
			err:      nil,
			expected: makeJSONLog(`"error":null`),
		},
		"plain error": {
			err:      errors.New("plain error"),
			expected: makeJSONLog(`"error":"plain error"`),
		},
		"serrors.New": {
			err:      serrors.New("api error", "statusCode", 500),
			expected: makeJSONLog(`"error":"api error","statusCode":500`),
		},
		"serrors.New with slog.Attr": {
			err:      serrors.New("api error", slog.Int("statusCode", 500)),
			expected: makeJSONLog(`"error":"api error","statusCode":500`),
		},
		"serrors.Wrap": {
			err: serrors.Wrap(
				errors.New("timeout"),
				"failed to fetch data",
				"endpoint", "/api/data",
			),
			expected: makeJSONLog(`"error":"failed to fetch data: timeout","endpoint":"/api/data"`),
		},
		"serrors.Wrap with slog.Attr": {
			err: serrors.Wrap(
				errors.New("timeout"),
				"failed to fetch data",
				slog.String("endpoint", "/api/data"),
			),
			expected: makeJSONLog(`"error":"failed to fetch data: timeout","endpoint":"/api/data"`),
		},
		"serrors.Wrap nil": {
			err:      serrors.Wrap(nil, "no error occurred", "info", "none"),
			expected: makeJSONLog(`"error":"no error occurred","info":"none"`),
		},
		"multiple wrapping": {
			err: serrors.Wrap(
				serrors.Wrap(
					errors.New("connection refused"),
					"failed to connect",
					"host", "localhost",
					"port", 5432,
				),
				"database error",
				"dbName", "users",
			),
			expected: makeJSONLog(`"error":"database error: failed to connect: connection refused","dbName":"users","host":"localhost","port":5432`),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			run(t, tc)
		})
	}
}
