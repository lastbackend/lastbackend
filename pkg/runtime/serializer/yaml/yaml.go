package yaml

import (
	"github.com/lastbackend/lastbackend/pkg/runtime"
	"gopkg.in/yaml.v2"
)

type yamlSerializer struct {
	// the nested serializer
	runtime.Serializer
}

func NewSerializer() runtime.Serializer {
	return &yamlSerializer{}
}

func (yamlSerializer) Encode(obj runtime.Object) ([]byte, error) {
	return yaml.Marshal(obj)
}

func (yamlSerializer) Decode(buf []byte, into runtime.Object) error {
	return yaml.Unmarshal(buf, into)
}
