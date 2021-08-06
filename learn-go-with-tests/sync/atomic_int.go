package main

import "sync"

type AtomicInt struct {
	mu  sync.RWMutex
	val int
}

func NewAtomicInt() *AtomicInt {
	return &AtomicInt{}
}

// Get 获取值
func (a *AtomicInt) Get() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.val
}

// IncrAndGet 修改值
func (a *AtomicInt) IncrAndGet(val int) int {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.val += +val
	return a.val
}
