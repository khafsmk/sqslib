package mqueue

import (
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

func TestMultiHandlers(t *testing.T) {
	ctx := context.Background()
	check := checker(t)
	kc := NewKinesisHandler(KinesisOption{StreamName: "mercury"}, kinesis.New(kinesis.Options{
		HTTPClient: nopClient,
		Region:     "us-west-2",
	}))
	eb := NewEventBridgeHandler(eventbridge.New(eventbridge.Options{
		HTTPClient: nopClient,
		Region:     "us-west-2",
	}))
	jh := NewJSONHandler(io.Discard)

	c := NewSequenceHandlers(kc, eb, jh)
	err := c.Handle(ctx, Record{})
	check(err)

	c = NewFanOutHandlers(kc, eb, jh)
	err = c.Handle(ctx, Record{})
	check(err)
}
