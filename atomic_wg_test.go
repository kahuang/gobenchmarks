package main

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

var result int64
var wgResult sync.WaitGroup

func BenchmarkAtomicAddSingle(b *testing.B) {
	a := int64(0)
	for n := 0; n < b.N; n++ {
		atomic.AddInt64(&a, int64(1))
	}
	for n := 0; n < b.N; n++ {
		atomic.AddInt64(&a, int64(-1))
	}
	result = a
}

func BenchmarkWaitGroup(b *testing.B) {
	wg := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		wg.Add(1)
	}
	for n := 0; n < b.N; n++ {
		wg.Done()
	}
	wgResult = wg
}

func BenchmarkAtomicAddMultiple(b *testing.B) {
	a := int64(0)
	numCPUs := runtime.NumCPU()
	nPerWorker := b.N / numCPUs
	runner := func(wg *sync.WaitGroup) {
		for n := 0; n < nPerWorker; n++ {
			atomic.AddInt64(&a, int64(1))
			atomic.AddInt64(&a, int64(-1))
		}
		result = a
		wg.Done()
	}
	wg := sync.WaitGroup{}
	wg.Add(numCPUs)
	for i := 0; i < numCPUs; i++ {
		go runner(&wg)
	}
	wg.Wait()
}

func BenchmarkWaitGroupMultiple(b *testing.B) {
	numCPUs := runtime.NumCPU()
	nPerWorker := b.N / numCPUs
	runner2 := func(wg *sync.WaitGroup) {
		for n := 0; n < nPerWorker; n++ {
			wg.Add(1)
		}
		for n := 0; n < nPerWorker; n++ {
			wg.Done()
		}
		wgResult = *wg
		wg.Done()
	}
	wg := sync.WaitGroup{}
	wg.Add(numCPUs)
	for i := 0; i < numCPUs; i++ {
		go runner2(&wg)
	}
	wg.Wait()
}
