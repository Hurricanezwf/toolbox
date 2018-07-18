package mutex

import (
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("key not found")
	ErrBadKey   = errors.New("bad key")
)

var mutex sync.Map

func Lock(key string) error {
	if len(key) <= 0 {
		return ErrBadKey
	}

	v, _ := mutex.LoadOrStore(key, &sync.RWMutex{})
	l := v.(*sync.RWMutex)
	l.Lock()
	return nil
}

func Unlock(key string) error {
	v, ok := mutex.Load(key)
	if !ok {
		return ErrNotFound
	}
	l := v.(*sync.RWMutex)
	l.Unlock()
	return nil
}

func RLock(key string) error {
	if len(key) <= 0 {
		return ErrBadKey
	}
	v, _ := mutex.LoadOrStore(key, &sync.RWMutex{})
	l := v.(*sync.RWMutex)
	l.RLock()
	return nil
}

func RUnlock(key string) error {
	v, ok := mutex.Load(key)
	if !ok {
		return ErrNotFound
	}
	l := v.(*sync.RWMutex)
	l.RUnlock()
	return nil
}

type TryLock struct {
	lock chan struct{}
}

func NewTryLock() TryLock {
	return TryLock{
		lock: make(chan struct{}, 1),
	}
}

// If lock successfully, true will be returned
// Or false will be returned
func (l *TryLock) Lock() bool {
	select {
	case l.lock <- struct{}{}:
		return true
	default:
		return false
	}
	return false
}

func (l *TryLock) Unlock() {
	select {
	case <-l.lock:
		return
	default:
		return
	}
	return
}
