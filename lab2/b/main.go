package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var sum int
var wg sync.WaitGroup

type Queue struct {
	lock int32
	pos  int
	arr  []int
}

func (q *Queue) put(value int) bool {
	for !atomic.CompareAndSwapInt32(&q.lock, 0, 1) {
	}

	if q.pos == len(q.arr)-1 {
		atomic.StoreInt32(&q.lock, 0)
		return false
	}

	q.pos++
	// fmt.Printf("put %d\n", value)
	q.arr[q.pos] = value

	atomic.StoreInt32(&q.lock, 0)

	return true
}

func (q *Queue) take() (int, bool) {
	for !atomic.CompareAndSwapInt32(&q.lock, 0, 1) {
	}

	if q.pos == -1 {
		atomic.StoreInt32(&q.lock, 0)
		return 0, false
	}

	value := q.arr[q.pos]
	// fmt.Printf("take %d\n", value)
	q.pos--

	atomic.StoreInt32(&q.lock, 0)

	return value, true
}

func NewQueue(size int) *Queue {
	return &Queue{
		arr: make([]int, size),
		pos: -1,
	}
}

func Ivanchuk(property []int, toPetrov *Queue) {
	for _, val := range property {
		for !toPetrov.put(val) {
		}

		fmt.Printf("Ivanchuk put %d\n", val)
		time.Sleep(10 * time.Millisecond)
	}

	for !toPetrov.put(-1) {
	}
}

func Petrov(toNechiporchuk *Queue, fromIvanchuk *Queue) {

	for {
		res, suc := fromIvanchuk.take()
		for !suc {
			res, suc = fromIvanchuk.take()
		}

		fmt.Printf("Petrov put %d\n", res)
		time.Sleep(10 * time.Millisecond)

		for !toNechiporchuk.put(res) {
		}

		if res == -1 {
			return
		}
	}
}

func Nechiporchuk(fromPetrov *Queue) {
	defer wg.Done()

	for {
		res, suc := fromPetrov.take()
		for !suc {
			res, suc = fromPetrov.take()
		}

		fmt.Printf("Nechiporchuk received %d\n", res)
		time.Sleep(10 * time.Millisecond)

		if res == -1 {
			return
		}

		sum += res
	}
}

func main() {
	wg.Add(1)
	property := []int{1, 2, 5, 1, 2}

	q1 := NewQueue(1)
	q2 := NewQueue(1)

	go Ivanchuk(property, q1)
	go Petrov(q2, q1)
	go Nechiporchuk(q2)

	wg.Wait()

	fmt.Println(sum)
}
