package json

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/runtime"
)

type jsonSerializer struct {
	// the nested serializer
	runtime.Serializer
}

func NewSerializer() runtime.Serializer {
	return &jsonSerializer{}
}

func (jsonSerializer) Encode(obj runtime.Object) ([]byte, error) {
	return json.Marshal(obj)
}

func (jsonSerializer) Decode(buf []byte, into runtime.Object) error {
	return json.Unmarshal(buf, into)
}
