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

		mySvc.runScheduledForever(ctx)
	}()

	//time.Sleep(2250 * time.Millisecond)
	time.Sleep(8000 * time.Millisecond)

	cancel() // cancel the context, shutdown the thing

	log.Printf("wait for goroutine to finish...")

	wg.Wait()

	log.Printf("main done")
}

type myService struct {
	loopTicker *time.Ticker

	stopCh chan struct{}
	doneCh chan struct{}
}

func newMyService() *myService {
	return &myService{
		loopTicker: time.NewTicker(30 * time.Second),
		stopCh:     make(chan struct{}),
		doneCh:     make(chan struct{}),
	}
}

func (s *myService) runScheduledForever(ctx context.Context) {
	defer close(s.doneCh)

	for {
		s.uploadBatch(ctx)

		select {
		case <-ctx.Done():
			log.Printf("context done while sleeping, terminate...")

			return
		case <-s.loopTicker.C:
		}
	}
}

func (s *myService) uploadBatch(ctx context.Context) {
	log.Printf("batch upload start")

	for i := range 10 {
		select {
		case <-ctx.Done():
			log.Printf("context done in the middle of batch upload, terminate...")

			return
		default:
			log.Printf("upload file %d", i)
			time.Sleep(500 * time.Millisecond)
			log.Printf("upload file %d done", i)
		}
	}

	log.Printf("batch upload done")
}
