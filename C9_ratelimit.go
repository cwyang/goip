package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

//client

func Open() *APIConnection {
	return &APIConnection{
		rateLimiter: rate.NewLimiter(rate.Limit(1), 5), // 1 events/sec, burst 5
	}
}

type APIConnection struct {
	rateLimiter *rate.Limiter
}

// ReadFile and ResolveAddress is not relevant
func (a *APIConnection) ReadFile(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// do some work
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// do some work
	return nil
}

func main() {
	defer log.Printf("done")

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := Open()
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("readfile err: %v", err)
			}
			log.Printf("readfile")
		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("resolve err: %v", err)
			}
			log.Printf("resolve")
		}()
	}

	wg.Wait()
}
