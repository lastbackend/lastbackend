package runtime

import (
	"bytes"
)

// codec binds an encoder and decoder.
type serializer struct {
	Encoder
	Decoder
}

// NewCodec creates a Codec from an Encoder and Decoder.
func NewCodec(e Encoder, d Decoder) Serializer {
	return serializer{e, d}
}

// Encode is a convenience wrapper for encoding to a []byte from an Encoder
func Encode(e Encoder, obj Object) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := e.Encode(obj, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode is a convenience wrapper for decoding data into an Object.
func Decode(d Decoder, data []byte) (Object, error) {
	obj, err := d.Decode(data, nil)
	return obj, err
}
