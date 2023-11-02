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
		var name string
		fields := map[string]string{}
		if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
			name = t.Elem().Name()
		} else {
			name = t.Name()
		}
		Fields(v, fields)
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

func Fields(x any, f map[string]string) {
	ts := reflect.TypeOf(x)
	vs := reflect.ValueOf(x)
	for i := 0; i < ts.NumField(); i++ {
		switch ts.Field(i).Type.String() {
		case "gorm.Model":
			Fields(vs.Field(i).Interface(), f)
		default:
			f[ts.Field(i).Name] = ts.Field(i).Type.String()
		}
	}
}
