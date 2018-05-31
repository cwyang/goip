package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func doWork(
	done <-chan interface{},
	id int,
	wg *sync.WaitGroup,
	result chan<- int,
) {
	started := time.Now()
	defer wg.Done()

	// simulate random load
	simulatedLoadTime := time.Duration(1+rand.Intn(5)) * time.Second
	select {
	case <-done:
	case <-time.After(simulatedLoadTime):
	}

	select {
	case <-done:
	case result <- id:
	}
	took := time.Since(started)
	if took < simulatedLoadTime { // :-(
		took = simulatedLoadTime
	}
	fmt.Printf("%v took %v\n", id, took)
}

func main() {
	done := make(chan interface{})
	result := make(chan int)

	const N int = 10
	var wg sync.WaitGroup
	wg.Add(N)

	for i := 0; i < N; i++ {
		go doWork(done, i, &wg, result)
	}
	firstReturned := <-result
	close(done)
	wg.Wait()

	fmt.Printf("#%v is the winner\n", firstReturned)
}
