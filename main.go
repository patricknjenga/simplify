package main

import (
	"fmt"

	"github.com/patricknjenga/simplify/errors"
)

func main() {
	errors.Log(fmt.Errorf("hello world"))
}
