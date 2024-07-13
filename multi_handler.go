package mqueue

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

// NewSequenceHandlers returns the new handler that will emit sequently to multiple handlers.
func NewSequenceHandlers(handlers ...Handler) *multiHandler {
	return &multiHandler{
		handlers: handlers,
	}
}

// NewFanOutHandlers returns a new handler that will fan-out the records to multiple handlers.
func NewFanOutHandlers(handlers ...Handler) *multiHandler {
	return &multiHandler{
		handlers:    handlers,
		concurrency: true,
	}
}

type multiHandler struct {
	handlers    []Handler
	concurrency bool
}

// Handle publish its argument [Record] to multiple handlers.
func (h *multiHandler) Handle(ctx context.Context, record Record) error {
	if h.concurrency {
		return h.handleConcurrency(ctx, record)
	}
	return h.handleLinear(ctx, record)
}

func (h *multiHandler) handleLinear(ctx context.Context, record Record) error {
	var errs []error
	for _, fn := range h.handlers {
		fn := fn
		if err := fn.Handle(ctx, record); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (h *multiHandler) handleConcurrency(ctx context.Context, record Record) error {
	var g errgroup.Group
	for _, fn := range h.handlers {
		fn := fn
		g.Go(func() error {
			return fn.Handle(ctx, record)
		})
	}
	return g.Wait()
}
