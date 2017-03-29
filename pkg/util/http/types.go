package http

import "net/http"

type Route struct {
	Path string
	Handler func(w http.ResponseWriter, r *http.Request)
	Middleware []Middleware
	Method string
}

type Middleware func(http.HandlerFunc) http.HandlerFunc
