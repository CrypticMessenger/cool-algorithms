package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

// ref: https://karlsims.com/rd.html
const (
	width  int     = 320
	height int     = 240
	D_a    float64 = 0.21
	D_b    float64 = 0.11
	feed   float64 = 0.014
	k      float64 = 0.054
	dt     float64 = 1
)

var fav_feed []FavoriteValues = []FavoriteValues{
	FavoriteValues{0.0367, 0.0649, 1, 0.5}, // mitosis
	FavoriteValues{0.029, 0.057, 0.21, 0.13},
	FavoriteValues{0.026, 0.052, 0.26, 0.13}, // favorite values
	FavoriteValues{0.014, 0.054, 0.21, 0.11},
}

var pixels_t1 [width][height]Pixel
var pixels_t2 [width][height]Pixel
var initial_amt Pixel = Pixel{1, 0}
var laplacian_kernel [3][3]float64 = [3][3]float64{{0.05, 0.2, 0.05}, {0.2, -1, 0.2}, {0.05, 0.2, 0.05}}
var orientation Orientation = Orientation{width / 2, height / 2, 0, 0} // symmetric shape

func (g *Game) Update() error {
	swap(&pixels_t1, &pixels_t2)
	for i := 1; i < width-1; i++ {
		for j := 1; j < height-1; j++ {
			var a float64 = pixels_t1[i][j].conc_a
			var b float64 = pixels_t1[i][j].conc_b
			var laplacian_a float64 = laplacian(true, i, j)
			var laplacian_b float64 = laplacian(false, i, j)
			var newD_a float64 = D_a + math.Abs(float64(i-orientation.center_x))*orientation.dx //  inreases diffusion rate by dx per pixel going away from center
			var newD_b float64 = D_b + math.Abs(float64(j-orientation.center_y))*orientation.dy
			pixels_t2[i][j].conc_a = a + (newD_a*laplacian_a-a*b*b+feed*(1-a))*dt
			pixels_t2[i][j].conc_b = b + (newD_b*laplacian_b+a*b*b-(k+feed)*b)*dt
			clip(&pixels_t2[i][j].conc_a, 0, 1)
			clip(&pixels_t2[i][j].conc_b, 0, 1)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := 1; i < width-1; i++ {
		for j := 1; j < height-1; j++ {
			var a float64 = pixels_t2[i][j].conc_a
			var b float64 = pixels_t2[i][j].conc_b
			var temp_c int = int(math.Floor((a - b) * 255))
			clip(&temp_c, 0, 255)
			var c uint8 = uint8(temp_c)
			screen.Set(i, j, color.RGBA{c, c, c, 255})
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Reaction Diffusion Algorithm")
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			pixels_t1[i][j] = initial_amt
			pixels_t2[i][j] = initial_amt
		}
	}
	var side int = 20
	for i := (width / 2) - side; i < (width/2)+side+1; i++ {
		for j := (height / 2) - side; j < (height/2)+side+1; j++ {
			pixels_t1[i][j].conc_b = 1
			pixels_t1[i][j].conc_a = 0
			pixels_t2[i][j].conc_b = 1
			pixels_t2[i][j].conc_a = 0
		}
	}

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
