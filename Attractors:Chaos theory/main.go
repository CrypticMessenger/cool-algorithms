package main

import (
	"chaos/pkg"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	attractor pkg.Attractor
}

// 0 -> WIDTH
// |
// v
// HEIGHT
func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.attractor.Draw(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.attractor.GetTotalWidth(), g.attractor.GetTotalHeight()
}

func main() {
	var g *Game = &Game{
		attractor: pkg.NewReccurencePlots(),
	}
	ebiten.SetWindowSize(g.attractor.GetTotalWidth(), g.attractor.GetTotalHeight())
	ebiten.SetWindowTitle("Chaos theory: Attractors")

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
