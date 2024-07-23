package mqueue

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/smithy-go/ptr"
)

// NewEventBridgeHandler returns a new EventBridge handler.
// It's better to allow the client to test by exposing the eventbridge.Options
// testing by HTTPClient is better than using extra libraries for mocking it.
// This is also good for using with localstack. See the all_test.go
func NewEventBridgeHandler(busName string, config aws.Config, optFns ...func(*eventbridge.Options)) *EventBridgeHandler {
	return &EventBridgeHandler{
		BusName: busName,
		client:  eventbridge.NewFromConfig(config, optFns...),
	}
}

// EventBridgeHandler is a handler that sends events to AWS EventBridge.
type EventBridgeHandler struct {
	BusName string
	client  *eventbridge.Client
}

// Handle sends the record to AWS EventBridge.
func (h *EventBridgeHandler) Handle(ctx context.Context, record Record) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(record)
	if err != nil {
		return err
	}
	_, err = h.client.PutEvents(ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				EventBusName: ptr.String(h.BusName),
				Detail:       ptr.String(buf.String()),
			},
		},
	})
	return err
}
