// context package

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//	done := make(chan interface{})
	//	defer close(done)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(ctx); err != nil {
			fmt.Printf("printGreeting Err: %v\n", err)
			cancel()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(ctx); err != nil {
			fmt.Printf("printFarewell Err: %v\n", err)
		}
	}()
	wg.Wait()
}

func printGreeting(ctx context.Context) error {
	greeting, e := genGreeting(ctx)
	if e != nil {
		return e
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}
func printFarewell(ctx context.Context) error {
	greeting, e := genFarewell(ctx)
	if e != nil {
		return e
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}
func genGreeting(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch locale, e := locale(ctx); {
	case e != nil:
		return "err", e
	case locale == "UTF-8":
		return "안녕", nil
	}
	return "??", fmt.Errorf("unsupported locale")
}
func genFarewell(ctx context.Context) (string, error) {
	switch locale, e := locale(ctx); {
	case e != nil:
		return "err", e
	case locale == "UTF-8":
		return "잘 가", nil
	}
	return "??", fmt.Errorf("unsupported locale")
}

func locale(ctx context.Context) (string, error) {
	if deadline, ok := ctx.Deadline(); ok {
		if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
			return "", context.DeadlineExceeded
		}
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(1 * time.Minute):
	}
	return "UTF-8", nil
}
