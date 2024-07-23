package mqueue

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Record contains optional metadata fields for the event
type Record struct {
	Source             string    `json:"source"`
	Destination        string    `json:"destination,omitempty"`
	Time               time.Time `json:"time"`
	TraceID            string    `json:"trace_id,omitempty"`
	IdempotencyKey     string    `json:"idempotency_key,omitempty"`
	SequenceID         string    `json:"sequence_id,omitempty"`
	TenantID           string    `json:"tenant_id,omitempty"`
	DataClassification string    `json:"data_classification,omitempty"`
	Data               any       `json:"data"`
}

func (c *Client) newRecord(ctx context.Context, source string, data any) Record {
	var traceID string
	if span, ok := tracer.SpanFromContext(ctx); ok {
		traceID = fmt.Sprintf("%d", span.Context().TraceID())
	}
	if c.timeNow == nil {
		c.timeNow = time.Now
	}
	if c.newUUID == nil {
		c.newUUID = uuid.NewString
	}
	return Record{
		Data:           data,
		Source:         source,
		Time:           c.timeNow().UTC(),
		TraceID:        traceID,
		IdempotencyKey: c.newUUID(),
	}
}
