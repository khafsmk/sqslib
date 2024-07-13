package mqueue

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

// Handler is the interface that wraps the basic Handle method.
func NewSequenceHandlers(handlers ...Handler) *multiHandler {
	return &multiHandler{
		handlers: handlers,
	}
}

// NewFanInHandlers returns a new handler that will fan-in the records to the given handlers.
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
