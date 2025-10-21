package serrors

import (
	"context"
	"log/slog"
)

type Handler struct {
	slog.Handler
}

var _ slog.Handler = (*Handler)(nil)

func NewHandler(h slog.Handler) *Handler {
	return &Handler{h}
}

// Handle implements slog.Handler.
func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	record.Attrs(func(a slog.Attr) bool {
		if a.Key != "error" {
			return true
		}

		v := a.Value.Any()
		if v == nil {
			return true
		}

		err, ok := v.(LogAttrsProviderError)
		if !ok {
			return true
		}

		record.Add(err.LogAttrs()...)

		return true
	})

	return h.Handler.Handle(ctx, record)
}
