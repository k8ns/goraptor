package bg

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	StatusWaiting  = "waiting"
	StatusRunning  = "running"
	StatusStopping = "stopping"
	StatusStopped  = "stopped"
	StatusError    = "error"
)

type Counter interface {
	Inc()
}

type Process struct {
	Name       string
	Job        func() error
	Interval   time.Duration
	OkCounter  Counter
	ErrCounter Counter
	status     string
	stop       func()
	clock      *time.Ticker
	muProc     sync.Mutex
}

func (p *Process) Start(ctx context.Context, wg *sync.WaitGroup) {
	p.muProc.Lock()
	defer p.muProc.Unlock()

	if p.status != StatusStopped {
		return
	}
	p.status = StatusWaiting

	c, stop := context.WithCancel(ctx)
	p.stop = stop

	p.clock = time.NewTicker(p.Interval)

	log.Println("bg proc start", p.Name)
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.doJob()
		for {
			select {
			case <-c.Done():
				p.status = StatusStopping
				p.clock.Stop()

				p.status = StatusStopped
				log.Println("bg proc stopped", p.Name)
				return
			case <-p.clock.C:
				p.doJob()
			}
		}
	}()
}

func (p *Process) doJob() {
	if p.status == StatusWaiting {
		p.status = StatusRunning
		log.Printf("bg proc running %s", p.Name)
		err := p.Job()
		if err != nil {
			p.status = StatusError
			log.Println("bg proc error", p.Name, err)
			p.ErrCounter.Inc()
		} else {
			p.OkCounter.Inc()
		}
		p.status = StatusWaiting
		log.Printf("bg proc waiting %s", p.Name)
	}
}

func (p *Process) Stop() {
	p.stop()
}

func (p *Process) Status() string {
	return p.status
}

func (p *Process) String() string {
	return fmt.Sprintf("%s - %s", p.Name, p.status)
}
