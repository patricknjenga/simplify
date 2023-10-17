package model

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

func Schema(r *mux.Router, s ...any) error {
	var res = map[string]map[string]string{}
	for _, v := range s {
		var fields = map[string]string{}
		for i := 0; i < reflect.TypeOf(v).NumField(); i++ {
			field := reflect.TypeOf(v).Field(i)
			fields[field.Name] = field.Type.String()
		}
		var name string
		if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
			name = t.Elem().Name()
		} else {
			name = t.Name()
		}
		res[name] = fields
	}
	r.HandleFunc("/Schema", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	return nil
}
