package errors

import (
	"fmt"
	"log"
	"runtime"
)

func Format(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("%s:%d %w", file, line, err)
}

func Log(err error) {
	log.Println(err)
}
