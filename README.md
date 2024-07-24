# A quick demo

```go
package main

func main() {

	var (
		squadName   = "squad"
		serviceName = "service"
		domain      = "domain"
	)
	h := mqueue.NewEventBridgeHandler("bus-name", mqueue.MKAWSConfig) // or mqueue.LocalStackConfig
	client := mqueue.New(
		h,
		mqueue.WithSquadName(squadName), // optional
		mqueue.WithServiceName(serviceName), // optional
		mqueue.WithDomain(domain), // optional
	)

	input := map[string]string{"key": "value"}
	err := client.Publish(context.Background(), input)
	if err != nil {
		panic(err)
	}

	// or we can set the default client to avoid passing the client around
	mqueue.SetDefault(client)
	mqueue.Publish(input)
	// or with context
	mqueue.PublishContext(context.Background(), input)
}
```

# Design

The `mqueue` package contains three main types:

1. `Client` is the frontend, providing output methods like `Publish`.
2. Each call to a `Client` output method creates a `Record`.
3. The `Record` is passed to a `Handler` for output.

```
                        +-------------+
                        |   Record    |
+---------------+       +-------------+        +-------------+
|               |                              |             |
|     Client    +----------------------------->|   Handler   |
|               |                              |             |
+---------------+                              +-------------+
```

## The Handler

A Handler describes the queue backend. It handles log records produced by a
Client.

```go
// Handler defines the interface for sending records to a backend.
type Handler interface {
	// Handle sends the record to the backend.
	Handle(ctx context.Context, record Record) error
}
```

## The record

A Record holds information about a message event.

```go
type Record struct {
	Source             string    `json:"source"`
	Destination        string    `json:"destination,omitempty"`
	Time               time.Time `json:"time"`
	TraceID            string    `json:"trace_id,omitempty"`
	IdempotencyKey     string    `json:"idempotency_key,omitempty"`
	SequenceID         string    `json:"sequence_id,omitempty"`
	TenantID           string    `json:"tenant_id,omitempty"`
	DataClassification string    `json:"data_classification,omitempty"`
	EventName          string    `json:"event_name,omitempty"`
	Data               any       `json:"data"`
}
```

## Provided Handlers

```go
// It's better to allow the client to test by exposing the aws.Config
// testing by HTTPClient is better than using extra libraries for mocking it.
// This is also good for using with unit testing and development environment such as local stack
// See the all_test.go
type EventBridgeHandler struct{ ... }
	func NewEventBridgeHandler(busName string, config aws.Config, optFns ...func(*eventbridge.Options)) *EventBridgeHandler

type KinesisHandler struct{ ... }
	func NewKinesisHandler(streamName string, config aws.Config, optFns ...func(*kinesis.Options)) *KinesisHandler

type SQSHandler struct{ ... }
	func NewSQSHandler(queueURL string, config aws.Config, optFns ...func(*sqs.Options)) *SQSHandler
```

## Preconfiguration environment

```go
var FSAWSConfig = func() aws.Config { ... }()
var LocalStackConfig = func() aws.Config { ... }()
var MKAWSConfig = func() aws.Config { ... }()
var MSMConfig = func() aws.Config { ... }()
```

## Testing

It's simple to implements the Handler for testing. You can turn the Client into
the fake struct without mocking it.

Source: (fakes vs mocks by Martin
Fowler)[https://martinfowler.com/articles/mocksArentStubs.html]

Fake: objects actually have working implementations, but usually take some
shortcut which makes them not suitable for production (an in memory database is
a good example).

Mocks: are what we are talking about here: objects pre-programmed with
expectations which form a specification of the calls they are expected to
receive.

```go
package main

import (
	mq "github.com/khafsmk/mqueue"
)

func TestClient(t *testing.T) {
	client := &mq.Client{
		Handler: mq.HandlerFunc(func(ctx context.Context, record mq.Record) error {
			return nil
		}),
	}
	err := client.Publish(context.Background(), map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
```

# Appendix: API

```
package mqueue // import "github.com/khafsmk/mqueue"

var FSAWSConfig = func() aws.Config { ... }()
var LocalStackConfig = func() aws.Config { ... }()
var MKAWSConfig = func() aws.Config { ... }()
var MSMConfig = func() aws.Config { ... }()


func Publish(input any) error
func PublishContext(ctx context.Context, input any) error
func SetDefault(c *Client)


type Record struct{ ... }

type Client struct{ ... }
	func Default() *Client
	func New(name string, h Handler) *Client

type Handler interface{ ... }

// It's better to allow the client to test by exposing the aws.Config
// testing by HTTPClient is better than using extra libraries for mocking it.
// This is also good for using with unit testing and development environment such as local stack
// See the all_test.go
type EventBridgeHandler struct{ ... }
	func NewEventBridgeHandler(busName string, config aws.Config, optFns ...func(*eventbridge.Options)) *EventBridgeHandler

type KinesisHandler struct{ ... }
	func NewKinesisHandler(streamName string, config aws.Config, optFns ...func(*kinesis.Options)) *KinesisHandler

type SQSHandler struct{ ... }
	func NewSQSHandler(queueURL string, config aws.Config, optFns ...func(*sqs.Options)) *SQSHandler

type JSONHandler struct{ ... }
	func NewJSONHandler(w io.Writer) *JSONHandler

func NewFanOutHandlers(handlers ...Handler) *MultiHandler
func NewSequenceHandlers(handlers ...Handler) *MultiHandler
```

