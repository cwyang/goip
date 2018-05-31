package main1

import (
	"fmt"
	"time"
)

func doWork(
	done <-chan interface{},
	pulseInterval time.Duration,
) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)
	go func() {
		defer close(heartbeat)
		defer close(results)
		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(pulseInterval * 2)

		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default: // if no one listens heartbeat
			}
		}
		sendResult := func(r time.Time) {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case results <- r:
					return
				}
			}
		}
		for {
			//		for i := 0; i < 2; i++ {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			}
		}

	}()
	return heartbeat, results
}

func main2() {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				fmt.Println("hearbeat false")
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				fmt.Println("results false")
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			fmt.Println("timeout")
			return
		}
	}
}
