package bonus

import (
	"errors"
	"sync"
)

type BonusRegistry struct {
	freeRides map[string]int
	lock      sync.RWMutex
}

func New() *BonusRegistry {
	return &BonusRegistry{
		freeRides: map[string]int{},
	}
}

func (b *BonusRegistry) Add(id string) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, exist := b.freeRides[id]; exist == true {
		b.freeRides[id]++
		return nil
	}
	b.freeRides[id] = 1
	return nil
}

func (b *BonusRegistry) Use(id string) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.freeRides[id] > 0 {
		b.freeRides[id]--
		return nil
	}
	return errors.New("No bonus games for this player")
}
