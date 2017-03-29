package serializer

import (
	"bytes"
)

type Codec Serializer

// serializer binds an encoder and decoder.
type serializer struct {
	Encoder
	Decoder
}

// NewSerializer creates a NewSerializer from an Encoder and Decoder.
func NewSerializer(e Encoder, d Decoder) Codec {
	return serializer{e, d}
}

// Encode is a convenience wrapper for encoding to a []byte from an Encoder
func Encode(e Encoder, obj interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := e.Encode(obj, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode is a convenience wrapper for decoding data into an Object.
func Decode(d Decoder, data []byte, obj interface{}) error {
	return d.Decode(data, obj)
}
