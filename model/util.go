package model

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

type Action struct {
	Fields []Field
	Struct interface{}
	Url    string
}

type Field struct {
	Name string
	Tag  reflect.StructTag
	Type string
}

type Schema struct {
	Fields []Field
	Name   string
	Struct interface{}
}

func ActionRoute(r *mux.Router, s ...Action) {
	var res = []Action{}
	for _, v := range s {
		if v.Struct != nil {
			v.Fields = Fields(v.Struct)
		}
		res = append(res, v)
	}
	r.HandleFunc("/Actions", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods(http.MethodGet)
}

func Fields(x interface{}) []Field {
	var (
		result  []Field
		typeOf  = reflect.TypeOf(x)
		valueOf = reflect.ValueOf(x)
	)
	for i := 0; i < typeOf.NumField(); i++ {
		switch typeOf.Field(i).Type.String() {
		case "gorm.Model":
			result = append(result, Fields(valueOf.Field(i).Interface())...)
		default:
			result = append(result, Field{
				Name: typeOf.Field(i).Name,
				Tag:  typeOf.Field(i).Tag,
				Type: typeOf.Field(i).Type.String(),
			})
		}
	}
	return result
}
