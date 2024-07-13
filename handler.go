package mqueue

import (
	"context"
)

// Handler is the interface that wraps the basic Publish method.
type Handler interface {
	Handle(ctx context.Context, record Record) error
}
