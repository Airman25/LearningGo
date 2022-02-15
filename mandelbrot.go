package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"math/cmplx"
	"os"

	"github.com/cpmech/gosl/utl"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func plus_yellow(x int) uint8 {
	if x/2 > 155 {
		x = 155
	} else {
		x /= 2
	}
	return uint8(x)
}

func minus_blue(x int) uint8 {
	if x*2 > 255 {
		x = 255
	} else {
		x *= 2
	}
	return uint8(x)
}

func draw_mandelbrot(interval [4]float64, i int, m int) {
	amin, amax, bmin, bmax := interval[0], interval[1], interval[2], interval[3]
	apoints, bpoints := m+1, m+1
	center := math.Round(float64(apoints) / 2)
	max_iterations := i
	m_border := m

	//blue := [3]uint8{0, 0, 255}
	yellow := [3]uint8{100, 100, 0}

	//plus_yellow = lambda x: 155 if x//2>155 else x//2
	myImage := image.NewRGBA(image.Rect(0, 0, apoints, bpoints))
	//var img [401][401][3]int

	for bi, b := range utl.LinSpace(bmin, bmax, bpoints) {
		for ai, a := range utl.LinSpace(amin, amax, apoints) {
			c := complex(a, b)
			z := complex(0, 0)
			outside_mandelbrot := true
			distance_from_center_to_point := cmplx.Abs(complex(float64(center-float64(ai)), float64(center-float64(bi))))

			for k := 0; k < max_iterations; k++ {
				z = z*z + c
				if cmplx.Abs(z) > float64(m_border) {
					tm := plus_yellow(int(distance_from_center_to_point))
					newcolor := [3]uint8{yellow[0] + tm, yellow[1] + tm, yellow[2]}
					//myImage.Pix[ai+bi*apoints]=uint8(float64(originalColor.R)*0.21 + float64(originalColor.G)*0.72 + float64(originalColor.B)*0.07)
					myImage.Set(ai, bi, color.RGBA{newcolor[0], newcolor[1], newcolor[2], 0})
					//img[ai][bi] = newcolor

					outside_mandelbrot = false
					break
				}
			}
			if outside_mandelbrot {

				//newcolor := [3]uint8{blue[0], blue[1], blue[2] - minus_blue(int(distance_from_center_to_point))}

				//img[ai][bi] = newcolor
				myImage.Set(ai, bi, color.RGBA{0, 0, 0, 0})
			}
		}
	}

	//myImage.At(1, 1).RGBA(col.A, col.B, col.G, col.R)
	//myImage.Pix[0] = blue
	//fmt.Print(myImage)

	o, _ := os.Create("out.jpeg")
	defer o.Close()
	writer := bufio.NewWriter(o)
	jpeg.Encode(writer, myImage, nil)
	writer.Flush()
	fmt.Println("DONE")
}

var globalimg *ebiten.Image

type Game struct{}

func (g *Game) Update() error {
	// Write your game's logical update.
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	//screen.Fill(color.White)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(globalimg, op)
	zoom()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 400, 400
}

func main() {
	interv := [4]float64{-2, 2, -2, 2}
	//interv := [4]float64{-1.4, -0.6, -0.1, 0.5}
	//interv := [4]float64{-1, 1, -1, 1}
	//draw_mandelbrot(interv, 800, 400)
	draw_mandelbrot(interv, 800, 400)
	globalimg, _, _ = Image("out.jpeg")
	game := &Game{}
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Mandelbrot aplication")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func Image(filename string) (*ebiten.Image, image.Image, error) {
	playImg, image, err := ebitenutil.NewImageFromFile(`D:\testVSCodeGo\` + filename)
	if err != nil {
		log.Fatalf("Load image error: %v", err)
		return nil, nil, err
	}
	return playImg, image, err
}

var (
	click          = true
	secondclick    = false
	thirdclickable = false
	savex          = 0.0
	savey          = 0.0

	dimensionx = 4.0
	dimensiony = 4.0
	decreasex  = 400.0 / dimensionx
	decreasey  = 400.0 / dimensiony
	leftx      = 0.0
	lefty      = 0.0
)

func zoom() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if click {
			mx, my := ebiten.CursorPosition()
			click = false
			newx := float64(mx)/decreasex - dimensionx/2 - leftx
			newy := float64(my)/decreasey - dimensiony/2 - lefty

			fmt.Println(newx, newy)
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
				interv := [4]float64{savex, newx, savey, newy}

				leftx = dimensionx/2 - (math.Abs(savex)+math.Abs(newx))/2
				lefty = dimensiony/2 - (math.Abs(savey)+math.Abs(newy))/2

				/*
					if leftx != 0 {
						newx -= (math.Abs(newx) + dimensionx/2) / 2
					}
					if lefty != 0 {
						newy -= (math.Abs(newy) + dimensiony/2) / 2
					}
				*/
				//leftx = 1
				//lefty = 1
				/*
					leftx = -savex
					lefty = -savey
				*/
				/*
					if savex < 0 {
						leftx = 1.0
					} else {
						leftx = -1.0
					}
					if savey < 0 {
						lefty = 1.0
					} else {
						lefty = -1.0
					}
				*/
				//-dimension/2 = -2.0
				//-diminsion/2 = -1.0
				//if savex = -2.0
				//-1.0 -dimesion/2 = -2.0
				//if savex = 0.0
				//-1.0 -dimension/2= 0.0
				//if savex = 2.0
				//1.0 - dimension/2= 2.0
				dimensionx = math.Abs(savex) + math.Abs(newx)
				dimensiony = math.Abs(savey) + math.Abs(newy)
				decreasex = 400.0 / dimensionx
				decreasey = 400.0 / dimensiony

				draw_mandelbrot(interv, 800, 400)
				globalimg, _, _ = Image("out.jpeg")
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
}
