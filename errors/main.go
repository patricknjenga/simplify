package errors

import (
	"encoding/json"
	"os"
	"runtime"
	"time"
)

func ArrDo(functions []func() error) error {
	for _, function := range functions {
		err := function()
		if err != nil {
			return err
		}
	}
	return nil
}

func Log(err error) {
	_, file, line, _ := runtime.Caller(1)
	var e = struct {
		File    string
		Line    int
		Message string
		Time    string
	}{
		File:    file,
		Line:    line,
		Message: err.Error(),
		Time:    time.Now().Format("20060102150405"),
	}
	json.NewEncoder(os.Stderr).Encode(&e)
}
