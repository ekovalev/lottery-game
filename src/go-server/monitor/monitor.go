package monitor

import (
	"time"

	"go-server/stack"

	log "github.com/sirupsen/logrus"
)

type Monitor struct {
	stack           *stack.BytePairStack
	popTimeout      time.Duration
	checkPopTimeout time.Duration
	refillTimeout   time.Duration
}

func New(st *stack.BytePairStack) *Monitor {
	return &Monitor{
		stack:           st,
		popTimeout:      10 * time.Second,
		checkPopTimeout: 1 * time.Second,
		refillTimeout:   60 * time.Second,
	}
}

func (m *Monitor) RemoveStalePair() {
	for {
		time.Sleep(m.checkPopTimeout)
		t := time.Now()
		if t.Sub(m.stack.GetTimePoppedLast()) >= m.popTimeout {
			log.WithFields(log.Fields{"time": t}).Info("[Monitor::RemoveStalePair] Popping a pair of bytes from the stack")
			_, _ = m.stack.Pop()
		}
	}
}

func (m *Monitor) RefillStack() {
	for {
		time.Sleep(m.refillTimeout)
		log.WithFields(log.Fields{"time": time.Now()}).Info("[Monitor::RefillStack] Refilling stack")
		m.stack.FillUp()
	}
}
