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
}

func newMyService() *myService {
	return &myService{
		loopTicker: time.NewTicker(1 * time.Second),
	}
}

func (s *myService) uploadForever() {
	for {
		s.uploadSingle()

		<-s.loopTicker.C
	}
}

func (s *myService) uploadSingle() {
	log.Printf("uploading...")
	time.Sleep(500 * time.Millisecond)
	log.Printf("upload done")
}

func (s *myService) onDestroy() {
	log.Printf("destroying...")

	log.Printf("destroyed")
}
