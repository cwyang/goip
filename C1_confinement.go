// Lexical confinement
package main

import (
	"bytes"
	"fmt"
	"sync"
)

func ex1() {
	// write in confined inside chanOwner
	// only read is exported
	chanOwner := func() <-chan int {
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Printf("got: %d\n", result)
		}
		fmt.Println("done")
	}

	results := chanOwner()
	consumer(results)
}

// data is confined inside printData
func ex2() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buf bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buf, "%c", b)
		}
		fmt.Println(buf.String())
	}
	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("hello")
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])
	wg.Wait()
}
func main() {
	ex2()
}
