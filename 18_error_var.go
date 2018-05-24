// Error variables
// package scoped
package main

import (
	"errors"
	"fmt"
	"math/rand"
)

var ErrTimeout = errors.New("timed out")
var ErrRejected = errors.New("rejected")

var random = rand.New(rand.NewSource(35))

func main() {
	r, err := SendRequest("Hello")
	for err == ErrTimeout {
		fmt.Println("timeout, retrying")
		r, err = SendRequest("Hello")
	}
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func SendRequest(req string) (string, error) {
	switch random.Int() % 3 {
	case 0:
		return "ok", nil
	case 1:
		return "", ErrRejected
	default:
		return "", ErrTimeout
	}
}
