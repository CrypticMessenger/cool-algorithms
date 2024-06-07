package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WIDTH  = 320
	HEIGHT = 320
)

type Game struct {
	noise       [WIDTH][HEIGHT]float64
	permutation [256]uint8
}

func fade(t float64) float64 {
	return 6*math.Pow(t, 5) - 15*math.Pow(t, 4) + 10*math.Pow(t, 3)
}

func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func gradient(h uint8, x, y float64) float64 {
	var vectors [][]int = [][]int{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	var g []int = vectors[h%8]
	return float64(g[0])*x + float64(g[1])*y
}

func perlinNoise(x, y float64, permutation [256]uint8) float64 {
	var x0 int = int(math.Floor(x)) & 255
	var y0 int = int(math.Floor(y)) & 255
	var x1 int = (x0 + 1) & 255
	var y1 int = (y0 + 1) & 255

	var sx = x - math.Floor(x)
	var sy = y - math.Floor(y)

	var ii int = int(permutation[x0]) + y0
	var jj int = int(permutation[x1]) + y0
	var kk int = int(permutation[x0]) + y1
	var ll int = int(permutation[x1]) + y1

	// wrap around values of ii, jj, kk, ll to 0-255
	ii = ii & 255
	jj = jj & 255
	kk = kk & 255
	ll = ll & 255

	var g0 = gradient(permutation[ii], sx, sy)
	var g1 = gradient(permutation[jj], sx-1, sy)
	var g2 = gradient(permutation[kk], sx, sy-1)
	var g3 = gradient(permutation[ll], sx-1, sy-1)

	var u = fade(sx)
	var v = fade(sy)

	var nx0 = lerp(u, g0, g1)
	var nx1 = lerp(u, g2, g3)
	var nxy = lerp(v, nx0, nx1)

	return nxy

}

func (g *Game) Update() error {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	generate(&g.permutation, generator)
	for i := 0; i < WIDTH; i++ {
		for j := 0; j < HEIGHT; j++ {
			g.noise[i][j] = perlinNoise(float64(i)*0.05, float64(j)*0.05, g.permutation)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Find min and max
	minVal := g.noise[0][0]
	maxVal := g.noise[0][0]
	for i := range g.noise {
		for j := range g.noise[i] {
			if g.noise[i][j] < minVal {
				minVal = g.noise[i][j]
			}
			if g.noise[i][j] > maxVal {
				maxVal = g.noise[i][j]
			}
		}
	}

	// Normalize values
	for i := range g.noise {
		for j := range g.noise[i] {
			g.noise[i][j] = (g.noise[i][j] - minVal) / (maxVal - minVal)
		}
	}

	for i := 0; i < WIDTH; i++ {
		for j := 0; j < HEIGHT; j++ {
			value := g.noise[i][j]
			red := uint8(255)
			green := uint8(255 * (1 - value)) // Interpolate between 255 (white) and 0 (yellow)
			blue := uint8(0)
			color := color.RGBA{red, green, blue, 255}
			screen.Set(i, j, color)

		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func generate(permutation *[256]uint8, generator *rand.Rand) {
	for i := 0; i < 256; i++ {
		permutation[i] = uint8(i)
	}
	for i := 0; i < 256; i++ {
		j := generator.Intn(i + 1)
		permutation[i], permutation[j] = permutation[j], permutation[i]
	}
}

func main() {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Perlin noise")
	source := rand.NewSource(42)
	generator := rand.New(source)

	var noise [WIDTH][HEIGHT]float64
	var permutation [256]uint8
	for i := 0; i < WIDTH; i++ {
		for j := 0; j < HEIGHT; j++ {
			noise[i][j] = 0
			noise[i][j] = 0
		}
	}
	// arrange values 0 to 255 in random order
	generate(&permutation, generator)

	if err := ebiten.RunGame(&Game{noise: noise, permutation: permutation}); err != nil {
		panic(err)
	}
}
