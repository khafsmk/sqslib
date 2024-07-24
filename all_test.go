package mqueue

import (
	"context"
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
	fake := HandlerFunc(func(ctx context.Context, r Record) error {
		return nil
	})

	c := NewSequenceHandlers(kc, eb, fake, sq)
	err := c.Handle(ctx, Record{})
	check(err)

	c = NewFanOutHandlers(kc, eb, fake)
	err = c.Handle(ctx, Record{})
	check(err)
}
