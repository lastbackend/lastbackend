package runtime

import (
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/storage"
)


type Runtime struct {
	routes    *mux.Router
	storage   storage.IStorage
}

func New (storage storage.IStorage) *Runtime {

	runtime := new(Runtime)
	runtime.storage = storage

	return runtime
}


