package main

import "sync"

// AtomicBool is simple
// wrapper fro synchronized
// access to bool
type AtomicBool struct {
	mtx   sync.Mutex
	value bool
}

// NewAtomicBool returns pointer to fully
// initialized AtomicBool struct
func NewAtomicBool() *AtomicBool {
	return &AtomicBool{
		sync.Mutex{},
		false,
	}
}

// Get returns actual
// value of AtomicBool
func (a *AtomicBool) Get() bool {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	return a.value
}

// Set sets the bool value
// of AtomicBool
func (a *AtomicBool) Set(value bool) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	a.value = value
}
