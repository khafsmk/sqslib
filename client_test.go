package mqueue

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestDefaultHandler(t *testing.T) {
	nowf := func() time.Time { return time.Time{} }
	uuidf := func() string { return "" }

	cases := []struct {
		name    string
		value   any
		event   Event
		want    Record
		wantErr bool
	}{
		{
			name:    "unknown event",
			value:   map[string]string{"key": "value"},
			event:   Event("unknown"),
			wantErr: true,
		},
		{
			name:  "map",
			value: map[string]string{"key": "value"},
			event: EventLoanCreate,
			want: Record{
				EventName: string(EventLoanCreate),
				Data:      map[string]string{"key": "value"},
			},
		},
		{
			name:  "struct",
			value: struct{ Key string }{Key: "value"},
			event: EventLoanCreate,
			want: Record{
				EventName: string(EventLoanCreate),
				Data:      struct{ Key string }{Key: "value"},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &Client{
				Handler: HandlerFunc(func(ctx context.Context, r Record) error {
					if diff := cmp.Diff(r, tc.want); diff != "" {
						t.Errorf("diff record: %s", diff)
					}
					return nil
				}),
				timeNow: nowf,
				newUUID: uuidf,
			}

			err := p.Publish(context.TODO(), tc.event, tc.value)
			if tc.wantErr && err == nil {
				t.Fatalf("want error for %t", tc.value)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func BenchmarkClient(b *testing.B) {
	b.ReportAllocs()
	c := &Client{
		Handler: HandlerFunc(func(ctx context.Context, record Record) error {
			return nil
		}),
		newUUID: func() string { return "" },
	}
	ctx := context.Background()
	input := map[string]string{"key": "value"}
	for i := 0; i < b.N; i++ {
		err := c.Publish(ctx, EventLoanCreate, input)
		if err != nil {
			b.Fatal(err)
		}
	}
}
