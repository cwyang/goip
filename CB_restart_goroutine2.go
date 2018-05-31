package main

import (
	"fmt"
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

func doWorkFn(done <-chan interface{}, intList ...int) (startGoroutineFn, <-chan interface{}) {
	intChanStream := make(chan (<-chan interface{}))
	intStream := bridge(done, intChanStream)
	doWork := func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) <-chan interface{} {
		intStream := make(chan interface{})
		heartbeat := make(chan interface{})
		go func() {
			defer close(intStream)
			select {
			case intChanStream <- intStream:
			case <-done:
				return
			}
			pulse := time.Tick(pulseInterval)

			for {
			valueLoop:
				for _, intVal := range intList {
					if intVal < 0 {
						log.Printf("negative: %v\n", intVal)
						return
					}
					for {
						select {
						case <-pulse:
							select {
							case heartbeat <- struct{}{}:
							default:
							}
						case intStream <- intVal:
							continue valueLoop
						case <-done:
							return
						}
					}
				}
			}
		}()
		return heartbeat
	}

	return doWork, intStream
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	done := make(chan interface{})
	defer close(done)

	doWork, intStream := doWorkFn(done, 1, 2, -1, 3, 4, 5)
	doWorkWithMaster := newMaster(1*time.Millisecond, doWork)
	doWorkWithMaster(done, 1*time.Hour)
	for intVal := range take(done, intStream, 6) {
		fmt.Printf("recv: %v\n", intVal)
	}
}
