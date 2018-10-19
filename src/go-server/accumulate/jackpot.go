package accumulate

import (
	"sync"
)

type JackpotType struct {
	value float64
	lock  sync.RWMutex
}

func New() *JackpotType {
	return &JackpotType{
		value: 0.0,
	}
}

func (j *JackpotType) Add(v float64) {
	j.lock.Lock()
	defer j.lock.Unlock()
	j.value += v
}

func (j *JackpotType) Redeem() (res float64) {
	j.lock.Lock()
	defer j.lock.Unlock()
	res = j.value
	j.value = 0
	return
}

func (j *JackpotType) GetCurrent() float64 {
	j.lock.RLock()
	defer j.lock.RUnlock()
	return j.value
}
