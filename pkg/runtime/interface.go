package runtime

import (
	"io"
)

// All API types registered with Scheme must support the Object interface.
type Object interface {
}

// Encoders write objects to a serialized form
type Encoder interface {
	Encode(obj Object, w io.Writer) error
}

// Decoders attempt to load an object from data.
type Decoder interface {
	Decode(data []byte, into Object) (Object, error)
}

// Serializer is the core interface for transforming objects into a serialized format and back.
// Implementations may choose to perform conversion of the object, but no assumptions should be made.
type Serializer interface {
	Encoder
	Decoder
}
