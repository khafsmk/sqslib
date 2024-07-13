package mqueue

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/smithy-go/ptr"
)

type KinesisOption struct {
	StreamName string
}

func NewKinesisHandler(options KinesisOption, kc *kinesis.Client) *KinesisHandler {
	return &KinesisHandler{
		kc: kc,
	}
}

type KinesisHandler struct {
	kc      *kinesis.Client
	options KinesisOption
}

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
