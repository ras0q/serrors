package serrors_test

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"
	"testing/synctest"

	"github.com/ras0q/serrors"
)

func TestLogValue(t *testing.T) {
	t.Parallel()

	type testCase struct {
		err      error
		expected string
	}

	run := func(t *testing.T, tc testCase) {
		synctest.Test(t, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buf, nil))

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
			expected: makeJSONLog(`"error":{"msg":"api error","statusCode":500}`),
		},
		"serrors.New with slog.Attr": {
			err: serrors.New(
				"api error",
				"userID", 12345,
				slog.Group("request",
					slog.String("method", "GET"),
					slog.String("url", "/api/data"),
				),
				slog.Group("response",
					slog.Int("statusCode", 500),
					"status", "Internal Server Error",
				),
			),
			expected: makeJSONLog(
				`"error":{"msg":"api error","userID":12345,"request":{"method":"GET","url":"/api/data"},"response":{"statusCode":500,"status":"Internal Server Error"}}`,
			),
		},
		"serrors.Wrap": {
			err: serrors.Wrap(
				errors.New("timeout"),
				"failed to fetch data",
				"endpoint", "/api/data",
			),
			expected: makeJSONLog(`"error":{"msg":"failed to fetch data","endpoint":"/api/data","cause":"timeout"}`),
		},
		"serrors.Wrap with slog.Attr": {
			err: serrors.Wrap(
				errors.New("timeout"),
				"failed to fetch data",
				slog.String("endpoint", "/api/data"),
			),
			expected: makeJSONLog(`"error":{"msg":"failed to fetch data","endpoint":"/api/data","cause":"timeout"}`),
		},
		"serrors.Wrap nil": {
			err:      serrors.Wrap(nil, "no error occurred", "info", "none"),
			expected: makeJSONLog(`"error":{"msg":"no error occurred","info":"none"}`),
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
			expected: makeJSONLog(
				`"error":{"msg":"database error","dbName":"users","cause":{"msg":"failed to connect","host":"localhost","port":5432,"cause":"connection refused"}}`,
			),
		},
		"each layer keeps its own keys": {
			err: serrors.Wrap(
				serrors.New("low-level error", "code", 1001, "detail", "disk full"),
				"high-level error",
				"code", 2002,
				"detail", "network issue",
			),
			expected: makeJSONLog(
				`"error":{"msg":"high-level error","code":2002,"detail":"network issue","cause":{"msg":"low-level error","code":1001,"detail":"disk full"}}`,
			),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			run(t, tc)
		})
	}
}
