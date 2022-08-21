package utils

import "sync"

type Locker struct {
	mutexes sync.Map // Zero value is empty and ready for use
}

func (m *Locker) Lock(key string) func() {
	value, _ := m.mutexes.LoadOrStore(key, &sync.Mutex{})

	mtx, ok := value.(*sync.Mutex)
	if ok {
		mtx.Lock()

		return func() { mtx.Unlock() }
	}

	return func() {}
}
