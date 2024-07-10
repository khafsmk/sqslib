package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// EventMetadata contains optional metadata fields for the event
type EventMetadata struct {
	TraceID            string `json:"trace_id,omitempty"`
	Domain             string `json:"domain,omitempty"`
	IdempotencyKey     string `json:"idempotency_key,omitempty"`
	SequenceID         string `json:"sequence_id,omitempty"`
	TenantID           string `json:"tenant_id,omitempty"`
	DataClassification string `json:"data_classification,omitempty"`
}

// EventInput contains all the information needed to publish an event
type EventInput struct {
	DetailType string
	Source     string
	EventMetadata
	Data interface{}
}

// Publisher interface for publishing events
type Publisher interface {
	PublishEvent(ctx context.Context, input EventInput) error
}

// BasePublisher provides common functionality for publishers
type BasePublisher struct{}

func (b *BasePublisher) ensureTraceAndIdempotencyKey(ctx context.Context, input *EventInput) {
	if input.TraceID == "" {
		if span, ok := tracer.SpanFromContext(ctx); ok {
			input.TraceID = fmt.Sprintf("%d", span.Context().TraceID())
		}
	}
	if input.IdempotencyKey == "" {
		input.IdempotencyKey = uuid.New().String()
	}
}

// EventBridgePublisher implements Publisher for EventBridge
type EventBridgePublisher struct {
	BasePublisher
	client       *eventbridge.Client
	eventBusName string
}

// SQSPublisher implements Publisher for SQS
type SQSPublisher struct {
	BasePublisher
	client   *sqs.Client
	queueURL string
}

// NewEventBridgePublisher creates a new EventBridgePublisher
func NewEventBridgePublisher(ctx context.Context, eventBusName string) (*EventBridgePublisher, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := eventbridge.NewFromConfig(cfg)
	return &EventBridgePublisher{
		client:       client,
		eventBusName: eventBusName,
	}, nil
}

// NewSQSPublisher creates a new SQSPublisher
func NewSQSPublisher(ctx context.Context, queueURL string) (*SQSPublisher, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := sqs.NewFromConfig(cfg)
	return &SQSPublisher{
		client:   client,
		queueURL: queueURL,
	}, nil
}

// PublishEvent publishes an event to EventBridge
func (p *EventBridgePublisher) PublishEvent(ctx context.Context, input EventInput) error {
	p.ensureTraceAndIdempotencyKey(ctx, &input)

	eventDetail := struct {
		Metadata EventMetadata `json:"metadata"`
		Data     interface{}   `json:"data"`
	}{
		Metadata: input.EventMetadata,
		Data:     input.Data,
	}

	detailBytes, err := json.Marshal(eventDetail)
	if err != nil {
		return err
	}

	_, err = p.client.PutEvents(ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				EventBusName: aws.String(p.eventBusName),
				Source:       aws.String(input.Source),
				DetailType:   aws.String(input.DetailType),
				Detail:       aws.String(string(detailBytes)),
				Time:         aws.Time(time.Now().UTC()),
			},
		},
	})
	return err
}

// PublishEvent publishes an event to SQS
func (p *SQSPublisher) PublishEvent(ctx context.Context, input EventInput) error {
	p.ensureTraceAndIdempotencyKey(ctx, &input)

	event := struct {
		DetailType string        `json:"detail-type"`
		Source     string        `json:"source"`
		Metadata   EventMetadata `json:"metadata"`
		Data       interface{}   `json:"data"`
	}{
		DetailType: input.DetailType,
		Source:     input.Source,
		Metadata:   input.EventMetadata,
		Data:       input.Data,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(string(eventBytes)),
	})
	return err
}
