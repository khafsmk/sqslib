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
	ctx := context.Background()

	var buf bytes.Buffer
	h := NewJSONHandler(&buf)

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
		want    Record
		wantErr bool
	}{
		{
			name:  "map",
			value: map[string]string{"key": "value"},
			want: Record{
				Data: map[string]string{"key": "value"},
			},
		},
		{
			name:  "struct",
			value: struct{ Key string }{Key: "value"},
			want: Record{
				Data: map[string]string{"Key": "value"},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &Client{
				handler: h,
				timeNow: nowf,
				newUUID: uuidf,
			}
			err := p.Publish(ctx, tc.value)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %t", tc.value)
			}
			if err != nil {
				t.Fatal(err)
			}
			check(tc.want)
		})
	}
}
