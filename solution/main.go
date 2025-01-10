package main

import (
	"log"
	"time"
)

func main() {
	mySvc := newMyService()

	go mySvc.uploadForever()

	time.Sleep(2250 * time.Millisecond)

	mySvc.onDestroy()

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

func (s *myService) uploadForever() {
	for {
		s.uploadSingle()

		select {
		case <-s.stopCh:
			log.Printf("stop received, terminate...")

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

func (s *myService) onDestroy() {
	log.Printf("destroying...")

	select {
	case s.stopCh <- struct{}{}:
	case <-time.After(10 * time.Second):
		log.Printf("couldn't stop file upload after 10 seconds, give up (timeout)")
	}

	log.Printf("destroyed")
}
