// for-select
package main

import (
	"bytes"
	"fmt"
	"sync"
)

// Sending iter vars out on a channel
/*
func ex1() {
	for _, s := range []string{"a", "b", "c"} {
		select {
		case <-done:
			return
		case stringStream <- s:
			return
		}
	}
}
*/

// Looping infinitely waiting to be stopped
/*
func ex1 () {
	for {
		select {
		case <- done:
			return
			default:
		}
		// Do non-preemtable work
	}
}
*/
func main() {
	ex1()
}
