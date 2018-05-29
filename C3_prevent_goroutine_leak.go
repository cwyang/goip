// for-select
package main

import (
	"fmt"
	"time"
)

func doWork(done <-chan interface{},
	strings <-chan string,
) <-chan interface{} {
	completed := make(chan interface{})
	go func() {
		defer fmt.Println("doWork exited")
		defer close(completed)
		for {
			select {
			case s := <-strings:
				// Do job
				fmt.Println(s)
			case <-done:
				return
			}
		}
	}()
	return completed
}

func main() {
	done := make(chan interface{})
	terminated := doWork(done, nil)
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("canceling")
		close(done)
	}()
	<-terminated

	fmt.Println("done")
}
