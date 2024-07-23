API

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

## The caller uses this API by

```go
package main
func main() {
	h := mqueue.NewEventBridgeHandler("bus-name", mqueue.MKAWSConfig) // or mqueue.FSAWSConfig, mqueue.LocalStackConfig


	client := mqueue.New("facility-service", h)
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
