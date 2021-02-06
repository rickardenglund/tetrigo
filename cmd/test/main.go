package main

import (
	"Tetrigo/fonts"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/text"
	"github.com/faiface/pixel/imdraw"
	"image"
	_ "image/png"
	"os"
	"time"
)

func run() {
	fmt.Printf("run\n")
	wd, _ := os.Getwd()
	fmt.Printf("wd: %v\n", wd)
	//pic, err := loadPicture("assets/Pixel_Mart/bacon.png")
	pic, err := loadPicture("/Users/rickard/code/Tetrigo/assets/emoji-3x.png")
	if err != nil {
		panic(err)
	}

	spriteSizeX := pic.Bounds().W() / 7
	spriteSizeY := pic.Bounds().H() / 5
	sprite := pixel.NewSprite(pic, pixel.R(0, 0, spriteSizeX, spriteSizeY))

	cfg := pixelgl.WindowConfig{
		Title:  "TetriGo",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	angle := 0.0
	last := time.Now()
	//start := last

	gophers := []pixel.Vec{}
	n := 5
	dx := win.Bounds().W() / float64(n)
	dy := win.Bounds().H() / float64(n)
	for ix := 0; ix < n; ix++ {
		for iy := 0; iy < n; iy++ {
			gophers = append(gophers, pixel.V(float64(ix)*dx+0.5*dx, float64(iy)*dy+0.5*dy))
		}
	}

	seconds := time.Tick(time.Second)
	frames := 0
	fps := -1

	camPos := win.Bounds().Center()
	camSpeed := 500.0

	batch := pixel.NewBatch(&pixel.TrianglesData{}, pic)
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)
		//
		//angle += 3 * dt
		//size := math.Sin(time.Since(start).Seconds())
		//pos := win.MousePosition()
		font := fonts.GetFont()
		atlas := text.NewAtlas(
			font,
			text.ASCII,
		)
		size := 1.0

		win.Clear(colornames.Darkgray)
		boxpos := pixel.V(win.Bounds().H(), win.Bounds().Center().X)

		for i := range gophers {
			pos := gophers[i]
			mat := pixel.IM
			mat = mat.Moved(pos) //.Add(sprite.Frame().Center()))
			mat = mat.Rotated(pos, angle)
			mat = mat.Scaled(pos, size)
			sprite.Draw(batch, mat)
		}
		batch.Draw(win)

		// Draw box
		if true {
			//boxpos = win.Bounds().Center()
			const bw = 100
			imd := imdraw.New(sprite.Picture())
			//imd.SetColorMask(colornames.Green)
			imd.Color = pixel.RGB(1, 0, 0)
			imd.Push(pixel.V(boxpos.X, boxpos.Y))
			//imd.Push(pixel.V(boxpos.X+bw, boxpos.Y))
			imd.Color = pixel.RGB(0,1, 0)
			imd.Push(pixel.V(boxpos.X+bw, boxpos.Y-bw))
			//imd.Push(pixel.V(boxpos.X, boxpos.Y-bw))
			//imd.Polygon(0)
			imd.Rectangle(0)
			imd.Draw(win)
		} else {
			imd := imdraw.New(nil)
			imd.Color = pixel.RGB(1, 0, 0)
			imd.Push(pixel.V(200, 100))
			imd.Color = pixel.RGB(0, 1, 0)
			imd.Push(pixel.V(800, 100))
			imd.Color = pixel.RGB(0, 0, 1)
			imd.Push(pixel.V(500, 700))
			imd.Polygon(0)
			imd.Draw(win)
		}

		//win.Clear(colornames.Black)
		pos := cam.Unproject(pixel.V(100, 100))
		basicTxt := text.New(pos, atlas)
		basicTxt.Color = colornames.Purple
		fmt.Fprintf(basicTxt, "FPS: %d", fps)
		mat := pixel.IM
		//mat = mat.Scaled(pos, 6)

		basicTxt.Draw(win, mat)

		s := win.Typed()
		if len(s) > 0 {
			fmt.Printf("typed: %v\n", s)
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			gopherpos := cam.Unproject(win.MousePosition())
			gophers = append(gophers, gopherpos)
		}

		if win.JustPressed(pixelgl.KeySpace) {
			camPos = win.Bounds().Center()
		}

		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		win.Update()
		frames++
		select {
		case <-seconds:
			fps = frames
			win.SetTitle(fmt.Sprintf("FPS: %d", fps))
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
	pixelgl.Run(run)
}
