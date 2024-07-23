package mqueue

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/smithy-go/ptr"
)

// NewSQSHandler returns a new SQS handler.
// It's better to allow the client to test by exposing the eventbridge.Options
// testing by HTTPClient is better than using extra libraries for mocking it.
// This is also good for using with localstack.
func NewSQSHandler(queueURL string, cfg aws.Config, optFns ...func(*sqs.Options)) *SQSHandler {
	return &SQSHandler{
		QueueURL: queueURL,
		client:   sqs.NewFromConfig(cfg, optFns...),
	}
}

// SQSHandler sends records to an AWS SQS queue.
type SQSHandler struct {
	QueueURL string
	client   *sqs.Client
}

// Handle sends the record to AWS SQS.
// it implements the Handler interface.
func (h *SQSHandler) Handle(ctx context.Context, record Record) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(record)
	if err != nil {
		return err
	}
	_, err = h.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    ptr.String(h.QueueURL),
		MessageBody: ptr.String(buf.String()),
	})
	return err
}
