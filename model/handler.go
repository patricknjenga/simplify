package model

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type IHandler[T any] interface {
	RegisterRoutes()
	Create(w http.ResponseWriter, r *http.Request)
	CreateBatch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	DeleteAll(w http.ResponseWriter, r *http.Request)
	DeleteBatch(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Query(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type MuxHandler[T any] struct {
	Router  *mux.Router
	Service IService[T]
}

func NewModelHandler[T any](router *mux.Router, service IService[T]) IHandler[T] {
	return &MuxHandler[T]{router, service}
}

func (h MuxHandler[T]) RegisterRoutes() {
	h.Router.HandleFunc("/", h.Create).Methods(http.MethodPost)
	h.Router.HandleFunc("/batch", h.CreateBatch).Methods(http.MethodPost)

	h.Router.HandleFunc("/all", h.DeleteAll).Methods(http.MethodDelete)
	h.Router.HandleFunc("/batch", h.DeleteBatch).Methods(http.MethodDelete)
	h.Router.HandleFunc("/{id:[0-9]+}", h.Delete).Methods(http.MethodDelete)

	h.Router.HandleFunc("/", h.Query).Methods(http.MethodGet)
	h.Router.HandleFunc("/{id:[0-9]+}", h.Get).Methods(http.MethodGet)

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
	err = h.Service.Create(r.Context(), t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h MuxHandler[T]) CreateBatch(w http.ResponseWriter, r *http.Request) {
	var (
		t []T
	)
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, v := range t {
		err = validator.New().Struct(&v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	}
	err = h.Service.CreateBatch(r.Context(), t, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h MuxHandler[T]) Delete(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)
	)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.Service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h MuxHandler[T]) DeleteAll(w http.ResponseWriter, r *http.Request) {
	err := h.Service.DeleteAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h MuxHandler[T]) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	var (
		ids []int
	)
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.Service.DeleteBatch(r.Context(), ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h MuxHandler[T]) Query(w http.ResponseWriter, r *http.Request) {
	var query Query
	if r.URL.Query().Get("q") != "" {
		err := json.Unmarshal([]byte(r.URL.Query().Get("q")), &query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	data, err := h.Service.Query(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h MuxHandler[T]) Get(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)
	)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := h.Service.Get(r.Context(), id)
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
	err = h.Service.Update(r.Context(), id, t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
