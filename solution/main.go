package main

import (
	"log"
	"time"
)

func main() {
	mySvc := newMyService()

	go mySvc.runScheduledForever()

	time.Sleep(2250 * time.Millisecond)
	//time.Sleep(8000 * time.Millisecond)

	mySvc.onDestroy()

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

func (s *myService) runScheduledForever() {
	defer close(s.doneCh)

	for {
		if stopped := s.uploadBatch(); stopped {
			return
		}

		select {
		case <-s.stopCh:
			log.Printf("stop received while sleeping, terminate...")

			return
		case <-s.loopTicker.C:
		}
	}
}

func (s *myService) uploadBatch() (stopped bool) {
	log.Printf("batch upload start")

	for i := range 10 {
		select {
		case <-s.stopCh:
			log.Printf("stop received in the middle of batch upload, terminate...")

			return true
		default:
			log.Printf("upload file %d", i)
			time.Sleep(500 * time.Millisecond)
			log.Printf("upload file %d done", i)
		}
	}

	log.Printf("batch upload done")

	return false
}

func (s *myService) onDestroy() {
	log.Printf("destroying...")

	select {
	case s.stopCh <- struct{}{}:
		log.Printf("stop signal sent")
		<-s.doneCh
		log.Printf("done signal received")
	case <-time.After(10 * time.Second):
		log.Printf("couldn't stop file upload after 10 seconds, give up (timeout)")
	}

	log.Printf("destroyed")
}
