package main

import (
	"sync"
	"testing"
)

func TestAtomicInt(t *testing.T) {
	t.Run("同步修改值", func(t *testing.T) {
		atomicInt := NewAtomicInt()
		atomicInt.IncrAndGet(1)
		atomicInt.IncrAndGet(1)
		atomicInt.IncrAndGet(1)

		want := 3
		got := atomicInt.Get()

		assertEquals(t, got, want)
	})

	t.Run("并发修改值", func(t *testing.T) {
		atomicInt := NewAtomicInt()
		want := 100

		// 闭锁
		wg := sync.WaitGroup{}
		wg.Add(want)

		for i := 0; i < want; i++ {
			go func() {
				atomicInt.IncrAndGet(1)

				wg.Done()
			}()
		}

		wg.Wait()

		got := atomicInt.Get()

		if got != want {
			t.Errorf("want: %v, got: %v", want, got)
		}
	})
}

func assertEquals(t *testing.T, got int, want int) {
	if got != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
