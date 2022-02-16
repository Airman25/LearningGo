package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"math/cmplx"

	"github.com/cpmech/gosl/utl"
	"github.com/hajimehoshi/ebiten/v2"
)

func plus_color(x int) uint8 {
	if x/2 > 155 {
		x = 155
	} else {
		x /= 2
	}
	return uint8(x)
}

func draw_mandelbrot(interval [4]float64, max_iterations int, m_border int) image.Image {
	amin, amax, bmin, bmax := interval[0], interval[1], interval[2], interval[3]
	apoints, bpoints := m_border+1, m_border+1
	center := math.Round(float64(apoints) / 2)
	yellow := [3]uint8{100, 100, 0}
	myImage := image.NewRGBA(image.Rect(0, 0, apoints, bpoints))
	for bi, b := range utl.LinSpace(bmin, bmax, bpoints) {
		for ai, a := range utl.LinSpace(amin, amax, apoints) {
			c := complex(a, b)
			z := complex(0, 0)
			outside_mandelbrot := true
			distance_from_center_to_point := cmplx.Abs(complex(float64(center-float64(ai)), float64(center-float64(bi))))
			for k := 0; k < max_iterations; k++ {
				z = z*z + c
				if cmplx.Abs(z) > float64(m_border) {
					tm := plus_color(int(distance_from_center_to_point))
					newcolor := [3]uint8{yellow[0] + tm, yellow[1] + tm, yellow[2]}
					myImage.Set(ai, bi, color.RGBA{newcolor[0], newcolor[1], newcolor[2], 0})
					outside_mandelbrot = false
					break
				}
			}
			if outside_mandelbrot {
				myImage.Set(ai, bi, color.RGBA{0, 0, 0, 0})
			}
		}
	}
	var reader = bytes.NewBuffer(nil)
	jpeg.Encode(reader, myImage, nil)
	imgtest, _ := jpeg.Decode(reader)
	fmt.Print(reader)
	return imgtest
}

var globalimg *ebiten.Image

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(globalimg, op)
	zoom()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 400, 400
}

var (
	click          = true
	secondclick    = false
	thirdclickable = false
	rightclick     = false
	savex          = 0.0
	savey          = 0.0

	stackarr = [][4]float64{}
	top      = 0
	newx     = 0.0
	newy     = 0.0
)

func main() {
	interv := [4]float64{-2, 2, -2, 2}
	stackarr = append(stackarr, interv)
	img := draw_mandelbrot(stackarr[top], 800, 400)
	globalimg = ebiten.NewImageFromImage(img)
	ebiten.SetWindowIcon([]image.Image{img})
	game := &Game{}
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Mandelbrot aplication")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func zoom() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if click {
			mx, my := ebiten.CursorPosition()
			for i, val := range utl.LinSpace(stackarr[top][0], stackarr[top][1], 101) {
				if i == mx/4 {
					newx = val
				}
			}
			for i, val := range utl.LinSpace(stackarr[top][2], stackarr[top][3], 101) {
				if i == my/4 {
					newy = val
				}
			}
			click = false
			if secondclick == false {
				savex = newx
				savey = newy
			} else {
				if savex > newx {
					savex, newx = newx, savex
				}
				if savey > newy {
					savey, newy = newy, savey
				}
				top++
				interv := [4]float64{savex, newx, savey, newy}
				stackarr = append(stackarr, interv)
				img := draw_mandelbrot(stackarr[top], 800, 400)
				globalimg = ebiten.NewImageFromImage(img)
				thirdclickable = true
			}
		}
	} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !secondclick && !click {
			secondclick = true
			click = true
		}
		if thirdclickable {
			click = true
			secondclick = false
			thirdclickable = false
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		rightclick = true
	} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if rightclick {
			rightclick = false
			if top > 0 {
				stackarr = stackarr[:top]
				top--
				img := draw_mandelbrot(stackarr[top], 800, 400)
				globalimg = ebiten.NewImageFromImage(img)
				click = true
				secondclick = false
				thirdclickable = false
			}
		}
	}
}
