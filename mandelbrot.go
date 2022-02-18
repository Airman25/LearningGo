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

	"github.com/cpmech/gosl/utl"       //used because of LinSpace function
	"github.com/hajimehoshi/ebiten/v2" //used to render image and make interactions with mouse possible
)

func plus_color(x int) uint8 { //making gradient
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
					myImage.Set(ai, bi, color.RGBA{yellow[0] + tm, yellow[1] + tm, yellow[2], 0})
					outside_mandelbrot = false
					break
				}
			}
			if outside_mandelbrot {
				myImage.Set(ai, bi, color.RGBA{0, 0, 0, 0})
			}
		}
	}
	//Ebiten won't render *image.RGBA so I convert it to image.Image
	var reader = bytes.NewBuffer(nil)
	jpeg.Encode(reader, myImage, nil)
	imgtest, _ := jpeg.Decode(reader)
	fmt.Print(reader)
	return imgtest
}

var globalimg *ebiten.Image //mandelbrot image to render

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(globalimg, op) //rendering image on screen
	zoom()                          //zoom in the image
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 400, 400
}

var (
	click          = true           //is it your first click
	secondclick    = false          //is it your second click
	thirdclickable = false          //is it possible to click again, used to prevent trying to zoom while generating image with new interval
	rightclick     = false          //right clicked
	savex          = 0.0            //previous x coordinate
	savey          = 0.0            //         y
	stackarr       = [][4]float64{} //array of intervals
	top            = 0              //the number of the current used interval
	newx           = 0.0            //new x coordinate for interval
	newy           = 0.0            //    y
)

func main() {
	interv := [4]float64{-2, 2, -2, 2} //the interval int he beginning
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
			newx = utl.LinSpace(stackarr[top][0], stackarr[top][1], 101)[mx/4] //converting coordinates clicked to coordinate for interval
			newy = utl.LinSpace(stackarr[top][2], stackarr[top][3], 101)[my/4]
			click = false
			if secondclick == false {
				savex = newx
				savey = newy
			} else {
				if savex > newx { //preventing image from being rotated
					savex, newx = newx, savex
				}
				if savey > newy {
					savey, newy = newy, savey
				}
				top++
				interv := [4]float64{savex, newx, savey, newy}
				stackarr = append(stackarr, interv)
				globalimg = ebiten.NewImageFromImage(draw_mandelbrot(stackarr[top], 800, 400))
				thirdclickable = true
			}
		}
	} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !secondclick && !click { //when it is second click
			secondclick = true
			click = true
		}
		if thirdclickable { //when it is clickable again
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
				globalimg = ebiten.NewImageFromImage(draw_mandelbrot(stackarr[top], 800, 400))
				click = true //setting that it would count as first click even if it was clicked before
				secondclick = false
				thirdclickable = false
			}
		}
	}
}
