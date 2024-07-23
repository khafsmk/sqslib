package mqueue

import (
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestMultiHandlers(t *testing.T) {
	ctx := context.Background()
	check := checker(t)
	kc := NewKinesisHandler("stream-name", aws.Config{
		HTTPClient: nopClient,
		Region:     "us-west-2",
	})
	eb := NewEventBridgeHandler("bus-name", aws.Config{
		HTTPClient: nopClient,
		Region:     "us-west-2",
	})

	sq := NewSQSHandler("queue-url", aws.Config{
		HTTPClient: nopClient,
		Region:     "us-west-2",
	})
	jh := NewJSONHandler(io.Discard)

	c := NewSequenceHandlers(kc, eb, jh, sq)
	err := c.Handle(ctx, Record{})
	check(err)

	c = NewFanOutHandlers(kc, eb, jh)
	err = c.Handle(ctx, Record{})
	check(err)
}
