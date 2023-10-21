package logger

import (
	"context"
	"fmt"
	"log/slog"
)

type SLogTeeHandler struct {
	handler1 slog.Handler
	handler2 slog.Handler
}

func NewSLogTeeHandler(h1, h2 slog.Handler) slog.Handler {
	return &SLogTeeHandler{
		handler1: h1,
		handler2: h2,
	}
}

func (h *SLogTeeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	h1 := h.handler1.Enabled(ctx, level)
	h2 := h.handler2.Enabled(ctx, level)
	switch {
	case h1 && !h2:
		panicf("SLog handler 1 is enabled while handler 2 is not.")
	case h2 && !h1:
		panicf("SLog handler 2 is enabled while handler 1 is not.")
	}
	return h1 && h2
}

func (h *SLogTeeHandler) Handle(ctx context.Context, r slog.Record) (err error) {
	if err = h.handler1.Handle(ctx, r); err != nil {
		err = fmt.Errorf("slog handler 1 failed to handle error; %w", err)
		goto end
	}
	if err = h.handler2.Handle(ctx, r); err != nil {
		err = fmt.Errorf("slog handler 2 failed to handle error; %w", err)
		goto end
	}
end:
	return err
}

func (h *SLogTeeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewSLogTeeHandler(
		h.handler1.WithAttrs(attrs),
		h.handler2.WithAttrs(attrs),
	)
}

func (h *SLogTeeHandler) WithGroup(name string) slog.Handler {
	return NewSLogTeeHandler(
		h.handler1.WithGroup(name),
		h.handler2.WithGroup(name),
	)
}
