package main

import (
	"fmt"
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
var barrier = NewBarrier(SIZE + 1)
var pBarrier []chan struct{}
var started bool
var wg sync.WaitGroup

type Game struct {
}

func getColor(i, j int) byte {
	var last int
	last += int(status[(i*SIZE+j)*4]) * 256 * 256
	last += int(status[(i*SIZE+j)*4+1]) * 256
	last += int(status[(i*SIZE+j)*4+2])

	switch last {
	case 0:
		return 0
	case 255:
		return 1
	case 255 * 256:
		return 2
	case 255 * 256 * 256:
		return 3
	case (255*256+255)*256 + 255:
		return 4
	}
	return 0
}

func getCount(i, j int) (byte, byte, byte, byte) {
	red := 0
	green := 0
	blue := 0
	white := 0

	// last := getColor(i, j)

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

			if getColor(i2, j2) == 1 {
				blue++
			}
			if getColor(i2, j2) == 2 {
				green++
			}
			if getColor(i2, j2) == 3 {
				red++
			}
			if getColor(i2, j2) == 4 {
				white++
			}

			// c += int(status[i2*SIZE+j2])
		}
	}
	return byte(white), byte(red), byte(green), byte(blue)
}

func updateStatus(i, j int) byte {
	prev := getColor(i, j)

	w, r, g, b := getCount(i, j)

	// if prev != 0 {
	// 	if r >= 2 && r <= 3 {
	// 		return 2
	// 	}
	// 	if g >= 2 && g <= 3 {
	// 		return 1
	// 	}
	// 	if w >= 2 && w <= 3 {
	// 		return 3
	// 	}

	// 	return 0
	// } else {
	// 	if nb == 3 {
	// 		return 0xff
	// 	}
	// }
	if prev != 0 {
		if prev == 4 && w >= 2 && w <= 3 {
			return 4
		}

		if prev == 3 && r >= 2 && r <= 3 {
			return 3
		}
		if prev == 2 && g >= 2 && g <= 3 {
			return 2
		}

		if prev == 1 && b >= 2 && b <= 3 {
			return 1
		}
	} else {
		if w == 3 {
			return 4
		}
		if r == 3 {
			return 3
		}
		if g == 3 {
			return 2
		}
		if b == 3 {
			return 1
		}
	}

	return 0
}

func updateLine(line int, screen *ebiten.Image) {
	wg.Wait()
	for {
		for j := 0; j < SIZE; j++ {
			nv := updateStatus(line, j)
			fill(newStatus, line, j, nv)
		}
		barrier.Await()

		barrier.Await()
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	// screen.Fill(color.White)

	screen.ReplacePixels(status)

	if !started {
		// screen.Fill(color.White)
		started = true
		for i := 0; i < SIZE; i++ {
			go updateLine(i, screen)
		}
		wg.Done()
	}

	barrier.Await()

	screen.ReplacePixels(newStatus)
	tmp := status
	status = newStatus
	newStatus = tmp

	// time.Sleep(1 * time.Millisecond)

	barrier.Await()

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SIZE, SIZE
}

func fill(arr []byte, i, j int, value byte) {
	if value == 4 {
		arr[(i*SIZE+j)*4] = 255
		arr[(i*SIZE+j)*4+1] = 255
		arr[(i*SIZE+j)*4+2] = 255
		arr[(i*SIZE+j)*4+3] = 255
	} else if value == 3 {
		arr[(i*SIZE+j)*4] = 255
		arr[(i*SIZE+j)*4+1] = 0
		arr[(i*SIZE+j)*4+2] = 0
		arr[(i*SIZE+j)*4+3] = 255
	} else if value == 2 {
		arr[(i*SIZE+j)*4] = 0
		arr[(i*SIZE+j)*4+1] = 255
		arr[(i*SIZE+j)*4+2] = 0
		arr[(i*SIZE+j)*4+3] = 255
	} else if value == 1 {
		arr[(i*SIZE+j)*4] = 0
		arr[(i*SIZE+j)*4+1] = 0
		arr[(i*SIZE+j)*4+2] = 255
		arr[(i*SIZE+j)*4+3] = 255
	} else {
		arr[(i*SIZE+j)*4] = value
		arr[(i*SIZE+j)*4+1] = value
		arr[(i*SIZE+j)*4+2] = value
		arr[(i*SIZE+j)*4+3] = value
	}
}

func main() {
	// squares = make([][]*ebiten.Image, SIZE)
	// positions = make([][]*ebiten.DrawImageOptions, SIZE)
	status = make([]byte, SIZE*SIZE*4)
	newStatus = make([]byte, SIZE*SIZE*4)
	pBarrier = make([]chan struct{}, SIZE)
	wg.Add(1)

	for i := 0; i < SIZE; i++ {
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
		var c byte

		_, err := fmt.Fscanf(file, "%d %d %d\n", &x, &y, &c)
		if err == io.EOF {
			break
		}

		fill(status, x, y, c)
	}

	// if err := ebiten.RunGame(&Game{}); err != nil {
	// 	panic(err)
	// }

	scale := 1.0
	if err := ebiten.Run((&Game{}).Update, SIZE, SIZE, scale, "Game of life"); err != nil {
		log.Fatal(err)
	}
}
