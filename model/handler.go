package model

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type IHandler interface {
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
}

type Handler[T interface{}] struct {
	Router  *mux.Router
	Service IService[T]
}

func NewModelHandler[T interface{}](n string, rt *mux.Router, svc IService[T]) IHandler {
	var h = &Handler[T]{rt, svc}
	h.Router.Path(fmt.Sprintf("/%s", n)).Methods(http.MethodDelete).HandlerFunc(h.Delete)
	h.Router.Path(fmt.Sprintf("/%s", n)).Methods(http.MethodGet).HandlerFunc(h.Get)
	h.Router.Path(fmt.Sprintf("/%s", n)).Methods(http.MethodPost).HandlerFunc(h.Post)
	h.Router.Path(fmt.Sprintf("/%s", n)).Methods(http.MethodPut).HandlerFunc(h.Put)
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
		return
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
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"Count": count,
		"Data":  data,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
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
		return
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
		return
	}
}
