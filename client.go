package mqueue

import (
	"context"
	"io"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// Client is a frontend that callers will use to publish events.
type Client struct {
	Name    string
	handler Handler
	timeNow func() time.Time
	newUUID func() string
}

func (c *Client) publish(ctx context.Context, input any) error {
	if ctx == nil {
		ctx = context.Background()
	}
	record := c.newRecord(ctx, c.Name, input)
	return c.handler.Handle(ctx, record)
}

// Publish publishes an event to the backend queue.
func (c *Client) Publish(ctx context.Context, input any) error {
	return c.publish(ctx, input)
}

var defaultClient atomic.Pointer[Client]

const noName = "no-name"

func init() {
	defaultClient.Store(New(noName, NewJSONHandler(io.Discard)))
}

// New initializes a new client with the given publisher.
func New(name string, p Handler) *Client {
	if p == nil {
		panic("nil publisher")
	}
	return &Client{
		Name:    name,
		handler: p,
		timeNow: time.Now,
		newUUID: uuid.NewString,
	}
}

// Default returns the default publisher.
func Default() *Client { return defaultClient.Load() }

// SetDefault sets the default publisher.
func SetDefault(c *Client) {
	defaultClient.Store(c)
}

// Publish writes an event to the default publisher.
func Publish(input any) error {
	return Default().publish(context.Background(), input)
}

// PublishContext publishes an event to the default publisher with context.
func PublishContext(ctx context.Context, input any) error {
	return Default().publish(ctx, input)
}
