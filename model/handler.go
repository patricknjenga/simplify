package model

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type IHandler[T any] interface {
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	RegisterRoutes()
}

type Handler[T any] struct {
	Router   *mux.Router
	Service  IService[T]
	Validate *validator.Validate
}

func NewModelHandler[T any](router *mux.Router, service IService[T]) IHandler[T] {
	return &Handler[T]{router, service, validator.New()}
}

func (h Handler[T]) RegisterRoutes() {
	h.Router.Path("/").Methods(http.MethodDelete).HandlerFunc(h.Delete)
	h.Router.Path("/").Methods(http.MethodGet).HandlerFunc(h.Get)
	h.Router.Path("/").Methods(http.MethodPost).HandlerFunc(h.Post)
	h.Router.Path("/").Methods(http.MethodPut).HandlerFunc(h.Put)
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
	var q Query
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
	err = json.NewEncoder(w).Encode(map[string]any{"Count": count, "Data": data})
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
