package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func hitDuck(d *Duck, x, y int) bool {
	if x > d.x && x < d.x+d.frameWidth && y > d.y && y < d.y+d.frameHeight {
		return true
	}

	return false
}

type Duck struct {
	frameFlyHorizontalOX int
	frameFlyHorizontalOY int
	framesHorizontal     int

	frameFlyDiagonalOX int
	frameFlyDiagonalOY int
	framesDiagonal     int

	frameVertialOX int
	frameVertialOY int
	framesVertical int

	frameDieOX int
	frameDieOY int
	framesDie  int

	frameFallOX int
	frameFallOY int
	framesFall  int

	frameWidth  int
	frameHeight int

	currentFrame int
	x            int
	y            int
	direction    int

	isFree bool
	isDead bool
}

var (
	BlueDuck = Duck{
		frameFlyHorizontalOY: 118,
		frameFlyDiagonalOY:   155,
		frameVertialOY:       192,
		frameDieOY:           229,
		frameFallOY:          229,
		frameFallOX:          38,

		frameWidth:       38,
		frameHeight:      37,
		framesHorizontal: 3,
		framesDiagonal:   3,
		framesVertical:   3,
		framesDie:        1,
		framesFall:       1,
	}
)

const (
	screenWidth  = 500
	screenHeight = 500

	frameOX     = 0
	frameOY     = 118
	frameWidth  = 38
	frameHeight = 37
	frameNum    = 3
)

var duckStarted bool = false
var objects []chan struct{}
var ducks []*Duck
var nestWaiter chan struct{}

var (
	runnerImage     *ebiten.Image
	backgroundImage *ebiten.Image
)
var score int
var ammo int

func nest(screen *ebiten.Image, connector chan struct{}) {
	var timer int
	for {
		<-connector

		for i := range objects {
			if ducks[i].isFree {
				continue
			}
			objects[i] <- struct{}{}
		}

		for i := range objects {
			if ducks[i].isFree {
				continue
			}
			<-objects[i]
		}

		deleted := true
		for deleted {
			deleted = false
			for i := range objects {
				if ducks[i].isFree {
					deleted = true

					ducks[i] = ducks[len(ducks)-1]
					ducks = ducks[:len(ducks)-1]

					objects[i] = objects[len(objects)-1]
					objects = objects[:len(objects)-1]

					break
				}
			}
		}

		if timer == 40 {
			duck := BlueDuck
			ducks = append(ducks, &duck)
			objects = append(objects, make(chan struct{}))

			go duckFly(screen, ducks[len(ducks)-1], objects[len(objects)-1])

			timer = 0
		}

		if ammo < 2 {
			if timer%30 == 0 {
				ammo++
			}
		}

		timer++

		connector <- struct{}{}
	}
}

func duckFly(screen *ebiten.Image, duck *Duck, connector chan struct{}) {
	duck.direction = 3
	duck.y = 380
	duck.x = screenWidth/4 + (screenWidth/4)*rand.Intn(3)

	speed := 1 + rand.Intn(3)

	for {
		<-connector

		swich := rand.Intn(3)
		if swich == 0 {
			if Abs(duck.direction) == 1 {
				duck.direction = 2 * duck.direction
			} else if Abs(duck.direction) == 2 {
				if rand.Intn(2) == 1 {
					duck.direction = 3 * duck.direction / 2
				} else {
					duck.direction = duck.direction / 2
				}
			} else if Abs(duck.direction) == 3 {
				if rand.Intn(2) == 1 {
					duck.direction = -2 * duck.direction / 3
				} else {
					duck.direction = 2 * duck.direction / 3
				}
			}
		}

		op := &ebiten.DrawImageOptions{}

		if duck.direction < 0 {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(frameWidth, 0)
		}

		op.GeoM.Translate(float64(duck.x), float64(duck.y))

		if Abs(duck.direction) == 1 {
			i := (duck.currentFrame) % duck.framesHorizontal
			sx, sy := duck.frameFlyHorizontalOX+i*duck.frameWidth, duck.frameFlyHorizontalOY
			screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

			duck.x += duck.direction * 5 * speed
			duck.y -= speed

		} else if Abs(duck.direction) == 2 {
			i := (duck.currentFrame) % duck.framesDiagonal
			sx, sy := duck.frameFlyDiagonalOX+i*duck.frameWidth, duck.frameFlyDiagonalOY
			screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

			duck.x += duck.direction * speed
			duck.y -= 2 * speed
		} else if Abs(duck.direction) == 3 {
			i := (duck.currentFrame) % duck.framesVertical
			sx, sy := duck.frameVertialOX+i*duck.frameWidth, duck.frameVertialOY
			screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

			// duck.x += duck.direction * 5 / 2
			duck.y -= 3 * speed
		} else if Abs(duck.direction) == 0 {
			i := (duck.currentFrame) % duck.framesDie
			sx, sy := duck.frameDieOX+i*duck.frameWidth, duck.frameDieOY
			screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

			duck.direction = 10
		} else if Abs(duck.direction) == 10 {
			i := (duck.currentFrame) % duck.framesFall
			sx, sy := duck.frameFallOX+i*duck.frameWidth, duck.frameFallOY
			screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

			duck.y += 15
			duck.direction *= -1
		}

		duck.currentFrame++
		connector <- struct{}{}

		if (duck.x < 0) || duck.x > screenWidth || duck.y < 0 || duck.y > 380 {
			duck.isFree = true
			return
		}

		// i := (duck. / 10) % frameNum
		// sx, sy := frameOX+i*frameWidth, frameOY
		// screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
	}
}

type Game struct {
	count int
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !duckStarted {
		nestWaiter = make(chan struct{})
		go nest(screen, nestWaiter)

		duckStarted = true
	}

	mx, my := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if ammo > 0 {
			ammo--
			for i := range ducks {
				if ducks[i].isFree || ducks[i].isDead {
					continue
				}
				if hitDuck(ducks[i], mx, my) {
					ducks[i].direction = 0
					ducks[i].isDead = true
					score++
					break
				}
			}
		}
	}

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(backgroundImage.SubImage(image.Rect(220, 0, 720, 500)).(*ebiten.Image), op)

	nestWaiter <- struct{}{}

	time.Sleep(100 * time.Millisecond)

	<-nestWaiter

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Score: %d\nAmmo: %d", score, ammo))
	// screen.Clear()

	// op.GeoM.Translate(screenWidth/2, screenHeight/2)
	// i := (g.count / 10) % frameNum
	// sx, sy := frameOX+i*frameWidth, frameOY
	// screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	sprites, err := os.Open("../duckhunt_various_sheet.png")
	if err != nil {
		log.Fatal(err)
	}

	bg, err := os.Open("../bg.png")
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(sprites)
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(img)

	bgi, _, err := image.Decode(bg)
	if err != nil {
		log.Fatal(err)
	}
	backgroundImage = ebiten.NewImageFromImage(bgi)

	scale := 2
	ebiten.SetWindowSize(scale*screenWidth, scale*screenHeight)
	ebiten.SetWindowTitle("Duckhunt")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
