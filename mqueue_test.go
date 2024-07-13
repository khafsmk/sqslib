package mqueue_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/khafsmk/mqueue"
)

func TestMQueueTest(t *testing.T) {
	var buf bytes.Buffer
	client := mqueue.New("facility-service", mqueue.NewJSONHandler(&buf))
	err := client.Publish(context.Background(), map[string]string{"key": "value"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(buf.String())
}

func checkPublished(t *testing.T, buf *bytes.Buffer, want *mqueue.Record) {
	t.Helper()
	got := new(mqueue.Record)
	err := json.Unmarshal(buf.Bytes(), got)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(mqueue.Record{}, "IdempotencyKey", "Time")); diff != "" {
		t.Errorf("unexpected write (-want +got):\n%s", diff)
	}
	buf.Reset()
}
func check(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetDefault(t *testing.T) {
	var buf bytes.Buffer
	currentHandler := mqueue.Default()
	name := "test-service"
	mqueue.SetDefault(mqueue.New(name, mqueue.NewJSONHandler(&buf)))
	t.Cleanup(func() {
		mqueue.SetDefault(currentHandler)
	})

	err := mqueue.Publish(map[string]string{"a": "1"})
	check(t, err)
	checkPublished(t, &buf, &mqueue.Record{
		Source: name,
		Data:   map[string]any{"a": "1"},
	})

	err = mqueue.PublishContext(context.Background(), map[string]string{"b": "2"})
	check(t, err)
	checkPublished(t, &buf, &mqueue.Record{
		Source: name,
		Data:   map[string]any{"b": "2"},
	})

}
