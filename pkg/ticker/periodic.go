package ticker

import (
	"log"
	"sync"
	"time"
)

var (
	PeriodicExecution5min  = 60 * 6 * time.Second
	PeriodicExecution30sec = 5 * time.Second
)

type Handler func() error

type Scheduler struct {
	mx   sync.Mutex
	stop chan struct{}
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		mx:   sync.Mutex{},
		stop: make(chan struct{}),
	}
}

func (s *Scheduler) Run(name string, duration time.Duration, handler Handler) {

	log.Printf("added job for %q scheduler", name)
	go s.run(name, duration, handler)
}

func (s *Scheduler) run(name string, duration time.Duration, handler Handler) {
	ticker := time.NewTicker(duration)

	running := true
	for running {
		select {
		case <-s.stop:
			ticker.Stop()
			s.mx.Lock()
			running = false
			s.mx.Unlock()
		case <-ticker.C:
			log.Printf("executing %q scheduled job", name)
			err := handler()
			if err != nil {
				log.Printf("error while executing %q scheduled job: %q", name, err)
			}
			ticker.Stop()
			ticker = time.NewTicker(duration)
		}
	}

	log.Printf("ending %q scheduled job", name)
}

func (s *Scheduler) Stop() {
	s.stop <- struct{}{}
	close(s.stop)
	log.Println("scheduled job dispatcher closing")
}

func (s *Scheduler) Wait() {
	<-s.stop
}

func GetExtractionPeriod(name string) (time.Duration, bool) {
	switch name {
	case "30s":
		return PeriodicExecution30sec, true
	case "5m":
		return PeriodicExecution5min, true
	default:
		return 0, false
	}
}
