package mqueue

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/smithy-go/ptr"
)

// NewEventBridgeHandler returns a new EventBridge handler.
func NewEventBridgeHandler(client *eventbridge.Client) *EventBridgeHandler {
	return &EventBridgeHandler{
		client: client,
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
