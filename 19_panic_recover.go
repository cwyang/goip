// Recovering from a panic
package main

import (
	"errors"
	"fmt"
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("Trapped panic: %s (%T)\n", e, e)
		}
	}()
	yikes()
}

func yikes() {
	panic(errors.New("bad, bad!"))
}
