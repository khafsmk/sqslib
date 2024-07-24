package mqueue

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestDefaultHandler(t *testing.T) {
	var buf bytes.Buffer
	h := HandlerFunc(func(ctx context.Context, record Record) error {
		return json.NewEncoder(&buf).Encode(record)
	})

	nowf := func() time.Time { return time.Time{} }
	uuidf := func() string { return "" }

	check := func(want Record) {
		t.Helper()
		b, err := json.Marshal(want)
		if err != nil {
			t.Fatal(err)
		}
		got := buf.Bytes()
		if diff := cmp.Diff(b, got, transformJSON); diff != "" {
			t.Errorf("unexpected write (-want +got):\n%s", diff)
		}
		buf.Reset()
	}

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
				Data:      map[string]string{"Key": "value"},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &Client{
				Handler: h,
				timeNow: nowf,
				newUUID: uuidf,
			}
			err := p.Publish(context.TODO(), tc.event, tc.value)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for %t", tc.value)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				check(tc.want)
			}
		})
	}
}
