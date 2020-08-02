package middleware

import (
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
)

const (
	logPrefix = "http:middleware"
)

type Middleware struct {
	storage storage.IStorage
	secret  string
}

func New(stg storage.IStorage, token string) Middleware {
	return Middleware{
		storage: stg,
		secret:  token,
	}
}
