package main

import "sync"

type AtomicBool struct {
	mtx   sync.Mutex
	value bool
}

func NewAtomicBool() *AtomicBool {
	return &AtomicBool{
		sync.Mutex{},
		false,
	}
}

func (a *AtomicBool) Get() bool {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	return a.value
}

func (a *AtomicBool) Set(value bool) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	a.value = value
}
