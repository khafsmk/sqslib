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
func NewFanOutHandlers(handlers ...Handler) *multiHandler
func NewSequenceHandlers(handlers ...Handler) *multiHandler

type Client struct{ ... }
    func Default() *Client
    func New(name string, h Handler) *Client
type EventBridgeHandler struct{ ... }
    func NewEventBridgeHandler(busname string, config aws.Config, optFns ...func(*eventbridge.Options)) *EventBridgeHandler
type Handler interface{ ... }
type JSONHandler struct{ ... }
    func NewJSONHandler(w io.Writer) *JSONHandler
type KinesisHandler struct{ ... }
    func NewKinesisHandler(streamName string, config aws.Config, optFns ...func(*kinesis.Options)) *KinesisHandler
type Record struct{ ... }
type SQSHandler struct{ ... }
    func NewSQSHandler(queueURL string, cfg aws.Config, optFns ...func(*sqs.Options)) *SQSHandler

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
