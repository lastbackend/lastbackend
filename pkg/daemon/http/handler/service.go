package handler

import (
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"net/http"
)

func ServiceListH(w http.ResponseWriter, r *http.Request) {

	var (
		er  error
		ctx = context.Get()
	)

	w.WriteHeader(200)
	_, er = w.Write([]byte{})
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

func ServiceInfoH(w http.ResponseWriter, r *http.Request) {

	var (
		er  error
		ctx = context.Get()
	)

	w.WriteHeader(200)
	_, er = w.Write([]byte{})
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

func ServiceCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		er  error
		ctx = context.Get()
	)

	w.WriteHeader(200)
	_, er = w.Write([]byte{})
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

func ServiceUpdateH(w http.ResponseWriter, r *http.Request) {

	var (
		er  error
		ctx = context.Get()
	)

	w.WriteHeader(200)
	_, er = w.Write([]byte{})
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}

func ServiceRemoveH(w http.ResponseWriter, r *http.Request) {

	var (
		er  error
		ctx = context.Get()
	)

	w.WriteHeader(200)
	_, er = w.Write([]byte{})
	if er != nil {
		ctx.Log.Error("Error: write response", er.Error())
		return
	}
}
