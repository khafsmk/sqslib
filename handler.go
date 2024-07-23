package mqueue

import (
	"context"
)

// Handler defines the interface for sending records to a backend.
type Handler interface {
	// Handle sends the record to the backend.
	Handle(ctx context.Context, record Record) error
}

// HandlerFunc is an adapter to allow the use of ordinary functions as handlers.
type HandlerFunc func(context.Context, Record) error

// Handle implements the Handler interface.
func (f HandlerFunc) Handle(ctx context.Context, record Record) error {
	return f(ctx, record)
}
