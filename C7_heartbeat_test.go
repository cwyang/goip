package main1

import (
	"testing"
	"time"
)

// export IntStream to Test
func IntStream(done <-chan interface{}, nums ...int,
) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(intStream)
		//

		time.Sleep(2 * time.Second)

		for _, n := range nums {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
			select {
			case <-done:
				return
			case intStream <- n:
			}
		}
	}()
	return heartbeat, intStream
}

func TestDoWork_GenerateAllNumbersOK(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	heartbeat, results := IntStream(done, intSlice...)

	<-heartbeat

	i := 0
	for r := range results {
		if expected := intSlice[i]; r != expected {
			t.Errorf("index %v: expected %v, but received %v,",
				i, expected, r)
		}
		i++
	}
}

func TestDoWork_GenerateAllNumbersBad(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	_, results := IntStream(done, intSlice...)

	for i, expected := range intSlice {
		select {
		case r := <-results:
			if r != expected {
				t.Errorf("index %v: expected %v, but received %v,",
					i, expected, r)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("timeout")
		}
	}
}
