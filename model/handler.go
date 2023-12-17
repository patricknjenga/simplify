package model

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type IHandler interface {
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
}

type Handler[T any] struct {
	Router  *mux.Router
	Service IService[T]
}

func NewModelHandler[T any](rt *mux.Router, svc IService[T]) IHandler {
	var h = &Handler[T]{rt, svc}
	h.Router.Path("/").Methods(http.MethodDelete).HandlerFunc(h.Delete)
	h.Router.Path("/").Methods(http.MethodGet).HandlerFunc(h.Get)
	h.Router.Path("/").Methods(http.MethodPost).HandlerFunc(h.Post)
	h.Router.Path("/").Methods(http.MethodPut).HandlerFunc(h.Put)
	return h
}

func (h Handler[T]) Delete(w http.ResponseWriter, r *http.Request) {
	var t []T
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.Service.Delete(r.Context(), t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (h Handler[T]) Get(w http.ResponseWriter, r *http.Request) {
	var (
		t T
		q Query
	)
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	count, data, err := h.Service.Get(r.Context(), q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	err = json.NewEncoder(w).Encode(map[string]any{
		"Count":  count,
		"Data":   data,
		"Fields": Fields(t),
		"Struct": t,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (h Handler[T]) Post(w http.ResponseWriter, r *http.Request) {
	var t []T
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.Service.Post(r.Context(), t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (h Handler[T]) Put(w http.ResponseWriter, r *http.Request) {
	var t []T
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.Service.Put(r.Context(), t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}
