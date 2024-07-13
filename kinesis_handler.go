package mqueue

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/smithy-go/ptr"
)

// KinesisOption is the configuration for the Kinesis handler.
type KinesisOption struct {
	StreamName string
}

// NewKinesisHandler returns a new Kinesis handler.
func NewKinesisHandler(options KinesisOption, kc *kinesis.Client) *KinesisHandler {
	return &KinesisHandler{
		kc: kc,
	}
}

// KinesisHandler sends records to a Kinesis stream.
type KinesisHandler struct {
	kc      *kinesis.Client
	options KinesisOption
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
		StreamName:   ptr.String(h.options.StreamName),
	})
	return err
}
