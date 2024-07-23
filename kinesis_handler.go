package mqueue

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/smithy-go/ptr"
)

// NewKinesisHandler returns a new Kinesis handler.
// It's better to allow the client to test by exposing the eventbridge.Options
// testing by HTTPClient is better than using extra libraries for mocking it.
// This is also good for using with localstack.
func NewKinesisHandler(streamName string, config aws.Config, optFns ...func(*kinesis.Options)) *KinesisHandler {
	return &KinesisHandler{
		StreamName: streamName,
		kc:         kinesis.NewFromConfig(config, optFns...),
	}
}

// KinesisHandler sends records to a Kinesis stream.
type KinesisHandler struct {
	kc         *kinesis.Client
	StreamName string
}

// Handle sends the record to the Kinesis stream.
func (h *KinesisHandler) Handle(ctx context.Context, record Record) error {
	buf, err := json.Marshal(record.Data)
	if err != nil {
		return err
	}
	_, err = h.kc.PutRecord(ctx, &kinesis.PutRecordInput{
		Data:         buf,
		PartitionKey: ptr.String(record.IdempotencyKey),
		StreamName:   ptr.String(h.StreamName),
	})
	return err
}
