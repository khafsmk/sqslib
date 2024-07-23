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
	SquadName   string
	ServiceName string
	Domain      string
	Handler     Handler

	timeNow func() time.Time
	newUUID func() string
}

func (c *Client) publish(ctx context.Context, input any) error {
	if ctx == nil {
		ctx = context.Background()
	}
	record := c.newRecord(ctx, c.ServiceName, input)
	return c.Handler.Handle(ctx, record)
}

// Publish publishes an event to the backend queue.
func (c *Client) Publish(ctx context.Context, input any) error {
	return c.publish(ctx, input)
}

var defaultClient atomic.Pointer[Client]

func init() {
	defaultClient.Store(New(NewJSONHandler(io.Discard)))
}

type ClientOption struct {
	SquadName   string
	ServiceName string
	Domain      string
}

type ClientOptionFn func(*ClientOption)

// WithSquadName sets the squad name for the client.
func WithSquadName(name string) ClientOptionFn {
	return func(o *ClientOption) {
		o.SquadName = name
	}
}

// WithServiceName sets the service name for the client.
func WithServiceName(name string) ClientOptionFn {
	return func(o *ClientOption) {
		o.ServiceName = name
	}
}

// WithDomain sets the domain for the client.
func WithDomain(domain string) ClientOptionFn {
	return func(o *ClientOption) {
		o.Domain = domain
	}
}

// New initializes a new client with the given publisher.
func New(h Handler, fns ...ClientOptionFn) *Client {
	if h == nil {
		panic("nil handler")
	}

	conf := &ClientOption{}
	for _, fn := range fns {
		fn(conf)
	}
	return &Client{
		SquadName:   conf.ServiceName,
		ServiceName: conf.ServiceName,
		Domain:      conf.Domain,
		Handler:     h,
		timeNow:     time.Now,
		newUUID:     uuid.NewString,
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
