package model

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

type Field struct {
	Name string
	Type string
}

func Schema(r *mux.Router, s ...any) error {
	var res = map[string][]Field{}
	for _, v := range s {
		var name string
		if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
			name = t.Elem().Name()
		} else {
			name = t.Name()
		}
		res[name] = Fields(v)
	}
	r.HandleFunc("/Schema", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	return nil
}

func Fields(x any) []Field {
	var (
		res = []Field{}
		ts  = reflect.TypeOf(x)
		vs  = reflect.ValueOf(x)
	)
	for i := 0; i < ts.NumField(); i++ {
		if vs.Field(i).Kind() == reflect.Struct {
			res = append(res, Fields(vs.Field(i).Interface())...)
		} else {
			res = append(res, Field{
				Name: ts.Field(i).Name,
				Type: ts.Field(i).Type.String(),
			})
		}
	}
	return res
}
