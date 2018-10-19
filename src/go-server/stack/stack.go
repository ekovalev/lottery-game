package stack

import (
	_ "encoding/binary"
	"errors"
	"math/rand"
	"sync"
	"time"
)

type BytePairStack struct {
	data           []uint16
	capacity       uint
	top            uint
	randomGen      *rand.Rand
	lock           sync.RWMutex
	timePoppedLast time.Time
}

func New() *BytePairStack {
	src := rand.NewSource(time.Now().UnixNano())
	return &BytePairStack{
		data:           make([]uint16, 100),
		capacity:       100,
		top:            0,
		randomGen:      rand.New(src),
		timePoppedLast: time.Now(),
	}
}

func (s *BytePairStack) Push(v uint16) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.top >= s.capacity {
		return errors.New("Stack is full")
	}
	s.data[s.top] = v
	s.top++
	return nil
}

func (s *BytePairStack) Pop() (uint16, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	defer func() {
		s.timePoppedLast = time.Now()
	}()

	if s.top == 0 {
		return 0, errors.New("Stack is empty")
	}
	s.top--
	return s.data[s.top], nil
}

func (s *BytePairStack) FillUp() {
	s.lock.Lock()
	defer s.lock.Unlock()
	for s.top < s.capacity {
		s.data[s.top] = uint16(s.randomGen.Intn(0xffff))
		s.top++
	}
}

func (s *BytePairStack) GetTimePoppedLast() time.Time {
	return s.timePoppedLast
}
