package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Model[T interface{}] struct {
	Fields     []Field
	Handler    IHandler[T]
	Name       string
	Repository IRepository[T]
	Service    IService[T]
	Struct     T
}

func (m Model[T]) New(rt *mux.Router, db *gorm.DB) *Model[T] {
	var (
		t T
		f = Fields(t)
		n = reflect.TypeOf(t).Name()
		r = NewGormRepository[T](1000, db)
		s = NewModelService[T](r)
		h = NewModelHandler[T](rt.PathPrefix(fmt.Sprintf("/%s", n)).Subrouter(), s)
	)
	return &Model[T]{f, h, n, r, s, t}
}

func New(rt *mux.Router, db *gorm.DB, models ...Model[interface{}]) []*Model[interface{}] {
	var (
		res    []*Model[interface{}]
		schema []Schema
	)
	for k, v := range models {
		res = append(res, v.New(rt, db))
		schema = append(schema, Schema{
			Fields: res[k].Fields,
			Name:   res[k].Name,
			Struct: res[k].Struct,
		})
		res[k].Handler.RegisterRoutes()
	}
	rt.Path("/Schema").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(schema)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	return res
}
