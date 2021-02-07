package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"math/rand"
	"os"
	"time"
)

func run() {
	fmt.Printf("run\n")
	wd, _ := os.Getwd()
	fmt.Printf("wd: %v\n", wd)
	//spriteSheet, err := loadPicture("assets/Pixel_Mart/bacon.png")
	spriteSheet, err := loadPicture("/Users/rickard/code/Tetrigo/assets/emoji-3x.png")
	if err != nil {
		panic(err)
	}

	//sprite := pixel.NewSprite(spriteSheet, spriteSheet.Bounds())
	spriteRects := make([]pixel.Rect, 0, 7*5)
	sw := spriteSheet.Bounds().W() / 7
	sh := spriteSheet.Bounds().H() / 5
	x := 0.0
	y := 0.0
	for x+sw < spriteSheet.Bounds().W() + sw{
		for y+sh < spriteSheet.Bounds().H()  + sh{
			spriteRects = append(spriteRects, pixel.R(x, y, x+sw, y+sh))
			fmt.Printf("x: %f y: %f\n", x, y)
			y += sh
		}
		y = 0
		x += sw
	}

	cfg := pixelgl.WindowConfig{
		Title:  "TetriGo",
		Bounds: pixel.R(0, 0, 1024, 768),
		//VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	frames := 0
	seconds := time.Tick(time.Second)
	batch := pixel.NewBatch(&pixel.TrianglesData{}, spriteSheet)
	ticks := time.Tick(time.Millisecond * 500)
	var positions []pixel.Vec
	var ids []int
	for !win.Closed() {
		select {
		case <-ticks:
			const n = 10
			ids = make([]int, n)
			positions = make([]pixel.Vec, n)
			for i := 0; i < n; i++ {
				ids[i] = rand.Intn(len(spriteRects))

				positions[i] = pixel.V(rand.Float64()*win.Bounds().W(), rand.Float64()*win.Bounds().H())
			}
		default:
		}
		win.Clear(colornames.Darkgray)

		batch.Clear()
		for i := range positions {
			gopher := pixel.NewSprite(spriteSheet, spriteRects[ids[i]])
			mat := pixel.IM.Moved(positions[i])
			gopher.Draw(batch, mat)
		}

		batch.Draw(win)

		win.Update()
		frames++
		select {
		case <-seconds:
			win.SetTitle(fmt.Sprintf("FPS: %d", frames))
			frames = 0
		default:
		}
	}

}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func main() {
	rand.Seed(time.Now().Unix())
	pixelgl.Run(run)
}
