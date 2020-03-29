package main

import (
	"./noise"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"math/rand"
	"time"
)

const (
	WinSizeW = 800
	WinSizeH = 800

	SizeW    = 400
	SizeH    = 400
	OneSize  = 3
	CellSize = 2
	OctavNum = 5
)

func Draw(imd *imdraw.IMDraw, Perlin2D noise.Perlin2D) {
	var n, x, y float64
	min, max := 1., 0.

	for i := 0; i < SizeW; i++ {
		for j := 0; j < SizeH; j++ {
			x, y = float64(i)*OneSize/SizeW, float64(j)*OneSize/SizeH
			n = Perlin2D.Noise(x, y, OctavNum, 0.5)/2 + 0.5
			if n < min {
				min = n
			}
			if n > max {
				max = n
			}
			if n < 0.6 {
				//imd.Color = pixel.RGB(n, n, 1)
				imd.Push(pixel.V(float64(i)*CellSize, float64(j)*CellSize), pixel.V(float64(i+1)*CellSize, float64(j+1)*CellSize))
				imd.Rectangle(0)
			}

		}
	}
	fmt.Println(min, max)
}

func run() {
	rand.Seed(time.Now().UnixNano())

	imd := imdraw.New(nil)

	Perlin2D := noise.Perlin2D{}
	Perlin2D.Init(time.Now().UnixNano(), SizeW+1, SizeH+1, OctavNum)
	Draw(imd, Perlin2D)

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Perlin noise",
		Bounds: pixel.R(0, 0, float64(WinSizeW), float64(WinSizeH)),
	})
	if err != nil {
		panic(err)
	}

	imd.Draw(win)

	for !win.Closed() {
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
