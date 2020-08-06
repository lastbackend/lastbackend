package middleware

import (
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/server/config"
	"net/http"
)

const (
	logPrefix = "http:middleware"
)

type Middleware struct {
	storage storage.IStorage
	secret  string
	items   []func(h http.Handler, cfg config.Config) http.Handler
}

func New(stg storage.IStorage, token string) Middleware {
	return Middleware{
		storage: stg,
		secret:  token,
		items:   make([]func(h http.Handler, cfg config.Config) http.Handler, 0),
	}
}

func (m *Middleware) Add(h ...func(h http.Handler, cfg config.Config) http.Handler) {
	m.items = append(m.items, h...)
}

func (m *Middleware) Apply(h http.Handler, cfg config.Config) http.Handler {

	if len(m.items) < 1 {
		return h
	}

	wrapped := h

	for i := len(m.items) - 1; i >= 0; i-- {
		wrapped = m.items[i](wrapped, cfg)
	}

	return wrapped
}
