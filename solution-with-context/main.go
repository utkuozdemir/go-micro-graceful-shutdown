package main

import (
	"context"
	"log"
	"sync"
	"time"
)

func main() {
	mySvc := newMyService()

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		mySvc.uploadForever(ctx)
	}()

	time.Sleep(2250 * time.Millisecond)

	cancel() // cancel the context, shutdown the thing

	log.Printf("wait for goroutine to finish...")

	wg.Wait()

	log.Printf("main done")
}

type myService struct {
	loopTicker *time.Ticker

	stopCh chan struct{}
}

func newMyService() *myService {
	return &myService{
		loopTicker: time.NewTicker(1 * time.Second),
		stopCh:     make(chan struct{}),
	}
}

func (s *myService) uploadForever(ctx context.Context) {
	for {
		s.uploadSingle()

		select {
		case <-ctx.Done():
			log.Printf("context done, terminate...")

			return
		case <-s.loopTicker.C:
		}
	}
}

func (s *myService) uploadSingle() {
	log.Printf("uploading...")
	time.Sleep(500 * time.Millisecond)
	log.Printf("upload done")
}
