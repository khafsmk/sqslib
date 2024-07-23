package mqueue

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

// NewSequenceHandlers returns the new handler that will emit sequently to multiple handlers.
func NewSequenceHandlers(handlers ...Handler) *MultiHandler {
	return &MultiHandler{
		handlers: handlers,
	}
}

// NewFanOutHandlers returns a new handler that will fan-out the records to multiple handlers.
func NewFanOutHandlers(handlers ...Handler) *MultiHandler {
	return &MultiHandler{
		handlers:    handlers,
		concurrency: true,
	}
}

type MultiHandler struct {
	handlers    []Handler
	concurrency bool
}

// Handle publish its argument [Record] to multiple handlers.
func (h *MultiHandler) Handle(ctx context.Context, record Record) error {
	if h.concurrency {
		return h.handleConcurrency(ctx, record)
	}
	return h.handleLinear(ctx, record)
}

func (h *MultiHandler) handleLinear(ctx context.Context, record Record) error {
	var errs []error
	for _, fn := range h.handlers {
		fn := fn
		if err := fn.Handle(ctx, record); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (h *MultiHandler) handleConcurrency(ctx context.Context, record Record) error {
	var g errgroup.Group
	for _, fn := range h.handlers {
		fn := fn
		g.Go(func() error {
			return fn.Handle(ctx, record)
		})
	}
	return g.Wait()
}
