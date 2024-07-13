
API 

```
package mqueue // import "github.com/khafsmk/mqueue"

func Publish(input any) error
func PublishContext(ctx context.Context, input any) error
func SetDefault(c *Client)
func NewJSONHandler(w io.Writer) *jsonHandler
func NewFanOutHandlers(handlers ...Handler) *multiHandler
func NewSequenceHandlers(handlers ...Handler) *multiHandler
type Client struct{ ... }
    func Default() *Client
    func New(name string, p Handler) *Client
type EventBridgeHandler struct{ ... }
    func NewEventBridgeHandler(client *eventbridge.Client) *EventBridgeHandler
type Handler interface{ ... }
type KinesisHandler struct{ ... }
    func NewKinesisHandler(options KinesisOption, kc *kinesis.Client) *KinesisHandler
type KinesisOption struct{ ... }
type Record struct{ ... }
```