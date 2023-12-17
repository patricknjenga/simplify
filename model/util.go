package model

import (
	"encoding/json"
	"fmt"
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

type Route struct {
	Method string
	Path   string
}

type Schema struct {
	Fields []Field
	Name   string
	Struct interface{}
}

func Actions(r *mux.Router, s ...Action) {
	var res = []Action{}
	for _, v := range s {
		if v.Struct != nil {
			v.Fields = Fields(v.Struct)
		}
		res = append(res, v)
	}
	r.Path("/Action").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
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

func Routes(rt *mux.Router) error {
	var routes []Route
	err := rt.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return fmt.Errorf("%s %w", path, err)
		}
		for _, v := range methods {
			routes = append(routes, Route{v, path})
		}
		return nil
	})
	if err != nil {
		return err
	}
	rt.Path("/").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(routes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	return nil
}

func Schemas(r *mux.Router, s ...interface{}) {
	var res = []Schema{}
	for _, v := range s {
		if t := reflect.TypeOf(v); t.Kind() != reflect.Ptr {
			res = append(res, Schema{Fields(v), t.Name(), s})
		} else {
			res = append(res, Schema{Fields(v), t.Elem().Name(), s})
		}
	}
	r.Path("/Schema").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
