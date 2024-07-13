package mqueue

import (
	"context"
)

// Handler defines the interface for sending records to a backend.
type Handler interface {
	// Handle sends the record to the backend.
	Handle(ctx context.Context, record Record) error
}
