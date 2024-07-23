package mqueue

import (
	"context"
	"encoding/json"
	"io"
)

type JSONHandler struct {
	w io.Writer
}

// NewJSONHandler returns a new JSON handler that writes the record to the given writer.
func NewJSONHandler(w io.Writer) *JSONHandler {
	return &JSONHandler{w}
}

// Handle formats its argument [Record] as a as a JSON object on a single line.
func (p *JSONHandler) Handle(ctx context.Context, record Record) error {
	b, err := json.Marshal(record)
	if err != nil {
		return err
	}
	_, err = p.w.Write(b)
	return err
}
