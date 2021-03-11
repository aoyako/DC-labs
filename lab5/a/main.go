package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup
var changed int32

type Barrier struct {
	sync    *sync.Mutex
	size    int32
	counter int32
	wg      *sync.WaitGroup
}

func NewBarrier(n int) *Barrier {
	waiter := &sync.WaitGroup{}
	waiter.Add(n)

	return &Barrier{
		sync: &sync.Mutex{},
		size: int32(n),
		wg:   waiter,
	}
}

func (b *Barrier) Await() {
	b.sync.Lock()

	atomic.AddInt32(&b.counter, 1)
	if b.size == b.counter {
		// fmt.Println(b.counter)
		b.wg.Done()
		b.wg.Wait()

		for !atomic.CompareAndSwapInt32(&b.counter, 1, 0) {
		}

		// fmt.Println("added")
		b.wg.Add(int(b.size))

	} else {
		b.sync.Unlock()
		b.wg.Done()
		// time.Sleep(100 * time.Microsecond)
		b.wg.Wait()

		atomic.AddInt32(&b.counter, -1)
		return
	}

	b.sync.Unlock()
}

func Recruit(r []int, barrier *Barrier) {
	var tmp []int
	for i := range r {
		tmp = append(tmp, r[i])
	}

	for changed == 1 {
		barrier.Await()
		atomic.CompareAndSwapInt32(&changed, 1, 0)
		barrier.Await()

		firstChanged := false
		for i := range r {
			tmp[i] = r[i]
		}

		if (r[0] == 1) && (r[1] == -1) {
			tmp[0] = -1
			tmp[1] = 1
			firstChanged = true
			atomic.StoreInt32(&changed, 1)
		}

		barrier.Await()

		for i := range r[:len(r)-1] {
			if (r[i] == 1) && (r[i+1] == -1) {
				tmp[i] = -1
				tmp[i+1] = 1
				atomic.StoreInt32(&changed, 1)
			}
		}

		// barrier.Await()

		for i := range r[1:] {
			r[i] = tmp[i]
		}

		barrier.Await()

		if firstChanged {
			r[0] = tmp[0]
		}
	}
	wg.Done()
}

func main() {
	wg.Add(10)
	changed = 1
	barrier := NewBarrier(10)

	var arr []int
	for i := 0; i < 1000; i++ {
		arr = append(arr, ((rand.Intn(2)+2)%3)-1)
	}
	fmt.Println(arr)

	go Recruit(arr[0:100], barrier)
	go Recruit(arr[99:200], barrier)
	go Recruit(arr[199:300], barrier)
	go Recruit(arr[299:400], barrier)
	go Recruit(arr[399:500], barrier)
	go Recruit(arr[499:600], barrier)
	go Recruit(arr[599:700], barrier)
	go Recruit(arr[699:800], barrier)
	go Recruit(arr[799:900], barrier)
	go Recruit(arr[899:], barrier)

	wg.Wait()
	// time.Sleep(100 * time.Second)
	fmt.Println(arr)
}
