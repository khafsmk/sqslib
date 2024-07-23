package mqueue_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	mq "github.com/khafsmk/mqueue"
)

func TestHandlerFunc(t *testing.T) {
	client := &mq.Client{
		Handler: mq.HandlerFunc(func(ctx context.Context, record mq.Record) error {
			return nil
		}),
	}
	err := client.Publish(context.Background(), map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMQueueTest(t *testing.T) {
	var buf bytes.Buffer
	h := mq.NewJSONHandler(&buf)
	client := &mq.Client{Handler: h}
	err := client.Publish(context.Background(), map[string]string{"key": "value"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(buf.String())
}

func checkRecord(t *testing.T, buf *bytes.Buffer, want *mq.Record) {
	t.Helper()
	got := new(mq.Record)
	err := json.Unmarshal(buf.Bytes(), got)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(mq.Record{}, "IdempotencyKey", "Time")); diff != "" {
		t.Errorf("unexpected write (-want +got):\n%s", diff)
	}
	buf.Reset()
}

func TestSetDefault(t *testing.T) {
	var buf bytes.Buffer
	currentHandler := mq.Default()
	serviceName := "test-service"
	client := mq.New("", "", serviceName, mq.NewJSONHandler(&buf))
	mq.SetDefault(client)

	t.Cleanup(func() {
		mq.SetDefault(currentHandler)
	})

	err := mq.Publish(map[string]string{"a": "1"})
	check(t, err)
	checkRecord(t, &buf, &mq.Record{
		Source: serviceName,
		Data:   map[string]any{"a": "1"},
	})

	err = mq.PublishContext(context.Background(), map[string]string{"b": "2"})
	check(t, err)
	checkRecord(t, &buf, &mq.Record{
		Source: serviceName,
		Data:   map[string]any{"b": "2"},
	})
}

func check(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
