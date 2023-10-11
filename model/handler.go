package model

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type MuxHandler[T any] struct {
	Router  *mux.Router
	Service IService[T]
}

func NewModelHandler[T any](router *mux.Router, service IService[T]) IHandler[T] {
	return &MuxHandler[T]{router, service}
}

func (h MuxHandler[T]) RegisterRoutes() {
	h.Router.HandleFunc("/", h.Create).Methods(http.MethodPost)
	h.Router.HandleFunc("/", h.Index).Methods(http.MethodGet)
	h.Router.HandleFunc("/{id:[0-9]+}", h.Destroy).Methods(http.MethodDelete)
	h.Router.HandleFunc("/{id:[0-9]+}", h.Show).Methods(http.MethodGet)
	h.Router.HandleFunc("/{id:[0-9]+}", h.Update).Methods(http.MethodPut)
}

func (h MuxHandler[T]) Create(w http.ResponseWriter, r *http.Request) {
	var (
		t T
	)
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validator.New().Struct(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	err = h.Service.Create(r.Context(), &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h MuxHandler[T]) Destroy(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)
	)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.Service.Destroy(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h MuxHandler[T]) Index(w http.ResponseWriter, r *http.Request) {
	var (
		query Query
	)
	if r.URL.Query().Get("q") != "" {
		err := json.Unmarshal([]byte(r.URL.Query().Get("q")), &query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	data, err := h.Service.Index(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h MuxHandler[T]) Show(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)
	)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := h.Service.Show(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h MuxHandler[T]) Update(w http.ResponseWriter, r *http.Request) {
	var (
		t    T
		vars = mux.Vars(r)
	)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validator.New().Struct(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	err = h.Service.Update(r.Context(), id, &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
