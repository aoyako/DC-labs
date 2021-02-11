package main

import (
	"fmt"
	"sync"
	"time"
)

type Customer struct {
	ID     int
	finish chan struct{}
}

func (c *Customer) process() <-chan struct{} {
	reply := make(chan struct{})

	go func() {
		time.Sleep(10 * time.Millisecond)
		fmt.Printf("Customer %d sits\n", c.ID)

		reply <- struct{}{}
	}()

	return reply
}

func (c *Customer) walkToBarber(queue chan<- *Customer) {
	for {
		queue <- c
		<-c.finish

		time.Sleep(100 * time.Millisecond)
	}
}

func Barberer(queue <-chan *Customer) {
	for c := range queue {
		res := c.process()
		<-res

		fmt.Printf("Barber started %d\n", c.ID)
		time.Sleep(1000 * time.Millisecond)
		fmt.Printf("Barber finished %d\n", c.ID)

		c.finish <- struct{}{}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	c1 := Customer{finish: make(chan struct{}), ID: 1}
	c2 := Customer{finish: make(chan struct{}), ID: 2}
	c3 := Customer{finish: make(chan struct{}), ID: 3}
	c4 := Customer{finish: make(chan struct{}), ID: 4}

	queue := make(chan *Customer)
	go Barberer(queue)

	go c1.walkToBarber(queue)
	go c2.walkToBarber(queue)
	go c3.walkToBarber(queue)
	go c4.walkToBarber(queue)

	wg.Wait()
}
