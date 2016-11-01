package runtime

type Serializer interface {
	Encode(obj Object) ([]byte, error)
	Decode(data []byte, into Object) error
}

type Object interface {
}
