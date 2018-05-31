package main

import (
	"log"
	"os"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} { // <1>
	switch len(channels) {
	case 0: // <2>
		return nil
	case 1: // <3>
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() { // <4>
		defer close(orDone)

		switch len(channels) {
		case 2: // <5>
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default: // <6>
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or(append(channels[3:], orDone)...): // <6>
			}
		}
	}()
	return orDone
}

type startGoroutineFn func(
	done <-chan interface{},
	pulseInterval time.Duration,
) (heartbeat <-chan interface{})

func newMaster(
	timeout time.Duration,
	startGoroutine startGoroutineFn,
) startGoroutineFn { // master is also monitorable
	return func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) <-chan interface{} {
		heartbeat := make(chan interface{})
		go func() {
			defer close(heartbeat)

			var slaveDone chan interface{}
			var slaveHeartbeat <-chan interface{}
			startSlave := func() {
				slaveDone = make(chan interface{})
				slaveHeartbeat = startGoroutine(or(slaveDone, done), timeout/2)
			}
			startSlave()
			pulse := time.Tick(pulseInterval)
		monitorLoop:
			for {
				timeoutSignal := time.After(timeout)
				for {
					select {
					case <-pulse:
						select {
						case heartbeat <- struct{}{}:
						default:
						}
					case <-slaveHeartbeat:
						continue monitorLoop
					case <-timeoutSignal:
						log.Println("master: slave unresponding; restarting")
						close(slaveDone)
						startSlave()
						continue monitorLoop
					case <-done:
						return
					}
				}
			}
		}()
		return heartbeat
	}
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	doWork := func(done <-chan interface{}, _ time.Duration) <-chan interface{} {
		log.Println("slave: Hello~ argh~~")
		go func() {
			<-done
			log.Println("slave:I'm halting")
		}()
		return nil
	}
	doWorkWithMaster := newMaster(4*time.Second, doWork)

	done := make(chan interface{})
	time.AfterFunc(9*time.Second, func() {
		log.Println("main: halting master and slave")
		close(done)
	})

	for range doWorkWithMaster(done, 4*time.Second) {
	}
	log.Println("Done")
}
