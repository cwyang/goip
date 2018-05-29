// repeat, take, orDone and tee
package main

import (
	"fmt"
)

func orDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valStream <- v:
				}
			}
		}
	}()
	return valStream
}

func take(done, valStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valStream:
			}
		}
	}()
	return takeStream
}

func tee(done, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range orDone(done, in) {
			var o1, o2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case o1 <- val:
					o1 = nil
				case o2 <- val:
					o2 = nil
				}
			}
		}
	}()
	return out1, out2
}

// bridge-channel: channel of channels
func bridge(
	done <-chan interface{},
	chanStream <-chan <-chan interface{},
) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range orDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func test1() {
	done := make(chan interface{})

	doWork := func(done, c <-chan interface{}) {
		go func() {
			for val := range orDone(done, c) { // simple orDone
				fmt.Println(val)
			}
		}()
	}

	myChan := make(chan interface{})
	defer close(myChan)
	doWork(done, myChan)
	for i := 1; i <= 10; i++ {
		myChan <- i
	}
	close(done)
	fmt.Println("done")
}

func test2() {
	done := make(chan interface{})
	defer close(done)

	o1, o2 := tee(done, take(done, repeat(done, "hello", "world"), 5))

	for v := range o1 {
		fmt.Printf("o1: %v o2: %v\n", v, <-o2)
	}
}

func test3() {
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}
	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v", v)
	}
}

func main() {
	test3()
}
