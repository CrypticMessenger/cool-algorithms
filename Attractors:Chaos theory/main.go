package main

import (
	"chaos/pkg"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	attractor pkg.Attractor
}

// 0 -> WIDTH |
// |   Board  |
// v		  |
// HEIGHT	  |
func (g *Game) Update() error {
	g.attractor.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.attractor.Draw(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.attractor.GetTotalWidth(), g.attractor.GetTotalHeight()
}

func f_xy(c complex128) float64 {
	var x float64 = real(c)
	var y float64 = imag(c)
	return 50 * x * math.Cos(y/50)
}
func g_xy(c complex128) float64 {
	var x float64 = real(c)
	var y float64 = imag(c)
	return 50 * y * math.Sin(x/50)
}

func main() {
	var g *Game = &Game{
		// attractor: pkg.NewReccurencePlots(),
		attractor: pkg.NewFieldSimulation(600, 600, 150, 40, 40, 100, f_xy, g_xy, 2, 1, false, color.RGBA{255, 153, 0, 255}, color.RGBA{58, 28, 68, 255}),
	}
	ebiten.SetWindowSize(g.attractor.GetTotalWidth(), g.attractor.GetTotalHeight())
	ebiten.SetWindowTitle("Chaos theory: Attractors")

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
