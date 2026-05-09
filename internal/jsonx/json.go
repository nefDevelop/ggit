package jsonx

import (
	"bytes"
	"io"

	"github.com/goccy/go-json"
)

// Marshal is a drop-in replacement for json.Marshal.
func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent is a drop-in replacement for json.MarshalIndent.
func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// Unmarshal is a drop-in replacement for json.Unmarshal.
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *json.Encoder {
	return json.NewEncoder(w)
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}

// Valid reports whether data is a valid JSON encoding.
func Valid(data []byte) bool {
	return json.Valid(data)
}

// Indent is a drop-in replacement for json.Indent.
func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error {
	return json.Indent(dst, src, prefix, indent)
}

// HTMLEscape is a drop-in replacement for json.HTMLEscape.
func HTMLEscape(dst *bytes.Buffer, src []byte) {
	json.HTMLEscape(dst, src)
}

// RawMessage is a raw encoded JSON value.
type RawMessage = json.RawMessage
