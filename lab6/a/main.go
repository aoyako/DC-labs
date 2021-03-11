package main

import (
	"fmt"
	"image/color"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten"
)

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

var status []byte
var newStatus []byte
var SIZE int = 1000
var barrier = NewBarrier(SIZE/rowsPerThread + 1)
var pBarrier []chan struct{}
var started bool
var wg sync.WaitGroup
var rowsPerThread int = 100

type Game struct {
}

func getCount(i, j int) int {
	c := 0
	for a := -1; a <= 1; a++ {
		for b := -1; b <= 1; b++ {
			if a == 0 && b == 0 {
				continue
			}
			i2 := i + a
			j2 := j + b
			if i2 < 0 || j2 < 0 || SIZE <= i2 || SIZE <= j2 {
				continue
			}

			if status[(i2*SIZE+j2)*4] != 0 {
				c++
			}

			// c += int(status[i2*SIZE+j2])
		}
	}
	return c
}

func updateStatus(i, j int) byte {
	prev := status[(i*SIZE+j)*4]

	nb := getCount(i, j)

	if prev == 0xff {
		if nb < 2 {
			return 0
		}
		if nb >= 2 && nb <= 3 {
			return 0xff
		}
		return 0
	} else {
		if nb == 3 {
			return 0xff
		}
		return 0
	}
}

func updateLine(line int, screen *ebiten.Image) {
	wg.Wait()
	for {
		for currLine := 0; currLine < rowsPerThread; currLine++ {

			normalizedLine := line*rowsPerThread + currLine
			for j := 0; j < SIZE; j++ {
				nv := updateStatus(normalizedLine, j)
				// fill(newStatus, line, j, nv)
				newStatus[(normalizedLine*SIZE+j)*4] = nv
				newStatus[(normalizedLine*SIZE+j)*4+1] = nv
				newStatus[(normalizedLine*SIZE+j)*4+2] = nv
				newStatus[(normalizedLine*SIZE+j)*4+3] = nv
			}
		}

		// if line != 0 {
		// 	<-pBarrier[line-1]
		// }
		// pBarrier[line] <- struct{}{}
		barrier.Await()

		// for j := 0; j < SIZE; j++ {
		// 	if status[line*SIZE+j] == 0 {
		// 		// squares[line][j].Fill(color.White)
		// 		newStatus[line*SIZE+j] = 0
		// 	} else {
		// 		// squares[line][j].Fill(color.Black)
		// 		newStatus[line*SIZE+j] = 0xff
		// 	}

		// screen.ReplacePixels(status)
		// screen.DrawImage(squares[line][j], positions[line][j])
		// status[line][j] = newStatus[line][j]
		// }

		// if line != 0 {
		// 	<-pBarrier[line-1]
		// }
		// pBarrier[line] <- struct{}{}

		barrier.Await()
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	// screen.Fill(color.White)

	screen.ReplacePixels(status)

	if !started {
		screen.Fill(color.White)
		started = true
		for i := 0; i < SIZE/rowsPerThread; i++ {
			go updateLine(i, screen)
		}
		wg.Done()
	}

	barrier.Await()

	tmp := status
	status = newStatus
	newStatus = tmp

	barrier.Await()

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SIZE, SIZE
}

func fill(arr []byte, i, j int, value byte) {
	status[(i*SIZE+j)*4] = value
	status[(i*SIZE+j)*4+1] = value
	status[(i*SIZE+j)*4+2] = value
	status[(i*SIZE+j)*4+3] = value
}

func main() {
	// squares = make([][]*ebiten.Image, SIZE)
	// positions = make([][]*ebiten.DrawImageOptions, SIZE)
	status = make([]byte, SIZE*SIZE*4)
	newStatus = make([]byte, SIZE*SIZE*4)
	pBarrier = make([]chan struct{}, SIZE/rowsPerThread)
	wg.Add(1)

	for i := 0; i < SIZE/rowsPerThread; i++ {
		// 	squares[i] = make([]*ebiten.Image, SIZE)
		// 	positions[i] = make([]*ebiten.DrawImageOptions, SIZE)
		// 	status[i] = make([]int, SIZE)
		// 	newStatus[i] = make([]int, SIZE)
		pBarrier[i] = make(chan struct{})

		// 	for j := range squares[i] {
		// 		squares[i][j], _ = ebiten.NewImage(5, 5, ebiten.FilterNearest)
		// 		positions[i][j] = &ebiten.DrawImageOptions{}
		// 		positions[i][j].GeoM.Translate(float64(i*5), float64(j*5))
		// 	}
	}

	// fill(status, SIZE/2, SIZE/2)
	// fill(status, SIZE/2-1, SIZE/2+1)
	// fill(status, SIZE/2, SIZE/2+2)
	// fill(status, SIZE/2+1, SIZE/2+1)

	file, _ := os.Open("field.txt")

	for {
		var x int
		var y int

		_, err := fmt.Fscanf(file, "%d %d\n", &x, &y)
		if err == io.EOF {
			break
		}

		fill(status, x, y, 0xff)
	}

	// if err := ebiten.RunGame(&Game{}); err != nil {
	// 	panic(err)
	// }

	scale := 1.0
	if err := ebiten.Run((&Game{}).Update, SIZE, SIZE, scale, "Game of life"); err != nil {
		log.Fatal(err)
	}
}
