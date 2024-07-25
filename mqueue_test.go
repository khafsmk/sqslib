package mqueue_test

import (
	"bytes"
	"context"
	"testing"

	mq "github.com/khafsmk/mqueue"
)

func TestHandlerFunc(t *testing.T) {
	client := &mq.Client{
		Handler: mq.HandlerFunc(func(ctx context.Context, record mq.Record) error {
			return nil
		}),
	}
	err := client.Publish(context.Background(), mq.EventLoanCreate, map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMQueueTest(t *testing.T) {
	var buf bytes.Buffer
	h := mq.NewJSONHandler(&buf)
	client := &mq.Client{Handler: h}
	err := client.Publish(context.Background(), mq.EventLoanCreate, map[string]string{"key": "value"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(buf.String())
}

func TestSetDefault(t *testing.T) {
	currentHandler := mq.Default()
	serviceName := "test-service"
	client := mq.New(
		mq.HandlerFunc(func(ctx context.Context, record mq.Record) error {
			return nil
		}),
		mq.WithServiceName(serviceName),
		mq.WithDomain("test-domain"),
		mq.WithSquadName("squad-name"),
	)
	mq.SetDefault(client)

	t.Cleanup(func() {
		mq.SetDefault(currentHandler)
	})

	err := mq.Publish(mq.EventLoanUpdate, map[string]string{"a": "1"})
	check(t, err)

	err = mq.PublishContext(context.Background(), mq.EventLoanCreate, map[string]string{"b": "2"})
	check(t, err)
}

func check(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkClient(b *testing.B) {
	b.ReportAllocs()
	c := mq.New(
		mq.HandlerFunc(func(ctx context.Context, record mq.Record) error {
			return nil
		}),
	)
	ctx := context.Background()
	input := map[string]string{"key": "value"}
	for i := 0; i < b.N; i++ {
		err := c.Publish(ctx, mq.EventLoanCreate, input)
		if err != nil {
			b.Fatal(err)
		}
	}
}
