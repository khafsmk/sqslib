package mqueue

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

func TestMultiHandlers(t *testing.T) {
	ctx := context.Background()
	check := checker(t)
	kc := kinesis.New(kinesis.Options{
		HTTPClient: nopClient,
		Region:     "us-west-2",
	})
	eb := eventbridge.New(eventbridge.Options{
		HTTPClient: nopClient,
		Region:     "us-west-2",
	})

	c := NewSequenceHandlers(NewKinesisHandler(KinesisOption{StreamName: "mercury"}, kc), NewEventBridgeHandler(eb))
	err := c.Handle(ctx, Record{})
	check(err)

	c = NewFanOutHandlers(NewKinesisHandler(KinesisOption{StreamName: "mercury"}, kc), NewKinesisHandler(KinesisOption{StreamName: "venus"}, kc))
	err = c.Handle(ctx, Record{})
	check(err)
}
